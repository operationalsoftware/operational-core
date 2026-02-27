package components

import (
	"strconv"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

// Deprecated: Build number input markup directly with h.Input and CSS classes instead of this props wrapper.
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

// Deprecated: Build number input markup directly and apply "input-container" and size CSS classes.
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
		h.Label(h.For(p.Name), g.Text(p.Label)),
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
	)
}
