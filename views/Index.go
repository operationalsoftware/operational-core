package views

import (
	"operationalcore/layout"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

var indexCrumb layout.Crumb = layout.Crumb{
	Text: "Home",
	UrlToken: "",
}

func Index() g.Node {
	crumbs := []layout.Crumb{
		indexCrumb,
		{
			Text: "Test",
			UrlToken: "test",
		},
	}

	indexContent := H1(g.Text("Operational Core Home"))

	return layout.Page(layout.PageParams{
		Title: "Home",
		Content: indexContent,
		Crumbs: crumbs,
	})
}
