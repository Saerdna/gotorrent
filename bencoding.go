package gotorrent

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Marshal takes the given Go datastructure and converts it to a bencoded string.
func Marshal(source interface{}) (bencoded_string string, err error) {
	value := reflect.ValueOf(source)
	switch value.Kind() {
	case reflect.Int:
		i := source.(int)
		return "i" + strconv.Itoa(i) + "e", nil
	case reflect.String:
		s := source.(string)
		return strconv.Itoa(len(s)) + ":" + s, nil
	case reflect.Array, reflect.Slice:
		listString := "l"
		for i := 0; i < value.Len(); i++ {
			token, err := Marshal(value.Index(i).Interface())
			if err != nil {
				return "", err
			}
			listString += token
		}
		listString += "e"
		return listString, nil
	case reflect.Struct:
		dictString := "d"
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if field.CanInterface() {
				fieldName := value.Type().Field(i).Name
				marshalledFieldName, err := Marshal(lowerCaseWithSpaces(fieldName))
				if err != nil {
					return "", err
				}
				dictString += marshalledFieldName

				marshalledValue, err := Marshal(field.Interface())
				if err != nil {
					return "", err
				}
				dictString += marshalledValue
			}
		}
		dictString += "e"
		return dictString, nil
	case reflect.Map:
		marshalledMap := map[string]string{}
		marshalledKeys := []string{}
		for _, keyValue := range value.MapKeys() {
			marshalledKey, err := Marshal(keyValue.Interface())
			if err != nil {
				return "", err
			}
			marshalledKeys = append(marshalledKeys, marshalledKey)
			marshalledValue, err := Marshal(value.MapIndex(keyValue).Interface())
			if err != nil {
				return "", err
			}
			marshalledMap[marshalledKey] = marshalledValue
		}
		sort.Strings(marshalledKeys)
		dictString := "d"
		for _, marshalledKey := range marshalledKeys {
			dictString += marshalledKey + marshalledMap[marshalledKey]
		}
		dictString += "e"
		return dictString, nil
	default:
		return "", fmt.Errorf("Can't marshal type %v, value %v", value.Kind(), source)
	}
}

// Unmarshal takes a bencoded string and a target object, and fills out the target object
// with the values from the bencoded string.  The structure of the target object must match
// the structure of the string.  Slices will be automatically sized.
// TODO(apm): Lots of string copies in here, look into optimizations.
func Unmarshal(s string, v interface{}) error {
	ptrValue := reflect.ValueOf(v)
	switch ptrValue.Kind() {
	case reflect.Interface, reflect.Ptr:
		break
	default:
		return fmt.Errorf("Must pass a pointer or struct to Unmarshal, received %v", ptrValue)
	}

	value := ptrValue.Elem()
	if !value.CanSet() {
		return fmt.Errorf("Received unsettable value %v", v)
	}
	switch value.Kind() {
	case reflect.Int:
		if s[0] != 'i' || s[len(s)-1] != 'e' {
			return fmt.Errorf("Expected integer for %v, found %v", v, s)
		}
		s = s[1 : len(s)-1]
		i, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Error parsing %v as int: %v", s, err)
		}
		value.SetInt(int64(i))
		return nil
	case reflect.String:
		substrings := strings.SplitN(s, ":", 2)
		if len(substrings) < 2 {
			return fmt.Errorf("Invalid string, missing colon: %v", s)
		}
		length, err := strconv.Atoi(substrings[0])
		if err != nil {
			return fmt.Errorf("Unable to parse length for string %v: err %v", s, err)
		}
		if len(substrings[1]) != length {
			return fmt.Errorf("Invalid length specification for string %v", s)
		}
		value.SetString(substrings[1])
		return nil
	case reflect.Array, reflect.Slice:
		if s[0] != 'l' || s[len(s)-1] != 'e' {
			return fmt.Errorf("Expected list for %v, found %v", v, s)
		}
		s = s[1 : len(s)-1]
		tokens := []string{}
		for len(s) > 0 {
			token, leftovers, err := getOneToken(s)
			if err != nil {
				return fmt.Errorf("Unable to tokenize string %v: err %v", s, err)
			}
			s = leftovers
			tokens = append(tokens, token)
		}

		if value.Len() != len(tokens) {
			if value.Kind() == reflect.Array {
				return fmt.Errorf("Length mismatch on array %v: tokens %v", value, tokens)
			}
			value.Set(reflect.MakeSlice(value.Type(), len(tokens), len(tokens)))
		}

		for i := 0; i < value.Len(); i++ {
			err := Unmarshal(tokens[i], value.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}
		if s != "" {
			return fmt.Errorf("Unconsumed inputs in %v when unmarshaling to %v, %v", s, value, v)
		}
		return nil
	case reflect.Struct:
		if s[0] != 'd' || s[len(s)-1] != 'e' {
			return fmt.Errorf("Expected dict for %v, found %v", v, s)
		}
		s = s[1 : len(s)-1]
		for len(s) > 0 {
			token, leftovers, err := getOneToken(s)
			if err != nil {
				return fmt.Errorf("Unable to tokenize string %v: err %v", s, err)
			}
			s = leftovers
			var unmarshalledFieldName string
			err = Unmarshal(token, &unmarshalledFieldName)
			if err != nil {
				return fmt.Errorf("Unable to unmarshall field name %v: %v", token, err)
			}
			fieldName := camelCase(unmarshalledFieldName)
			field := value.FieldByName(fieldName)
			if !field.IsValid() {
				return fmt.Errorf("Struct contains no field named %v", fieldName)
			}
			if !field.CanSet() {
				return fmt.Errorf("Dict contained value for unsettable field %v", token)
			}

			token, leftovers, err = getOneToken(s)
			if err != nil {
				return fmt.Errorf("Unable to tokenize string %v: err %v", s, err)
			}
			s = leftovers
			err = Unmarshal(token, field.Addr().Interface())
			if err != nil {
				return err
			}
		}
		if s != "" {
			return fmt.Errorf("Unconsumed inputs in %v when unmarshaling to %v, %v", s, value, v)
		}
		return nil
	case reflect.Map:
		if s[0] != 'd' || s[len(s)-1] != 'e' {
			return fmt.Errorf("Expected map for %v, found %v", v, s)
		}
		s = s[1 : len(s)-1]
		for len(s) > 0 {
			key, leftovers, err := getOneToken(s)
			if err != nil {
				return fmt.Errorf("Unable to tokenize string %v: err %v", s, err)
			}
			s = leftovers
			key_val := reflect.New(value.Type().Key())
			err = Unmarshal(key, key_val.Interface())
			if err != nil {
				return err
			}

			elem, leftovers, err := getOneToken(s)
			if err != nil {
				return fmt.Errorf("Unable to tokenize string %v: err %v", s, err)
			}
			s = leftovers
			elem_val := reflect.New(value.Type().Elem())
			err = Unmarshal(elem, elem_val.Interface())
			if err != nil {
				return err
			}

			value.SetMapIndex(key_val.Elem(), elem_val.Elem())
		}
		return nil
	default:
		return fmt.Errorf("Can't unmarshal to type %v, value %v", value.Kind(), v)
	}
}

