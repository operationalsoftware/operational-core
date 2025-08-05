package components

import (
	"strings"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

// SearchSelectOption represents <li data-value="...">Label</li>
type SearchSelectOption struct {
	Label string
	Value string
	Nodes []g.Node
}

type SearchSelectProps struct {
	Name          string
	Placeholder   string
	Mode          string // "single", "multi"
	Options       []SearchSelectOption
	Selected      string
	ShowOnlyLabel bool
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

		displayText := o.Value + " - " + o.Label
		if p.ShowOnlyLabel {
			displayText = o.Label
		}

		listOptions = append(listOptions,
			h.Div(
				h.Class(classes),
				h.DataAttr("value", o.Value),
				g.Group(o.Nodes),
				g.Text(displayText),
			),
		)

	}

	var inputText string
	if p.Selected != "" {
		for _, o := range p.Options {
			if o.Value == p.Selected {
				inputText = o.Label
				break
			}
		}
	} else {
		inputText = p.Placeholder
	}

	return h.Div(
		h.Class("search-select"),
		g.Attr("data-mode", p.Mode),
		g.Attr("data-name", p.Name),
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
	)
}
