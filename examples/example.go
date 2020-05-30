package examples

//go:generate annotation-gen -i .

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
	B1 *ComponentB `autowired:"true"`
	B2 *ComponentB `autowired:"true"`
	B3 *ComponentB
}

// Annotation@Component={"type": "Singleton"}
type ComponentB struct {
}
