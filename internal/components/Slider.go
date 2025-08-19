package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
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
		InlineScript("/components/Slider.js"),
	)
}
