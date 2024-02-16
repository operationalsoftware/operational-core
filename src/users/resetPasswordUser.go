package users

import (
	"app/components"
	"app/db"
	"app/layout"
	userModel "app/src/users/model"
	"app/utils"
	"fmt"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type resetPasswordUserViewProps struct {
	Id  int
	Ctx utils.Context
}

func resetPasswordUserView(p *resetPasswordUserViewProps) g.Node {
	dbInstance := db.UseDB()
	user, err := userModel.ByID(dbInstance, p.Id)

	if err != nil {
		fmt.Println("Error:", err)
		return g.Text("Error")
	}

	resetPasswordContent := g.Group([]g.Node{
		g.If(
			user.UserID != 1,
			components.Form(
				h.ID("reset-password-form"),
				hx.Post(""),
				hx.Target("#submission-error"),
				hx.Swap("outerHTML"),

				h.Div(
					passwordInput(&passwordInputProps{}),
				),

				h.Div(
					confirmPasswordInput(&confirmPasswordInputProps{}),
				),

				h.Div(
					h.ID("submission-error"),
					components.InputHelper(&components.InputHelperProps{
						Label: "",
						Type:  components.InputHelperTypeError,
					}),
				),

				components.Button(&components.ButtonProps{}, h.Type("submit"), g.Text("Submit")),
			),
		),
		components.InlineStyle("/src/users/resetPasswordUser.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Reset Password",
		Content: resetPasswordContent,
		Ctx:     p.Ctx,
	})
}
