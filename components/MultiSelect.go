package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type MultiSelectProps struct {
	Options []Option
	Name    string
}

func MultiSelectOptions(o Option) g.Node {
	return h.Div(
		h.Class("option"),
		h.Li(
			h.Role("option"),
			Checkbox(&CheckboxProps{
				Label: o.Label,
				Value: o.Value,
			}),
		),
		Icon("check"),
	)
}

func MultiSelect(p *MultiSelectProps) g.Node {
	classes := c.Classes{
		"custom-multi-select": true,
	}

	if p.Name == "" {
		p.Name = "multi-select"
	}

	return h.Div(
		classes,
		h.Div(
			h.Class("multi-select-button"),
			h.Role("combobox"),
			h.Aria("labelledby", "multi-select button"),
			h.Aria("haspopup", "listbox"),
			h.Aria("expanded", "false"),
			h.Aria("controls", "multi-select-dropdown"),
			h.Input(
				h.Type("hidden"),
				h.Class("hidden-input"),
				h.Name(p.Name),
				h.Value(p.Options[0].Label),
			),
			h.Div(
				h.Class("selected-values"),
			),
			h.Span(h.Class("arrow")),
		),
		h.Ul(
			h.Class("multi-select-dropdown"),
			h.Role("listbox"),
			h.ID("multi-select-dropdown"),
			g.Group(g.Map(p.Options, MultiSelectOptions)),
		),
		InlineStyle(Assets, "/MultiSelect.css"),
		InlineScript(Assets, "/MultiSelect.js"),
	)
}
