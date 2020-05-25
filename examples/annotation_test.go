package examples

import (
	"testing"

	"github.com/u2takey/go-annotation/pkg/lib"
)

func TestAnnotation(t *testing.T) {
	for _, a := range lib.GetAllAnnotations() {
		t.Logf("%+v\n", a)
	}
}
