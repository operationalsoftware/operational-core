package teamview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"fmt"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type TeamPageProps struct {
	Ctx  reqcontext.ReqContext
	Team model.Team
}

func TeamPage(p *TeamPageProps) g.Node {

	team := p.Team

	content := g.Group([]g.Node{
		h.Div(
			h.Class("button-container"),
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Classes: c.Classes{
					"edit-button": true,
				},
				Link: fmt.Sprintf("/teams/%d/edit", team.TeamID),
			},
				components.Icon(&components.IconProps{
					Identifier: "pencil",
				}),
			),
		),
		h.Div(
			h.H3(g.Text("Team Details")),
			h.Div(
				h.Class("properties-grid"),
				g.Group([]g.Node{
					h.Span(
						h.Strong(g.Text("Team Name")),
					),
					h.Span(
						g.Text(team.TeamName),
					),

					h.Span(
						h.Strong(g.Text("Is Archived?")),
					),
					h.Span(
						g.If(team.IsArchived, g.Text("Yes")),
						g.If(!team.IsArchived, g.Text("No")),
					),
				}),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Title: "Team: " + team.TeamName,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			teamsBreadCrumb,
			{
				IconIdentifier: "account-group",
				Title:          team.TeamName,
			},
		},
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/teamview/team_page.css"),
		},
	})
}
