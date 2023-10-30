package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type InputLabelProps struct {
	For string
}

func InputLabel(p *InputLabelProps, children ...g.Node) g.Node {
	return h.Label(
		h.For(p.For),
		g.Group(children),
	)
}
