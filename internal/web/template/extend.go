package template

import (
	"errors"
	"reflect"

	"github.com/nikolalohinski/gonja/v2/exec"
)

// Directive creates a new global function from a given function that returns string
func Directive(name string, fn any) any {
	return AsCtx(func() *exec.Context {
		fnType := reflect.TypeOf(fn)
		if fnType.Kind() != reflect.Func {
			panic("directive expects a function")
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

				result := fnValue.Call(args)
				return exec.AsValue(result[0].String())
			},
		})
	})
}
