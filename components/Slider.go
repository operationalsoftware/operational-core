package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type SliderProps struct {
	Min     string
	Max     string
	Classes c.Classes
}

func Slider(p *SliderProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	p.Classes["slider"] = true
	return h.Div(
		p.Classes,
		h.Input(
			h.Type("range"),
			h.Min(p.Min),
			h.Max(p.Max),
		),
		h.Div(
			h.Class("show-range-value"),
		),
		InlineStyle("/components/Slider.css"),
		InlineScript("/components/Slider.js"),
	)
}
