package layout

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
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

			classes := c.Classes{
				"breadcrumb": true,
			}
			return Li(
				classes,
				g.If(index != 1, Span(Class("divider"), g.Text("/"))),
				g.If(index == len(crumbs), Span(g.Text(p.Title))),
				g.If(index != len(crumbs), A(
					h.Href(link),
					g.If(p.Icon != "", o.Icon(&o.IconProps{
						Identifier: p.Icon,
					})),
					Span(g.Text(p.Title)),
				)),
			)
		}),
	)
}

func breadcrumbs(items []Crumb) g.Node {
	return Nav(
		Aria("label", "breadcrumbs"),
		Ol(
			Class("breadcrumbs"),
			renderCrumb(items),
		),
		o.InlineStyle(
			Assets, "/breadcrumbs.css",
		),
	)
}
