package usersviews

import (
	"app/components"
	"app/layout"
	"app/models/usermodel"
	"app/reqcontext"
	"fmt"

	g "github.com/maragudk/gomponents"
)

type EditUserPageProps struct {
	User usermodel.User
	Ctx  reqcontext.ReqContext
}

func EditUserPage(p *EditUserPageProps) g.Node {

	content := g.Group([]g.Node{
		editUserForm(&editUserFormProps{
			user: p.User,
		}),
	})

	return layout.Page(layout.PageProps{
		Title:   fmt.Sprintf("Edit User: %s", p.User.Username),
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/src/users/editUser.css"),
		},
	})
}
