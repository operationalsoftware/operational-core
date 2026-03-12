package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

// Deprecated: Build radio markup directly with h.Div/h.Label/h.Input and existing radio CSS classes.
type RadioProps struct {
	Label   string
	Name    string
	Classes c.Classes
}

// Deprecated: Use direct HTML markup and apply "radio-container"/"radio" classes instead of this wrapper component.
func Radio(p *RadioProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}
	classes := c.Classes{
		"radio": true,
	}

	p.Classes["radio-container"] = true
	return h.Div(
		p.Classes,
		h.Label(h.For(p.Name), g.Text(p.Label)),
		h.Input(
			classes,
			h.Type("radio"),
			h.Name(p.Name),
			h.ID(p.Name),
			g.Group(children),
		),
	)
}
