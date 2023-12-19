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

var viewUserCrumb layout.Crumb = layout.Crumb{
	Text:     "View User",
	UrlToken: "view",
}

type UserRenderer interface {
	Render() map[string]components.RenderedCell
}

type CustomUser model.User

func (u CustomUser) Render() map[string]components.RenderedCell {
	return map[string]components.RenderedCell{
		"userId": {
			Content: g.Text(fmt.Sprint(u.UserId)),
			Attributes: []g.Node{
				h.StyleEl(
					h.StyleAttr("text-align: right"),
				),
			},
			Classes: c.Classes{
				"table-cell": true,
			},
		},
		"username": {
			Content: g.Text(u.Username),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
		"firstName": {
			Content: g.Text(u.FirstName),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
		"lastName": {
			Content: g.Text(u.LastName),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
		"email": {
			Content: g.Text(u.Email),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
	}
}

func ViewUser() g.Node {
	crumbs := []layout.Crumb{
		viewUserCrumb,
	}

	dbInstance, _ := sql.Open("sqlite3", "./db/operationalcore.db")
	users := model.GetUsers(dbInstance)

	var data []components.TableRowRenderer
	for _, user := range users {
		data = append(data, CustomUser(user))
	}

	var columns = []components.TableColumn{
		{
			Name: "User ID",
			Key:  "userId",
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
		{
			Name: "Username",
			Key:  "username",
		},
	}

	viewUserContent := g.Group([]g.Node{
		h.H1(g.Text("View User")),
		components.Table(&components.TableProps{
			Columns: columns,
			Data:    data,
		}),
	})

	return layout.Page(layout.PageParams{
		Title:   "View User",
		Content: viewUserContent,
		Crumbs:  crumbs,
	})
}
