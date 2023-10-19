package views

import (
	"operationalcore/layout"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

var contactsCrumb layout.Crumb = layout.Crumb{
	Text:     "Contacts",
	UrlToken: "contacts",
}

func Contacts() g.Node {
	crumbs := []layout.Crumb{
		contactsCrumb,
	}

	contactsContent := g.Group([]g.Node{
		h.H1(
			g.Text("Contacts"),
		),
	})

	return layout.Page(
		layout.PageParams{
			Title:   "Contacts",
			Crumbs:  crumbs,
			Content: contactsContent,
		},
	)
}
