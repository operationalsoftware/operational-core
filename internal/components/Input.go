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
	HelperText  string
	InputType   string
	HelperType  InputHelperType
	InputProps  []g.Node
	Classes     c.Classes
}

func Input(p *InputProps, children ...g.Node) g.Node {
	classes := c.Classes{}

	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	if p.InputType == "" {
		p.InputType = "text"
	}

	if p.InputProps == nil {
		p.InputProps = []g.Node{}
	}

	if p.Size == "" {
		p.Size = InputSizeMedium
	}

	p.Classes["input-container"] = true
	classes[string(p.Size)] = true

	return h.Div(
		p.Classes,
		g.If(
			p.Label != "",
			h.Label(h.For(p.Name), g.Text(p.Label)),
		),
		h.Input(
			classes,
			h.Name(p.Name),
			h.ID(p.Name),
			h.Placeholder(p.Placeholder),
			h.Type(p.InputType),
			g.Group(p.InputProps),
		),
		g.If(p.HelperText != "", InputHelper(&InputHelperProps{
			Label: p.HelperText,
			Type:  p.HelperType,
		})),
		g.Group(children),
	)
}
