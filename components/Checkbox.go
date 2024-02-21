package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type CheckboxProps struct {
	Name    string
	Label   string
	Value   string
	Checked bool
}

func Checkbox(p *CheckboxProps, children ...g.Node) g.Node {
	return h.Div(
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
		InlineStyle("/components/Checkbox.css"),
	)
}
