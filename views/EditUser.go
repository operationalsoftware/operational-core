package views

import (
	"fmt"
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
	Id  int
	Ctx utils.Context
}

func createEditUserCrumbs(userId int) []layout.Crumb {
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

	editUserContent := o.Form(
		h.ID("edit-user-form"),
		ghtmx.Post(""),

		g.If(
			!user.IsAPIUser,
			h.Div(
				partials.CreateUserFirstNameInput(&partials.CreateUserFirstNameInputProps{
					Value: user.FirstName.String,
				}),
			),
		),

		g.If(
			!user.IsAPIUser,
			h.Div(
				partials.CreateUserLastNameInput(&partials.CreateUserLastNameInputProps{
					Value: user.LastName.String,
				}),
			),
		),

		h.Div(
			partials.CreateUserUsernameInput(&partials.CreateUserUsernameInputProps{
				Value: user.Username,
			}),
		),

		g.If(
			!user.IsAPIUser,
			h.Div(
				partials.CreateUserEmailInput(&partials.CreateUserEmailInputProps{
					Value: user.Email.String,
				}),
			),
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
	)

	return layout.Page(layout.PageProps{
		Title:   fmt.Sprintf("Edit User: %s", user.Username),
		Content: editUserContent,
		Ctx:     p.Ctx,
		Crumbs:  crumbs,
	})
}
