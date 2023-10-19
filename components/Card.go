package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func Card(children ...g.Node) g.Node {
	return h.Div(
		InlineStyle(Assets, "/Card.css"),
		h.Class("card"),
		g.Group(children),
	)
}
