package src

import (
	o "app/components"
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
		o.InlineScript("/src/index.js"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Home",
		Content: indexContent,
		Ctx:     p.Ctx,
	})
}
