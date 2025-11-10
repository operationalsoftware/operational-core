package serviceview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
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
					h.Class("link small"),
					h.Href("/services/metrics"),
					g.Text("Manage Metrics"),
				),

				h.A(
					h.Class("link small"),
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

		h.Form(
			h.Method("GET"),

			// resourceServicingFilters(&resourceServicingFiltersProps{
			// 	availableFilters: p.AvailableFilters,
			// 	activeFilters:    p.ActiveFilters,
			// }),

			resourceServicingTable(&resourceServicingProps{
				sort:      p.Sort,
				resources: p.Resources,
				count:     p.Count,
				pageSize:  p.PageSize,
				page:      p.Page,
			}),
			h.Div(
				h.Class("status-legend"),
				h.Div(
					h.Span(h.Class("status-dot threshold-80")),
					g.Text("> 80%"),
				),
				h.Div(
					h.Span(h.Class("status-dot threshold-90")),
					g.Text("> 90%"),
				),
				h.Div(
					h.Span(h.Class("status-dot is-due")),
					g.Text("Is Due"),
				),
			),
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
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/serviceview/servicing_page.js"),
		},
	})
}

func servicingNav() g.Node {
	return h.Nav(
		h.Class("servicing-nav"),

		// h.A(
		// 	h.Class("button primary small"),
		// 	h.Href(fmt.Sprintf("/resources/%d/services/new", p.resourceID)),
		// 	components.Icon(&components.IconProps{
		// 		Identifier: "plus",
		// 	}),
		// 	g.Text("Service"),
		// ),
	)
}

type resourceServicingFiltersProps struct {
	availableFilters model.AndonAvailableFilters
	activeFilters    model.AndonFilters
}

func resourceServicingFilters(p *resourceServicingFiltersProps) g.Node {
	type selectDef struct {
		label            string
		name             string
		availableFilters []string
		activeFilters    []string
	}

	return g.Group{
		h.Div(
			h.Class("resources-filters"),

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
		{TitleContents: g.Text("Metric")},
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

		cells := []components.TableCell{
			{Contents: g.Text(r.Reference)},
			{Contents: g.Text(serviceOwnershipTeamLabel(r.ServiceOwnershipTeamName))},
			{Contents: g.Text(r.MetricName)},
			{Contents: g.Text(r.CurrentValue.String())},
			{Contents: g.Text(r.Threshold.String())},
			{Contents: g.Text(r.NormalisedPercentage.String())},
			{Contents: g.Text(lastRecordedAt)},
			{Contents: g.Text(lastServicedAt)},
			{Contents: g.Group([]g.Node{

				h.Div(
					h.Class("andon-actions"),

					g.If(
						!r.HasWIPService,
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
