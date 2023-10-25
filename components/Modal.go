package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type ModalProps struct {
	Title         string
	FooterContent g.Node
}

func Modal(p *ModalProps, children ...g.Node) g.Node {
	classes := c.Classes{
		"modal":  true,
		"hidden": true,
	}
	return h.Dialog(
		classes,
		h.Header(
			h.Class("modal-header"),
			h.H3(g.Text(p.Title)),
			h.Span(
				h.ID("close-btn"),
				h.Class("primary close-btn"),
				g.Text("X"),
			),
		),
		h.Div(
			h.Class("modal-content"),
			g.Group(children),
		),
		g.If(p.FooterContent != nil,
			h.Footer(
				h.Class("modal-footer"),
				p.FooterContent,
			),
		),
		InlineStyle(Assets, "/Modal.css"),
		InlineScript(Assets, "/Modal.js"),
	)
}
