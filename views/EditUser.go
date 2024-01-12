package views

import (
	o "operationalcore/components"
	"operationalcore/db"
	"operationalcore/layout"
	"operationalcore/model"
	"operationalcore/partials"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

func EditUser(id string) g.Node {

	dbInsance := db.UseDB()
	user := model.GetUser(dbInsance, id)

	editUserContent := g.Group([]g.Node{

		o.Form(
			h.ID("edit-user-form"),
			ghtmx.Post("/users/edit/"+id),

			h.Div(
				partials.CreateUserFirstNameInput(&partials.CreateUserFirstNameInputProps{
					Value: user.FirstName,
				}),
			),

			h.Div(
				partials.CreateUserLastNameInput(&partials.CreateUserLastNameInputProps{
					Value: user.LastName,
				}),
			),

			h.Div(
				partials.CreateUserUsernameInput(&partials.CreateUserUsernameInputProps{
					Value: user.Username,
				}),
			),

			h.Div(
				partials.CreateUserEmailInput(&partials.CreateUserEmailInputProps{
					Value: user.Email,
				}),
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
		Title:   "Edit User",
		Content: editUserContent,
	})
}
