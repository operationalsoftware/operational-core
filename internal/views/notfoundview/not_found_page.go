package notfoundview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type NotFoundPageProps struct {
	Ctx reqcontext.ReqContext
}

func NotFoundPage(p *NotFoundPageProps) g.Node {
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
			components.InlineStyle("/internal/views/notfoundview/not_found_page.css"),
		},
	})
}
