package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type CheckboxProps struct {
	Name    string
	Label   string
	Value   string
	Checked bool
	Classes c.Classes
}

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
