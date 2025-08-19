package components

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Form(children ...g.Node) g.Node {
	return h.FormEl(
		h.Class("form"),
		g.Group(children),
	)
}
