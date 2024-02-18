package users

import (
	"app/components"
	"app/layout"
	userModel "app/src/users/model"
	"app/utils"

	g "github.com/maragudk/gomponents"
)

type resetPasswordViewProps struct {
	User userModel.User
	Ctx  utils.Context
}

func resetPasswordView(p *resetPasswordViewProps) g.Node {

	resetPasswordContent := g.Group([]g.Node{
		resetPasswordForm(&resetPasswordFormProps{
			userID: p.User.UserID,
		}),

		components.InlineStyle("/src/users/resetPassword.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Reset Password: " + p.User.Username,
		Content: resetPasswordContent,
		Ctx:     p.Ctx,
	})
}
