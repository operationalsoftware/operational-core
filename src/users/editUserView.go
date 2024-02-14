package users

import (
	"fmt"
	"app/components"
	"app/db"
	"app/layout"
	userModel "app/src/users/model"
	"app/utils"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type editUserViewProps struct {
	Id  int
	Ctx utils.Context
}

func editUserView(p *editUserViewProps) g.Node {

	dbInsance := db.UseDB()
	user, err := userModel.ByID(dbInsance, p.Id)
	options := []components.CheckboxOption{}

	for _, role := range user.Roles {
		if role != "" {
			options = append(options, components.CheckboxOption{
				Value:   role,
				Label:   role,
				Checked: true,
			})
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

			g.If(
				!user.IsAPIUser,
				h.Div(
					components.CheckboxGroup(&components.CheckboxGroupProps{
						Name:    "roles",
						Label:   "Roles",
						Options: options,
					}),
				),
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
