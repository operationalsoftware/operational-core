package views

import (
	o "operationalcore/components"
	"operationalcore/layout"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func Page404() g.Node {
	content404 := g.Group([]g.Node{
		h.H1(g.Text("404 Page")),
		o.Card(
			h.Div(
				h.Class("row"),
				h.Div(
					h.Class("col col-1"),
					h.Div(
						h.Class("img-container"),
						h.Img(
							h.Src("/static/img/logo.png"),
						),
					),
				),
				h.Div(
					h.Class("col col-2"),
					h.P(h.Class("text-404"), g.Text("The page you are looking for does not exist.")),
				),
			),

			h.A(
				h.Class("home-btn"),
				h.Href("/"),
				o.Button(&o.ButtonProps{
					ButtonType: o.ButtonSecondary,
				}, g.Text("Go to home page")),
			),
		),
		o.InlineStyle(Assets, "/404.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "404 Page",
		Content: content404,
	})
}
