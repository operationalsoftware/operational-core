package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func Form(children ...g.Node) g.Node {
	return h.FormEl(
		g.Group(children),
		InlineStyle("/components/Form.css"),
	)
}
