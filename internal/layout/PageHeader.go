package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PageHeaderProps struct {
	Title   g.Node
	Actions []g.Node
}

func PageHeader(p *PageHeaderProps) g.Node {
	if p == nil {
		return nil
	}

	return h.Header(
		h.Class("page-header"),
		h.Div(
			h.Class("page-header-title"),
			p.Title,
		),
		h.Div(
			h.Class("page-header-actions"),
			g.Group(p.Actions),
		),
	)
}
