package template

type Option struct {
	Value interface{}
	Valid bool
}

var (
	Directive = newDirective
	CustomIf  = newIf
)
