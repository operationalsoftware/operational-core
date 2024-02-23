package users

import (
	"app/components"
	"app/db"
	"app/layout"
	reqContext "app/reqcontext"
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
	Ctx reqContext.ReqContext
}

func indexView(p *indexViewProps) g.Node {

	db := db.UseDB()

	sort := utils.Sort{}
	sort.ParseQueryParam(p.Ctx.Req.URL.Query().Get("sort"), userModel.ListSortableKeys)
	users, err := userModel.List(db, userModel.ListQuery{
		Sort: sort,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return g.Text("Error")
	}

	var data []components.TableRowRenderer
	for _, user := range users {
		data = append(data, userWithRender(user))
	}

	var columns = components.TableColumns{
		{
			Name: "Username",
			Key:  "Username",
		},
		{
			Name: "First Name",
			Key:  "FirstName",
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
		components.Table(&components.TableProps{
			Columns:      columns,
			SortableKeys: userModel.ListSortableKeys,
			Data:         data,
			HXGetPath:    "/users",
			HXSelect:     ".table-container > *",
			UrlQuery:     p.Ctx.Req.URL.Query(),
		},
			h.ID("users-table"),
		),
		components.InlineStyle("/src/users/index.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Users",
		Content: content,
		Ctx:     p.Ctx,
	})
}
