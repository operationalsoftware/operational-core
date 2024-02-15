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
	"golang.org/x/exp/slices"
)

type editUserViewProps struct {
	Id  int
	Ctx utils.Context
}

var userRoleOptions = []components.CheckboxOption{
	{
		Value: "User Admin",
		Label: "User Admin",
	},
}

func editUserView(p *editUserViewProps) g.Node {

	dbInsance := db.UseDB()
	user, err := userModel.ByID(dbInsance, p.Id)

	// update user role options, checked if user has the role
	for idx, option := range userRoleOptions {
		// check if user.Roles contains option.Value
		if slices.Contains(user.Roles, option.Value) {
			userRoleOptions[idx].Checked = true
		}
	}

	if err != nil {
		fmt.Println("Error:", err)
		return g.Text("Error")
	}

	editUserContent := g.Group([]g.Node{
		components.Form(
			h.ID("edit-user-form"),
			hx.Post(""),

			g.If(
				!user.IsAPIUser,
				h.Div(
					firstNameInput(&firstNameInputProps{
						Value: user.FirstName.String,
					}),
				),
			),

			g.If(
				!user.IsAPIUser,
				h.Div(
					lastNameInput(&lastNameInputProps{
						Value: user.LastName.String,
					}),
				),
			),

			h.Div(
				usernameInput(&usernameInputProps{
					Value: user.Username,
				}),
			),

			g.If(
				!user.IsAPIUser,
				h.Div(
					emailInput(&emailInputProps{
						Value: user.Email.String,
					}),
				),
			),

			h.Div(
				components.CheckboxGroup(&components.CheckboxGroupProps{
					Name:    "roles",
					Label:   "Roles",
					Options: userRoleOptions,
				}),
			),

			h.Div(
				h.ID("submission-error"),
				components.InputHelper(&components.InputHelperProps{
					Label: "",
					Type:  components.InputHelperTypeError,
				},
				),
			),

			components.Button(&components.ButtonProps{}, h.Type("submit"), g.Text("Submit")),
		),

		components.InlineStyle("/src/users/editUser.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   fmt.Sprintf("Edit User: %s", user.Username),
		Content: editUserContent,
		Ctx:     p.Ctx,
	})
}
