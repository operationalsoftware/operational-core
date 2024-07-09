package usersviews

import (
	"app/components"
	"app/layout"
	"app/reqcontext"

	g "github.com/maragudk/gomponents"
)

type AddUserPageProps struct {
	Ctx reqcontext.ReqContext
}

func AddUserPage(p *AddUserPageProps) g.Node {
	content := g.Group([]g.Node{

		addUserForm(&addUserFormProps{}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add New User",
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/src/users/addUser.css"),
		},
	})
}
