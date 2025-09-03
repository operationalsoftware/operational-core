package andonview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"fmt"
	"strconv"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
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

	var outstandingColumns = components.TableColumns{
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

	var outstandingTableRows []components.TableRow
	for _, a := range p.OutstandingAndons {

		namePathStr := strings.Join(a.NamePath, " > ")

		twoMinutesPassed := time.Since(a.RaisedAt) > 2*time.Minute
		fiveMinutesPassed := time.Since(a.RaisedAt) > 5*time.Minute

		cells := []components.TableCell{
			{
				Contents: g.Text(a.Location),
			},
			{
				Contents: g.Text(a.IssueDescription),
			},
			{
				Contents: g.Text(namePathStr),
			},
			{
				Contents: g.Text(a.AssignedTeamName),
			},
			{
				Contents: g.Text(string(a.Severity)),
			},
			{
				Contents: g.Text(a.Status),
			},
			{
				Contents: g.Text(a.RaisedByUsername),
			},
			{
				Contents: g.Text(a.RaisedAt.Format("2006-01-02 15:04:05")),
			},
			{
				Contents: g.Group([]g.Node{

					h.Div(
						h.Class("andon-actions"),

						g.If(
							a.CanUserAcknowledge,

							components.Button(&components.ButtonProps{
								Size:       "small",
								ButtonType: "button",
							},
								g.Attr("onclick", "updateAndon(event)"),
								g.Attr("data-id", strconv.Itoa(a.AndonEventID)),
								g.Attr("data-action", "acknowledge"),
								g.Attr("title", "Acknowledge"),

								components.Icon(&components.IconProps{
									Identifier: "gesture-tap-hold",
								}),
							),
						),

						g.If(
							a.Severity == "Self-resolvable" && a.CanUserResolve,
							components.Button(&components.ButtonProps{
								Size:       "small",
								ButtonType: "button",
							},
								g.Attr("onclick", "updateAndon(event)"),
								g.Attr("data-id", strconv.Itoa(a.AndonEventID)),
								g.Attr("data-action", "resolve"),
								g.Attr("title", "Resolve"),

								components.Icon(&components.IconProps{
									Identifier: "check",
								}),
							),
						),

						g.If(a.CanUserCancel,
							components.Button(&components.ButtonProps{
								Size: "small",
							},
								g.Attr("onclick", "updateAndon(event)"),
								g.Attr("data-id", strconv.Itoa(a.AndonEventID)),
								g.Attr("data-action", "cancel"),
								g.Attr("title", "Cancel"),

								components.Icon(&components.IconProps{
									Identifier: "cancel",
								}),
							),
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

		outstandingTableRows = append(outstandingTableRows, components.TableRow{
			Cells: cells,
			HREF:  fmt.Sprintf("/andons/%d", a.AndonEventID),
		})
	}

	var acknowledgedTableRows components.TableRows
	for _, a := range p.AcknowledgedAndons {
		namePathStr := strings.Join(a.NamePath, " > ")

		if a.Severity == "Info" {
			continue
		}

		cells := []components.TableCell{
			{
				Contents: g.Text(a.Location),
			},
			{
				Contents: g.Text(a.IssueDescription),
			},
			{
				Contents: g.Text(namePathStr),
			},
			{
				Contents: g.Text(a.AssignedTeamName),
			},
			{
				Contents: g.Text(string(a.Severity)),
			},
			{
				Contents: g.Text(a.Status),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(a.AcknowledgedByUsername != nil, g.Text(nilsafe.Str(a.AcknowledgedByUsername))),
				}),
			},
			{
				Contents: g.Text(a.AcknowledgedAt.Format("2006-01-02 15:04:05")),
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
							g.Attr("data-id", strconv.Itoa(a.AndonEventID)),
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
							g.Attr("data-id", strconv.Itoa(a.AndonEventID)),
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
			HREF:  fmt.Sprintf("/andons/%d", a.AndonEventID),
		})
	}

	content := g.Group([]g.Node{
		h.Nav(
			h.Class("andon-nav"),

			components.Button(&components.ButtonProps{
				Size: "small",
				Classes: c.Classes{
					"primary": true,
				},
				Link: "/andons/add",
			},
				components.Icon(&components.IconProps{
					Identifier: "plus",
				}),
				g.Text("New Andon"),
			),
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

			h.Form(
				g.Attr("method", "GET"),
				g.Attr("id", "team-form"),

				h.Div(
					g.Attr("data-selected", strings.Join(p.SelectedTeams, ",")),
					g.Attr("id", "search-select-wrapper"),

					components.SearchSelect(&components.SearchSelectProps{
						Name:        "AndonTeams",
						Placeholder: "Select a team",
						Mode:        "multi",
						Options:     MapTeamsToOptions(p.Teams, p.SelectedTeams),
					},
						g.Attr("onchange", "handleTeamSelectChange(event)"),
					),
				),
			),
		),

		h.Form(
			h.ID("andon-wip-table-form"),
			g.Attr("method", "GET"),

			h.H3(
				h.Class("table-title"),
				g.Text("New"),
			),
			h.Hr(),
			components.Table(&components.TableProps{
				Columns: outstandingColumns,
				Sort:    p.Sort,
				Rows:    outstandingTableRows,
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

		h.Div(
			h.Class("status-legend"),

			h.Div(
				h.Span(
					h.Class("status-dot two-minutes-passed"),
				),
				g.Text("Outstanding (> 2 minutes)"),
			),
			h.Div(
				h.Span(
					h.Class("status-dot five-minutes-passed"),
				),
				g.Text("Outstanding (> 5 minutes)"),
			),
			h.Div(
				h.Span(
					h.Class("status-dot status-acknowledged"),
				),
				g.Text("Acknowledged"),
			),
			h.Div(
				h.Span(
					h.Class("status-dot status-cancelled"),
				),
				g.Text("Cancelled"),
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

func MapTeamsToOptions(teams []model.Team, selectedValues []string) []components.SearchSelectOption {
	out := make([]components.SearchSelectOption, len(teams))
	for i, v := range teams {
		isSelected := false
		for _, selectedTeam := range selectedValues {
			if v.TeamName == selectedTeam {
				isSelected = true
				break
			}
		}

		out[i] = components.SearchSelectOption{
			Text:     v.TeamName,
			Value:    v.TeamName,
			Selected: isSelected,
		}
	}
	return out
}
