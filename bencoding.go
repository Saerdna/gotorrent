package gotorrent

import (
	"fmt"
	"reflect"
	"strconv"
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
	default:
		return "", fmt.Errorf("Can't marshal type %v, value %v", value.Kind(), v)
	}
}
