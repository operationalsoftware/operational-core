package userview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/models"
	"app/pkg/reqcontext"
	"fmt"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type UserPageProps struct {
	Id   int
	Ctx  reqcontext.ReqContext
	User models.User
}

func UserPage(p *UserPageProps) g.Node {

	userContent := g.Group([]g.Node{
		g.If(
			p.User.UserID != 1,
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
					!p.User.IsAPIUser,
					g.Group([]g.Node{

						h.Span(
							h.Strong(g.Text("First Name")),
						),
						h.Span(
							g.Text(*p.User.FirstName),
						),
						h.Span(
							h.Strong(g.Text("Last Name")),
						),
						h.Span(
							g.Text(*p.User.LastName),
						),
						h.Span(
							h.Strong(g.Text("Email")),
						),
						h.Span(
							g.Text(*p.User.Email),
						),
					}),
				),
				h.Span(
					h.Strong(g.Text("Username")),
				),
				h.Span(
					g.Text(p.User.Username),
				),
			),
		),

		permissionsCheckboxesPartial(p.User.Permissions),
	})

	return layout.Page(layout.PageProps{
		Title:   "View User",
		Content: userContent,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/routes/users/usersviews/user.css"),
		},
	})
}
