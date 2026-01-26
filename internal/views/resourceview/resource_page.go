package resourceview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	serviceview "app/internal/views/serviceview"
	"app/pkg/appsort"
	"app/pkg/format"
	"app/pkg/reqcontext"
	"fmt"

	"github.com/shopspring/decimal"
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type ResourcePageProps struct {
	Ctx            reqcontext.ReqContext
	Resource       model.Resource
	Services       []model.ResourceService
	CurrentMetrics []model.ResourceServiceMetricStatus
	LifetimeTotals []model.ServiceMetricLifetimeTotal
	ServiceCount   int
	Sort           appsort.Sort
	Page           int
	PageSize       int
	CanManage      bool
}

func ResourcePage(p *ResourcePageProps) g.Node {

	lastServicedAtStr := "\u2013"
	if p.Resource.LastServicedAt != nil {
		lastServicedAtStr = p.Resource.LastServicedAt.Format("2006-01-02 15:04:05")
	}

	type attribute struct {
		label string
		value g.Node
	}
	statusNode := components.Badge(&components.BadgeProps{
		Type: components.BadgeSuccess,
		Size: components.BadgeSm,
	}, g.Text("Active"))
	if p.Resource.IsArchived {
		statusNode = components.Badge(&components.BadgeProps{
			Type: components.BadgeSecondary,
			Size: components.BadgeSm,
		}, g.Text("Archived"))
	}

	attributes := []attribute{
		{label: "Type", value: g.Text(p.Resource.Type)},
		{label: "Reference", value: g.Text(p.Resource.Reference)},
		{
			label: "Service Ownership Team",
			value: g.Text(serviceOwnershipTeamLabel(p.Resource.ServiceOwnershipTeamName)),
		},
		{label: "Last Serviced At", value: g.Text(lastServicedAtStr)},
		{label: "Status", value: statusNode},
	}

	content := g.Group([]g.Node{

		h.Div(
			h.Class("header"),
			h.H3(g.Textf("%s \u2013 %s", p.Resource.Type, p.Resource.Reference)),

			resourceNav(&resourceNavProps{
				resourceID: p.Resource.ResourceID,
				canManage:  p.CanManage,
				isArchived: p.Resource.IsArchived,
			}),
		),

		h.Div(
			h.Class("attributes-list"),

			g.Group{g.Map(attributes, func(a attribute) g.Node {
				return h.Li(
					components.Icon(&components.IconProps{
						Identifier: "arrow-right-thin",
					}),
					h.Strong(g.Textf("%s: ", a.label)),
					h.Span(a.value),
				)
			})},
		),

		h.Div(
			h.Class("service-schedules"),

			h.Div(
				h.Class("service-schedule-header"),

				h.H3(g.Text("Service Schedules")),
			),

			currentMetricsTable(&currentMetricsTableProps{
				records:    p.CurrentMetrics,
				count:      p.ServiceCount,
				sort:       p.Sort,
				pageSize:   p.PageSize,
				page:       p.Page,
				resourceID: p.Resource.ResourceID,
			}),
			serviceview.StatusLegend(),

			g.If(len(p.LifetimeTotals) > 0,
				g.Group([]g.Node{
					h.H3(g.Text("Lifetime Totals")),
					lifetimeTotalsTable(&lifetimeTotalsTableProps{
						records: p.LifetimeTotals,
					}),
				}),
			),

			h.H3(g.Text("Service History")),

			serviceHistoryTable(&serviceHistoryTableProps{
				sort:     p.Sort,
				services: p.Services,
				count:    p.ServiceCount,
				pageSize: p.PageSize,
				page:     p.Page,
			}),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   fmt.Sprintf("Resource - %s", p.Resource.Reference),
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "cube-scan",
				Title:          "Resources",
				URLPart:        "resources",
			},
			{
				Title: p.Resource.Reference,
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/resourceview/resource_page.css"),
			components.InlineStyle("/internal/views/serviceview/components.css"),
			components.InlineScript("/internal/views/resourceview/resource_page.js"),
		},
	})
}

func serviceOwnershipTeamLabel(name *string) string {
	if name == nil || *name == "" {
		return "Unassigned"
	}
	return *name
}

type resourceNavProps struct {
	resourceID int
	canManage  bool
	isArchived bool
}

func resourceNav(p *resourceNavProps) g.Node {
	actions := []g.Node{}

	if !p.isArchived {
		if p.canManage {
			actions = append(actions,
				h.A(
					h.Class("button primary"),
					h.Href(fmt.Sprintf("/resources/%d/services/new", p.resourceID)),
					components.Icon(&components.IconProps{
						Identifier: "plus",
					}),
					g.Text("Start Service"),
				),
				h.A(
					h.Class("button primary"),
					h.Href(fmt.Sprintf("/services/resource/%d/schedules/add", p.resourceID)),
					components.Icon(&components.IconProps{
						Identifier: "plus",
					}),
					g.Text("Assign Service Schedule"),
				),

				h.A(
					h.Class("button primary"),
					h.Href(fmt.Sprintf("/resources/%d/metric-recording/add", p.resourceID)),
					components.Icon(&components.IconProps{
						Identifier: "plus",
					}),
					g.Text("Recording"),
				),

				h.A(
					h.Class("button primary"),
					h.Href(fmt.Sprintf("/resources/%d/edit", p.resourceID)),
					components.Icon(&components.IconProps{
						Identifier: "pencil",
					}),
					g.Text("Edit"),
				),
			)
		}
	}

	return h.Nav(
		h.Class("resource-nav"),
		g.Group(actions),
	)
}

type currentMetricsTableProps struct {
	records []model.ResourceServiceMetricStatus
	count   int

	sort       appsort.Sort
	pageSize   int
	page       int
	resourceID int
}

