package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type RadioProps struct {
	Label   string
	Name    string
	Classes c.Classes
}

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
