package components

import (
	"strconv"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type InputNumberProps struct {
	Size        InputSize
	Name        string
	Label       string
	Placeholder string
	Min         int
	Max         int
	HelperType  InputHelperType
	HelperText  string
	Classes     c.Classes
}

func InputNumber(p *InputNumberProps, children ...g.Node) g.Node {
	classes := c.Classes{}
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	if p.Size == "" {
		p.Size = InputSizeMedium
	}

	p.Classes["input-container"] = true
	classes[string(p.Size)] = true

	return h.Div(
		p.Classes,
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
			h.Type("number"),
			g.If(p.Min != 0, h.Min(strconv.Itoa(p.Min))),
			g.If(p.Max != 0, h.Max(strconv.Itoa(p.Max))),
			g.Group(children),
		),
		g.If(p.HelperText != "", InputHelper(&InputHelperProps{
			Label: p.HelperText,
			Type:  p.HelperType,
		})),
		InlineStyle("/components/Input.css"),
	)
}
