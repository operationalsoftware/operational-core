package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type TextareaProps struct {
	Name        string
	Label       string
	Placeholder string
}

func Textarea(p *TextareaProps, children ...g.Node) g.Node {
	return h.Div(
		InputLabel(&InputLabelProps{
			For: p.Name,
		},
			g.Text(p.Label),
		),
		h.Textarea(
			h.Name(p.Name),
			h.ID(p.Name),
			h.Placeholder(p.Placeholder),
			g.Group(children),
		),
		InlineStyle(Assets, "/Textarea.css"),
	)
}
