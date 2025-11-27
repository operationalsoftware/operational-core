package userview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"fmt"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type UsersHomePageProps struct {
	Ctx                 reqcontext.ReqContext
	Users               []model.User
	UserCount           int
	ActiveUserCountLast int
	Sort                appsort.Sort
	Page                int
	PageSize            int
	MyFilter            string
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
		{
			TitleContents: g.Text("Last Active"),
			SortKey:       "LastActive",
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
					Contents: g.Group([]g.Node{
						g.If(u.FirstName == nil, g.Text("\u2013")),
						g.If(u.FirstName != nil, g.Text(nilsafe.Str(u.FirstName))),
					}),
				},
				{
					Contents: g.Group([]g.Node{
						g.If(u.LastName == nil, g.Text("\u2013")),
						g.If(u.LastName != nil, g.Text(nilsafe.Str(u.LastName))),
					}),
				},
				{
					Contents: g.Group([]g.Node{
						g.If(u.Email == nil, g.Text("\u2013")),
						g.If(u.Email != nil, g.Text(nilsafe.Str(u.Email))),
					}),
				},
				{
					Contents: g.Text(u.Created.Format("2006-01-02 15:04:05")),
				},
				{
					Contents: g.Group([]g.Node{
						g.If(u.LastLogin == nil, g.Text("\u2013")),
						g.If(u.LastLogin != nil, g.Text(nilsafe.Time(u.LastLogin).Format("2006-01-02 15:04:05"))),
					}),
				},
				{
					Contents: g.Group([]g.Node{
						g.If(u.LastActive == nil, g.Text("\u2013")),
						g.If(u.LastActive != nil, g.Text(nilsafe.Time(u.LastActive).Format("2006-01-02 15:04:05"))),
					}),
				},
			},
			HREF: fmt.Sprintf("/users/%d", u.UserID),
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

		h.P(
			h.Class("user-stats"),
			g.Text(fmt.Sprintf(
				"Total users: %d, Active in last 30 days: %d",
				p.UserCount,
				p.ActiveUserCountLast,
			)),
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
		Ctx:     p.Ctx,
		Title:   "Users",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			usersBreadCrumb,
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/userview/users_home_page.css"),
		},
	})
}

var usersBreadCrumb = layout.Breadcrumb{
	IconIdentifier: "account-multiple",
	Title:          "Users",
	URLPart:        "users",
}
