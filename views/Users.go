package views

import (
	"database/sql"
	"fmt"
	"operationalcore/components"
	"operationalcore/layout"
	"operationalcore/model"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

var usersCrumb layout.Crumb = layout.Crumb{
	Text:     "View Users",
	UrlToken: "users",
}

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
			Content: g.Text(u.FirstName),
		},
		"lastName": {
			Content: g.Text(u.LastName),
		},
		"email": {
			Content: g.Text(u.Email),
		},
	}
}

func Users() g.Node {
	crumbs := []layout.Crumb{
		usersCrumb,
	}

	dbInstance, _ := sql.Open("sqlite3", "./db/operationalcore.db")
	users := model.GetUsers(dbInstance)

	var data []components.TableRowRenderer
	for _, user := range users {
		data = append(data, CustomUser(user))
	}

	var columns = []components.TableColumn{
		{
			Name: "Username",
			Key:  "username",
		},
		{
			Name: "First Name",
			Key:  "firstName",
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
		h.H1(
			g.Text("View Users"),
		),
		components.Table(&components.TableProps{
			Columns: columns,
			Data:    data,
		}),
	})

	return layout.Page(layout.PageParams{
		Title:   "View Users",
		Content: viewUserContent,
		Crumbs:  crumbs,
	})
}
