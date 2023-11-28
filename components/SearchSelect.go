package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type SearchSelectProps struct {
	Name    string
	Options []Option
}

func SearchSelect(p *SearchSelectProps) g.Node {
	classes := c.Classes{
		"search-select": true,
	}

	if p.Name == "" {
		p.Name = "search-select"
	}

	return h.Div(
		classes,
		h.Role("combobox"),
		h.Aria("labelledby", "search-select button"),
		h.Aria("haspopup", "listbox"),
		h.Aria("expanded", "false"),
		h.Aria("controls", "search-select-dropdown"),
		h.Input(
			h.Type("hidden"),
			h.Class("hidden-input"),
			h.Name(p.Name),
			h.Value(""),
		),
		h.Div(
			h.Class("content-container"),
			h.Div(
			// h.Class("selected-values"),
			// h.Div(
			// 	h.Class("selected-value"),
			// 	h.Span(
			// 		g.Text("One"),
			// 	),
			// ),
			),
			h.Input(
				h.Type("search"),
				h.Aria("autocomplete", "off"),
				h.Aria("label", "search-select button"),
				h.Aria("controls", "search-select-dropdown"),
				h.Aria("haspopup", "listbox"),
				h.Aria("expanded", "false"),
				h.Aria("role", "combobox"),
			),
		),
		h.Ul(
			h.Class("search-select-dropdown"),
			h.Role("listbox"),
			h.ID("search-select-dropdown"),
			g.Group(g.Map(p.Options, MultiSelectOptions)),
		),
		InlineStyle(Assets, "/SearchSelect.css"),
		InlineScript(Assets, "/SearchSelect.js"),
	)

}
