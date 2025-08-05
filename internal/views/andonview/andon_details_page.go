package andonview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type AndonDetailsPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	AndonID          int
	AndonEvent       model.AndonEvent
	AndonChanges     []model.AndonChange
	AndonComments    []model.Comment
}

var changelogFieldDefs = []components.ChangelogFieldDefinition{
	{Name: "IssueDescription", Label: "Issue Description"},
	{Name: "IssueID", Label: "IssueID"},
	{Name: "Location", Label: "Location"},
	{Name: "Status", Label: "Status"},
	{Name: "RaisedByUsername", Label: "Raised By"},
	{Name: "AcknowledgedByUsername", Label: "Acknowledged By"},
	{Name: "ResolvedByUsername", Label: "Resolved By"},
	{Name: "CancelledByUsername", Label: "Cancelled By"},
}

func AndonDetailsPage(p *AndonDetailsPageProps) g.Node {

	andonEvent := p.AndonEvent
	namePathStr := strings.Join(andonEvent.NamePath, " > ")
	// hasInfoSeverity := andonEvent.Severity == "Info"

	// isSelfResolvable := false
	// if andonEvent.Severity == "Self-resolvable" && andonEvent.IsTeamMate && andonEvent.Status == "Outstanding" {
	// 	isSelfResolvable = true
	// }
	// isAckBtnEnabled := true
	// if andonEvent.Severity == "Requires Intervention" && !andonEvent.IsTeamMate {
	// 	isAckBtnEnabled = false
	// }

	// isResolveBtnEnabled := true
	// if hasInfoSeverity {
	// 	isResolveBtnEnabled = false
	// }
	// if andonEvent.Severity == "Requires Intervention" && !andonEvent.IsTeamMate {
	// 	isResolveBtnEnabled = false
	// }

	var changelogEntries []components.ChangelogEntry
	for _, change := range p.AndonChanges {
		entry := components.ChangelogEntry{
			ChangedAt:         change.ChangeAt,
			ChangedByUsername: change.ChangeByUsername,
			Changes: map[string]interface{}{
				"IssueDescription":       change.IssueDescription,
				"IssueID":                change.IssueID,
				"Location":               change.Location,
				"Status":                 change.Status,
				"RaisedByUsername":       change.RaisedByUsername,
				"AcknowledgedByUsername": change.AcknowledgedByUsername,
				"ResolvedByUsername":     change.ResolvedByUsername,
				"CancelledByUsername":    change.CancelledByUsername,
			},
		}
		changelogEntries = append(changelogEntries, entry)
	}

	andonComments := []g.Node{}
	for _, comment := range p.AndonComments {
		commentNode := h.Div(
			h.Class("comment"),

			h.P(
				g.Text(comment.Comment),
			),

			h.Div(
				h.Class("date"),

				h.Strong(
					g.Text(comment.CommentedAt.Format("2006-01-02 15:04:05")),
				),

				h.Div(
					h.Class("username"),

					g.Text(comment.CommentedByUsername),
				),
			),
		)
		andonComments = append(andonComments, commentNode)
	}

	content := g.Group([]g.Node{

		h.Div(
			h.Class("details"),

			h.Div(

				h.H3(
					g.Text(fmt.Sprintf("%s @ %s", namePathStr, andonEvent.Location)),
				),

				h.Div(
					h.Class("detail-list"),

					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Location: "),
						),

						h.Span(
							g.Text(andonEvent.Location),
						),
					),
					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Issue: "),
						),

						h.Span(
							g.Text(namePathStr),
						),
					),
					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Issue Description: "),
						),

						h.Span(
							g.Text(andonEvent.IssueDescription),
						),
					),
					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Status: "),
						),

						h.Span(
							g.Text(andonEvent.Status),
						),
					),
					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Source: "),
						),

						h.Span(
							g.Text(andonEvent.Source),
						),
					),
					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Assigned Team: "),
						),

						h.Span(
							g.Text(andonEvent.AssignedTeam),
						),
					),
					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Raised By: "),
						),

						h.Span(
							g.Text(andonEvent.RaisedByUsername),
						),
					),
				),
			),

			h.Div(
				h.Class("actions"),

				g.If(andonEvent.Status == "Outstanding" && andonEvent.CanUserAcknowledge,
					components.Button(&components.ButtonProps{
						Size:       "small",
						ButtonType: "button",
					},
						g.Attr("onclick", "updateAndon(event)"),
						g.Attr("data-id", strconv.Itoa(p.AndonID)),
						g.Attr("data-action", "acknowledge"),
						g.Attr("title", "Acknowledge"),

						components.Icon(&components.IconProps{
							Identifier: "gesture-tap-hold",
						}),

						g.Text("Acknowledge"),
					),
				),
				g.If(andonEvent.Status == "Acknowledged" && andonEvent.CanUserResolve,
					components.Button(&components.ButtonProps{
						Size:       "small",
						ButtonType: "button",
					},
						g.Attr("onclick", "updateAndon(event)"),
						g.Attr("data-id", strconv.Itoa(p.AndonID)),
						g.Attr("data-action", "resolve"),
						g.Attr("title", "Resolve"),

						components.Icon(&components.IconProps{
							Identifier: "check",
						}),

						g.Text("Resolve"),
					),
				),
				g.If(andonEvent.Status == "Outstanding" && andonEvent.Severity == "Self-resolvable" && andonEvent.CanUserResolve,
					components.Button(&components.ButtonProps{
						Size:       "small",
						ButtonType: "button",
					},
						g.Attr("onclick", "updateAndon(event)"),
						g.Attr("data-id", strconv.Itoa(p.AndonID)),
						g.Attr("data-action", "resolve"),
						g.Attr("title", "Resolve"),

						components.Icon(&components.IconProps{
							Identifier: "check",
						}),

						g.Text("Resolve"),
					),
				),
				g.If(andonEvent.Status == "Cancelled",
					components.Button(&components.ButtonProps{
						Size:       "small",
						ButtonType: "button",
					},
						g.Attr("onclick", "updateAndon(event)"),
						g.Attr("data-id", strconv.Itoa(p.AndonID)),
						g.Attr("data-action", "reopen"),
						g.Attr("title", "Reopen"),

						components.Icon(&components.IconProps{
							Identifier: "restore",
						}),

						g.Text("Reopen"),
					),
				),

				components.Button(&components.ButtonProps{
					Size: "small",
				},
					g.Attr("onclick", "updateAndon(event)"),
					g.Attr("data-id", strconv.Itoa(p.AndonID)),
					g.Attr("data-action", "cancel"),
					g.Attr("title", "Cancel"),

					components.Icon(&components.IconProps{
						Identifier: "cancel",
					}),

					g.Text("Cancel"),
				),
			),
		),

		h.Div(
			h.Class("history-section"),

			h.Div(
				h.Class("comment-section"),

				h.H3(
					g.Text("Comments"),
				),

				h.Div(
					h.Class("comments"),

					g.Group(andonComments),
				),

				h.FormEl(
					h.Action("/andons/add/comment"),
					h.Method("POST"),

					h.Div(
						h.Class("comment-box"),

						h.Textarea(
							h.Class("new-comment"),
							h.Name("Comment"),

							h.Placeholder("Enter Comment"),
						),

						h.Input(
							h.Name("AndonID"),
							h.Type("hidden"),
							h.Value(fmt.Sprintf("%d", p.AndonID)),
						),

						components.Button(&components.ButtonProps{
							Classes: c.Classes{
								"add-comment-btn": true,
							},
							ButtonType: "submit",
						},
							components.Icon(&components.IconProps{
								Identifier: "comment-text-outline",
							}),
							g.Text(" Comment"),
						),
					),
				),
			),

			h.Br(),
			h.Br(),

			h.Div(
				h.Class("change-log"),
				components.Changelog(changelogEntries, changelogFieldDefs),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Add New Andon Issue",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			andonIssuesBreadCrumb,
			{Title: "Details"},
		},
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonview/andon_details_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/andonview/andon_details_page.js"),
		},
	})
}
