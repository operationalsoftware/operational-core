package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

// Deprecated: Use h.Div directly and apply the "divider" CSS class instead of this wrapper component.
func Divider(children ...g.Node) g.Node {

	children = ensureClasses(children, c.Classes{
		"divider": true,
	})

	return h.Div(children...)
}
