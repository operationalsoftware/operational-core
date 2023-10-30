package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type InputSize string

const (
	InputSizeSmall  InputSize = "sm"
	InputSizeMedium InputSize = "md"
	InputSizeLarge  InputSize = "lg"
)

type InputProps struct {
	Size        InputSize
	Name        string
	Label       string
	Placeholder string
	HelperType  InputHelperType
	HelperText  string
}

func Input(p *InputProps, children ...g.Node) g.Node {
	classes := c.Classes{}

	if p.Size == "" {
		p.Size = InputSizeMedium
	}

	classes[string(p.Size)] = true

	return h.Div(
		InputLabel(&InputLabelProps{
			For: p.Name,
		},
			g.Text(p.Label),
		),
		h.Input(
			classes,
			h.Name(p.Name),
			h.ID(p.Name),
			h.Placeholder(p.Placeholder),
			h.Type("text"),
			g.Group(children),
		),
		g.If(p.HelperText != "", InputHelper(&InputHelperProps{
			Label: p.HelperText,
			Type:  p.HelperType,
		})),
		InlineStyle(Assets, "/Input.css"),
	)
}
