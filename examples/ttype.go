package examples

import "github.com/u2takey/go-annotation/pkg/plugins"

var _ = plugins.Description{}

// Annotation@Description={"body":"a"}
type A struct {
	FieldA string
}

// Annotation@Description={"body":"b"}
type B struct {
	FieldB string
}
