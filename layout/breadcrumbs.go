package layout

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type Crumb struct {
	Title    string
	LinkPart string
	Icon     string
}

func breadcrumbs(items []Crumb) g.Node {
	item := Crumb{
		Title:    "Login",
		LinkPart: "/login/password",
		Icon:     "",
	}

	return Nav(
		Aria("label", "breadcrumbs"),
		Ol(
			Li(
				A(
					Href(item.LinkPart),
					g.Text(item.Title),
				),
			),
		),
		o.InlineStyle(
			Assets, "/breadcrumbs.css",
		),
	)
}
