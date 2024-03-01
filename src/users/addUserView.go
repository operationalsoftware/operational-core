package users

import (
	"app/components"
	"app/layout"
	reqContext "app/reqcontext"

	g "github.com/maragudk/gomponents"
)

type addUserViewProps struct {
	Ctx reqContext.ReqContext
}

func addUserView(p *addUserViewProps) g.Node {
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
