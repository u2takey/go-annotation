package examples

//go:generate annotation-gen -i . -v 8

import (
	"github.com/u2takey/go-annotation/pkg/plugin"
)

var _ = plugin.Description{}

// Annotation@Description={"body":"a"}
type A struct {
	FieldA string
}

// Annotation@Description={"body":"b"}
type B struct {
	FieldB string
}

// Annotation@Description={"body":"b"}
type C interface {
}

// Annotation@Component
type ComponentA struct {
	B1 *ComponentB `autowired:"true"` // Will populate with new(ComponentB)
	B2 *ComponentB `autowired:"true"` // Will populate with new(ComponentB)
	B3 *ComponentB
}

// Annotation@Component={"type": "Singleton"}
type ComponentB struct {
	C *ComponentC `autowired:"true"` // Will populate with NewComponentC()
}

// Annotation@Component
type ComponentC struct {
	D        *ComponentD `autowired:"true"` // Will populate with NewComponentD()
	IntValue int
}

func NewComponentC() *ComponentC {
	return &ComponentC{IntValue: 1}
}

// Annotation@Component
type ComponentD struct {
	IntValue int
}

func NewComponentD() (*ComponentD, error) {
	return &ComponentD{IntValue: 2}, nil
}
