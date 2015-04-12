package bencoding

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// Marshal takes the given Go datastructure and converts it to a bencoded string.
// See https://wiki.theory.org/BitTorrentSpecification#Bencoding for details about bencoding.
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
				marshalledFieldName, err := Marshal(ToLowerCaseWithSpaces(fieldName))
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
