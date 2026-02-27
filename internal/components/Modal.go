package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

// Deprecated: Use native browser dialogs (window.confirm/window.alert/window.prompt) instead of this props wrapper.
type ModalProps struct {
	Title         string
	FooterContent g.Node
	Classes       c.Classes
}

// Deprecated: Use native browser dialogs (window.confirm/window.alert/window.prompt) instead of this wrapper component.
func Modal(p *ModalProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	p.Classes["modal"] = true
	p.Classes["hidden"] = true
	return h.Dialog(
		p.Classes,
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
	)
}
