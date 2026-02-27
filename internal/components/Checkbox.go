package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

// Deprecated: Build checkbox markup directly with h.Label/h.Input and the "checkbox" CSS class.
type CheckboxProps struct {
	Name    string
	Label   string
	Value   string
	Checked bool
	Classes c.Classes
}

// Deprecated: Use h.Label and h.Input directly and apply the "checkbox" CSS class instead of this wrapper component.
func Checkbox(p *CheckboxProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}
	p.Classes["checkbox"] = true

	return h.Label(
		p.Classes,
		g.Text(p.Label),
		h.Input(
			h.Type("checkbox"),
			h.Name(p.Name),
			h.Value(p.Value),
			g.If(p.Checked, h.Checked()),
			g.Group(children),
		),
	)
}
