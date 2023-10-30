package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type RadioProps struct {
	Label string
	Name  string
}

func Radio(p *RadioProps, children ...g.Node) g.Node {
	return h.Div(
		InputLabel(&InputLabelProps{
			For: p.Name,
		},
			g.Text(p.Label),
		),
		h.Input(
			h.Type("radio"),
			h.Name(p.Name),
			h.ID(p.Name),
		),
		InlineStyle(Assets, "/Radio.css"),
	)
}
