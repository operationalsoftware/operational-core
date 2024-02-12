package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type InputHelperType string

const (
	InputHelperTypeSuccess InputHelperType = "success"
	InputHelperTypeWarning InputHelperType = "warning"
	InputHelperTypeError   InputHelperType = "error"
)

type InputHelperProps struct {
	Label string
	Type  InputHelperType
}

func InputHelper(p *InputHelperProps) g.Node {
	classes := c.Classes{}

	if p.Type == "" {
		p.Type = InputHelperTypeWarning
	}

	classes[string(p.Type)] = true

	return h.Div(
		classes,
		g.Text(p.Label),
		InlineStyle("/components/InputHelper.css"),
	)
}
