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

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AndonPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	AndonID          int
	Andon            model.Andon
	AndonChangelog   []model.AndonChange
	AndonComments    []model.Comment
}

var changelogFieldDefs = []components.ChangelogProperty{
	{FieldKey: "Description", Label: g.Text("Description")},
	{FieldKey: "RaisedByUsername", Label: g.Text("Raised By")},
	{FieldKey: "AcknowledgedByUsername", Label: g.Text("Acknowledged By")},
	{FieldKey: "ResolvedByUsername", Label: g.Text("Resolved By")},
	{FieldKey: "CancelledByUsername", Label: g.Text("Cancelled By")},
}

func AndonPage(p *AndonPageProps) g.Node {

	andon := p.Andon
	namePathStr := strings.Join(andon.NamePath, " > ")

	var changelogEntries []components.ChangelogEntry
	for _, change := range p.AndonChangelog {
		entry := components.ChangelogEntry{
			ChangedAt:        change.ChangeAt,
			ChangeByUsername: change.ChangeByUsername,
			IsCreation:       change.IsCreation,
			Changes: map[string]any{
				"Description":            change.Description,
				"RaisedByUsername":       change.RaisedByUsername,
				"AcknowledgedByUsername": change.AcknowledgedByUsername,
				"ResolvedByUsername":     change.ResolvedByUsername,
				"CancelledByUsername":    change.CancelledByUsername,
			},
		}
		changelogEntries = append(changelogEntries, entry)
	}

	acknowledgedByUsername := "\u2013"
	if andon.AcknowledgedByUsername != nil {
		acknowledgedByUsername = *andon.AcknowledgedByUsername
	}
	acknowledgedAtStr := "\u2013"
	if andon.AcknowledgedAt != nil {
		acknowledgedAtStr = andon.AcknowledgedAt.Format("2006-01-02 15:04:05")
	}
	resolvedByUsername := "\u2013"
	if andon.ResolvedByUsername != nil {
		resolvedByUsername = *andon.ResolvedByUsername
	}
	resolvedAtStr := "\u2013"
	if andon.ResolvedAt != nil {
		resolvedAtStr = andon.ResolvedAt.Format("2006-01-02 15:04:05")
	}

	type attribute struct {
		label string
		value g.Node
	}
	attributes := []attribute{
		{label: "Location", value: g.Text(andon.Location)},
		{label: "Issue", value: g.Text(namePathStr)},
		{label: "Description", value: g.Text(andon.Description)},
		{label: "Source", value: g.Text(andon.Source)},
		{label: "Assigned Team", value: g.Text(andon.AssignedTeamName)},
		{label: "Raised By", value: g.Text(andon.RaisedByUsername)},
		{label: "Raised At", value: g.Text(andon.RaisedAt.Format("2006-01-02 15:04:05"))},
		{label: "Acknowledged By", value: g.Text(acknowledgedByUsername)},
		{label: "Acknowledged At", value: g.Text(acknowledgedAtStr)},
		{label: "Resolved By", value: g.Text(resolvedByUsername)},
		{label: "Resolved At", value: g.Text(resolvedAtStr)},
	}

	if andon.IsCancelled {
		attributes = append(attributes,
			attribute{
				label: "Cancelled By",
				value: g.Text(nilsafe.Str(andon.CancelledByUsername)),
			},
			attribute{
				label: "Cancelled At",
				value: g.Text(nilsafe.Time(andon.CancelledAt).Format("2006-01-02 15:04:05")),
			},
		)
	}

	content := g.Group([]g.Node{

		h.Div(
			h.Class("header"),

			h.Div(
				h.Class("title"),
				h.H3(g.Text(fmt.Sprintf("%s @ %s", namePathStr, andon.Location))),

				severityBadge(andon.Severity, "large"),

				statusBadge(andon.Status, "large"),
			),

			h.Div(
				h.Class("actions"),

				g.If(andon.CanUserAcknowledge,
					components.Button(&components.ButtonProps{},
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
				g.If(andon.CanUserResolve,
					components.Button(&components.ButtonProps{},
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
				g.If(andon.CanUserReopen,
					components.Button(&components.ButtonProps{},
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

				g.If(andon.CanUserCancel,
					components.Button(&components.ButtonProps{},
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
			components.InlineStyle("/internal/views/andonview/components.css"),
			components.InlineStyle("/internal/views/andonview/andon_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/andonview/andon_page.js"),
		},
	})
}
