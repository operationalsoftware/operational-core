package views

import (
	"fmt"
	"operationalcore/components"
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

	editUserContent := g.Group([]g.Node{
		o.Form(
			h.ID("edit-user-form"),
			ghtmx.Post(""),

			g.If(
				!user.IsAPIUser,
				h.Div(
					partials.UserFormFirstNameInput(&partials.UserFormFirstNameInputProps{
						Value: user.FirstName.String,
					}),
				),
			),

			g.If(
				!user.IsAPIUser,
				h.Div(
					partials.UserFormLastNameInput(&partials.UserFormLastNameInputProps{
						Value: user.LastName.String,
					}),
				),
			),

			h.Div(
				partials.UserFormUsernameInput(&partials.UserFormUsernameInputProps{
					Value: user.Username,
				}),
			),

			g.If(
				!user.IsAPIUser,
				h.Div(
					partials.UserFormEmailInput(&partials.UserFormEmailInputProps{
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
		),

		components.InlineStyle(Assets, "/AddUser.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   fmt.Sprintf("Edit User: %s", user.Username),
		Content: editUserContent,
		Ctx:     p.Ctx,
		Crumbs:  crumbs,
	})
}
