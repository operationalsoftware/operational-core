package serviceview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/reqcontext"
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ServiceMetricsPageProps struct {
	Ctx              reqcontext.ReqContext
	ShowArchived     bool
	Metrics          []model.ServiceMetric
	MetricsCount     int
	AvailableFilters model.AndonAvailableFilters
	ActiveFilters    model.AndonFilters
	Sort             appsort.Sort
	Page             int
	PageSize         int
}

func ServiceMetricsPage(p *ServiceMetricsPageProps) g.Node {

	content := g.Group([]g.Node{

		h.Div(
			h.Class("header"),

			h.H3(g.Text("Service Metrics")),
			metricsNav(),
		),

		h.Form(
			h.ID("service-metrics-form"),
			h.Method("GET"),

			components.Checkbox(
				&components.CheckboxProps{
					Name:    "ShowArchived",
					Label:   "Show Archived",
					Checked: p.ShowArchived,
					Value:   "true",
				},
				g.Attr("onchange", "submitTableForm(this.form)"),
			),

			metricsTable(&metricsProps{
				sort:     p.Sort,
				metrics:  p.Metrics,
				count:    p.MetricsCount,
				pageSize: p.PageSize,
				page:     p.Page,
			}),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Service Metrics",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "account-wrench",
				Title:          "Services",
				URLPart:        "services",
			},
			{
				IconIdentifier: "speedometer",
				Title:          "Metrics",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/serviceview/metrics_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/serviceview/metrics_page.js"),
		},
	})
}

func metricsNav() g.Node {
	return h.Nav(
		h.Class("metrics-nav"),

		h.A(
			h.Class("button primary"),
			h.Href("/services/metrics/add"),
			components.Icon(&components.IconProps{
				Identifier: "plus",
			}),
			g.Text("Metric"),
		),
	)

}

type metricsProps struct {
	metrics []model.ServiceMetric
	count   int

	sort     appsort.Sort
	pageSize int
	page     int
}

func metricsTable(p *metricsProps) g.Node {
	var columns = components.TableColumns{
		{TitleContents: g.Text("Name")},
		{TitleContents: g.Text("Description")},
		{TitleContents: g.Text("Is Cumulative?")},
		{TitleContents: g.Text("Status")},
		{TitleContents: g.Text("Actions")},
	}

	var tableRows components.TableRows
	for _, m := range p.metrics {

		isCumulative := g.Text("No")
		if m.IsCumulative {
			isCumulative = g.Text("Yes")
		}
		status := components.Badge(&components.BadgeProps{
			Type: components.BadgeSuccess,
			Size: components.BadgeSm,
		}, g.Text("Active"))
		if m.IsArchived {
			status = components.Badge(&components.BadgeProps{
				Type: components.BadgeSecondary,
				Size: components.BadgeSm,
			}, g.Text("Archived"))
		}

		cells := []components.TableCell{
			{Contents: g.Text(m.Name)},
			{Contents: g.Text(m.Description)},
			{Contents: isCumulative},
			{Contents: status},
			{
				Contents: components.Button(&components.ButtonProps{
					Size:       "small",
					ButtonType: "primary",
					Link:       fmt.Sprintf("/services/metrics/%d/edit", m.ServiceMetricID),
				},
					components.Icon(&components.IconProps{
						Identifier: "pencil",
					}),
					g.Text("Edit"),
				),
			},
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
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
