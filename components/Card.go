package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func Card(children ...g.Node) g.Node {
	return h.Div(
		h.Class("card"),
		g.Group(children),
		InlineStyle("/components/Card.css"),
	)
}
