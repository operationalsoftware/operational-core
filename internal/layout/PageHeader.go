package layout

import (
	"app/internal/components"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
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
				components.Button(&components.ButtonProps{
					ButtonType: components.ButtonPrimary,
					Link:       p.BackToLink,
					Classes: c.Classes{
						"page-header-back-button": true,
					},
				},
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
