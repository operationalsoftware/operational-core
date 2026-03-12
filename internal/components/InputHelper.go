package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

// Deprecated: Use string CSS class names directly on helper text markup.
type InputHelperType string

const (
	InputHelperTypeSuccess InputHelperType = "success"
	InputHelperTypeWarning InputHelperType = "warning"
	InputHelperTypeError   InputHelperType = "error"
	InputHelperTypeNone    InputHelperType = "none"
)

// Deprecated: Build helper text markup directly and apply "input-helper" and state CSS classes.
type InputHelperProps struct {
	Label string
	Type  InputHelperType
}

// Deprecated: Use h.Div directly and apply "input-helper" and state classes instead of this wrapper component.
func InputHelper(p *InputHelperProps) g.Node {
	classes := c.Classes{
		"input-helper": true,
	}

	if p.Type == "" {
		p.Type = InputHelperTypeWarning
	}

	classes[string(p.Type)] = true

	return h.Div(
		classes,
		g.Text(p.Label),
	)
}
