package serviceview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/format"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type ResourceServicingPageProps struct {
	Ctx              reqcontext.ReqContext
	ShowArchived     bool
	Resources        []model.ResourceServiceMetricStatus
	Count            int
	AvailableFilters model.AndonAvailableFilters
	ActiveFilters    model.AndonFilters
	Sort             appsort.Sort
	Page             int
	PageSize         int
	Teams            []model.Team
	SelectedTeamIDs  []string
}

func ResourceServicingPage(p *ResourceServicingPageProps) g.Node {

	content := g.Group([]g.Node{

		h.Div(
			h.Class("header"),

			h.Div(
				h.Class("actions"),
				h.A(
					h.Class("link"),
					h.Href("/services/metrics"),
					g.Text("Manage Metrics"),
				),

				h.A(
					h.Class("link"),
					h.Href("/services/schedules"),
					g.Text("Manage Schedules"),
				),

				h.A(
					h.Class("link"),
					h.Href("/services/all"),
					g.Text("All Services"),
				),
			),

			h.Div(
				h.Class("title"),
				h.H3(g.Text("Resource Servicing Priority List")),
			),
		),

		servicingTeamFilter(&servicingTeamFilterProps{
			teams:           p.Teams,
			selectedTeamIDs: p.SelectedTeamIDs,
		}),
		StatusLegend(),

		h.Form(
			h.Method("GET"),

			resourceServicingTable(&resourceServicingProps{
				sort:      p.Sort,
				resources: p.Resources,
				count:     p.Count,
				pageSize:  p.PageSize,
				page:      p.Page,
			}),
			StatusLegend(),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Resource Servicing Priority List",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "account-wrench",
				Title:          "Services",
				URLPart:        "services",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/serviceview/servicing_page.css"),
			components.InlineStyle("/internal/views/serviceview/components.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/serviceview/servicing_page.js"),
		},
	})
}

type servicingTeamFilterProps struct {
	teams           []model.Team
	selectedTeamIDs []string
}

func servicingTeamFilter(p *servicingTeamFilterProps) g.Node {
	selectedAttr := strings.Join(p.selectedTeamIDs, ",")

	return h.Div(
		h.Class("servicing-team-select"),
		h.Form(
			h.Method("GET"),
			h.ID("service-team-form"),
			h.Div(
				h.ID("service-team-select-wrapper"),
				g.Attr("data-selected", selectedAttr),
				components.SearchSelect(&components.SearchSelectProps{
					Name:        "ServiceOwnershipTeamIDs",
					Placeholder: "Service Ownership Team",
					Mode:        "multi",
					Options:     mapServiceTeamsToOptions(p.teams, p.selectedTeamIDs),
				},
					g.Attr("onchange", "handleServiceTeamSelectChange(event)"),
				),
			),
		),
	)
}

func mapServiceTeamsToOptions(teams []model.Team, selectedValues []string) []components.SearchSelectOption {
	out := make([]components.SearchSelectOption, len(teams))
	for i, team := range teams {
		value := strconv.Itoa(team.TeamID)
		out[i] = components.SearchSelectOption{
			Text:     team.TeamName,
			Value:    value,
			Selected: slices.Contains(selectedValues, value),
		}
	}
	return out
}

type resourceServicingProps struct {
	resources []model.ResourceServiceMetricStatus
	count     int

	sort     appsort.Sort
	pageSize int
	page     int
}

func resourceServicingTable(p *resourceServicingProps) g.Node {
	var columns = components.TableColumns{
		{TitleContents: g.Text("Reference")},
		{TitleContents: g.Text("Service Ownership Team")},
		{TitleContents: g.Text("Schedule")},
		{TitleContents: g.Text("Current Value")},
		{TitleContents: g.Text("Threshold")},
		{TitleContents: g.Text("Threshold Utilisation (%)")},
		{TitleContents: g.Text("Last Recorded At")},
		{TitleContents: g.Text("Last Serviced At")},
		{TitleContents: g.Text("Actions")},
	}

	var tableRows components.TableRows
	for _, r := range p.resources {

		lastRecordedAt := "\u2013"
		if r.LastRecordedAt != nil {
			lastRecordedAt = r.LastRecordedAt.Format("2006-01-02 15:04:05")
		}

		lastServicedAt := "\u2013"
		if r.LastServicedAt != nil {
			lastServicedAt = r.LastServicedAt.Format("2006-01-02 15:04:05")
		}

		scheduleCell := h.Div(
			h.Class("schedule-cell"),
			h.Span(
				h.Class("schedule-name"),
				g.Text(r.ServiceScheduleName),
			),
			h.Span(
				h.Class("metric-name"),
				g.Text(r.MetricName),
			),
		)

		cells := []components.TableCell{
			{Contents: g.Text(r.Reference)},
			{Contents: g.Text(serviceOwnershipTeamLabel(r.ServiceOwnershipTeamName))},
			{Contents: scheduleCell},
			{Contents: g.Text(format.DecimalWithCommas(r.CurrentValue.String())), Classes: c.Classes{
				"text-right": true,
			}},
			{Contents: g.Text(format.DecimalWithCommas(r.Threshold.String())), Classes: c.Classes{
				"text-right": true,
			}},
			{Contents: g.Text(r.NormalisedPercentage.String()), Classes: c.Classes{
				"text-right": true,
			}},
			{Contents: g.Text(lastRecordedAt)},
			{Contents: g.Text(lastServicedAt)},
			{Contents: g.Group([]g.Node{

				h.Div(
					h.Class("andon-actions"),

					g.If(
						!r.HasWIPService && r.CanUserManage,
						components.Button(&components.ButtonProps{
							Size:       "small",
							ButtonType: "primary",
							Link: fmt.Sprintf("/resources/%d/services/new",
								r.ResourceID),
						},
							g.Attr("title", "Start Service"),

							components.Icon(&components.IconProps{
								Identifier: "plus",
							}),
							components.Icon(&components.IconProps{
								Identifier: "account-wrench",
							}),
						),
					),
					g.If(
						r.HasWIPService,
						components.Button(&components.ButtonProps{
							Size:       "small",
							ButtonType: "primary",
							Link: fmt.Sprintf("/services/%d",
								nilsafe.Int(r.WIPServiceID)),
						},
							g.Attr("title", "Service In Progress"),

							components.Icon(&components.IconProps{
								Identifier: "clock-outline",
							}),
						),
					),
				),
			}),
			},
		}

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
		for i := 0; i < len(cells)-1; i++ {
			if cells[i].Classes == nil {
				cells[i].Classes = c.Classes{}
			}
			cells[i].Classes[colourClass] = true
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
			HREF: fmt.Sprintf("/resources/%d",
				r.ResourceID),
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

func serviceOwnershipTeamLabel(name *string) string {
	if name == nil || *name == "" {
		return "Unassigned"
	}
	return *name
}
