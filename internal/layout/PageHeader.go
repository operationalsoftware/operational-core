package layout

import (
	"app/internal/components"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PageHeaderProps struct {
	BackToText string
	BackToLink string
	Title      g.Node
	Actions    []g.Node
}

func PageHeader(p *PageHeaderProps) g.Node {
	if p == nil {
		return nil
	}

	return h.Header(
		h.Class("page-header"),
		g.If(
			p.BackToText != "" && p.BackToLink != "",
			h.Div(
				h.Class("page-header-back"),
				h.A(
					h.Class("button primary page-header-back-button"),
					h.Href(p.BackToLink),
					components.Icon(&components.IconProps{
						Identifier: "arrow-left-thin",
					}),
					g.Text(p.BackToText),
				),
			),
		),
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
