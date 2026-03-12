package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

// Deprecated: Build textarea markup directly with h.Div/h.Label/h.Textarea and existing CSS classes.
type TextareaProps struct {
	Name        string
	Label       string
	Placeholder string
	Classes     c.Classes
}

// Deprecated: Use h.Textarea directly and apply the "textarea-container" CSS class instead of this wrapper component.
func Textarea(p *TextareaProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	p.Classes["textarea-container"] = true
	return h.Div(
		p.Classes,
		h.Label(h.For(p.Name), g.Text(p.Label)),
		h.Textarea(
			h.Name(p.Name),
			h.ID(p.Name),
			h.Placeholder(p.Placeholder),
			g.Group(children),
		),
	)
}
