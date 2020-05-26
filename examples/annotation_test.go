package examples

import (
	"testing"

	"github.com/u2takey/go-annotation/pkg/lib"
	"github.com/u2takey/go-annotation/pkg/plugins"
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
		annotation := lib.GetAnnotation(tCase.target, &plugins.Description{})
		if annotation == nil {
			t.Errorf("Expect annotation %q exsits", "Description")
		}
		anno := annotation.(*plugins.Description)
		if anno.Body != tCase.pluginBody {
			t.Errorf("Expect annotation body %q, Got %q", tCase.pluginBody, anno.Body)
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
