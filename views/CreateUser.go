package views

import (
	o "operationalcore/components"
	"operationalcore/layout"
	"operationalcore/partials"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

var createUserCrumb layout.Crumb = layout.Crumb{
	Text:     "Create",
	UrlToken: "create",
}

func CreateUser() g.Node {
	crumbs := []layout.Crumb{
		createUserCrumb,
	}

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

	return layout.Page(layout.PageParams{
		Title:   "Create User",
		Content: createUserContent,
		Crumbs:  crumbs,
	})
}
