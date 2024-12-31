package template

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/nikolalohinski/gonja/v2/exec"
)

// newDirective creates a new global function from a given function
func newDirective(name string, fn any) *exec.Context {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("directive expects a function")
	}

	returnFnType := Option{
		Valid: false,
	}

	if fnType.NumOut() > 1 {
		panic("function must return 1 thing")
	}

	if fnType.NumOut() == 1 {
		returnFnType.Valid = true
	}

	if returnFnType.Valid {
		switch fnType.Out(0).Kind() {
		case reflect.String:
			returnFnType.Value = "string"
		default:
			panic(fmt.Sprintf("function return %T, not supported", fnType.Out(0).Kind()))
		}
	}

	if fnType.NumOut() != 1 || fnType.Out(0).Kind() != reflect.String {
		panic("function must return a string")
	}

	return exec.NewContext(map[string]interface{}{
		name: func(_ *exec.Evaluator, params *exec.VarArgs) *exec.Value {
			fnValue := reflect.ValueOf(fn)
			numIn := fnType.NumIn()

			if len(params.Args) != numIn {
				return exec.AsValue(
					exec.ErrInvalidCall(errors.New("invalid number of arguments")),
				)
			}

			args := make([]reflect.Value, numIn)
			for i := 0; i < numIn; i++ {
				arg := params.Args[i]
				expectedType := fnType.In(i)

				var value reflect.Value
				switch expectedType.Kind() {
				case reflect.String:
					value = reflect.ValueOf(arg.String())
				case reflect.Int, reflect.Int64:
					value = reflect.ValueOf(arg.Integer())
				case reflect.Float64:
					value = reflect.ValueOf(arg.Float())
				case reflect.Bool:
					value = reflect.ValueOf(arg.Bool())
				default:
					return exec.AsValue(
						exec.ErrInvalidCall(errors.New("unsupported argument type")),
					)
				}
				args[i] = value
			}

			if returnFnType.Valid {
				result := fnValue.Call(args)
				return exec.AsValue(result[0].String())
			} else {
				fnValue.Call(args)
				return exec.AsValue("")
			}
		},
	})
}
