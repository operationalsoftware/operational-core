package layout

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type Crumb struct {
	UrlToken string
	Text     string
	IconPath string
}

func breadcrumbsGenerator() func(Crumb, bool) g.Node {
	url := ""
	return func(item Crumb, active bool) g.Node {
		url = url + "/" + item.UrlToken

		return Li(
			g.If(url != "/", Span(Class("separator"), g.Text("/"))),
			g.If(active, Span(g.Text(item.Text))),
			g.If(!active, A(Href(url), g.Text(item.Text))),
		)
	}
}

func breadcrumbs(items []Crumb) g.Node {

	breadcrumb := breadcrumbsGenerator()

	listItems := []g.Node{}
	for i, item := range items {
		listItems = append(listItems, breadcrumb(item, i == len(items)-1))
	}

	return Nav(
		Aria("label", "breadcrumbs"),
		o.InlineStyle(
			Assets, "/breadcrumbs.css",
		),
		Ol(listItems...),
	)
}
