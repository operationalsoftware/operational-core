package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type SearchSelectOptionData struct {
	Text  string
	Value string
	Nodes []g.Node
}

type SearchSelectOptionProps struct {
	Classes string
	Value   string
	Nodes   []g.Node
	Text    string
}

func SearchSelectOption(p *SearchSelectOptionProps) g.Node {

	return h.Div(
		h.Class(p.Classes),
		h.DataAttr("value", p.Value),
		g.Group(p.Nodes),
		g.Text(p.Text),
	)
}