// Converts a CamelCase string to a lower case string with spaces
func lowerCaseWithSpaces(s string) (lower string) {
	for _, r := range s {
		if unicode.IsLower(r) {
			lower += string(r)
		} else if unicode.IsUpper(r) {
			lower += " " + string(unicode.ToLower(r))
		}
	}
	return strings.TrimSpace(lower)
}

func camelCase(s string) (camel string) {
	previousRuneWasSpace := true
	for _, r := range s {
		if previousRuneWasSpace {
			camel += string(unicode.ToUpper(r))
			previousRuneWasSpace = false
		} else if unicode.IsSpace(r) {
			previousRuneWasSpace = true
		} else {
			camel += string(r)
		}
	}
	return camel
}

// TODO(apm): This would be a lot cleaner if we built a syntax tree.
func getOneToken(s string) (token, leftovers string, err error) {
	switch s[0] {
	case 'i':
		substrings := strings.SplitAfterN(s, "e", 2)
		if len(substrings) < 2 {
			return "", "", fmt.Errorf("No termination for leading integer in %v", s)
		}
		return substrings[0], substrings[1], nil
	case 'l', 'd':
		token := s[0:1]
		leftovers := s[1:]
		for len(leftovers) > 0 {
			if leftovers[0] == 'e' {
				token += leftovers[0:1]
				leftovers = leftovers[1:]
				return token, leftovers, nil
			}
			subtoken, new_leftovers, err := getOneToken(leftovers)
			if err != nil {
				return "", "", fmt.Errorf("Subtoken error: %v", err)
			}
			token += subtoken
			leftovers = new_leftovers
		}
		return "", "", fmt.Errorf("No termination for token in %v", s)
	default:
		colonIndex := strings.Index(s, ":")
		if colonIndex < 0 {
			return "", "", fmt.Errorf("Couldn't parse token in %v", s)
		}
		length, err := strconv.Atoi(s[:colonIndex])
		if err != nil {
			return "", "", fmt.Errorf("Couldn't parse length of string %v: %v", s, err)
		}
		tokenLength := colonIndex + length + 1
		return s[:tokenLength], s[tokenLength:], nil
	}
}
