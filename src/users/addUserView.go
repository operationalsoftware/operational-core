package users

import (
	"app/components"
	"app/layout"
	"app/utils"

	g "github.com/maragudk/gomponents"
)

type addUserViewProps struct {
	Ctx utils.Context
}

func addUserView(p *addUserViewProps) g.Node {
	content := g.Group([]g.Node{

		addUserForm(&addUserFormProps{}),

		components.InlineStyle("/src/users/addUser.css"),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add New User",
		Content: content,
	})
}
