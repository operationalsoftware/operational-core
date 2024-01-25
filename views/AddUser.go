package views

import (
	"operationalcore/components"
	o "operationalcore/components"
	"operationalcore/layout"
	"operationalcore/partials"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type AddUserProps struct {
	Ctx utils.Context
}

func createAddUserCrumbs() []layout.Crumb {
	usersCrumbs := createUsersCrumbs()
	return append(usersCrumbs, layout.Crumb{
		LinkPart: "add",
		Icon:     "",
		Title:    "Add",
	})
}

func AddUser(p *AddUserProps) g.Node {

	crumbs := createAddUserCrumbs()

	addUserContent := g.Group([]g.Node{
		o.Form(
			h.ID("add-user-form"),
			ghtmx.Post(""),
			ghtmx.Target("#submission-error"),
			ghtmx.Swap("outerHTML"),

			h.Div(
				partials.UserFormFirstNameInput(&partials.UserFormFirstNameInputProps{}),
			),

			h.Div(
				partials.UserFormLastNameInput(&partials.UserFormLastNameInputProps{}),
			),

			h.Div(
				partials.UserFormUsernameInput(&partials.UserFormUsernameInputProps{}),
			),

			h.Div(
				partials.UserFormEmailInput(&partials.UserFormEmailInputProps{}),
			),

			h.Div(
				partials.UserFormPasswordInput(&partials.UserFormPasswordInputProps{}),
			),

			h.Div(
				partials.UserFormConfirmPasswordInput(&partials.UserFormConfirmPasswordInputProps{}),
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

		components.InlineStyle(Assets, "/AddUser.css"),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Crumbs:  crumbs,
		Title:   "Add User",
		Content: addUserContent,
	})
}
