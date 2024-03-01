package users

import (
	"app/components"
	"app/layout"
	reqContext "app/reqcontext"
	userModel "app/src/users/model"

	g "github.com/maragudk/gomponents"
)

type resetPasswordViewProps struct {
	User userModel.User
	Ctx  reqContext.ReqContext
}

func resetPasswordView(p *resetPasswordViewProps) g.Node {

	resetPasswordContent := g.Group([]g.Node{
		resetPasswordForm(&resetPasswordFormProps{
			userID: p.User.UserID,
		}),
	})

	return layout.Page(layout.PageProps{
		Title:   "Reset Password: " + p.User.Username,
		Content: resetPasswordContent,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/src/users/resetPassword.css"),
		},
	})
}
