package userview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/cookie"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"fmt"
	"strconv"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type UserPageProps struct {
	Ctx  reqcontext.ReqContext
	User model.User
}

func UserPage(p *UserPageProps) g.Node {

	user := p.User

	sessionDuration := nilsafe.Int(user.SessionDurationMinutes)
	if sessionDuration == 0 {
		sessionDuration = int(cookie.DefaultSessionDurationMinutes.Minutes())
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
					Link: fmt.Sprintf("/users/%d/edit", user.UserID),
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
					Link: fmt.Sprintf("/users/%d/reset-password", user.UserID),
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
							g.Text(nilsafe.Str(user.FirstName)),
						),
						h.Span(
							h.Strong(g.Text("Last Name")),
						),
						h.Span(
							g.Text(nilsafe.Str(user.LastName)),
						),
						h.Span(
							h.Strong(g.Text("Email")),
						),
						h.Span(
							g.Text(nilsafe.Str(user.Email)),
						),
						h.Span(
							h.Strong(g.Text("Session Duration In Minutes")),
						),
						h.Span(
							g.Text(strconv.Itoa(sessionDuration)),
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
	})

	return layout.Page(layout.PageProps{
		Title: "View User",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			usersBreadCrumb,
			{IconIdentifier: "account", Title: user.Username},
		},
		Content: userContent,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/userview/user_page.css"),
		},
	})
}
