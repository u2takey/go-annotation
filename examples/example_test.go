package examples

import (
	"testing"

	"github.com/u2takey/go-annotation/pkg/lib"
	"github.com/u2takey/go-annotation/pkg/plugin"
)

func TestAnnotationDescription(t *testing.T) {
	tests := []struct {
		target      interface{}
		pluginBody  string
		description string
	}{
		{
			target:      new(A),
			pluginBody:  "a",
			description: "A",
		},
		{
			target:      new(B),
			pluginBody:  "b",
			description: "B",
		},
	}

	type HasDescription interface {
		GetDescription() string
	}
	for _, tCase := range tests {
		annotation := lib.GetAnnotation(tCase.target, &plugin.Description{})
		if annotation == nil {
			t.Errorf("Expect annotation %q exsits", "Description")
		}
		anAnnotation := annotation.(*plugin.Description)
		if anAnnotation.Body != tCase.pluginBody {
			t.Errorf("Expect annotation body %q, Got %q", tCase.pluginBody, anAnnotation.Body)
		}

		des, ok := tCase.target.(HasDescription)
		if !ok {
			t.Errorf("Expect has GetDescription exsits")
			continue
		}
		if des.GetDescription() != tCase.description {
			t.Errorf("Expect GetDescription returns %q, Got %q", tCase.description, des.GetDescription())
		}

	}

}

func TestAnnotationComponent(t *testing.T) {
	ca, err := ProvideComponentA()
	if err != nil {
		t.Errorf("call ProvideComponentA failed: %s", err)
	}
	if ca.B1 == nil || ca.B2 == nil {
		t.Errorf("expect b1, b2 not nil b1: %v b2: %v", ca.B1, ca.B2)
	}
	if ca.B3 != nil {
		t.Errorf("expect b3 not nil")
	}

	if ca.B1 != ca.B2 {
		t.Errorf("expect b1 == b2")
	}

	if ca.B1.C == nil {
		t.Errorf("expect b1.c not nil")
	}

	if ca.B1.C.IntValue != 1 {
		t.Errorf("expect ca.b1.c.intvalue %q, got %q", 1, ca.B1.C.IntValue)
	}

	if ca.B1.C.D == nil {
		t.Errorf("expect b1.c.d not nil")
	}

	if ca.B1.C.D.IntValue != 2 {
		t.Errorf("expect ca.b1.c.d.intvalue %q, got %q", 2, ca.B1.C.D.IntValue)
	}
}
