package usersviews

import (
	"app/components"
	"app/layout"
	"app/models/usermodel"
	"app/reqcontext"

	g "github.com/maragudk/gomponents"
)

type ResetPasswordPageProps struct {
	User usermodel.User
	Ctx  reqcontext.ReqContext
}

func ResetPasswordPage(p *ResetPasswordPageProps) g.Node {

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
