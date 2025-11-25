package components

import (
	"slices"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SearchSelectOption struct {
	Text     string
	Value    string
	Selected bool
	Nodes    []g.Node
}

type SearchSelectProps struct {
	Name                 string
	Placeholder          string
	Mode                 string // "single", "multi"
	Options              []SearchSelectOption
	Selected             string
	OptionsEndpoint      string // optional: URL to fetch options
	SearchQueryParamName string // optional: query parameter name, default "SearchText"
}

func SearchSelectOptions(options []SearchSelectOption) g.Node {
	var listOptions []g.Node
	for _, o := range options {
		classes := "select-option"
		if o.Selected {
			classes += " selected"
		}

		selectOption := h.Div(
			h.Class(classes),
			h.Data("value", o.Value),
			g.Group(o.Nodes),
			g.Text(o.Text),
		)

		listOptions = append(listOptions, selectOption)

	}
	return h.Div(listOptions...)
}

func SearchSelect(p *SearchSelectProps, children ...g.Node) g.Node {
	for i := range p.Options {
		if p.Options[i].Value == p.Selected {
			p.Options[i].Selected = true
		}
	}

	listOptions := SearchSelectOptions(p.Options)

	searchQueryParamName := p.SearchQueryParamName
	if searchQueryParamName == "" && p.OptionsEndpoint != "" {
		searchQueryParamName = "SearchText"
	}

	var inputText string
	if p.Selected != "" {
		for _, o := range p.Options {
			if o.Value == p.Selected {
				inputText = o.Text
				break
			}
		}
	} else {
		inputText = p.Placeholder
	}

	return h.Div(
		h.Class("search-select"),
		h.Data("mode", p.Mode),
		h.Data("name", p.Name),
		g.If(
			p.OptionsEndpoint != "",
			g.Attr("data-options-endpoint", p.OptionsEndpoint),
		),
		g.If(
			searchQueryParamName != "",
			h.Data("search-query-param", searchQueryParamName),
		),

		h.Div(
			h.Class("select-input"),
			g.Attr("tabindex", "0"),

			h.Span(
				g.Text(inputText),
			),

			Icon(&IconProps{
				Identifier: "chevron-down",
			}),
		),
		h.Div(
			h.Class("select-dropdown"),

			h.Input(
				h.Class("select-search"),
				h.Type("text"),
				g.Attr("placeholder", "Search..."),
			),
			h.Ul(
				h.Class("select-options"),
				listOptions,
			),
		),
		h.Div(
			h.Class("select-hidden-inputs"),
		),
		g.Group(children),

		InlineScript("/internal/components/SearchSelect.js"),
	)
}

func MapStringsToOptions(vals []string, selectedValues []string) []SearchSelectOption {
	out := make([]SearchSelectOption, len(vals))
	for i, v := range vals {
		isSelected := slices.Contains(selectedValues, v)
		out[i] = SearchSelectOption{Text: v, Value: v, Selected: isSelected}
	}
	return out
}
