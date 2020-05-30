package plugin

// example plugin - plugin-description

import (
	"github.com/u2takey/go-annotation/pkg/lib"
)

func init() {
	lib.RegisterPlugin(new(Description))
}

type Description struct {
	Body string
}

func (p *Description) Template() string {
	return `
func (s *$.type|raw$) GetDescription() string {
	return "$.type|raw$"
}
`
}
