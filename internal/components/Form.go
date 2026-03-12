package components

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// Deprecated: Use h.FormEl directly and apply the "form" CSS class instead of this wrapper component.
func Form(children ...g.Node) g.Node {
	return h.FormEl(
		h.Class("form"),
		g.Group(children),
	)
}
