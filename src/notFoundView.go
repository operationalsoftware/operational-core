package src

import (
	"app/components"
	"app/layout"
	reqContext "app/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type notFoundViewProps struct {
	Ctx reqContext.ReqContext
}

func notFoundView(p *notFoundViewProps) g.Node {
	notFoundContent := g.Group([]g.Node{
		h.H1(g.Text("404 - Page Not Found")),
		components.Card(
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
				components.Button(&components.ButtonProps{
					ButtonType: components.ButtonSecondary,
				}, g.Text("Go to home page")),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Title:   "404 - Not Found",
		Content: notFoundContent,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/src/notFound.css"),
		},
	})
}
