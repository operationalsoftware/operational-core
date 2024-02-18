package src

import (
	"app/components"
	"app/layout"
	"app/utils"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type indexViewProps struct {
	Ctx utils.Context
}

func indexView(p *indexViewProps) g.Node {

	indexContent := g.Group([]g.Node{
		h.H1(g.Text("Operational Core Home")),
		components.InlineScript("/src/index.js"),
		components.Card(
			h.H1(g.Text("Card Title")),
		),
	})

	return layout.Page(layout.PageProps{
		Title:   "Home",
		Content: indexContent,
		Ctx:     p.Ctx,
	})
}
