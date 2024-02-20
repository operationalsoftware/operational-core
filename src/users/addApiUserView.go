package users

import (
	"app/components"
	"app/layout"
	"app/utils"

	g "github.com/maragudk/gomponents"
)

type addUserAPIViewProps struct {
	Ctx utils.Context
}

func addUserAPIView(p *addUserAPIViewProps) g.Node {
	content := g.Group([]g.Node{

		addApiUserForm(&addApiUserFormProps{}),
		components.InlineStyle("/src/users/addApiUser.css"),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add New API User",
		Content: content,
	})
}
