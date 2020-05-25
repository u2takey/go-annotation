package plugins

import (
	"github.com/u2takey/go-annotation/pkg/lib"
)

func init() {
	lib.RegisterPlugin(new(Description))
}

// plugin description

type Description struct {
	Body string
}

func (p *Description) GetTemplate() string {
	return `
func (s *$.type|raw$) GetDescription() string {
	return "$.type|raw$"
}
`
}
