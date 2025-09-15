package andonview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AllAndonsPageProps struct {
	Ctx              reqcontext.ReqContext
	ShowArchived     bool
	Andons           []model.Andon
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
	}

	var tableRows components.TableRows
	for _, a := range p.Andons {
		namePathStr := strings.Join(a.NamePath, " > ")

		cells := []components.TableCell{
			{
				Contents: g.Text(a.Location),
			},
			{
				Contents: g.Text(a.Description),
			},
			{
				Contents: g.Text(namePathStr),
			},
			{
				Contents: g.Text(a.AssignedTeamName),
			},
			{
				Contents: severityBadge(a.Severity, "small"),
			},
			{
				Contents: statusBadge(a.Status, "small"),
			},
			{
				Contents: g.Text(a.RaisedByUsername),
			},
			{
				Contents: g.Text(a.RaisedAt.Format("2006-01-02 15:04:05")),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(a.AcknowledgedByUsername != nil, g.Text(nilsafe.Str(a.AcknowledgedByUsername))),
				}),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(a.AcknowledgedAt != nil, g.Text(nilsafe.Time(a.AcknowledgedAt).Format("2006-01-02 15:04:05"))),
				}),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(a.ResolvedByUsername != nil, g.Text(nilsafe.Str(a.ResolvedByUsername))),
				}),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(a.ResolvedAt != nil, g.Text(nilsafe.Time(a.ResolvedAt).Format("2006-01-02 15:04:05"))),
				}),
			},
			{
				Contents: g.Group([]g.Node{
					g.If(a.LastUpdated != nil, g.Text(nilsafe.Time(a.LastUpdated).Format("2006-01-02 15:04:05"))),
				}),
			},
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
			HREF:  fmt.Sprintf("/andons/%d", a.AndonID),
		})
	}

	content := g.Group([]g.Node{

		h.H3(g.Text("All Andons")),

		h.Form(
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
						Options:     components.MapStringsToOptions(availableFilters.IssueIn, p.Filters.Issues),
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
						Options: components.MapStringsToOptions(availableFilters.
							SeverityIn, p.Filters.Severities),
						Selected: strings.Join(p.Filters.Severities, ","),
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
						Options:     components.MapStringsToOptions(availableFilters.TeamIn, p.Filters.Teams),
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
						Options:     components.MapStringsToOptions(availableFilters.LocationIn, p.Filters.Locations),
						Selected:    strings.Join(p.Filters.Locations, ","),
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
						Options:     components.MapStringsToOptions(availableFilters.RaisedByUsernameIn, p.Filters.RaisedByUsername),
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
						Options:     components.MapStringsToOptions(availableFilters.AcknowledgedByUsernameIn, p.Filters.AcknowledgedByUsername),
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
						Options:     components.MapStringsToOptions(availableFilters.ResolvedByUsernameIn, p.Filters.ResolvedByUsername),
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
			components.InlineStyle("/internal/views/andonview/components.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/andonview/all_andons_page.js"),
		},
	})
}
