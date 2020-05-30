package lib

import (
	"reflect"

	"k8s.io/gengo/types"
)

type TargetType = types.Kind

const (
	TargetDefault TargetType = TargetStruct
	TargetStruct  TargetType = types.Struct
	TargetMethod  TargetType = types.Func
)

// Annotation is simple annotation, it keeps annotation information for specific types.
type Annotation interface {
}

func GetAnnotationName(a Annotation) string {
	d := reflect.Indirect(reflect.ValueOf(a))
	dt := d.Type()
	return dt.Name()
}

type TargetedAnnotation interface {
	Target() TargetType
}

// CompileAnnotation, annotations which will registered by annotation_gen to gen additional codes
// CompileAnnotation MUST be registered as an Annotation Plugin for tool annotation-gen.
type CompileAnnotation interface {
	Annotation
	Template() string
}
