package components

import (
	"strings"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
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

type SearchSelectOptionsProps struct {
	Mode           string
	Selected       string
	SelectedValues map[string]bool
	Options        []SearchSelectOption
}

func SearchSelectOptions(p *SearchSelectOptionsProps) g.Node {
	var listOptions []g.Node
	for _, o := range p.Options {
		classes := "select-option"
		if o.Selected {
			classes += " selected"
		}

		selectOption := h.Div(
			h.Class(classes),
			h.DataAttr("value", o.Value),
			g.Group(o.Nodes),
			g.Text(o.Text),
		)

		listOptions = append(listOptions, selectOption)

	}
	return h.Div(listOptions...)
}

func SearchSelect(p *SearchSelectProps, children ...g.Node) g.Node {

	selectedValues := map[string]bool{}
	if p.Mode == "multi" && p.Selected != "" {
		for _, val := range strings.Split(p.Selected, ",") {
			selectedValues[val] = true
		}
	}

	listOptions := SearchSelectOptions(&SearchSelectOptionsProps{
		Mode:           p.Mode,
		Selected:       p.Selected,
		SelectedValues: selectedValues,
		Options:        p.Options,
	})

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

	if p.OptionsEndpoint != "" {
		attrs = append(attrs, g.Attr("data-options-endpoint", p.OptionsEndpoint))
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
