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

type ServiceSchedulesPageProps struct {
	Ctx              reqcontext.ReqContext
	ShowArchived     bool
	Schedules        []model.ServiceSchedule
	Count            int
	AvailableFilters model.AndonAvailableFilters
	ActiveFilters    model.AndonFilters
	Sort             appsort.Sort
	Page             int
	PageSize         int
}

func ServiceSchedulesPage(p *ServiceSchedulesPageProps) g.Node {

	content := g.Group([]g.Node{

		h.Div(
			h.Class("header"),

			h.H3(g.Text("Service Schedules")),
			schedulesNav(),
		),

		h.Form(
			h.ID("service-schedules-form"),
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

			schedulesTable(&schedulesProps{
				sort:      p.Sort,
				schedules: p.Schedules,
				count:     p.Count,
				pageSize:  p.PageSize,
				page:      p.Page,
			}),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Service Schedules",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "account-wrench",
				Title:          "Services",
				URLPart:        "services",
			},
			{
				IconIdentifier: "clock",
				Title:          "Schedules",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/serviceview/schedules_page.css"),
		},
	})
}

func schedulesNav() g.Node {
	return h.Nav(
		h.Class("schedules-nav"),

		h.A(
			h.Class("button primary"),
			h.Href("/services/schedules/add"),
			components.Icon(&components.IconProps{
				Identifier: "plus",
			}),
			g.Text("Schedule"),
		),
	)

}

type schedulesProps struct {
	schedules []model.ServiceSchedule
	count     int

	sort     appsort.Sort
	pageSize int
	page     int
}

func schedulesTable(p *schedulesProps) g.Node {
	var columns = components.TableColumns{
		{TitleContents: g.Text("Name")},
		{TitleContents: g.Text("Metric")},
		{TitleContents: g.Text("Threshold")},
		{TitleContents: g.Text("Status")},
		{TitleContents: g.Text("Actions")},
	}

	var tableRows components.TableRows
	for _, s := range p.schedules {

		status := components.Badge(&components.BadgeProps{
			Type: components.BadgeSuccess,
			Size: components.BadgeSm,
		}, g.Text("Active"))
		if s.IsArchived {
			status = components.Badge(&components.BadgeProps{
				Type: components.BadgeSecondary,
				Size: components.BadgeSm,
			}, g.Text("Archived"))
		}

		cells := []components.TableCell{
			{Contents: g.Text(s.Name)},
			{Contents: g.Text(s.MetricName)},
			{Contents: g.Text(s.Threshold.String())},
			{Contents: status},
			{
				Contents: components.Button(&components.ButtonProps{
					Size:       "small",
					ButtonType: "primary",
					Link:       fmt.Sprintf("/services/schedules/%d/edit", s.ServiceScheduleID),
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
