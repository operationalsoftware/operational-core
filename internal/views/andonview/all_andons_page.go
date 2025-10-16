package andonview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/format"
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

	content := g.Group([]g.Node{

		h.H3(g.Text("All Andons")),

		h.Form(
			h.ID("all-andon-table-form"),
			g.Attr("method", "GET"),

			allAndonsFilters(&allAndonsFiltersProps{
				availableFilters: p.AvailableFilters,
				activeFilters:    p.ActiveFilters,
			}),

			allAndonsTable(&allAndonsTableProps{
				sort:     p.Sort,
				andons:   p.Andons,
				count:    p.AndonsCount,
				pageSize: p.PageSize,
				page:     p.Page,
			}),
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

	type selectDef struct {
		label            string
		name             string
		availableFilters []string
		activeFilters    []string
	}

	var startDateValue, endDateValue g.Node
	if p.activeFilters.StartDate != nil {
		startDateValue = h.Value(p.activeFilters.StartDate.Format("2006-01-02"))
	}
	if p.activeFilters.EndDate != nil {
		endDateValue = h.Value(p.activeFilters.EndDate.Format("2006-01-02"))
	}

	return g.Group{
		h.Div(
			h.Class("andon-filters"),

			h.Div(
				h.Class("date-filters"),
				h.Label(
					g.Text("Start date"),
					h.Input(
						h.Name("StartDate"),
						h.Type("date"),
						startDateValue,
					),
				),
				h.Label(
					g.Text("End date"),
					h.Input(
						h.Name("EndDate"),
						h.Type("date"),
						endDateValue,
					),
				),
			),

			g.Map([]selectDef{
				{
					label:            "Location",
					name:             "LocationIn",
					availableFilters: p.availableFilters.LocationIn,
					activeFilters:    p.activeFilters.LocationIn,
				},
				{
					label:            "Issue",
					name:             "IssueIn",
					availableFilters: p.availableFilters.IssueIn,
					activeFilters:    p.activeFilters.IssueIn,
				},
				{
					label:            "Assigned Team",
					name:             "TeamIn",
					availableFilters: p.availableFilters.TeamIn,
					activeFilters:    p.activeFilters.TeamIn,
				},
				{
					label:            "Severity",
					name:             "SeverityIn",
					availableFilters: p.availableFilters.SeverityIn,
					activeFilters:    p.activeFilters.SeverityIn,
				},
				{
					label:            "Status",
					name:             "StatusIn",
					availableFilters: p.availableFilters.StatusIn,
					activeFilters:    p.activeFilters.StatusIn,
				},
				{
					label:            "Raised By",
					name:             "RaisedByUsernameIn",
					availableFilters: p.availableFilters.RaisedByUsernameIn,
					activeFilters:    p.activeFilters.RaisedByUsernameIn,
				},
				{
					label:            "Acknowledged By",
					name:             "AcknowledgedByUsernameIn",
					availableFilters: p.availableFilters.AcknowledgedByUsernameIn,
					activeFilters:    p.activeFilters.AcknowledgedByUsernameIn,
				},
				{
					label:            "Resolved By",
					name:             "ResolvedByUsernameIn",
					availableFilters: p.availableFilters.ResolvedByUsernameIn,
					activeFilters:    p.activeFilters.ResolvedByUsernameIn,
				},
			}, func(i selectDef) g.Node {
				return h.Label(
					g.Text(i.label),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        i.name,
						Placeholder: "-",
						Mode:        "multi",
						Options:     components.MapStringsToOptions(i.availableFilters, i.activeFilters),
						Selected:    strings.Join(i.activeFilters, ","),
					}),
				)
			}),
		),

		h.Button(
			h.Class("button primary"),
			h.Type("submit"),
			g.Text("GO"),
		),
	}
}

type allAndonsTableProps struct {
	sort     appsort.Sort
	andons   []model.Andon
	count    int
	pageSize int
	page     int
}

func allAndonsTable(p *allAndonsTableProps) g.Node {
	var columns = components.TableColumns{
		{TitleContents: g.Text("Location")},
		{TitleContents: g.Text("Issue Description")},
		{TitleContents: g.Text("Issue"), SortKey: "IssueName"},
		{TitleContents: g.Text("Assigned Team"), SortKey: "AssignedTeam"},
		{TitleContents: g.Text("Severity"), SortKey: "Severity"},
		{TitleContents: g.Text("Status"), SortKey: "Status"},
		{TitleContents: g.Text("Raised By"), SortKey: "RaisedByUsername"},
		{TitleContents: g.Text("Raised At"), SortKey: "RaisedAt"},
		{TitleContents: g.Text("Open Duration"), SortKey: "OpenDurationSeconds"},
		{TitleContents: g.Text("Acknowledged By"), SortKey: "AcknowledgedByUsername"},
		{TitleContents: g.Text("Acknowledged At"), SortKey: "AcknowledgedAt"},
		{TitleContents: g.Text("Resolved By"), SortKey: "ResolvedByUsername"},
		{TitleContents: g.Text("Resolved At"), SortKey: "ResolvedAt"},
		{TitleContents: g.Text("Updated At"), SortKey: "LastUpdated"},
	}

	var tableRows components.TableRows
	for _, a := range p.andons {
		namePathStr := strings.Join(a.NamePath, " > ")
		acknowledgedAt := "\u2013"
		if a.AcknowledgedAt != nil {
			acknowledgedAt = a.AcknowledgedAt.Format("2006-01-02 15:04:05")
		}
		acknowledgedBy := "\u2013"
		if a.AcknowledgedByUsername != nil {
			acknowledgedBy = *a.AcknowledgedByUsername
		}
		resolvedAt := "\u2013"
		if a.ResolvedAt != nil {
			resolvedAt = a.ResolvedAt.Format("2006-01-02 15:04:05")
		}
		resolvedBy := "\u2013"
		if a.ResolvedByUsername != nil {
			resolvedBy = *a.ResolvedByUsername
		}
		lastUpdated := "\u2013"
		if a.LastUpdated != nil {
			lastUpdated = a.LastUpdated.Format("2006-01-02 15:04:05")
		}

		cells := []components.TableCell{
			{Contents: g.Text(a.Location)},
			{Contents: g.Text(a.Description)},
			{Contents: g.Text(namePathStr)},
			{Contents: g.Text(a.AssignedTeamName)},
			{Contents: severityBadge(a.Severity, "small")},
			{Contents: statusBadge(a.Status, "small")},
			{Contents: g.Text(a.RaisedByUsername)},
			{Contents: g.Text(a.RaisedAt.Format("2006-01-02 15:04:05"))},
			{Contents: g.Text(format.FormatSecondsIntoDuration(a.OpenDurationSeconds))},
			{Contents: g.Text(acknowledgedBy)},
			{Contents: g.Text(acknowledgedAt)},
			{Contents: g.Text(resolvedBy)},
			{Contents: g.Text(resolvedAt)},
			{Contents: g.Text(lastUpdated)},
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
			HREF:  fmt.Sprintf("/andons/%d", a.AndonID),
		})
	}

	return components.Table(&components.TableProps{
		Columns: columns,
		Sort:    p.sort,
		Rows:    tableRows,
		Pagination: &components.TablePaginationProps{
			TotalRecords:        p.count,
			PageSize:            p.pageSize,
			CurrentPage:         p.page,
			CurrentPageQueryKey: "Page",
			PageSizeQueryKey:    "PageSize",
		},
	})

}
