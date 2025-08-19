package components

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Fieldset(children ...g.Node) g.Node {
	return h.Div(
		h.Class("fieldset"),
		g.Group(children),
	)
}
