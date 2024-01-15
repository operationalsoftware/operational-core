package views

import (
	"operationalcore/components"
	"operationalcore/db"
	"operationalcore/layout"
	"operationalcore/model"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type UserProps struct {
	Id  string
	Ctx utils.Context
}

func User(p *UserProps) g.Node {

	dbInstance := db.UseDB()
	user := model.GetUser(dbInstance, p.Id)

	userContent := g.Group([]g.Node{
		h.Div(
			h.Class("container"),
			h.Div(
				h.Class("grid-table"),
				h.Span(
					h.Strong(g.Text("First Name")),
				),
				h.Span(
					g.Text(user.FirstName),
				),
				h.Span(
					h.Strong(g.Text("Last Name")),
				),
				h.Span(
					g.Text(user.LastName),
				),
				h.Span(
					h.Strong(g.Text("Email")),
				),
				h.Span(
					g.Text(user.Email),
				),
				h.Span(
					h.Strong(g.Text("Username")),
				),
				h.Span(
					g.Text(user.Username),
				),
			),
			h.A(
				h.Class("edit-btn"),
				g.Attr("href", "/users/edit/"+p.Id),
				components.Icon(&components.IconProps{
					Identifier: "pencil",
				}),
			),
		),

		components.InlineStyle(Assets, "/User.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "View User",
		Content: userContent,
		Ctx:     p.Ctx,
	})
}
