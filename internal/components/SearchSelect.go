package components

import (
	c "github.com/maragudk/gomponents/components"
)

type SearchSelectProps struct {
	Name      string
	Options   []Option
	Classes   c.Classes
	ID        string
	OptionUrl string
	Multiple  bool
	Value     []string
}

// func SearchSelect(p *SearchSelectProps) g.Node {
// 	if p.Classes == nil {
// 		p.Classes = c.Classes{}
// 	}

// 	if p.Multiple {
// 		p.Classes["search-select"] = true
// 	} else {
// 		p.Classes["search-select-single"] = true
// 	}

// 	if p.Name == "" {
// 		p.Name = "search-select"
// 	}

// 	if p.ID == "" {
// 		p.ID = "search-select-dropdown"
// 	} else {
// 		p.ID = p.ID + "-search-select-dropdown"
// 	}

// 	p.Classes["search-select-container"] = true

// 	return h.Div(
// 		p.Classes,
// 		g.If(p.Multiple, h.DataAttr("multiple", "true")),
// 		h.Role("combobox"),
// 		h.Aria("labelledby", "search-select button"),
// 		h.Aria("haspopup", "listbox"),
// 		h.Aria("expanded", "false"),
// 		h.Aria("controls", "search-select-dropdown"),
// 		h.Input(
// 			h.Type("hidden"),
// 			h.Class("hidden-input"),
// 			h.Name(p.Name),
// 			h.Value(""),
// 		),
// 		h.Div(
// 			h.Class("content-container"),
// 			g.If(p.Multiple, h.Div(
// 				h.Class("selected-values"),
// 			)),
// 			h.Input(
// 				h.Type("search"),
// 				h.Aria("autocomplete", "off"),
// 				h.Aria("label", "search-select button"),
// 				h.Aria("controls", p.ID),
// 				h.Aria("haspopup", "listbox"),
// 				h.Aria("expanded", "false"),
// 				h.Aria("role", "combobox"),
// 				h.Name("search-value"),
// 				ghtmx.Post(p.OptionUrl),
// 				ghtmx.Trigger("keyup changed delay:500ms"),
// 				ghtmx.Target(".search-select-dropdown"),
// 				ghtmx.Swap("outerHTML"),
// 			),
// 			g.If(!p.Multiple, h.Span(h.Class("arrow"))),
// 		),
// 		MultiSelectOptions(p.Options, p.Value, p.Name, "search-select-dropdown", p.ID),
// 		InlineScript("/components/SearchSelect.js"),
// 	)

// }
