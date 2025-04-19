package searchview

import (
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type SearchPageProps struct {
	Ctx reqcontext.ReqContext
}

func SearchPage(p SearchPageProps) g.Node {
	return layout.Page(layout.PageProps{
		Title:   "Search",
		Content: g.Node(h.H1(g.Text("Search"))),
		Ctx:     p.Ctx,
	})
}
