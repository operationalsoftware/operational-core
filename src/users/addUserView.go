package users

import (
	"app/components"
	"app/layout"
	"app/utils"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type addUserViewProps struct {
	Ctx utils.Context
}

func addUserView(p *addUserViewProps) g.Node {
	content := g.Group([]g.Node{
		components.Form(
			h.ID("add-user-form"),
			hx.Post(""),
			hx.Target("#submission-error"),
			hx.Swap("outerHTML"),

			h.Div(
				firstNameInput(&firstNameInputProps{}),
			),

			h.Div(
				lastNameInput(&lastNameInputProps{}),
			),

			h.Div(
				usernameInput(&usernameInputProps{}),
			),

			h.Div(
				emailInput(&emailInputProps{}),
			),

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

		components.InlineStyle("/src/users/addUser.css"),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add User",
		Content: content,
	})
}
