package components

import (
	g "github.com/maragudk/gomponents"
)

type SearchSelectOptionsProps struct {
	Mode           string
	Selected       string
	SelectedValues map[string]bool
	Options        []SearchSelectOptionData
}

func SearchSelectOptions(p *SearchSelectOptionsProps) []g.Node {
	var listOptions []g.Node
	for _, o := range p.Options {
		classes := "select-option"
		if p.Mode == "multi" {
			if p.SelectedValues[o.Value] {
				classes += " selected"
			}
		} else {
			if o.Value == p.Selected {
				classes += " selected"
			}
		}

		listOptions = append(listOptions, SearchSelectOption(
			&SearchSelectOptionProps{
				Classes: classes,
				Value:   o.Value,
				Nodes:   o.Nodes,
				Text:    o.Text,
			},
		),
		)
	}
	return listOptions
}
