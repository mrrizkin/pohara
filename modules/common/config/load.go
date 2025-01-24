package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func load(cfg interface{}) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Handle nested structs and pointers to structs
		if field.Type.Kind() == reflect.Struct {
			if err := load(fieldValue.Addr().Interface()); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
			}
			continue
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(field.Type.Elem()))
			}
			if err := load(fieldValue.Interface()); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
			}
			continue
		}

		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		parts := strings.Split(envTag, ",")
		envVar := parts[0]
		var defaultValue string
		var isRequired bool

		for _, part := range parts[1:] {
			switch {
			case part == "required":
				isRequired = true
			case strings.HasPrefix(part, "default="):
				defaultValue = strings.TrimPrefix(part, "default=")
			}
		}

		value, exists := os.LookupEnv(envVar)
		if !exists {
			if defaultValue != "" {
				value = defaultValue
			} else if isRequired {
				return fmt.Errorf("required environment variable %q not set", envVar)
			} else {
				continue
			}
		}

		// Handle pointers to basic types
		if fieldValue.Kind() == reflect.Ptr {
			elemType := fieldValue.Type().Elem()
			switch elemType.Kind() {
			case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Bool, reflect.Float32, reflect.Float64:
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(elemType))
				}
				fieldValue = fieldValue.Elem()
			default:
				return fmt.Errorf("unsupported pointer type %q for field %q", elemType, field.Name)
			}
		}

		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set field %q", field.Name)
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("field %q: invalid int value %q: %w", field.Name, value, err)
			}
			fieldValue.SetInt(intValue)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintValue, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return fmt.Errorf("field %q: invalid uint value %q: %w", field.Name, value, err)
			}
			fieldValue.SetUint(uintValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("field %q: invalid bool value %q: %w", field.Name, value, err)
			}
			fieldValue.SetBool(boolValue)
		case reflect.Float32, reflect.Float64:
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("field %q: invalid float value %q: %w", field.Name, value, err)
			}
			fieldValue.SetFloat(floatValue)
		default:
			return fmt.Errorf("unsupported type %q for field %q", fieldValue.Type(), field.Name)
		}
	}

	return nil
}
