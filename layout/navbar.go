package layout

import (
	o "operationalcore/components"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func navbar() g.Node {
	return Nav(ID("navbar"),
		o.InlineStyle(
			Assets, "/navbar.css",
		),
		Div(
			Class("nav_links-container"),
			Div(
				Class("logo-container"),
				A(Href("/"),
					Img(
						Alt("Logo"),
						Src("/static/img/logo.png"),
					),
				),
			),
			Div(
				Class("nav_links"),
				o.AppGallery(),
			),
		),
	)
}
