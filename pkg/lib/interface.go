package lib

import "reflect"

// Annotation is simple annotation, it keeps annotation information for specific types.
type Annotation interface {
}

func GetAnnotationName(a Annotation) string{
	d := reflect.Indirect(reflect.ValueOf(a))
	dt := d.Type()
	return dt.Name()
}


// CompileAnnotation, annotations which will registered by annotation_gen to gen additional codes
// CompileAnnotation MUST be registered as an Annotation Plugin for tool annotation-gen.
type CompileAnnotation interface {
	Annotation
	GetTemplate() string
}
