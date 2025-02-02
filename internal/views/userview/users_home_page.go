package userview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/models"
	"app/pkg/appsort"
	"app/pkg/reqcontext"
	"fmt"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type UsersHomePageProps struct {
	Ctx       reqcontext.ReqContext
	Users     []models.User
	UserCount int
	Sort      appsort.Sort
	Page      int
	PageSize  int
	MyFilter  string
}

func UsersHomePage(p *UsersHomePageProps) g.Node {

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
	for _, u := range p.Users {
		tableRows = append(tableRows, components.TableRow{
			Cells: []components.TableCell{
				{
					Contents: h.A(
						g.Text(u.Username),
						g.Attr("href",
							fmt.Sprintf("/users/%d", u.UserID))),
				},
				{
					Contents: g.Text(*u.FirstName),
				},
				{
					Contents: g.Text(*u.LastName),
				},
				{
					Contents: g.Text(*u.Email),
				},
				{
					Contents: g.Text(u.Created.Format("2006-01-02 15:04:05")),
				},
				{
					Contents: g.If(u.LastLogin != nil, g.Text((*u.LastLogin).Format("2006-01-02 15:04:05"))),
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
				Columns: columns,
				Sort:    p.Sort,
				Rows:    tableRows,
				Pagination: &components.TablePaginationProps{
					TotalRecords:        p.UserCount,
					PageSize:            p.PageSize,
					CurrentPage:         p.Page,
					CurrentPageQueryKey: "Page",
					PageSizeQueryKey:    "PageSize",
				},
			},
				h.ID("users-table"),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Title:   "Users",
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/routes/users/usersviews/index.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/routes/users/usersviews/index.js"),
		},
	})
}
