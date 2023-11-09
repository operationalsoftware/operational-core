package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type SliderProps struct {
	Min string
	Max string
}

func Slider(p *SliderProps, children ...g.Node) g.Node {
	return h.Div(
		h.Input(
			h.Type("range"),
			h.Min(p.Min),
			h.Max(p.Max),
		),
		h.Div(
			h.Class("show-range-value"),
		),
		InlineStyle(Assets, "/Slider.css"),
		InlineScript(Assets, "/Slider.js"),
	)
}
