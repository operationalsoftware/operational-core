package views

import (
	o "operationalcore/components"
	"operationalcore/db"
	"operationalcore/layout"
	"operationalcore/model"
	"operationalcore/partials"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type EditUserProps struct {
	Id  string
	Ctx utils.Context
}

func createEditUserCrumbs(userId string) []layout.Crumb {
	existingCrumbs := createUserCrumbs(userId)
	return append(existingCrumbs, layout.Crumb{
		LinkPart: "edit",
		Icon:     "",
		Title:    "Edit",
	})
}

func EditUser(p *EditUserProps) g.Node {
	crumbs := createEditUserCrumbs(p.Id)

	dbInsance := db.UseDB()
	user := model.GetUser(dbInsance, p.Id)

	editUserContent := g.Group([]g.Node{

		o.Form(
			h.ID("edit-user-form"),
			ghtmx.Post(""),

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

	return layout.Page(layout.PageProps{
		Title:   "Edit User",
		Content: editUserContent,
		Ctx:     p.Ctx,
		Crumbs:  crumbs,
	})
}
