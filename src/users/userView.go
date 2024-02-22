package users

import (
	"app/components"
	"app/db"
	"app/layout"
	userModel "app/src/users/model"
	"app/utils"
	"fmt"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type userViewProps struct {
	Id  int
	Ctx utils.Context
}

func userView(p *userViewProps) g.Node {
	dbInstance := db.UseDB()
	user, err := userModel.ByID(dbInstance, p.Id)

	if err != nil {
		fmt.Println("Error:", err)
		return g.Text("Error")
	}

	userContent := g.Group([]g.Node{
		g.If(
			user.UserID != 1,
			h.Div(
				h.Class("button-container"),
				components.Button(&components.ButtonProps{
					ButtonType: "primary",
					Classes: c.Classes{
						"edit-button": true,
					},
					Link: fmt.Sprintf("/users/%d/edit", p.Id),
				},
					components.Icon(&components.IconProps{
						Identifier: "pencil",
					}),
				),
				components.Button(&components.ButtonProps{
					ButtonType: "primary",
					Classes: c.Classes{
						"reset-pw-button": true,
					},
					Link: fmt.Sprintf("/users/%d/reset-password", p.Id),
				},
					g.Text("Reset Password"),
				),
			),
		),
		h.Div(
			h.H3(g.Text("User Details")),
			h.Div(
				h.Class("properties-grid"),
				g.If(
					!user.IsAPIUser,
					g.Group([]g.Node{

						h.Span(
							h.Strong(g.Text("First Name")),
						),
						h.Span(
							g.Text(user.FirstName.String),
						),
						h.Span(
							h.Strong(g.Text("Last Name")),
						),
						h.Span(
							g.Text(user.LastName.String),
						),
						h.Span(
							h.Strong(g.Text("Email")),
						),
						h.Span(
							g.Text(user.Email.String),
						),
					}),
				),
				h.Span(
					h.Strong(g.Text("Username")),
				),
				h.Span(
					g.Text(user.Username),
				),
			),
		),

		h.Div(
			h.H3(g.Text("Permissions")),
			g.If(
				user.Permissions.UserAdmin.Access,
				h.Div(
					h.H4(h.Class("permission-group-title"), g.Text("User Admin")),
					h.Ul(
						h.Class("permission-group-list"),
						h.Li(g.Text(getPermissionDescription("UserAdmin", "Access"))),
					),
				),
			),
		),

		components.InlineStyle("/src/users/user.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "View User",
		Content: userContent,
		Ctx:     p.Ctx,
	})
}
