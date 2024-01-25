package views

import (
	o "operationalcore/components"
	"operationalcore/layout"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type IndexProps struct {
	Ctx utils.Context
}

func createIndexCrumbs() []layout.Crumb {
	var crumbs = []layout.Crumb{
		{
			Title:    "Home",
			LinkPart: "",
			Icon:     "home",
		},
	}

	return crumbs
}

func Index(p *IndexProps) g.Node {

	crumbs := createIndexCrumbs()

	indexContent := g.Group([]g.Node{
		h.H1(g.Text("Operational Core Home")),
		o.InlineScript(Assets, "/Index.js"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Home",
		Content: indexContent,
		Ctx:     p.Ctx,
		Crumbs:  crumbs,
	})
}
