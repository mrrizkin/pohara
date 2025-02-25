package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ParseQueryParams parses query parameters from a Fiber request into a struct
// using the "query" tag. The tag supports the following options:
// - fieldName: the query parameter name (required)
// - required: indicates that the field is required
// - default=value: default value if parameter is not provided
//
// Example usage:
//
//	type RequestQueryUserDatatable struct {
//	  Page  sql.Int64Nullable `query:"page,default=1"`
//	  Limit sql.Int64Nullable `query:"limit,default=10"`
//	  Name  sql.StringNullable `query:"name"`
//	}
func ParseQueryParams(c *fiber.Ctx, dest interface{}) error {
	// Get the reflected value of the destination struct
	destVal := reflect.ValueOf(dest)

	// Ensure dest is a pointer to a struct
	if destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Struct {
		return errors.New("destination must be a pointer to a struct")
	}

	// Get the struct value
	structVal := destVal.Elem()
	structType := structVal.Type()

	// Iterate through struct fields
	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		fieldType := structType.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Get query tag
		queryTag := fieldType.Tag.Get("query")
		if queryTag == "" {
			continue
		}

		// Parse tag options
		tagParts := strings.Split(queryTag, ",")
		if len(tagParts) == 0 {
			continue
		}

		paramName := tagParts[0]
		required := false
		defaultValue := ""

		// Process tag options
		for _, option := range tagParts[1:] {
			if option == "required" {
				required = true
			} else if strings.HasPrefix(option, "default=") {
				defaultValue = strings.TrimPrefix(option, "default=")
			}
		}

		// Get query parameter value
		paramValue := c.Query(paramName)

		// Check if required parameter is missing
		if paramValue == "" && required && defaultValue == "" {
			return fmt.Errorf("missing required query parameter: %s", paramName)
		}

		// Use default value if parameter is not provided
		if paramValue == "" && defaultValue != "" {
			paramValue = defaultValue
		}

		// Skip if no value to set (and not using default)
		if paramValue == "" {
			continue
		}

		// Try to set value using Scan method first (for SQL nullable types)
		scanMethod := field.Addr().MethodByName("Scan")
		if scanMethod.IsValid() {
			scanMethod.Call([]reflect.Value{reflect.ValueOf(paramValue)})
			continue
		}

		// Handle standard types if Scan wasn't available
		switch field.Kind() {
		case reflect.String:
			field.SetString(paramValue)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(paramValue, 10, 64)
			if err != nil {
				return fmt.Errorf(
					"invalid int value %q for parameter %s: %w",
					paramValue,
					paramName,
					err,
				)
			}
			field.SetInt(intValue)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintValue, err := strconv.ParseUint(paramValue, 10, 64)
			if err != nil {
				return fmt.Errorf(
					"invalid uint value %q for parameter %s: %w",
					paramValue,
					paramName,
					err,
				)
			}
			field.SetUint(uintValue)

		case reflect.Float32, reflect.Float64:
			floatValue, err := strconv.ParseFloat(paramValue, 64)
			if err != nil {
				return fmt.Errorf(
					"invalid float value %q for parameter %s: %w",
					paramValue,
					paramName,
					err,
				)
			}
			field.SetFloat(floatValue)

		case reflect.Bool:
			boolValue, err := strconv.ParseBool(paramValue)
			if err != nil {
				return fmt.Errorf(
					"invalid bool value %q for parameter %s: %w",
					paramValue,
					paramName,
					err,
				)
			}
			field.SetBool(boolValue)

		case reflect.Struct:
			// Try Unmarshal method as fallback
			unmarshalMethod := field.Addr().MethodByName("UnmarshalText")
			if unmarshalMethod.IsValid() {
				unmarshalMethod.Call([]reflect.Value{reflect.ValueOf([]byte(paramValue))})
				continue
			}

			return fmt.Errorf(
				"unsupported struct type %q for parameter %s",
				field.Type().String(),
				paramName,
			)

		default:
			return fmt.Errorf(
				"unsupported type %q for parameter %s",
				field.Type().String(),
				paramName,
			)
		}
	}

	return nil
}
