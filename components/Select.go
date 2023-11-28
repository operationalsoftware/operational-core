package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type Option struct {
	Label string
	Value string
}

type SelectProps struct {
	Options []Option
	Name    string
}

func renderOptions(o Option) g.Node {
	return h.Li(
		h.Role("option"),
		Radio(&RadioProps{
			Name:  o.Value,
			Label: o.Label,
		}),
	)
}

func Select(p *SelectProps) g.Node {
	classes := c.Classes{
		"custom-select": true,
	}

	if p.Name == "" {
		p.Name = "select"
	}

	return h.Div(
		classes,
		h.Div(
			h.Class("select-button"),
			h.Role("combobox"),
			h.Aria("labelledby", "select button"),
			h.Aria("haspopup", "listbox"),
			h.Aria("expanded", "false"),
			h.Aria("controls", "select-dropdown"),
			h.Input(
				h.Type("hidden"),
				h.Class("hidden-input"),
				h.Name(p.Name),
				h.Value(""),
			),
			h.Span(h.Class("default-value"), g.Text("")),
			h.Span(h.Class("arrow")),
		),
		h.Ul(
			h.Class("select-dropdown"),
			h.Role("listbox"),
			h.ID("select-dropdown"),
			g.Group(g.Map(p.Options, renderOptions)),
		),
		InlineStyle(Assets, "/Select.css"),
		InlineScript(Assets, "/Select.js"),
	)
}
