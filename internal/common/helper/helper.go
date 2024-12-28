package helper

import (
	"fmt"
	"reflect"
)

func StructToMap(in interface{}, tag string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("StructToMap only accept structs: got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tag != "" {
			if tagv := fi.Tag.Get(tag); tagv != "" {
				out[tagv] = v.Field(i).Interface()
			}
		} else {
			out[fi.Name] = v.Field(i).Interface()
		}
	}

	return out, nil
}
