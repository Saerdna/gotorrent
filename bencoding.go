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
	default:
		return "", fmt.Errorf("Unknown type for object %v", v)
	}
	return "", nil
}