func currentMetricsTable(p *currentMetricsTableProps) g.Node {
	var columns = components.TableColumns{
		{TitleContents: g.Text("Metric")},
		{TitleContents: g.Text("Service Ownership Team")},
		{TitleContents: g.Text("Current Value"), Classes: c.Classes{"text-right": true}},
		{TitleContents: g.Text("Threshold"), Classes: c.Classes{"text-right": true}},
		{TitleContents: g.Text("Threshold Utilisation (%)"), Classes: c.Classes{"text-right": true}},
		{TitleContents: g.Text("Is Due?")},
		{TitleContents: g.Text("Last Recorded At")},
		{TitleContents: g.Text("Actions")},
	}

	var tableRows components.TableRows
	for _, r := range p.records {

		isDue := g.Text("No")
		if r.IsDue {
			isDue = g.Text("Yes")
		}

		lastRecordedAt := "\u2013"
		if r.LastRecordedAt != nil {
			lastRecordedAt = r.LastRecordedAt.Format("2006-01-02 15:04:05")
		}

		cells := []components.TableCell{
			{Contents: h.Div(
				h.Class("metric-cell"),
				h.Span(
					h.Class("metric-name"),
					g.Text(r.MetricName),
				),
			)},
			{Contents: g.Text(serviceOwnershipTeamLabel(r.ServiceOwnershipTeamName))},
			{Contents: g.Text(format.DecimalWithCommas(r.CurrentValue.String())), Classes: c.Classes{"text-right": true}},
			{Contents: g.Text(format.DecimalWithCommas(r.Threshold.String())), Classes: c.Classes{"text-right": true}},
			{Contents: g.Text(format.DecimalWithCommas(r.NormalisedPercentage.String())), Classes: c.Classes{"text-right": true}},
			{Contents: isDue},
			{Contents: g.Text(lastRecordedAt)},
		}

		cells = append(cells, components.TableCell{
			Contents: h.Form(
				h.Method("POST"),
				h.Class("unassign-service-schedule-form"),
				h.Action(fmt.Sprintf("/services/resource/%d/schedules/%d/unassign", p.resourceID, r.ServiceScheduleID)),
				h.Button(
					g.If(
						!r.CanUserManage,
						g.Group{
							h.Disabled(),
							g.Attr("title", "You do not have permission to unassign this service schedule"),
						},
					),
					h.Class("button secondary small"),
					h.Type("submit"),
					g.Text("Unassign"),
				),
			),
		})

		one := decimal.NewFromInt(1)
		nineTenths := decimal.NewFromFloat(0.9)
		eightTenths := decimal.NewFromFloat(0.8)
		var colourClass string
		if r.NormalisedValue.GreaterThanOrEqual(one) {
			colourClass = "is-due"
		} else if r.NormalisedValue.GreaterThanOrEqual(nineTenths) {
			colourClass = "threshold-90"
		} else if r.NormalisedValue.GreaterThanOrEqual(eightTenths) {
			colourClass = "threshold-80"
		}

		endIdx := len(cells) - 1

		if colourClass != "" {
			for i := 0; i < endIdx; i++ {
				if cells[i].Classes == nil {
					cells[i].Classes = c.Classes{}
				}
				cells[i].Classes[colourClass] = true
			}
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
		})
	}

	return components.Table(&components.TableProps{
		Columns: columns,
		Sort:    p.sort,
		Rows:    tableRows,
	})
}

type lifetimeTotalsTableProps struct {
	records []model.ServiceMetricLifetimeTotal
}

func lifetimeTotalsTable(p *lifetimeTotalsTableProps) g.Node {
	var columns = components.TableColumns{
		{TitleContents: g.Text("Metric")},
		{TitleContents: g.Text("Lifetime Total"), Classes: c.Classes{"text-right": true}},
	}

	var tableRows components.TableRows
	for _, r := range p.records {

		lifetime := "\u2013"
		if !r.LifetimeTotal.IsZero() {
			lifetime = format.DecimalWithCommas(r.LifetimeTotal.String())
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: []components.TableCell{
				{Contents: g.Text(r.MetricName)},
				{Contents: g.Text(lifetime), Classes: c.Classes{"text-right": true}},
			},
		})
	}

	return components.Table(&components.TableProps{
		Columns: columns,
		Rows:    tableRows,
	})
}

type serviceHistoryTableProps struct {
	services []model.ResourceService
	count    int

	sort     appsort.Sort
	pageSize int
	page     int
}

func serviceHistoryTable(p *serviceHistoryTableProps) g.Node {
	var columns = components.TableColumns{
		{TitleContents: g.Text("Status")},
		{TitleContents: g.Text("Completed By")},
		{TitleContents: g.Text("Completed At")},
		{TitleContents: g.Text("Notes")},
	}

	var tableRows components.TableRows
	for _, s := range p.services {

		completedAt := "\u2013"
		if s.CompletedAt != nil {
			completedAt = s.CompletedAt.Format("2006-01-02 15:04:05")
		}

		completedBy := "\u2013"
		if s.CompletedByUsername != nil {
			completedBy = *s.CompletedByUsername
		}

		notes := "\u2013"
		if s.Notes != "" {
			notes = s.Notes
			if len(s.Notes) > 100 {
				notes = s.Notes[:100] + "..."
			}
		}

		cells := []components.TableCell{
			{Contents: g.Text(string(s.Status))},
			{Contents: g.Text(completedBy)},
			{Contents: g.Text(completedAt)},
			{Contents: h.Pre(
				h.Class("service-note"),
				g.Text(notes),
			)},
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
			HREF: fmt.Sprintf("/services/%d",
				s.ResourceServiceID),
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
