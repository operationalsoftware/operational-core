package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

func Divider(children ...g.Node) g.Node {

	children = ensureClasses(children, c.Classes{
		"divider": true,
	})

	return h.Div(children...)
}
