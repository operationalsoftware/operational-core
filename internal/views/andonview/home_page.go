package andonview

import (
	"slices"
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
	OutstandingAndons       []model.Andon
	OutstandingAndonsCount  int
	AcknowledgedAndons      []model.Andon
	AcknowledgedAndonsCount int
	Teams                   []model.Team
	SelectedTeams           []string
	IsSelfResolvable        bool
	OutstandingSort         appsort.Sort
	WIPSort                 appsort.Sort
}

func HomePage(p *HomePageProps) g.Node {

	var commonColumns = components.TableColumns{
		{
			TitleContents: g.Text("Issue"),
			SortKey:       "NamePath",
		},
		{
			TitleContents: g.Text("Location"),
		},
		{
			TitleContents: g.Text("Description"),
		},
		{
			TitleContents: g.Text("Assigned Team"),
			SortKey:       "AssignedTeamName",
		},
		{
			TitleContents: g.Text("Severity"),
			SortKey:       "Severity",
		},
	}

	var outstandingColumns = append(
		commonColumns,
		components.TableColumns{
			{
				TitleContents: g.Text("Raised By"),
				SortKey:       "RaisedByUsername",
			},
			{
				TitleContents: g.Text("Raised At"),
				SortKey:       "RaisedAt",
			}, {
				TitleContents: g.Text("Actions"),
			},
		}...,
	)

	var acknowledgedColumns = append(
		commonColumns,
		components.TableColumns{
			{
				TitleContents: g.Text("Acknowledged By"),
				SortKey:       "AcknowledgedByUsername",
			},
			{
				TitleContents: g.Text("Acknowledged At"),
				SortKey:       "AcknowledgedAt",
			}, {
				TitleContents: g.Text("Actions"),
			},
		}...,
	)

	var outstandingTableRows []components.TableRow
	for _, a := range p.OutstandingAndons {

		namePathStr := strings.Join(a.NamePath, " > ")

		twoMinutesPassed := time.Since(a.RaisedAt) > 2*time.Minute
		fiveMinutesPassed := time.Since(a.RaisedAt) > 5*time.Minute

		cells := []components.TableCell{
			{Contents: g.Text(namePathStr)},
			{Contents: g.Text(a.Location)},
			{Contents: g.Text(a.Description)},
			{Contents: g.Text(a.AssignedTeamName)},
			{Contents: severityBadge(a.Severity, "small")},
			{Contents: g.Text(a.RaisedByUsername)},
			{Contents: g.Text(a.RaisedAt.Format("2006-01-02 15:04:05"))},
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
								g.Attr("data-id", strconv.Itoa(a.AndonID)),
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
								g.Attr("data-id", strconv.Itoa(a.AndonID)),
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
								g.Attr("data-id", strconv.Itoa(a.AndonID)),
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
				"two-minutes-passed":  !fiveMinutesPassed && twoMinutesPassed,
				"five-minutes-passed": fiveMinutesPassed,
			}
		}

		outstandingTableRows = append(outstandingTableRows, components.TableRow{
			Cells: cells,
			HREF:  fmt.Sprintf("/andons/%d", a.AndonID),
		})
	}

	var acknowledgedTableRows components.TableRows
	for _, a := range p.AcknowledgedAndons {
		namePathStr := strings.Join(a.NamePath, " > ")

		cells := []components.TableCell{
			{
				Contents: g.Text(namePathStr),
			},
			{
				Contents: g.Text(a.Location),
			},
			{
				Contents: g.Text(a.Description),
			},
			{
				Contents: g.Text(a.AssignedTeamName),
			},
			{
				Contents: severityBadge(a.Severity, "small"),
			},

			{
				Contents: g.Text(nilsafe.Str(a.AcknowledgedByUsername)),
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
							g.Attr("data-id", strconv.Itoa(a.AndonID)),
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
							g.Attr("data-id", strconv.Itoa(a.AndonID)),
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
			HREF:  fmt.Sprintf("/andons/%d", a.AndonID),
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

			g.If(
				p.Ctx.User.Permissions.Andon.Admin,
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
				Columns:      outstandingColumns,
				SortQueryKey: "OutstandingSort",
				Sort:         p.OutstandingSort,
				Rows:         outstandingTableRows,
			},
				h.ID("new-andon-table"),
			),

			h.H3(
				h.Class("table-title wip-heading"),
				g.Text("WIP"),
			),
			h.Hr(),
			components.Table(&components.TableProps{
				Columns:      acknowledgedColumns,
				SortQueryKey: "WIPSort",
				Sort:         p.WIPSort,
				Rows:         acknowledgedTableRows,
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
			components.InlineStyle("/internal/views/andonview/components.css"),
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
		isSelected := slices.Contains(selectedValues, v.TeamName)

		out[i] = components.SearchSelectOption{
			Text:     v.TeamName,
			Value:    v.TeamName,
			Selected: isSelected,
		}
	}
	return out
}
