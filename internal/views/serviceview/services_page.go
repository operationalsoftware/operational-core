package serviceview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/reqcontext"
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ServicesPageProps struct {
	Ctx              reqcontext.ReqContext
	ShowArchived     bool
	Services         []model.ResourceService
	Count            int
	AvailableFilters model.AndonAvailableFilters
	ActiveFilters    model.AndonFilters
	Sort             appsort.Sort
	Page             int
	PageSize         int
}

func ServicesPage(p *ServicesPageProps) g.Node {

	content := g.Group([]g.Node{

		h.H3(g.Text("All Services")),

		h.Form(
			g.Attr("method", "GET"),

			// servicesFilters(&servicesFiltersProps{
			// 	availableFilters: p.AvailableFilters,
			// 	activeFilters:    p.ActiveFilters,
			// }),

			servicesTable(&servicesTableProps{
				sort:     p.Sort,
				services: p.Services,
				count:    p.Count,
				pageSize: p.PageSize,
				page:     p.Page,
			}),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "All Services",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "account-wrench",
				Title:          "Services",
				URLPart:        "services",
			},
			{
				Title: "All",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/serviceview/services_page.css"),
			components.InlineStyle("/internal/views/serviceview/components.css"),
		},
	})
}

type servicesFiltersProps struct {
	availableFilters model.AndonAvailableFilters
	activeFilters    model.AndonFilters
}

func servicesFilters(p *servicesFiltersProps) g.Node {
	type selectDef struct {
		label            string
		name             string
		availableFilters []string
		activeFilters    []string
	}

	return g.Group{
		h.Div(
			h.Class("services-filters"),

			g.Map([]selectDef{
				{
					label:            "Location",
					name:             "LocationIn",
					availableFilters: p.availableFilters.LocationIn,
					activeFilters:    p.activeFilters.LocationIn,
				},
				{
					label:            "Status",
					name:             "StatusIn",
					availableFilters: p.availableFilters.StatusIn,
					activeFilters:    p.activeFilters.StatusIn,
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

type servicesTableProps struct {
	services []model.ResourceService
	count    int

	sort     appsort.Sort
	pageSize int
	page     int
}

func servicesTable(p *servicesTableProps) g.Node {
	var columns = components.TableColumns{
		{TitleContents: g.Text("Resource")},
		{TitleContents: g.Text("Status")},
		{TitleContents: g.Text("Started By")},
		{TitleContents: g.Text("Started At")},
		{TitleContents: g.Text("Completed By")},
		{TitleContents: g.Text("Completed At")},
		{TitleContents: g.Text("Cancelled By")},
		{TitleContents: g.Text("Cancelled At")},
	}

	var tableRows components.TableRows
	for _, r := range p.services {

		hrefStr := fmt.Sprintf("/resources/%d", r.ResourceID)

		completedAt := "\u2013"
		if r.CompletedAt != nil {
			completedAt = r.CompletedAt.Format("2006-01-02 15:04:05")
		}

		cancelledAt := "\u2013"
		if r.CancelledAt != nil {
			cancelledAt = r.CancelledAt.Format("2006-01-02 15:04:05")
		}

		completedBy := "\u2013"
		if r.CompletedByUsername != nil {
			completedBy = *r.CompletedByUsername
		}

		cancelledBy := "\u2013"
		if r.CancelledByUsername != nil {
			cancelledBy = *r.CancelledByUsername
		}

		cells := []components.TableCell{
			{Contents: h.Div(
				h.A(
					h.Href(hrefStr),
					g.Text(r.ResourceReference),
				),
			)},
			{Contents: g.Text(string(r.Status))},
			{Contents: g.Text(r.StartedByUsername)},
			{Contents: g.Text(r.StartedAt.Format("2006-01-02 15:04:05"))},
			{Contents: g.Text(completedBy)},
			{Contents: g.Text(completedAt)},
			{Contents: g.Text(cancelledBy)},
			{Contents: g.Text(cancelledAt)},
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
			HREF: fmt.Sprintf("/services/%d",
				r.ResourceServiceID),
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
