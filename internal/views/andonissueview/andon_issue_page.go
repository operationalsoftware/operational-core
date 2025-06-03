package andonissueview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"fmt"
	"strings"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type AndonIssuePageProps struct {
	Ctx        reqcontext.ReqContext
	AndonIssue model.AndonIssue
}

func AndonIssuePage(p *AndonIssuePageProps) g.Node {

	andonIssue := p.AndonIssue

	content := g.Group([]g.Node{
		h.Div(
			h.Class("button-container"),
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Classes: c.Classes{
					"edit-button": true,
				},
				Link: fmt.Sprintf("/andon-issues/%d/edit", andonIssue.AndonIssueID),
			},
				components.Icon(&components.IconProps{
					Identifier: "pencil",
				}),
			),
		),
		h.Div(
			h.H3(g.Text("Andon Issue Details")),
			h.Div(
				h.Class("properties-grid"),
				g.Group([]g.Node{
					h.Span(
						h.Strong(g.Text("Issue Name")),
					),
					h.Span(
						g.Text(andonIssue.IssueName),
					),

					h.Span(
						h.Strong(g.Text("Issue Path")),
					),
					h.Span(
						g.Text(strings.Join(andonIssue.NamePath, " > ")),
					),

					h.Span(
						h.Strong(g.Text("Resolvable by Raiser?")),
					),
					h.Span(
						g.If(andonIssue.ResolvableByRaiser, g.Text("Yes")),
						g.If(!andonIssue.ResolvableByRaiser, g.Text("No")),
					),

					h.Span(
						h.Strong(g.Text("Will Stop Process?")),
					),
					h.Span(
						g.If(andonIssue.WillStopProcess, g.Text("Yes")),
						g.If(!andonIssue.WillStopProcess, g.Text("No")),
					),

					h.Span(
						h.Strong(g.Text("Is Archived?")),
					),
					h.Span(
						g.If(andonIssue.IsArchived, g.Text("Yes")),
						g.If(!andonIssue.IsArchived, g.Text("No")),
					),
				}),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Title: "Andon Issue: " + andonIssue.IssueName,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			andonIssuesBreadCrumb,
			{
				IconIdentifier: "alert-octagon-outline",
				Title:          andonIssue.IssueName,
			},
		},
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonissueview/andon_issue_page.css"),
		},
	})
}
