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

type UserRenderer interface {
	Render() map[string]components.RenderedCell
}

type CustomUser model.User

func (u CustomUser) Render() map[string]components.RenderedCell {
	return map[string]components.RenderedCell{
		"Username": {
			Content: h.A(
				g.Text(u.Username),
				g.Attr("href",
					fmt.Sprintf("/users/%d", u.UserId))),
		},
		"FirstName": {
			Content: g.Text(u.FirstName.String),
		},
		"LastName": {
			Content: g.Text(u.LastName.String),
		},
		"Email": {
			Content: g.Text(u.Email.String),
		},
		"Created": {
			Content: g.Text(u.Created.Format("2006-01-02 15:04:05")),
		},
		"LastLogin": {
			Content: g.If(u.LastLogin.Valid, g.Text(u.LastLogin.Time.Format("2006-01-02 15:04:05"))),
		},
	}
}

type UsersProps struct {
	Ctx utils.Context
}

func createUsersCrumbs() []layout.Crumb {
	indexCrumbs := createIndexCrumbs()
	return append(indexCrumbs, layout.Crumb{
		LinkPart: "users",
		Icon:     "",
		Title:    "Users",
	})
}

func Users(p *UsersProps) g.Node {

	db := db.UseDB()
	crumbs := createUsersCrumbs()

	users := model.GetUsers(db)

	var data []components.TableRowRenderer
	for _, user := range users {
		data = append(data, CustomUser(user))
	}

	var columns = []components.TableColumn{
		{
			Name:     "Username",
			Key:      "Username",
			Sortable: true,
		},
		{
			Name:     "First Name",
			Key:      "FirstName",
			Sortable: true,
		},
		{
			Name: "Last Name",
			Key:  "LastName",
		},
		{
			Name: "Email",
			Key:  "Email",
		},
		{
			Name: "Created",
			Key:  "Created",
		},
		{
			Name: "Last Login",
			Key:  "LastLogin",
		},
	}

	viewUserContent := g.Group([]g.Node{
		h.Div(
			h.Class("add-button-container"),
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Link:       "/users/add",
			},
				components.Icon(&components.IconProps{
					Identifier: "plus",
				}),
				g.Text("User"),
			),
		),
		components.Table(&components.TableProps{
			Columns: columns,
			Data:    data,
		}),
		components.InlineStyle(Assets, "/Users.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Users",
		Content: viewUserContent,
		Ctx:     p.Ctx,
		Crumbs:  crumbs,
	})
}
