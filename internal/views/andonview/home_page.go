package andonview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"strconv"
	"strings"
	"time"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type HomePageProps struct {
	Ctx                     reqcontext.ReqContext
	ShowArchived            bool
	OutstandingAndons       []model.AndonEvent
	AcknowledgedAndons      []model.AndonEvent
	NewAndonsCount          int
	AcknowledgedAndonsCount int
	Teams                   []model.Team
	SelectedTeams           []string
	IsSelfResolvable        bool
	Sort                    appsort.Sort
	Page                    int
	PageSize                int
}

func HomePage(p *HomePageProps) g.Node {

	var raisedColumns = components.TableColumns{
		{
			TitleContents: g.Text("Location"),
		},
		{
			TitleContents: g.Text("Issue Description"),
		},
		{
			TitleContents: g.Text("Issue"),
			SortKey:       "IssueName",
		},
		{
			TitleContents: g.Text("Assigned Team"),
			SortKey:       "AssignedTeam",
		},
		{
			TitleContents: g.Text("Severity"),
			SortKey:       "Severity",
		},
		{
			TitleContents: g.Text("Status"),
			SortKey:       "Status",
		},
		{
			TitleContents: g.Text("Raised By"),
			SortKey:       "RaisedBy",
		},
		{
			TitleContents: g.Text("Raised At"),
			SortKey:       "RaisedAt",
		},
		{
			TitleContents: g.Text("Actions"),
		},
	}

	var acknowledgedColumns = components.TableColumns{
		{
			TitleContents: g.Text("Location"),
		},
		{
			TitleContents: g.Text("Issue Description"),
		},
		{
			TitleContents: g.Text("Issue"),
			SortKey:       "IssueByName",
		},
		{
			TitleContents: g.Text("Assigned Team"),
			SortKey:       "AssignedTeam",
		},
		{
			TitleContents: g.Text("Severity"),
			SortKey:       "Severity",
		},
		{
			TitleContents: g.Text("Status"),
			SortKey:       "Status",
		},
		{
			TitleContents: g.Text("Acknowledged By"),
		},
		{
			TitleContents: g.Text("Acknowledged At"),
		},
		{
			TitleContents: g.Text("Actions"),
		},
	}

	var raisedTableRows components.TableRows
	for _, ai := range p.OutstandingAndons {

		namePathStr := strings.Join(ai.NamePath, " > ")

		twoMinutesPassed := time.Since(ai.RaisedAt) > 2*time.Minute
		fiveMinutesPassed := time.Since(ai.RaisedAt) > 5*time.Minute

		// isSelfResolvable := false
		// if ai.Severity == "Self-resolvable" && ai.IsTeamMate {
		// 	isSelfResolvable = true
		// }

		cells := []components.TableCell{
			{
				Contents: g.Text(ai.Location),
			},
			{
				Contents: g.Text(ai.IssueDescription),
			},
			{
				Contents: g.Text(namePathStr),
			},
			{
				Contents: g.Text(ai.AssignedTeam),
			},
			{
				Contents: g.Text(ai.Severity),
			},
			{
				Contents: g.Text(ai.Status),
			},
			{
				Contents: g.Text(ai.RaisedByUsername),
			},
			{
				Contents: g.Text(ai.RaisedAt.Format("2006-01-02 15:04:05")),
			},
			{
				Contents: g.Group([]g.Node{

					h.Div(
						h.Class("andon-actions"),

						g.If(
							ai.CanUserAcknowledge,

							components.Button(&components.ButtonProps{
								Size:       "small",
								ButtonType: "button",
							},
								g.Attr("onclick", "updateAndon(event)"),
								g.Attr("data-id", strconv.Itoa(ai.AndonEventID)),
								g.Attr("data-action", "acknowledge"),
								g.Attr("title", "Acknowledge"),

								components.Icon(&components.IconProps{
									Identifier: "gesture-tap-hold",
								}),
							),
						),

						g.If(
							ai.Severity == "Self-resolvable" && ai.CanUserResolve,
							components.Button(&components.ButtonProps{
								Size:       "small",
								ButtonType: "button",
							},
								g.Attr("onclick", "updateAndon(event)"),
								g.Attr("data-id", strconv.Itoa(ai.AndonEventID)),
								g.Attr("data-action", "resolve"),
								g.Attr("title", "Resolve"),

								components.Icon(&components.IconProps{
									Identifier: "check",
								}),
							),
						),

						components.Button(&components.ButtonProps{
							Size: "small",
						},
							g.Attr("onclick", "updateAndon(event)"),
							g.Attr("data-id", strconv.Itoa(ai.AndonEventID)),
							g.Attr("data-action", "cancel"),
							g.Attr("title", "Cancel"),

							components.Icon(&components.IconProps{
								Identifier: "cancel",
							}),
						),
					),
				}),
			},
		}

		for i := 0; i < len(cells)-1; i++ {
			cells[i].Classes = c.Classes{
				"amber":        twoMinutesPassed,
				"flashing-red": fiveMinutesPassed,
			}
		}

		raisedTableRows = append(raisedTableRows, components.TableRow{
			Cells: cells,
		})
	}

	var acknowledgedTableRows components.TableRows
	for _, ai := range p.AcknowledgedAndons {
		namePathStr := strings.Join(ai.NamePath, " > ")

		if ai.Severity == "Info" {
			continue
		}

		cells := []components.TableCell{
			{
				Contents: g.Text(ai.Location),
			},
			{
				Contents: g.Text(ai.IssueDescription),
			},
			{
				Contents: g.Text(namePathStr),
			},
			{
				Contents: g.Text(ai.AssignedTeam),
			},
			{
				Contents: g.Text(ai.Severity),
			},
			{
				Contents: g.Text(ai.Status),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(ai.AcknowledgedByUsername != nil, g.Text(nilsafe.Str(ai.AcknowledgedByUsername))),
				}),
			},
			{
				Contents: g.Text(ai.AcknowledgedAt.Format("2006-01-02 15:04:05")),
			},
			{
				Contents: g.Group([]g.Node{

					h.Div(
						h.Class("andon-actions"),

						components.Button(&components.ButtonProps{
							Size:       "small",
							ButtonType: "button",
						},
							g.Attr("onclick", "updateAndon(event)"),
							g.Attr("data-id", strconv.Itoa(ai.AndonEventID)),
							g.Attr("data-action", "resolve"),
							g.Attr("title", "Resolve"),

							components.Icon(&components.IconProps{
								Identifier: "check",
							}),
						),
						components.Button(&components.ButtonProps{
							Size: "small",
						},
							g.Attr("onclick", "updateAndon(event)"),
							g.Attr("data-id", strconv.Itoa(ai.AndonEventID)),
							g.Attr("data-action", "cancel"),
							g.Attr("title", "Cancel"),

							components.Icon(&components.IconProps{
								Identifier: "cancel",
							}),
						),
					),
				}),
			},
		}

		acknowledgedTableRows = append(acknowledgedTableRows, components.TableRow{
			Cells: cells,
		})
	}

	content := g.Group([]g.Node{
		h.Nav(
			h.Class("andon-nav"),

			h.A(
				h.Href("/andons/add"),

				components.Icon(&components.IconProps{
					Identifier: "plus",
					Classes: c.Classes{
						"icon": true,
					},
				},
				),
				g.Text("New Andon")),
			h.A(h.Href("/andons/all"), g.Text("All Andons")),
			h.A(
				h.Href("/andon-issues"),

				components.Icon(&components.IconProps{
					Identifier: "wrench-outline",
					Classes: c.Classes{
						"icon": true,
					},
				},
				),
				g.Text("Andon Issues")),
		),

		h.Div(
			h.Class("team-select"),

			h.FormEl(
				g.Attr("method", "GET"),
				g.Attr("id", "team-form"),

				h.Div(
					g.Attr("data-selected", strings.Join(p.SelectedTeams, ",")),
					g.Attr("id", "search-select-wrapper"),

					components.SearchSelect(&components.SearchSelectProps{
						Name:          "AndonTeams",
						Placeholder:   "Select a team",
						Mode:          "multi",
						Options:       MapTeamsToOptions(p.Teams),
						Selected:      strings.Join(p.SelectedTeams, ","),
						ShowOnlyLabel: true,
						// OnChange:
					},
						g.Attr("onchange", "handleTeamSelectChange(event)"),
					),
				),
			),
		),

		// h.Div(
		// 	h.Class("button-container"),

		// 	components.Button(&components.ButtonProps{
		// 		ButtonType: "primary",
		// 		Link:       "/andons/add?Location=Assembly Line 25&Source=Works Order Receipt",
		// 		Classes: c.Classes{
		// 			"add-andon-issue-btn": true,
		// 		},
		// 	},
		// 		components.Icon(&components.IconProps{
		// 			Identifier: "plus",
		// 		}),
		// 		g.Text("Andon Event"),
		// 	),
		// ),

		h.FormEl(
			h.ID("andon-wip-table-form"),
			g.Attr("method", "GET"),

			h.H3(
				h.Class("table-title"),
				g.Text("New"),
			),
			h.Hr(),
			components.Table(&components.TableProps{
				Columns: raisedColumns,
				Sort:    p.Sort,
				Rows:    raisedTableRows,
				Pagination: &components.TablePaginationProps{
					TotalRecords:        p.NewAndonsCount,
					PageSize:            p.PageSize,
					CurrentPage:         p.Page,
					CurrentPageQueryKey: "Page",
					PageSizeQueryKey:    "PageSize",
				},
			},
				h.ID("new-andon-table"),
			),

			h.H3(
				h.Class("table-title wip-heading"),
				g.Text("WIP"),
			),
			h.Hr(),
			components.Table(&components.TableProps{
				Columns: acknowledgedColumns,
				Sort:    p.Sort,
				Rows:    acknowledgedTableRows,
				Pagination: &components.TablePaginationProps{
					TotalRecords:        p.AcknowledgedAndonsCount,
					PageSize:            p.PageSize,
					CurrentPage:         p.Page,
					CurrentPageQueryKey: "Page",
					PageSizeQueryKey:    "PageSize",
				},
			},
				h.ID("andon-wip-table"),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Andons",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "alert-octagon-outline",
				Title:          "Andons",
				URLPart:        "andons",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonview/home_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/andonview/home_page.js"),
		},
	})
}

var andonIssuesBreadCrumb = layout.Breadcrumb{
	IconIdentifier: "alert-octagon-outline",
	Title:          "Andons",
	URLPart:        "andons",
}

func MapTeamsToOptions(teams []model.Team) []components.SearchSelectOption {
	out := make([]components.SearchSelectOption, len(teams))
	for i, v := range teams {
		out[i] = components.SearchSelectOption{
			Label: v.TeamName,
			Value: v.TeamName,
			// Value: strconv.Itoa(v.TeamID),
		}
	}
	return out
}
