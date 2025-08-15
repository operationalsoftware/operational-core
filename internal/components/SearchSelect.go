package components

import (
	"strings"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type SearchSelectOption struct {
	Text  string
	Value string
	Nodes []g.Node
}

type SearchSelectProps struct {
	Name                 string
	Placeholder          string
	Mode                 string // "single", "multi"
	Options              []SearchSelectOption
	Selected             string
	OptionsAPI           string // optional: URL to fetch options
	SearchQueryParamName string // optional: query parameter name, default "SearchText"
}

func SearchSelect(p *SearchSelectProps, children ...g.Node) g.Node {

	selectedValues := map[string]bool{}
	if p.Mode == "multi" && p.Selected != "" {
		for _, val := range strings.Split(p.Selected, ",") {
			selectedValues[val] = true
		}
	}

	var listOptions []g.Node
	for _, o := range p.Options {
		classes := "select-option"
		if p.Mode == "multi" {
			if selectedValues[o.Value] {
				classes += " selected"
			}
		} else {
			if o.Value == p.Selected {
				classes += " selected"
			}
		}

		listOptions = append(listOptions,
			h.Div(
				h.Class(classes),
				h.DataAttr("value", o.Value),
				g.Group(o.Nodes),
				g.Text(o.Text),
			),
		)

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

	attrs := []g.Node{
		h.Class("search-select"),
		g.Attr("data-mode", p.Mode),
		g.Attr("data-name", p.Name),
	}

	if p.OptionsAPI != "" {
		attrs = append(attrs, g.Attr("data-options-endpoint", p.OptionsAPI))
	}

	if p.SearchQueryParamName != "" {
		attrs = append(attrs, h.DataAttr("search-query-param", p.SearchQueryParamName))
	}

	return h.Div(
		g.Group(attrs),

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
				g.Group(listOptions),
			),
		),
		h.Div(
			h.Class("select-hidden-inputs"),
		),
		g.Group(children),

		InlineScript("/internal/components/SearchSelect.js"),
	)
}
