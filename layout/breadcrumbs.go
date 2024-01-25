package layout

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type Crumb struct {
	Title    string
	LinkPart string
	Icon     string
}

func renderCrumb(crumbs []Crumb) g.Node {
	index := 0
	linkParts := []string{}

	return g.Group(
		g.Map(crumbs, func(p Crumb) g.Node {
			index++
			linkParts = append(linkParts, p.LinkPart)
			link := ""

			for _, part := range linkParts {
				link += part + "/"
			}

			if len(link) > 1 {
				link = link[:len(link)-1]
			}

			var crumbContent = g.Group([]g.Node{
				g.If(p.Icon != "", o.Icon(&o.IconProps{
					Identifier: p.Icon,
				})),
				h.Span(g.Text(p.Title)),
			})

			return h.Li(
				g.If(index > 1, h.Span(h.Class("divider"), g.Text("/"))),
				g.If(index == len(crumbs), h.Div(crumbContent)),
				g.If(index < len(crumbs), h.A(
					h.Href(link),
					crumbContent,
				)),
			)
		}),
	)
}

func breadcrumbs(items []Crumb) g.Node {
	return h.Nav(
		h.Aria("label", "breadcrumbs"),
		h.Ol(
			h.Class("breadcrumbs"),
			renderCrumb(items),
		),
		o.InlineStyle(
			Assets, "/breadcrumbs.css",
		),
	)
}
