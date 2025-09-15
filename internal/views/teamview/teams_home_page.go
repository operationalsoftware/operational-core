package teamview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/reqcontext"
	"fmt"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type TeamsHomePageProps struct {
	Ctx          reqcontext.ReqContext
	ShowArchived bool
	Teams        []model.Team
	TeamCount    int
	Sort         appsort.Sort
	Page         int
	PageSize     int
}

func TeamsHomePage(p *TeamsHomePageProps) g.Node {
	var columns = components.TableColumns{
		{
			TitleContents: g.Text("Team Name"),
			SortKey:       "TeamName",
		},
	}

	if p.ShowArchived {
		columns = append(columns, components.TableColumn{
			TitleContents: g.Text("Is Archived?"),
			SortKey:       "IsArchived",
		})
	}

	var tableRows components.TableRows
	for _, t := range p.Teams {

		teamURL := fmt.Sprintf("/teams/%d", t.TeamID)
		cells := []components.TableCell{
			{
				Contents: h.A(
					g.Text(t.TeamName),
					g.Attr("href", teamURL),
				),
			},
		}

		if p.ShowArchived {
			cells = append(cells, components.TableCell{
				Contents: g.Group([]g.Node{
					g.If(t.IsArchived,
						g.Text("Archived"),
					),
				}),
			})
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
			HREF:  teamURL,
		})
	}

	content := g.Group([]g.Node{
		h.Div(
			h.Class("button-container"),
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Link:       "/teams/add",
				Classes: c.Classes{
					"add-team-btn": true,
				},
			},
				components.Icon(&components.IconProps{
					Identifier: "plus",
				}),
				g.Text("Team"),
			),
		),

		// form container for table interaction
		h.Form(
			h.ID("teams-table-form"),
			g.Attr("method", "GET"),

			components.Checkbox(
				&components.CheckboxProps{
					Name:    "ShowArchived",
					Label:   "Show Archived",
					Checked: p.ShowArchived,
					Value:   "true",
				},
				g.Attr("onchange", "submitTableForm(this.form)"),
			),

			components.Table(&components.TableProps{
				Columns: columns,
				Sort:    p.Sort,
				Rows:    tableRows,
				Pagination: &components.TablePaginationProps{
					TotalRecords:        p.TeamCount,
					PageSize:            p.PageSize,
					CurrentPage:         p.Page,
					CurrentPageQueryKey: "Page",
					PageSizeQueryKey:    "PageSize",
				},
			},
				h.ID("teams-table"),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Teams",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			teamsBreadCrumb,
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/teamview/teams_home_page.css"),
		},
	})
}

var teamsBreadCrumb = layout.Breadcrumb{
	IconIdentifier: "account-group",
	Title:          "Teams",
	URLPart:        "teams",
}
