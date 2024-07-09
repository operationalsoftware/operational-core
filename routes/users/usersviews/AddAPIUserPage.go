package usersviews

import (
	"app/components"
	"app/layout"
	"app/reqcontext"

	g "github.com/maragudk/gomponents"
)

type AddAPIUserPageProps struct {
	Ctx reqcontext.ReqContext
}

func AddAPIUserPage(p *AddAPIUserPageProps) g.Node {
	content := g.Group([]g.Node{

		addApiUserForm(&addApiUserFormProps{}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add New API User",
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/src/users/addApiUser.css"),
		},
	})
}
