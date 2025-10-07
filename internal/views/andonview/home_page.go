package andonview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"fmt"
	"slices"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type HomePageProps struct {
	Ctx            reqcontext.ReqContext
	NewAndons      []model.Andon
	NewAndonsCount int
	WIPAndons      []model.Andon
	WIPAndonsCount int
	Teams          []model.Team
	SelectedTeams  []string
	NewSort        appsort.Sort
	WIPSort        appsort.Sort
	ReturnTo       string
}

func HomePage(p *HomePageProps) g.Node {

	content := g.Group([]g.Node{
		andonsHomeNav(&andonsHomeNavProps{
			isUserAndonAdmin: p.Ctx.User.Permissions.Andon.Admin,
		}),

		h.Div(
			h.Class("team-select"),

			h.Form(
				g.Attr("method", "GET"),
				g.Attr("id", "team-form"),

				h.Div(
					g.Attr("data-selected", strings.Join(p.SelectedTeams, ",")),
					g.Attr("id", "search-select-wrapper"),

					components.SearchSelect(&components.SearchSelectProps{
						Name:        "AndonTeams",
						Placeholder: "Select a team",
						Mode:        "multi",
						Options:     mapTeamsToOptions(p.Teams, p.SelectedTeams),
					},
						g.Attr("onchange", "handleTeamSelectChange(event)"),
					),
				),
			),
		),

		h.Form(
			g.Attr("method", "GET"),

			h.H3(h.Class("table-title"),
				h.Title("New andons - require actioning or acknowledgement"),
				g.Text("New"),
			),
			h.Hr(),
			newAndonsTable(&newAndonsTableProps{
				andons:   p.NewAndons,
				sort:     p.NewSort,
				returnTo: p.ReturnTo,
			}),
			statusLegend(),

			h.H3(h.Class("table-title wip-heading"),
				h.Title("Recently acknowledged and open andons"),
				g.Text("Acknowledged"),
			),
			h.Hr(),
			wipAndonsTable(&wipAndonsTableProps{
				andons:   p.WIPAndons,
				sort:     p.WIPSort,
				returnTo: p.ReturnTo,
			}),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Andons",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "alert-octagon-outline",
				Title:          "Andons",
				URLPart:        "andons",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonview/components.css"),
			components.InlineStyle("/internal/views/andonview/home_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/andonview/home_page.js"),
			components.InlineScript("/internal/views/andonview/components.js"),
		},
	})
}

var andonIssuesBreadCrumb = layout.Breadcrumb{
	IconIdentifier: "alert-octagon-outline",
	Title:          "Andons",
	URLPart:        "andons",
}

type andonsHomeNavProps struct {
	isUserAndonAdmin bool
}

func andonsHomeNav(p *andonsHomeNavProps) g.Node {
	return h.Nav(
		h.Class("andon-nav"),

		h.A(
			h.Class("button primary small"),
			h.Href("/andons/add"),
			components.Icon(&components.IconProps{
				Identifier: "plus",
			}),
			g.Text("New Andon"),
		),

		h.A(h.Href("/andons/all"), g.Text("All Andons")),

		g.If(
			p.isUserAndonAdmin,
			h.A(
				h.Href("/andon-issues"),

				components.Icon(&components.IconProps{
					Identifier: "wrench-outline",
					Classes: c.Classes{
						"icon": true,
					},
				},
				),
				g.Text("Andon Issues")),
		),
	)

}

func mapTeamsToOptions(teams []model.Team, selectedValues []string) []components.SearchSelectOption {
	out := make([]components.SearchSelectOption, len(teams))
	for i, v := range teams {
		isSelected := slices.Contains(selectedValues, v.TeamName)

		out[i] = components.SearchSelectOption{
			Text:     v.TeamName,
			Value:    v.TeamName,
			Selected: isSelected,
		}
	}
	return out
}

var commonColumns = components.TableColumns{
	{TitleContents: g.Text("Issue"), SortKey: "NamePath"},
	{TitleContents: g.Text("Location")},
	{TitleContents: g.Text("Description")},
	{TitleContents: g.Text("Assigned Team"), SortKey: "AssignedTeamName"},
	{TitleContents: g.Text("Severity"), SortKey: "Severity"},
}

type newAndonsTableProps struct {
	andons   []model.Andon
	sort     appsort.Sort
	returnTo string
}

