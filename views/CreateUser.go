package views

import (
	o "operationalcore/components"
	"operationalcore/layout"
	"operationalcore/partials"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type CreateUserProps struct {
	Ctx utils.Context
}

func CreateUser(p *CreateUserProps) g.Node {

	createUserContent := g.Group([]g.Node{
		h.H1(g.Text("Form Page")),
		o.Form(
			h.ID("create-user-form"),
			ghtmx.Post("/users/create"),
			ghtmx.Target("#submission-error"),
			ghtmx.Swap("outerHTML"),

			h.Div(
				partials.CreateUserFirstNameInput(&partials.CreateUserFirstNameInputProps{}),
			),

			h.Div(
				partials.CreateUserLastNameInput(&partials.CreateUserLastNameInputProps{}),
			),

			h.Div(
				partials.CreateUserUsernameInput(&partials.CreateUserUsernameInputProps{}),
			),

			h.Div(
				partials.CreateUserEmailInput(&partials.CreateUserEmailInputProps{}),
			),

			h.Div(
				partials.CreateUserPasswordInput(&partials.CreateUserPasswordInputProps{}),
			),

			h.Div(
				partials.CreateUserConfirmPasswordInput(&partials.CreateUserConfirmPasswordInputProps{}),
			),

			h.Div(
				h.ID("submission-error"),
				o.InputHelper(&o.InputHelperProps{
					Label: "",
					Type:  o.InputHelperTypeError,
				},
				),
			),

			o.Button(&o.ButtonProps{}, h.Type("submit"), g.Text("Submit")),
		),
	})

	return layout.Page(layout.PageProps{
		Title:   "Create User",
		Content: createUserContent,
		Ctx:     p.Ctx,
	})
}
