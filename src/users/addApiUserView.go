package users

import (
	"app/components"
	"app/layout"
	reqContext "app/reqcontext"

	g "github.com/maragudk/gomponents"
)

type addUserAPIViewProps struct {
	Ctx reqContext.ReqContext
}

func addUserAPIView(p *addUserAPIViewProps) g.Node {
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
