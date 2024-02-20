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

type userWithRender userModel.User

func (u userWithRender) Render() map[string]components.RenderedCell {
	return map[string]components.RenderedCell{
		"Username": {
			Content: h.A(
				g.Text(u.Username),
				g.Attr("href",
					fmt.Sprintf("/users/%d", u.UserID))),
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

type indexViewProps struct {
	Ctx utils.Context
}

func indexView(p *indexViewProps) g.Node {

	db := db.UseDB()

	users, err := userModel.List(db)
	if err != nil {
		fmt.Println("Error:", err)
		return g.Text("Error")
	}

	var data []components.TableRowRenderer
	for _, user := range users {
		data = append(data, userWithRender(user))
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

	content := g.Group([]g.Node{
		h.Div(
			h.Class("button-container"),
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Link:       "/users/add",
				Classes: c.Classes{
					"add-user-btn": true,
				},
			},
				components.Icon(&components.IconProps{
					Identifier: "plus",
				}),
				g.Text("User"),
			),
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Link:       "/users/add-api-user",
				Classes: c.Classes{
					"add-api-user-btn": true,
				},
			},
				components.Icon(&components.IconProps{
					Identifier: "plus",
				}),
				g.Text("API User"),
			),
		),
		h.Div(
			h.ID("table"),
			components.Table(&components.TableProps{
				Columns: columns,
				Data:    data,
				HxGet:   "/users/table",
				Sort:    []components.SortItem{},
			}),
		),

		components.InlineStyle("/src/users/index.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Users",
		Content: content,
		Ctx:     p.Ctx,
	})
}
