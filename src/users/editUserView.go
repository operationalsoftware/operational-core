package users

import (
	"app/components"
	"app/layout"
	userModel "app/src/users/model"
	"app/utils"
	"fmt"

	g "github.com/maragudk/gomponents"
)

type editUserViewProps struct {
	User userModel.User
	Ctx  utils.Context
}

func editUserView(p *editUserViewProps) g.Node {

	content := g.Group([]g.Node{
		editUserForm(&editUserFormProps{
			user: p.User,
		}),
		components.InlineStyle("/src/users/editUser.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   fmt.Sprintf("Edit User: %s", p.User.Username),
		Content: content,
		Ctx:     p.Ctx,
	})
}
