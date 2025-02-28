package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func Divider(children ...g.Node) g.Node {

	children = ensureClasses(children, c.Classes{
		"divider": true,
	})

	return h.Div(children...)
}
