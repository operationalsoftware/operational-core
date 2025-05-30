package andonissueview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/reqcontext"
	"fmt"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type HomePageProps struct {
	Ctx             reqcontext.ReqContext
	ShowArchived    bool
	AndonIssues     []model.AndonIssue
	AndonIssueCount int
	Sort            appsort.Sort
	Page            int
	PageSize        int
}

func HomePage(p *HomePageProps) g.Node {
	var columns = components.TableColumns{
		{
			TitleContents: g.Text("Issue Name"),
			SortKey:       "IssueName",
		},
	}

	if p.ShowArchived {
		columns = append(columns, components.TableColumn{
			TitleContents: g.Text("Is Archived?"),
			SortKey:       "IsArchived",
		})
	}

	var tableRows components.TableRows
	for _, t := range p.AndonIssues {

		cells := []components.TableCell{
			{
				Contents: h.A(
					g.Text(t.IssueName),
					g.Attr("href",
						fmt.Sprintf("/andon-issues/%d", t.AndonIssueID))),
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
				g.Text("AndonIssue"),
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
			andonIssuesBreadCrumb,
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonissueview/home_page.css"),
		},
	})
}

var andonIssuesBreadCrumb = layout.Breadcrumb{
	IconIdentifier: "account-group",
	Title:          "AndonIssues",
	URLPart:        "andon-issues",
}
