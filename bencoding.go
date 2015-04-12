package gotorrent

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Marshal(v interface{}) (string, error) {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Int:
		i := v.(int)
		return "i" + strconv.Itoa(i) + "e", nil
	case reflect.String:
		s := v.(string)
		return strconv.Itoa(len(s)) + ":" + s, nil
	case reflect.Array, reflect.Slice:
		l := "l"
		for i := 0; i < value.Len(); i++ {
			s, err := Marshal(value.Index(i).Interface())
			if err != nil {
				return "", err
			}
			l += s
		}
		l += "e"
		return l, nil
	case reflect.Struct:
		l := "l"
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if field.CanInterface() {
				s, err := Marshal(field.Interface())
				if err != nil {
					return "", err
				}
				l += s
			}
		}
		l += "e"
		return l, nil
	case reflect.Map:
		d := "d"
		// TODO(apm): The iteration order is nondeterministic.  Ugh.
		for _, key := range value.MapKeys() {
			ks, err := Marshal(key.Interface())
			if err != nil {
				return "", err
			}
			vs, err := Marshal(value.MapIndex(key).Interface())
			if err != nil {
				return "", err
			}
			d += ks + vs
		}
		d += "e"
		return d, nil
	default:
		return "", fmt.Errorf("Can't marshal type %v, value %v", value.Kind(), v)
	}
}

// TODO(apm): Lots of string copies in here, would be an easy performance boost.
func Unmarshal(s string, v interface{}) error {
	ptr_value := reflect.ValueOf(v)
	switch ptr_value.Kind() {
	case reflect.Interface, reflect.Ptr:
		break
	default:
		return fmt.Errorf("Must pass a pointer or struct to Unmarshal, received %v", ptr_value)
	}

	value := ptr_value.Elem()
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
		for i := 0; i < value.Len(); i++ {
			token, leftovers, err := GetOneToken(s)
			if err != nil {
				return fmt.Errorf("Unable to tokenize string %v: err %v", s, err)
			}
			s = leftovers
			err = Unmarshal(token, value.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}
		if s != "" {
			return fmt.Errorf("Unconsumed inputs in %v when unmarshaling to %v", s, v)
		}
		return nil
	default:
		return fmt.Errorf("Can't unmarshal to type %v, value %v", value.Kind(), v)
	}
}

func GetOneToken(s string) (token, leftovers string, err error) {
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
			subtoken, new_leftovers, err := GetOneToken(leftovers)
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