func newAndonsTable(p *newAndonsTableProps) g.Node {
	var columns = append(
		commonColumns,
		components.TableColumns{
			{TitleContents: g.Text("Raised By"), SortKey: "RaisedByUsername"},
			{TitleContents: g.Text("Raised At"), SortKey: "RaisedAt"},
			{TitleContents: g.Text("Actions")},
		}...,
	)

	var rows []components.TableRow
	for _, a := range p.andons {

		namePathStr := strings.Join(a.NamePath, " > ")

		cells := []components.TableCell{
			{Contents: g.Text(namePathStr)},
			{Contents: g.Text(a.Location)},
			{Contents: g.Text(a.Description)},
			{Contents: g.Text(a.AssignedTeamName)},
			{Contents: severityBadge(a.Severity, "small")},
			{Contents: g.Text(a.RaisedByUsername)},
			{Contents: g.Text(a.RaisedAt.Format("2006-01-02 15:04:05"))},
			{
				Contents: g.Group([]g.Node{
					h.Div(
						h.Class("andon-actions"),
						acknowledgeButton(&acknowledgeButtonProps{
							andonID:  a.AndonID,
							returnTo: p.returnTo,
							isSmall:  true,
						}),
						cancelButton(&cancelButtonProps{
							andonID:  a.AndonID,
							returnTo: p.returnTo,
							isSmall:  true,
						}),
					),
				}),
			},
		}

		colourClass := "outstanding"
		if a.Severity == model.AndonSeverityInfo {
			colourClass = "requires-acknowledgement"
		} else if time.Since(a.RaisedAt) > 5*time.Minute {
			colourClass = "five-minutes-passed"
		}
		for i := 0; i < len(cells)-1; i++ {
			if cells[i].Classes == nil {
				cells[i].Classes = c.Classes{}
			}
			cells[i].Classes[colourClass] = true
		}

		rows = append(rows, components.TableRow{
			Cells: cells,
			HREF:  fmt.Sprintf("/andons/%d", a.AndonID),
		})
	}

	return components.Table(&components.TableProps{
		Columns:      columns,
		SortQueryKey: "NewSort",
		Sort:         p.sort,
		Rows:         rows,
	})
}

type wipAndonsTableProps struct {
	andons   []model.Andon
	sort     appsort.Sort
	returnTo string
}

func wipAndonsTable(p *wipAndonsTableProps) g.Node {
	var columns = append(
		commonColumns,
		components.TableColumns{
			{TitleContents: g.Text("Acknowledged By"), SortKey: "AcknowledgedByUsername"},
			{TitleContents: g.Text("Acknowledged At"), SortKey: "AcknowledgedAt"},
			{TitleContents: g.Text("Actions")},
		}...,
	)

	var rows components.TableRows
	for _, a := range p.andons {
		namePathStr := strings.Join(a.NamePath, " > ")

		cells := []components.TableCell{
			{Contents: g.Text(namePathStr)},
			{Contents: g.Text(a.Location)},
			{Contents: g.Text(a.Description)},
			{Contents: g.Text(a.AssignedTeamName)},
			{Contents: severityBadge(a.Severity, "small")},
			{Contents: g.Text(nilsafe.Str(a.AcknowledgedByUsername))},
			{Contents: g.Text(a.AcknowledgedAt.Format("2006-01-02 15:04:05"))},
			{
				Contents: g.Group([]g.Node{
					h.Div(
						h.Class("andon-actions"),

						resolveButton(&resolveButtonProps{
							andonID: a.AndonID,
							isSmall: true,
						}),

						cancelButton(&cancelButtonProps{
							andonID: a.AndonID,
							isSmall: true,
						}),
					),
				}),
			},
		}

		rows = append(rows, components.TableRow{
			Cells: cells,
			HREF:  fmt.Sprintf("/andons/%d", a.AndonID),
		})
	}

	return components.Table(&components.TableProps{
		Columns:      columns,
		SortQueryKey: "WIPSort",
		Sort:         p.sort,
		Rows:         rows,
	})
}

func statusLegend() g.Node {
	return h.Div(
		h.Class("status-legend"),
		h.Div(
			h.Span(h.Class("status-dot outstanding")),
			g.Text("Outstanding"),
		),
		h.Div(
			h.Span(h.Class("status-dot five-minutes-passed")),
			g.Text("Outstanding (> 5 minutes)"),
		),
		h.Div(
			h.Span(h.Class("status-dot requires-acknowledgement")),
			g.Text("Requires Acknowledgement"),
		),
	)
}
