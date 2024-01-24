package views

import (
	"fmt"
	"operationalcore/components"
	"operationalcore/db"
	"operationalcore/layout"
	"operationalcore/model"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type UserProps struct {
	Id  int
	Ctx utils.Context
}

func createUserCrumbs(userId int) []layout.Crumb {
	usersCrumbs := createUsersCrumbs()
	dbInstance := db.UseDB()
	user := model.GetUser(dbInstance, userId)
	return append(usersCrumbs, layout.Crumb{
		LinkPart: fmt.Sprintf("%d", user.UserId),
		Icon:     "",
		Title:    user.Username,
	})
}

func User(p *UserProps) g.Node {
	crumbs := createUserCrumbs(p.Id)

	dbInstance := db.UseDB()
	user := model.GetUser(dbInstance, p.Id)

	userContent := g.Group([]g.Node{
		g.If(
			user.UserId != 1,
			h.Div(
				h.Class("edit-button-container"),
				components.Button(&components.ButtonProps{
					ButtonType: "primary",
					Link:       fmt.Sprintf("/users/%d/edit", p.Id),
				},
					components.Icon(&components.IconProps{
						Identifier: "pencil",
					}),
				),
			),
		),
		h.Div(
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

		components.InlineStyle(Assets, "/User.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "View User",
		Content: userContent,
		Ctx:     p.Ctx,
		Crumbs:  crumbs,
	})
}
