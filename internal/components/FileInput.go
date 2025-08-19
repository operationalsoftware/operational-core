package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type FileInputProps struct {
	Name       string
	Label      string
	HelperText string
	HelperType InputHelperType
	InputProps []g.Node
	Classes    c.Classes
}

func FileInput(p *FileInputProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	if p.InputProps == nil {
		p.InputProps = []g.Node{}
	}

	p.Classes["input-container"] = true

	return h.Div(
		p.Classes,
		g.If(
			p.Label != "",
			h.Label(h.For(p.Name), g.Text(p.Label)),
		),
		h.Input(
			h.Type("file"),
			h.Name(p.Name),
			h.ID(p.Name),
			g.Group(p.InputProps),
		),
		g.If(p.HelperText != "", InputHelper(&InputHelperProps{
			Label: p.HelperText,
			Type:  p.HelperType,
		})),
		g.Group(children),
	)
}
