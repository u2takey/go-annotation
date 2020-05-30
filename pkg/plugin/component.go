package plugin

import "github.com/u2takey/go-annotation/pkg/lib"

func init() {
	lib.RegisterPlugin(new(Component))
}

type ComponentType string

const (
	Default   ComponentType = ""
	Singleton ComponentType = "Singleton"
)

type Component struct {
	Type ComponentType
}

func (p *Component) Template() string {
	// register a New Method
	return `

var New$.type|raw$Function = &lib.NewFunction{
	F: $.newFunction$,
	Singleton: $.newFunctionSingleton|print$,
}

func init() {
	lib.RegisterType(new($.type|raw$), New$.type|raw$Function) 
}

func Provide$.type|raw$ () (*$.type|raw$, error) {
	r, err := lib.Provide(new($.type|raw$))
	if err != nil{
		return nil, err
	}
	return r.(*$.type|raw$), nil
}
`
}
