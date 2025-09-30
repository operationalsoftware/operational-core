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

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type UserPageProps struct {
	Ctx  reqcontext.ReqContext
	User model.User
}

func UserPage(p *UserPageProps) g.Node {

	user := p.User

	userContent := g.Group([]g.Node{
		userActions(user.UserID),

		h.Div(
			h.Class("two-column-flex"),
			userDetails(user),
			userTeamsTable(user.Teams),
		),

		h.Div(
			h.H3(g.Text("Permissions")),
			permissionsDisplayPartial(user.Permissions),
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

func userActions(userID int) g.Node {
	return g.If(
		userID != 1,
		h.Div(
			h.Class("button-container"),
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Classes: c.Classes{
					"edit-button": true,
				},
				Link: fmt.Sprintf("/users/%d/edit", userID),
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
				Link: fmt.Sprintf("/users/%d/reset-password", userID),
			},
				g.Text("Reset Password"),
			),
		),
	)
}

func userDetails(user model.User) g.Node {

	sessionDuration := nilsafe.Int(user.SessionDurationMinutes)
	if sessionDuration == 0 {
		sessionDuration = int(cookie.DefaultSessionDurationMinutes.Minutes())
	}

	return h.Div(
		h.H3(g.Text("User Details")),
		h.Div(
			h.Class("properties-grid"),
			g.If(
				!user.IsAPIUser,
				g.Group([]g.Node{
					h.Span(h.Strong(g.Text("First Name"))),
					h.Span(g.Text(nilsafe.Str(user.FirstName))),

					h.Span(h.Strong(g.Text("Last Name"))),
					h.Span(g.Text(nilsafe.Str(user.LastName))),

					h.Span(h.Strong(g.Text("Email"))),
					h.Span(g.Text(nilsafe.Str(user.Email))),

					h.Span(h.Strong(g.Text("Session Duration In Minutes"))),
					h.Span(g.Text(strconv.Itoa(sessionDuration))),
				}),
			),

			h.Span(h.Strong(g.Text("Username"))),
			h.Span(g.Text(user.Username)),
		),
	)
}

func userTeamsTable(userTeams []model.UserTeam) g.Node {

	var columns = components.TableColumns{
		{TitleContents: g.Text("Team Name")},
		{TitleContents: g.Text("Role")},
	}

	var tableRows components.TableRows
	for _, ut := range userTeams {
		cells := []components.TableCell{
			{Contents: h.A(
				h.Href(fmt.Sprintf("/teams/%d", ut.TeamID)),
				g.Text(ut.TeamName),
			)},
			{Contents: g.Text(ut.Role)},
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
		})
	}

	return h.Div(
		h.H3(g.Text("Teams")),
		components.Table(&components.TableProps{
			Columns: columns,
			Rows:    tableRows,
		}),
	)
}
