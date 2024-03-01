package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
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

	return h.Div(
		p.Classes,
		InputLabel(&InputLabelProps{
			For: p.Name,
		}, g.Text(p.Label)),
		h.Input(
			h.Type("checkbox"),
			h.Name(p.Name),
			h.ID(p.Name),
			h.Value(p.Value),
			g.If(p.Checked, g.Attr("checked", "true")),
		),
		g.Group(children),
	)
}
