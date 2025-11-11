package resourceview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/reqcontext"
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type ResourcesPageProps struct {
	Ctx              reqcontext.ReqContext
	ShowArchived     bool
	Resources        []model.Resource
	ResourcesCount   int
	AvailableFilters model.AndonAvailableFilters
	ActiveFilters    model.AndonFilters
	Sort             appsort.Sort
	Page             int
	PageSize         int
}

func ResourcesPage(p *ResourcesPageProps) g.Node {

	content := g.Group([]g.Node{
		h.Div(
			h.Class("resources-header"),

			h.H3(g.Text("Resource Management")),
			resourcesHomeNav(),
		),

		h.Form(
			h.Method("GET"),

			components.Checkbox(
				&components.CheckboxProps{
					Name:    "IsArchived",
					Label:   "Show Archived",
					Value:   "true",
					Checked: p.ShowArchived,
					Classes: c.Classes{
						"filter-checkbox": true,
					},
				},
				g.Attr("onchange", "submitTableForm(this.form)"),
			),

			// resourcesFilters(&resourcesFiltersProps{
			// 	availableFilters: p.AvailableFilters,
			// 	activeFilters:    p.ActiveFilters,
			// }),

			resourcesTable(&resourcesProps{
				sort:         p.Sort,
				resources:    p.Resources,
				count:        p.ResourcesCount,
				pageSize:     p.PageSize,
				page:         p.Page,
				showArchived: p.ShowArchived,
			}),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Resource Management",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "cube-scan",
				Title:          "Resources",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/resourceview/resources_page.css"),
		},
	})
}

type resourcesFiltersProps struct {
	availableFilters model.AndonAvailableFilters
	activeFilters    model.AndonFilters
}

func resourcesHomeNav() g.Node {
	return h.Nav(
		h.Class("resources-nav"),

		h.A(
			h.Class("button primary"),
			h.Href("/resources/add"),
			components.Icon(&components.IconProps{
				Identifier: "plus",
			}),
			g.Text("Resource"),
		),
	)
}

func resourcesFilters(p *resourcesFiltersProps) g.Node {

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

type resourcesProps struct {
	resources []model.Resource
	count     int

	sort         appsort.Sort
	pageSize     int
	page         int
	showArchived bool
}

func resourcesTable(p *resourcesProps) g.Node {
	var columns = components.TableColumns{
		{TitleContents: g.Text("Reference"), SortKey: "Reference"},
		{TitleContents: g.Text("Type"), SortKey: "Type"},
		{TitleContents: g.Text("Service Ownership Team")},
		{TitleContents: g.Text("Last Serviced At"), SortKey: "LastServicedAt"},
	}
	if p.showArchived {
		columns = append(columns, components.TableColumn{
			TitleContents: g.Text("Status"),
		})
	}

	var tableRows components.TableRows
	for _, a := range p.resources {

		lastServicedAt := "\u2013"
		if a.LastServicedAt != nil {
			lastServicedAt = a.LastServicedAt.Format("2006-01-02 15:04:05")
		}
		teamName := "Unassigned"
		if a.ServiceOwnershipTeamName != nil && *a.ServiceOwnershipTeamName != "" {
			teamName = *a.ServiceOwnershipTeamName
		}

		cells := []components.TableCell{
			{Contents: g.Text(a.Reference)},
			{Contents: g.Text(a.Type)},
			{Contents: g.Text(teamName)},
			{Contents: g.Text(lastServicedAt)},
		}
		if p.showArchived {
			status := g.Text("")
			if a.IsArchived {
				status = components.Badge(&components.BadgeProps{
					Type: components.BadgeSecondary,
					Size: components.BadgeSm,
				}, g.Text("Archived"))
			}
			cells = append(cells, components.TableCell{Contents: status})
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: cells,
			HREF:  fmt.Sprintf("/resources/%d", a.ResourceID),
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
