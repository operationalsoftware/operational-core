package andonview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
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

var changelogFieldDefs = []components.ChangelogProperty{
	{FieldKey: "IssueDescription", Label: g.Text("Issue Description")},
	{FieldKey: "IssueID", Label: g.Text("IssueID")},
	{FieldKey: "Location", Label: g.Text("Location")},
	{FieldKey: "Status", Label: g.Text("Status")},
	{FieldKey: "RaisedByUsername", Label: g.Text("Raised By")},
	{FieldKey: "AcknowledgedByUsername", Label: g.Text("Acknowledged By")},
	{FieldKey: "ResolvedByUsername", Label: g.Text("Resolved By")},
	{FieldKey: "CancelledByUsername", Label: g.Text("Cancelled By")},
}

func AndonDetailsPage(p *AndonDetailsPageProps) g.Node {

	andonEvent := p.AndonEvent
	namePathStr := strings.Join(andonEvent.NamePath, " > ")

	isAcknowledged := andonEvent.Status == "Acknowledged"
	isResolved := andonEvent.Status == "Resolved"
	isCancelled := andonEvent.Status == "Cancelled"
	twoMinutesPassed := time.Since(andonEvent.RaisedAt) > 2*time.Minute
	fiveMinutesPassed := time.Since(andonEvent.RaisedAt) > 5*time.Minute

	var changelogEntries []components.ChangelogEntry
	for _, change := range p.AndonChanges {
		entry := components.ChangelogEntry{
			ChangedAt:        change.ChangeAt,
			ChangeByUsername: change.ChangeByUsername,
			IsCreation:       change.IsCreation,
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

	content := g.Group([]g.Node{

		h.Div(
			h.Class("details"),

			h.Div(

				h.Div(
					h.Class("title"),

					h.H3(
						g.Text(fmt.Sprintf("%s @ %s", namePathStr, andonEvent.Location)),
					),

					h.Span(
						g.Text(" \u2013 "),
					),

					h.H3(
						c.Classes{
							"status-badge":   true,
							"amber":          twoMinutesPassed,
							"flashing-red":   fiveMinutesPassed,
							"light-green":    isAcknowledged,
							"flashing-green": isResolved,
							"flashing-grey":  isCancelled,
						},
						g.Text(andonEvent.Status),
					),

					h.Span(
						g.Text(" \u2013 "),
					),

					h.H3(
						c.Classes{
							"severity-badge":        true,
							"info":                  andonEvent.Severity == model.AndonSeverityInfo,
							"self-resolvable":       andonEvent.Severity == model.AndonSeveritySelfResolvable,
							"requires-intervention": andonEvent.Severity == model.AndonSeverityRequiresIntervention,
						},

						h.Class("severity-badge"),
						g.Text(string(andonEvent.Severity)),
					),
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
							g.Text(andonEvent.AssignedTeamName),
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
					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Raised At: "),
						),

						h.Span(
							g.Text(andonEvent.RaisedAt.Format("2006-01-02 15:04:05")),
						),
					),
					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Acknowledged By: "),
						),

						h.Span(
							g.If(
								andonEvent.AcknowledgedByUsername == nil,
								g.Text("\u2013"),
							),
							g.If(
								andonEvent.AcknowledgedByUsername != nil,
								g.Text(nilsafe.Str(andonEvent.AcknowledgedByUsername)),
							),
						),
					),
					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Acknowledged At: "),
						),

						h.Span(
							g.If(andonEvent.AcknowledgedAt == nil,
								g.Text("\u2013"),
							),
							g.If(andonEvent.AcknowledgedAt != nil,
								g.Text(nilsafe.Time(andonEvent.AcknowledgedAt).Format("2006-01-02 15:04:05")),
							),
						),
					),
					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Resolved By: "),
						),

						h.Span(
							g.If(
								andonEvent.ResolvedByUsername == nil,
								g.Text("\u2013"),
							),
							g.If(
								andonEvent.ResolvedByUsername != nil,
								g.Text(nilsafe.Str(andonEvent.ResolvedByUsername)),
							),
						),
					),
					h.Li(
						components.Icon(&components.IconProps{
							Identifier: "arrow-right-thin",
						}),

						h.Strong(
							g.Text("Resolved At: "),
						),

						h.Span(
							g.If(andonEvent.ResolvedAt == nil,
								g.Text("\u2013"),
							),
							g.If(andonEvent.ResolvedAt != nil,
								g.Text(nilsafe.Time(andonEvent.ResolvedAt).Format("2006-01-02 15:04:05")),
							),
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

				g.If(andonEvent.CanUserCancel,
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
		),

		h.Div(
			h.Class("comments-and-changelog-container"),

			components.CommentsThread(&components.CommentsThreadProps{
				Comments: p.AndonComments,
				Entity:   "Andon",
				EntityID: p.AndonID,
			}),

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
			components.InlineStyle("/internal/views/andonview/andon_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/andonview/andon_page.js"),
		},
	})
}
