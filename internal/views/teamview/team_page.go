package teamview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/reqcontext"
	"fmt"
	"strconv"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type TeamPageProps struct {
	Ctx            reqcontext.ReqContext
	Team           model.Team
	TeamUsers      []model.TeamUser
	TeamUsersCount int
	Sort           appsort.Sort
	Page           int
	PageSize       int
}

func TeamPage(p *TeamPageProps) g.Node {

	team := p.Team

	var columns = components.TableColumns{
		{
			TitleContents: g.Text("Username"),
			SortKey:       "Username",
		},
		{
			TitleContents: g.Text("Role"),
			SortKey:       "Role",
		},
		{
			TitleContents: g.Text("Actions"),
		},
	}

	var tableRows components.TableRows
	for _, ai := range p.TeamUsers {

		cells := []components.TableCell{
			{
				Contents: g.Text(ai.Username),
			},
			{
				Contents: g.Text(ai.Role),
			},
			{
				Contents: g.Group([]g.Node{

					h.Div(
						h.Class("andon-actions"),

						components.Button(&components.ButtonProps{
							Size:       "small",
							ButtonType: "danger",
						},
							g.Attr("onclick", "updateAndon(event)"),
							g.Attr("data-id", strconv.Itoa(ai.UserID)),
							g.Attr("data-username", ai.Username),
							g.Attr("data-team-id", strconv.Itoa(ai.TeamID)),
							g.Attr("title", "Remove"),

							components.Icon(&components.IconProps{
								Identifier: "close",
							}),
						),
					),
				}),
			},
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
		})
	}

	content := g.Group([]g.Node{
		h.Div(
			h.Class("button-container"),
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Classes: c.Classes{
					"edit-button": true,
				},
				Link: fmt.Sprintf("/teams/%d/edit", team.TeamID),
			},
				components.Icon(&components.IconProps{
					Identifier: "pencil",
				}),
			),
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Classes: c.Classes{
					"assign-user-button": true,
				},
				Link: fmt.Sprintf("/teams/%d/assign-user", team.TeamID),
			},
				components.Icon(&components.IconProps{
					Identifier: "plus",
				}),
				g.Text("Assign User"),
			),
		),
		h.Div(
			h.H3(g.Text("Team Details")),
			h.Div(
				h.Class("properties-grid"),
				g.Group([]g.Node{
					h.Span(
						h.Strong(g.Text("Team Name")),
					),
					h.Span(
						g.Text(team.TeamName),
					),

					h.Span(
						h.Strong(g.Text("Is Archived?")),
					),
					h.Span(
						g.If(team.IsArchived, g.Text("Yes")),
						g.If(!team.IsArchived, g.Text("No")),
					),
				}),
			),
		),

		components.Table(&components.TableProps{
			Columns: columns,
			Sort:    p.Sort,
			Rows:    tableRows,
			Pagination: &components.TablePaginationProps{
				TotalRecords:        p.TeamUsersCount,
				PageSize:            p.PageSize,
				CurrentPage:         p.Page,
				CurrentPageQueryKey: "Page",
				PageSizeQueryKey:    "PageSize",
			},
		},
			h.ID("team-users-table"),
		),
	})

	return layout.Page(layout.PageProps{
		Title: "Team: " + team.TeamName,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			teamsBreadCrumb,
			{
				IconIdentifier: "account-group",
				Title:          team.TeamName,
			},
		},
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/teamview/team_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/teamview/team_page.js"),
		},
	})
}
