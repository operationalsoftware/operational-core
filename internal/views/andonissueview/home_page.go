package andonissueview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"fmt"
	"strings"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type HomePageProps struct {
	Ctx             reqcontext.ReqContext
	ShowArchived    bool
	AndonIssues     []model.AndonIssueNode
	AndonIssueCount int
	Sort            appsort.Sort
	Page            int
	PageSize        int
}

func HomePage(p *HomePageProps) g.Node {
	var columns = components.TableColumns{
		{
			TitleContents: g.Text("Issue Name"),
			SortKey:       "NamePath",
		},
		{
			TitleContents: g.Text("Assigned to Team"),
			SortKey:       "AssignedToTeamName",
		},
		{
			TitleContents: g.Text("Severity"),
			SortKey:       "Severity",
		},
		{
			TitleContents: g.Text("Created By"),
			SortKey:       "CreatedByUsername",
		},
		{
			TitleContents: g.Text("Created At"),
			SortKey:       "CreatedAt",
		},
		{
			TitleContents: g.Text("Last Updated"),
			SortKey:       "UpdatedAt",
		},
	}

	if p.ShowArchived {
		columns = append(columns, components.TableColumn{
			TitleContents: g.Text("Is Archived?"),
			SortKey:       "IsArchived",
		})
	}

	var tableRows components.TableRows
	for _, ai := range p.AndonIssues {

		namePathStr := strings.Join(ai.NamePath, " > ")

		cells := []components.TableCell{
			{
				Contents: h.Div(
					h.Class("flex"),

					h.A(
						h.Class("issue-title"),
						g.Text(namePathStr),
						g.Attr("href", fmt.Sprintf("/andon-issues/%d", ai.AndonIssueID)),
					),

					g.If(
						ai.IsGroup,
						h.Div(
							h.Class("badge primary"),

							g.Text("Group"),
						),
					),
				),
			},
			{
				Contents: g.Text(nilsafe.Str(ai.AssignedTeamName)),
			},
			{
				Contents: g.Text(nilsafe.Str((*string)(ai.Severity))),
			},
			{
				Contents: g.Text(string(ai.CreatedByUsername)),
			},
			{
				Contents: g.Text(ai.CreatedAt.Format("2006-01-02 15:04:05")),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(ai.UpdatedAt == nil, g.Text("\u2013")),
					g.If(ai.UpdatedAt != nil, g.Text(nilsafe.Time(ai.UpdatedAt).Format("2006-01-02 15:04:05"))),
				}),
			},
		}

		if p.ShowArchived {
			cells = append(cells, components.TableCell{
				Contents: g.Group([]g.Node{
					g.If(ai.IsArchived,
						g.Text("Archived"),
					),
				}),
			})
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
				Link:       "/andon-issues/add",
				Classes: c.Classes{
					"add-andon-issue-btn": true,
				},
			},
				components.Icon(&components.IconProps{
					Identifier: "plus",
				}),
				g.Text("Andon Issue"),
			),

			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Link:       "/andon-issues/add-group",
				Classes: c.Classes{
					"add-andon-issue-btn": true,
				},
			},
				components.Icon(&components.IconProps{
					Identifier: "plus",
				}),
				g.Text("Andon Issue Group"),
			),
		),

		// form container for table interaction
		h.FormEl(
			h.ID("andon-issues-table-form"),
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
					TotalRecords:        p.AndonIssueCount,
					PageSize:            p.PageSize,
					CurrentPage:         p.Page,
					CurrentPageQueryKey: "Page",
					PageSizeQueryKey:    "PageSize",
				},
			},
				h.ID("andonIssues-table"),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Andon Issues",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "alert-octagon-outline",
				Title:          "Andons",
				URLPart:        "andons",
			},
			{
				Title:   "Andon Issues",
				URLPart: "andon-issues",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonissueview/home_page.css"),
		},
	})
}
