package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func Fieldset(children ...g.Node) g.Node {
	return h.Div(
		h.Class("fieldset"),
		g.Group(children),
	)
}
