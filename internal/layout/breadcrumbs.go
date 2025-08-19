package layout

import (
	"app/internal/components"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type Breadcrumb struct {
	Title          string
	URLPart        string
	IconIdentifier string
}

func breadcrumbs(breadcrumbs []Breadcrumb) g.Node {

	index := 0
	link := "/"
	numBreadcrumbs := len(breadcrumbs)
	crumbNodes := g.Group(
		g.Map(breadcrumbs, func(b Breadcrumb) g.Node {
			index++

			if b.URLPart != "" {
				if link != "/" {
					link += "/"
				}
				link += b.URLPart
			}

			var crumbContent = g.Group([]g.Node{
				g.If(b.IconIdentifier != "", components.Icon(&components.IconProps{
					Identifier: b.IconIdentifier,
				})),
				h.Span(g.Text(b.Title)),
			})

			return h.Li(
				g.If(index > 1, h.Span(h.Class("divider"), g.Text("/"))),
				g.If(index == numBreadcrumbs, h.Div(crumbContent)),
				g.If(index < numBreadcrumbs, h.A(
					h.Href(link),
					crumbContent,
				)),
			)
		}),
	)

	return h.Nav(
		h.Class("breadcrumbs"),
		h.Aria("label", "breadcrumbs"),
		h.Ol(
			h.Class("breadcrumbs"),
			crumbNodes,
		),
	)
}
