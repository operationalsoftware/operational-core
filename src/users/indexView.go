package users

import (
	"app/components"
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

type indexViewProps struct {
	Ctx       reqContext.ReqContext
	users     []userModel.User
	userCount int
	sort      utils.Sort
	page      int
	pageSize  int
	myFilter  string
}

func indexView(p *indexViewProps) g.Node {

	var columns = components.TableColumns{
		{
			TitleContents: g.Text("Username"),
			SortKey:       "Username",
		},
		{
			TitleContents: g.Text("First Name"),
			SortKey:       "FirstName",
		},
		{
			TitleContents: g.Text("Last Name"),
			SortKey:       "LastName",
		},
		{
			TitleContents: g.Text("Email"),
			SortKey:       "Email",
		},
		{
			TitleContents: g.Text("Created"),
			SortKey:       "Created",
		},
		{
			TitleContents: g.Text("Last Login"),
			SortKey:       "LastLogin",
		},
	}

	var tableRows components.TableRows
	for _, u := range p.users {
		tableRows = append(tableRows, components.TableRow{
			Cells: []components.TableCell{
				{
					Contents: h.A(
						g.Text(u.Username),
						g.Attr("href",
							fmt.Sprintf("/users/%d", u.UserID))),
				},
				{
					Contents: g.Text(u.FirstName.String),
				},
				{
					Contents: g.Text(u.LastName.String),
				},
				{
					Contents: g.Text(u.Email.String),
				},
				{
					Contents: g.Text(u.Created.Format("2006-01-02 15:04:05")),
				},
				{
					Contents: g.If(u.LastLogin.Valid, g.Text(u.LastLogin.Time.Format("2006-01-02 15:04:05"))),
				},
			},
		})
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

		// form container for table interaction
		h.FormEl(
			h.ID("users-table-form"),
			g.Attr("method", "GET"),

			components.Table(&components.TableProps{
				Columns:      columns,
				SortableKeys: userModel.ListSortableKeys,
				Sort:         p.sort,
				Rows:         tableRows,
				OnChange:     "submitUsersTableForm()",
				Pagination: components.TablePaginationProps{
					TotalRecords:        p.userCount,
					PageSize:            p.pageSize,
					CurrentPage:         p.page,
					CurrentPageQueryKey: "Page",
					PageSizeQueryKey:    "PageSize",
				},
			},
				h.ID("users-table"),
			),
		),

		components.InlineScript("/src/users/index.js"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Users",
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/src/users/index.css"),
		},
	})
}
