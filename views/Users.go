package views

import (
	"fmt"
	"operationalcore/components"
	"operationalcore/db"
	"operationalcore/layout"
	"operationalcore/model"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type UserRenderer interface {
	Render() map[string]components.RenderedCell
}

type CustomUser model.User

func (u CustomUser) Render() map[string]components.RenderedCell {
	return map[string]components.RenderedCell{
		"username": {
			Content: h.A(
				g.Text(u.Username),
				g.Attr("href",
					fmt.Sprintf("/users/%d", u.UserId))),
			Attributes: []g.Node{},
			Classes: c.Classes{
				"table-link": true,
			},
		},
		"firstName": {
			Content: g.Text(u.FirstName.String),
		},
		"lastName": {
			Content: g.Text(u.LastName.String),
		},
		"email": {
			Content: g.Text(u.Email.String),
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
			Key:      "username",
			Sortable: true,
		},
		{
			Name:     "First Name",
			Key:      "firstName",
			Sortable: true,
		},
		{
			Name: "Last Name",
			Key:  "lastName",
		},
		{
			Name: "Email",
			Key:  "email",
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
