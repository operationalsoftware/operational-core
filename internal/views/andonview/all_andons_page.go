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

type AllAndonsPageProps struct {
	Ctx              reqcontext.ReqContext
	ShowArchived     bool
	Andons           []model.AndonEvent
	AndonsCount      int
	AvailableFilters model.AndonAvailableFilters
	Filters          model.AndonFilters
	Sort             appsort.Sort
	Page             int
	PageSize         int
}

type AvailableFilters struct {
	IssueIn                  []components.SearchSelectOption
	SeverityIn               []components.SearchSelectOption
	TeamIn                   []components.SearchSelectOption
	LocationIn               []components.SearchSelectOption
	StatusIn                 []components.SearchSelectOption
	RaisedByUsernameIn       []components.SearchSelectOption
	AcknowledgedByUsernameIn []components.SearchSelectOption
	ResolvedByUsernameIn     []components.SearchSelectOption
}

func AllAndonsPage(p *AllAndonsPageProps) g.Node {

	availableFilters := p.AvailableFilters

	var columns = components.TableColumns{
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
			TitleContents: g.Text("Acknowledged By"),
			SortKey:       "AcknowledgedBy",
		},
		{
			TitleContents: g.Text("Acknowledged At"),
			SortKey:       "AcknowledgedAt",
		},
		{
			TitleContents: g.Text("Resolved By"),
			SortKey:       "ResolvedBy",
		},
		{
			TitleContents: g.Text("Resolved At"),
			SortKey:       "ResolvedAt",
		},
		{
			TitleContents: g.Text("Updated At"),
			SortKey:       "LastUpdated",
		},
		{
			TitleContents: g.Text("Actions"),
		},
	}

	var tableRows components.TableRows
	for _, ai := range p.Andons {
		namePathStr := strings.Join(ai.NamePath, " > ")

		isAcknowledged := ai.Status == "Acknowledged"
		isResolved := ai.Status == "Resolved"
		isCancelled := ai.Status == "Cancelled"
		twoMinutesPassed := time.Since(ai.RaisedAt) > 2*time.Minute && !isResolved
		fiveMinutesPassed := time.Since(ai.RaisedAt) > 5*time.Minute && !isResolved

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
				Contents: g.Text(ai.AssignedTeamName),
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
					g.If(ai.AcknowledgedByUsername != nil, g.Text(nilsafe.Str(ai.AcknowledgedByUsername))),
				}),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(ai.AcknowledgedAt != nil, g.Text(nilsafe.Time(ai.AcknowledgedAt).Format("2006-01-02 15:04:05"))),
				}),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(ai.ResolvedByUsername != nil, g.Text(nilsafe.Str(ai.ResolvedByUsername))),
				}),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(ai.ResolvedAt != nil, g.Text(nilsafe.Time(ai.ResolvedAt).Format("2006-01-02 15:04:05"))),
				}),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(ai.LastUpdated != nil, g.Text(nilsafe.Time(ai.LastUpdated).Format("2006-01-02 15:04:05"))),
				}),
			},
			{
				Contents: g.Group([]g.Node{

					h.Div(
						h.Class("andon-actions"),

						g.If(ai.Status == "Outstanding" && ai.CanUserAcknowledge,
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
						g.If(ai.Status == "Acknowledged" && ai.CanUserResolve,
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
						g.If(ai.Status == "Outstanding" && ai.Severity == "Self-resolvable" && ai.CanUserResolve,
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
						g.If(ai.Status == "Cancelled",
							components.Button(&components.ButtonProps{
								Size:       "small",
								ButtonType: "button",
							},
								g.Attr("onclick", "updateAndon(event)"),
								g.Attr("data-id", strconv.Itoa(ai.AndonEventID)),
								g.Attr("data-action", "reopen"),
								g.Attr("title", "Reopen"),

								components.Icon(&components.IconProps{
									Identifier: "restore",
								},
								),
							),
						),

						g.If(ai.CanUserCancel,
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
					),
				}),
			},
		}

		for i := 0; i < len(cells)-1; i++ {
			cells[i].Classes = c.Classes{
				"amber":         twoMinutesPassed,
				"flashing-red":  fiveMinutesPassed,
				"light-green":   isAcknowledged,
				"flashing-grey": isCancelled,
			}
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
		})
	}

	issueOptions := make([]components.SearchSelectOption, len(availableFilters.IssueIn))
	for i, v := range availableFilters.IssueIn {
		isSelected := false
		for _, selectedTeam := range p.Filters.Issues {
			if v == selectedTeam {
				isSelected = true
				break
			}
		}

		issueOptions[i] = components.SearchSelectOption{
			Text:     v,
			Value:    v,
			Selected: isSelected,
		}
	}

	severityOptions := make([]components.SearchSelectOption, len(availableFilters.SeverityIn))
	for i, v := range availableFilters.SeverityIn {
		isSelected := false
		for _, selectedTeam := range p.Filters.Severities {
			if v == selectedTeam {
				isSelected = true
				break
			}
		}

		severityOptions[i] = components.SearchSelectOption{
			Text:     v,
			Value:    v,
			Selected: isSelected,
		}
	}

	teamOptions := make([]components.SearchSelectOption, len(availableFilters.TeamIn))
	for i, v := range availableFilters.TeamIn {
		isSelected := false
		for _, selectedTeam := range p.Filters.Teams {
			if v == selectedTeam {
				isSelected = true
				break
			}
		}

		teamOptions[i] = components.SearchSelectOption{
			Text:     v,
			Value:    v,
			Selected: isSelected,
		}
	}

	locationOptions := make([]components.SearchSelectOption, len(availableFilters.LocationIn))
	for i, v := range availableFilters.LocationIn {
		isSelected := false
		for _, selectedTeam := range p.Filters.Locations {
			if v == selectedTeam {
				isSelected = true
				break
			}
		}

		locationOptions[i] = components.SearchSelectOption{
			Text:     v,
			Value:    v,
			Selected: isSelected,
		}
	}

	statusOptions := make([]components.SearchSelectOption, len(availableFilters.StatusIn))
	for i, v := range availableFilters.StatusIn {
		isSelected := false
		for _, selectedTeam := range p.Filters.Statuses {
			if v == selectedTeam {
				isSelected = true
				break
			}
		}

		statusOptions[i] = components.SearchSelectOption{
			Text:     v,
			Value:    v,
			Selected: isSelected,
		}
	}

	raisedByUsernameOptions := make([]components.SearchSelectOption, len(availableFilters.RaisedByUsernameIn))
	for i, v := range availableFilters.RaisedByUsernameIn {
		isSelected := false
		for _, selectedTeam := range p.Filters.RaisedByUsername {
			if v == selectedTeam {
				isSelected = true
				break
			}
		}

		raisedByUsernameOptions[i] = components.SearchSelectOption{
			Text:     v,
			Value:    v,
			Selected: isSelected,
		}
	}

	acknowledgedByUsernameOptions := make([]components.SearchSelectOption, len(availableFilters.AcknowledgedByUsernameIn))
	for i, v := range availableFilters.AcknowledgedByUsernameIn {
		isSelected := false
		for _, selectedTeam := range p.Filters.AcknowledgedByUsername {
			if v == selectedTeam {
				isSelected = true
				break
			}
		}

		acknowledgedByUsernameOptions[i] = components.SearchSelectOption{
			Text:     v,
			Value:    v,
			Selected: isSelected,
		}
	}

	resolvedByUsernameOptions := make([]components.SearchSelectOption, len(availableFilters.ResolvedByUsernameIn))
	for i, v := range availableFilters.ResolvedByUsernameIn {
		isSelected := false
		for _, selectedTeam := range p.Filters.ResolvedByUsername {
			if v == selectedTeam {
				isSelected = true
				break
			}
		}

		resolvedByUsernameOptions[i] = components.SearchSelectOption{
			Text:     v,
			Value:    v,
			Selected: isSelected,
		}
	}

	content := g.Group([]g.Node{

		h.H3(g.Text("All Andons")),

		h.FormEl(
			h.ID("all-andon-table-form"),
			g.Attr("method", "GET"),

			h.Div(
				h.Class("andon-filters"),

				h.Div(
					h.Class("search-item date-section"),

					h.Div(

						h.Label(
							g.Text("Start date"),
						),
						h.Input(
							h.Name("StartDate"),
							h.Type("date"),

							func() g.Node {
								if p.Filters.StartDate != nil {
									return h.Value(p.Filters.StartDate.Format("2006-01-02"))
								}
								return nil
							}(),
						),
					),

					h.Div(

						h.Label(
							g.Text("End date"),
						),
						h.Input(
							h.Name("EndDate"),
							h.Type("date"),
							func() g.Node {
								if p.Filters.EndDate != nil {
									return h.Value(p.Filters.EndDate.Format("2006-01-02"))
								}
								return nil
							}(),
						),
					),
				),
				h.Div(
					h.Class("search-item"),

					h.Label(
						g.Text("Issues"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "IssueIn",
						Placeholder: "-",
						Mode:        "multi",
						Options:     issueOptions,
						Selected:    strings.Join(p.Filters.Issues, ","),
					}),
				),
				h.Div(
					h.Class("search-item"),

					h.Label(
						g.Text("Severity"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "SeverityIn",
						Placeholder: "-",
						Mode:        "multi",
						Options:     severityOptions,
						Selected:    strings.Join(p.Filters.Severities, ","),
					}),
				),
				h.Div(
					h.Class("search-item"),

					h.Label(
						g.Text("Assigned Teams"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "TeamIn",
						Placeholder: "-",
						Mode:        "multi",
						Options:     teamOptions,
						Selected:    strings.Join(p.Filters.Teams, ","),
					}),
				),
				h.Div(
					h.Class("search-item"),

					h.Label(
						g.Text("Locations"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "LocationIn",
						Placeholder: "-",
						Mode:        "multi",
						Options:     locationOptions,
						Selected:    strings.Join(p.Filters.Locations, ","),
					}),
				),
				h.Div(
					h.Class("search-item"),

					h.Label(
						g.Text("Statuses"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "StatusIn",
						Placeholder: "-",
						Mode:        "multi",
						Options:     statusOptions,
						Selected:    strings.Join(p.Filters.Statuses, ","),
					}),
				),
				h.Div(
					h.Class("search-item"),

					h.Label(
						g.Text("Raised By"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "RaisedByUsernameIn",
						Placeholder: "-",
						Mode:        "multi",
						Options:     raisedByUsernameOptions,
						Selected:    strings.Join(p.Filters.RaisedByUsername, ","),
					}),
				),
				h.Div(
					h.Class("search-item"),

					h.Label(
						g.Text("Acknowledged By"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "AcknowledgedByUsernameIn",
						Placeholder: "-",
						Mode:        "multi",
						Options:     acknowledgedByUsernameOptions,
						Selected:    strings.Join(p.Filters.AcknowledgedByUsername, ","),
					}),
				),
				h.Div(
					h.Class("search-item"),

					h.Label(
						g.Text("Resolved By"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "ResolvedByUsernameIn",
						Placeholder: "-",
						Mode:        "multi",
						Options:     resolvedByUsernameOptions,
						Selected:    strings.Join(p.Filters.ResolvedByUsername, ","),
					}),
				),
			),

			components.Button(&components.ButtonProps{
				ButtonType: components.ButtonPrimary,
				Size:       components.ButtonLg,
			},
				h.Type("submit"),
				h.ID("go-button"),
				g.Text("GO"),
			),

			components.Table(&components.TableProps{
				Columns: columns,
				Sort:    p.Sort,
				Rows:    tableRows,
				Pagination: &components.TablePaginationProps{
					TotalRecords:        p.AndonsCount,
					PageSize:            p.PageSize,
					CurrentPage:         p.Page,
					CurrentPageQueryKey: "Page",
					PageSizeQueryKey:    "PageSize",
				},
			},
				h.ID("andon-table"),
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
		Title:   "All Andons",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "alert-octagon-outline",
				Title:          "Andons",
				URLPart:        "andons",
			},
			{
				Title: "All",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonview/all_andons_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/andonview/all_andons_page.js"),
		},
	})
}
