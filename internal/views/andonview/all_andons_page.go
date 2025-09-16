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
	ActiveFilters    model.AndonFilters
	Sort             appsort.Sort
	Page             int
	PageSize         int
}

func AllAndonsPage(p *AllAndonsPageProps) g.Node {

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

			allAndonsFilters(&allAndonsFiltersProps{
				availableFilters: p.AvailableFilters,
				activeFilters:    p.ActiveFilters,
			}),

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

type allAndonsFiltersProps struct {
	availableFilters model.AndonAvailableFilters
	activeFilters    model.AndonFilters
}

func allAndonsFilters(p *allAndonsFiltersProps) g.Node {
	var startDateValue, endDateValue g.Node
	if p.activeFilters.StartDate != nil {
		startDateValue = h.Value(p.activeFilters.StartDate.Format("2006-01-02"))
	}
	if p.activeFilters.EndDate != nil {
		endDateValue = h.Value(p.activeFilters.EndDate.Format("2006-01-02"))
	}

	return h.Div(
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
					startDateValue,
				),
			),

			h.Div(

				h.Label(
					g.Text("End date"),
				),
				h.Input(
					h.Name("EndDate"),
					h.Type("date"),
					endDateValue,
				),
			),
		),

		h.Div(
			h.Class("search-item"),

			h.Label(
				g.Text("Location"),
			),
			components.SearchSelect(&components.SearchSelectProps{
				Name:        "LocationIn",
				Placeholder: "-",
				Mode:        "multi",
				Options:     components.MapStringsToOptions(p.availableFilters.LocationIn, p.activeFilters.LocationIn),
				Selected:    strings.Join(p.activeFilters.LocationIn, ","),
			}),
		),
		h.Div(
			h.Class("search-item"),

			h.Label(
				g.Text("Issue"),
			),
			components.SearchSelect(&components.SearchSelectProps{
				Name:        "IssueIn",
				Placeholder: "-",
				Mode:        "multi",
				Options:     components.MapStringsToOptions(p.availableFilters.IssueIn, p.activeFilters.IssueIn),
				Selected:    strings.Join(p.activeFilters.IssueIn, ","),
			}),
		),
		h.Div(
			h.Class("search-item"),

			h.Label(
				g.Text("Assigned Team"),
			),
			components.SearchSelect(&components.SearchSelectProps{
				Name:        "TeamIn",
				Placeholder: "-",
				Mode:        "multi",
				Options:     components.MapStringsToOptions(p.availableFilters.TeamIn, p.activeFilters.TeamIn),
				Selected:    strings.Join(p.activeFilters.TeamIn, ","),
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
				Options: components.MapStringsToOptions(
					p.availableFilters.SeverityIn,
					p.activeFilters.SeverityIn),
				Selected: strings.Join(p.activeFilters.SeverityIn, ","),
			}),
		),
		h.Div(
			h.Class("search-item"),

			h.Label(
				g.Text("Status"),
			),
			components.SearchSelect(&components.SearchSelectProps{
				Name:        "StatusIn",
				Placeholder: "-",
				Mode:        "multi",
				Options: components.MapStringsToOptions(
					p.availableFilters.StatusIn,
					p.activeFilters.StatusIn),
				Selected: strings.Join(p.activeFilters.StatusIn, ","),
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
				Options:     components.MapStringsToOptions(p.availableFilters.RaisedByUsernameIn, p.activeFilters.RaisedByUsernameIn),
				Selected:    strings.Join(p.activeFilters.RaisedByUsernameIn, ","),
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
				Options:     components.MapStringsToOptions(p.availableFilters.AcknowledgedByUsernameIn, p.activeFilters.AcknowledgedByUsernameIn),
				Selected:    strings.Join(p.activeFilters.AcknowledgedByUsernameIn, ","),
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
				Options:     components.MapStringsToOptions(p.availableFilters.ResolvedByUsernameIn, p.activeFilters.ResolvedByUsernameIn),
				Selected:    strings.Join(p.activeFilters.ResolvedByUsernameIn, ","),
			}),
		),
	)
}
