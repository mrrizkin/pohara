package template

import (
	"fmt"
	"reflect"

	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
)

type ifncs struct {
	name     string
	token    *tokens.Token
	params   [][]any
	wrappers []*nodes.Wrapper
	fn       any
}

func newIf(name string, fn any) *exec.ControlStructureSet {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("if expects a function")
	}

	if fnType.NumOut() != 1 || fnType.Out(0).Kind() != reflect.Bool {
		panic("function must return a bool")
	}

	// get all the argument types
	argsType := make([]reflect.Type, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		argsType[i] = fnType.In(i)
	}

	return exec.NewControlStructureSet(map[string]parser.ControlStructureParser{
		name: func(p *parser.Parser, args *parser.Parser) (nodes.ControlStructure, error) {
			i := &ifncs{
				token:  p.Current(),
				name:   name,
				params: make([][]any, 0),
				fn:     fn,
			}

			// parse the first if arguments
			index := 0
			i.params = append(i.params, make([]any, 0))
			for _, argType := range argsType {
				switch argType.Kind() {
				case reflect.String:
					i.params[index] = append(
						i.params[index],
						args.Match(tokens.Name, tokens.String),
					)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					i.params[index] = append(
						i.params[index],
						args.Match(tokens.Name, tokens.Integer),
					)
				case reflect.Float64, reflect.Float32:
					i.params[index] = append(i.params[index], args.Match(tokens.Name, tokens.Float))
				default:
					panic(fmt.Sprintf("function argument %s, not supported", argType.Kind()))
				}
			}

			if !args.End() {
				return nil, args.Error(
					fmt.Sprintf("Malformed %s controlStructure args.", name),
					nil,
				)
			}

			index++

			// Check the rest
			for {
				wrapper, tagArgs, err := p.WrapUntil(
					fmt.Sprintf("else%s", name),
					"else",
					fmt.Sprintf("end%s", name),
				)
				if err != nil {
					return nil, err
				}

				i.wrappers = append(i.wrappers, wrapper)

				if wrapper.EndTag == fmt.Sprintf("else%s", name) {
					// elif can take a condition
					i.params = append(i.params, make([]any, 0))
					for _, argType := range argsType {
						switch argType.Kind() {
						case reflect.String:
							i.params[index] = append(
								i.params[index],
								tagArgs.Match(tokens.Name, tokens.String),
							)
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							i.params[index] = append(
								i.params[index],
								tagArgs.Match(tokens.Name, tokens.Integer),
							)
						case reflect.Float64, reflect.Float32:
							i.params[index] = append(
								i.params[index],
								tagArgs.Match(tokens.Name, tokens.Float),
							)
						default:
							panic(
								fmt.Sprintf("function argument %s, not supported", argType.Kind()),
							)
						}
					}

					if !tagArgs.End() {
						return nil, tagArgs.Error(
							fmt.Sprintf("else%s-condition is malformed.", name),
							nil,
						)
					}
					index++
				} else {
					if !tagArgs.End() {
						// else/endif can't take any conditions
						return nil, tagArgs.Error("Arguments not allowed here.", nil)
					}
				}

				if wrapper.EndTag == fmt.Sprintf("end%s", name) {
					break
				}
			}

			return i, nil
		},
	})
}

func (i *ifncs) Position() *tokens.Token {
	return i.token
}

func (i *ifncs) String() string {
	t := i.Position()
	return fmt.Sprintf("IfControlStructure(Name=%s Line=%d Col=%d)", i.name, t.Line, t.Col)
}

func (i *ifncs) Execute(r *exec.Renderer, tag *nodes.ControlStructureBlock) error {
	for index, param := range i.params {
		args := make([]reflect.Value, len(param))
		for j, arg := range param {
			switch arg := arg.(type) {
			case *tokens.Token:
				switch arg.Type {
				case tokens.String:
					args[j] = reflect.ValueOf(arg.Val)
				case tokens.Integer:
					args[j] = reflect.ValueOf(arg.Val)
				case tokens.Float:
					args[j] = reflect.ValueOf(arg.Val)
				case tokens.Name:
					val, exist := r.Environment.Context.Get(arg.Val)
					if !exist {
						return fmt.Errorf("variable '%s' not found in context", arg.Val)
					}

					args[j] = reflect.ValueOf(val)
				}
			}
		}

		result := reflect.ValueOf(i.fn).Call(args)
		if result[0].Bool() {
			return r.ExecuteWrapper(i.wrappers[index])
		}

		// last condition is always true
		if len(i.params) == index+1 && len(i.wrappers) > index+1 {
			return r.ExecuteWrapper(i.wrappers[index+1])
		}
	}

	return nil
}
