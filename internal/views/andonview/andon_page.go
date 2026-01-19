package andonview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/format"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/url"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AndonPageProps struct {
	Ctx                    reqcontext.ReqContext
	Values                 url.Values
	ValidationErrors       validate.ValidationErrors
	IsSubmission           bool
	AndonID                int
	Andon                  model.Andon
	GalleryURL             string
	GalleryImageURLs       []string
	AndonChangelog         []model.AndonChange
	AndonComments          []model.Comment
	ReturnTo               string
	AddCommentHMACEnvelope string
}

func AndonPage(p *AndonPageProps) g.Node {

	andon := p.Andon
	namePathStr := strings.Join(andon.NamePath, " > ")

	content := g.Group([]g.Node{

		h.Div(
			h.Class("header"),

			h.Div(
				h.Class("title"),
				h.H3(g.Text(fmt.Sprintf("%s @ %s", namePathStr, andon.Location))),

				severityBadge(andon.Severity, "large"),

				statusBadge(andon.Status, "large"),
			),

			andonActions(&andonActionsProps{andon: andon, returnTo: p.ReturnTo}),
		),

		h.Div(
			h.Class("description"),
			h.Pre(g.Text(andon.Description)),
		),

		h.Div(
			h.Class("two-column-flex"),

			andonAttributesList(&andonAttributesListProps{
				andon: p.Andon,
			}),

			h.Div(
				h.Class("gallery-container"),

				components.Gallery(p.GalleryImageURLs),

				h.A(
					h.Class("button primary"),
					h.Href(p.GalleryURL),

					components.Icon(&components.IconProps{
						Identifier: "arrow-right-thin",
					}),
				),
			),
		),

		h.Div(
			h.Class("two-column-flex"),
			components.CommentsThread(&components.CommentsThreadProps{
				Comments:        p.AndonComments,
				CommentThreadID: p.Andon.CommentThreadID,
				HMACEnvelope:    p.AddCommentHMACEnvelope,
				CanAddComment:   true,
			}),
			andonChangeLog(&andonChangeLogProps{
				changeLog: p.AndonChangelog,
			}),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: fmt.Sprintf("Andon: %s", andon.Description),
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
			components.InlineScript("/internal/views/andonview/components.js"),
		},
	})
}

type andonActionsProps struct {
	andon    model.Andon
	returnTo string
}

func andonActions(p *andonActionsProps) g.Node {
	andon := p.andon
	return h.Div(
		h.Class("actions"),

		g.If(andon.CanUserAcknowledge, acknowledgeButton(&acknowledgeButtonProps{
			andonID:  andon.AndonID,
			showText: true,
			returnTo: p.returnTo,
		})),
		g.If(andon.CanUserResolve && !andon.CanUserAcknowledge, resolveButton(&resolveButtonProps{
			andonID:  andon.AndonID,
			showText: true,
			returnTo: p.returnTo,
		})),
		g.If(andon.CanUserCancel, cancelButton(&cancelButtonProps{
			andonID:  andon.AndonID,
			showText: true,
			returnTo: p.returnTo,
		})),
		g.If(andon.CanUserReopen, reopenButton(&reopenButtonProps{
			andonID:  andon.AndonID,
			showText: true,
		})),
	)
}

type andonAttributesListProps struct {
	andon model.Andon
}

func andonAttributesList(p *andonAttributesListProps) g.Node {

	andon := p.andon
	namePathStr := strings.Join(andon.NamePath, " > ")

	source := "\u2013"
	if andon.Source != "" {
		source = andon.Source
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

	renderDuration := func(display, tooltip string) g.Node {
		if tooltip == "" {
			return g.Text(display)
		}

		return h.Span(
			g.Attr("title", tooltip),
			g.Text(display),
		)
	}

	downtimeDisplay, downtimeTooltip := format.FormatOptionalSecondsIntoMinutes(andon.DowntimeDurationSeconds)
	openDurationDisplay, openDurationTooltip := format.FormatOptionalSecondsIntoMinutes(&andon.OpenDurationSeconds)

	type attribute struct {
		label string
		value g.Node
	}
	attributes := []attribute{
		{label: "Issue", value: g.Text(namePathStr)},
		{label: "Location", value: g.Text(andon.Location)},
		{label: "Source", value: g.Text(source)},
		{label: "Assigned Team", value: g.Text(andon.AssignedTeamName)},
		{label: "Raised By", value: g.Text(andon.RaisedByUsername)},
		{label: "Raised At", value: g.Text(andon.RaisedAt.Format("2006-01-02 15:04:05"))},
		{label: "Open Duration (m)", value: renderDuration(openDurationDisplay, openDurationTooltip)},
		{label: "Downtime (m)", value: renderDuration(downtimeDisplay, downtimeTooltip)},
		{label: "Acknowledged By", value: g.Text(acknowledgedByUsername)},
		{label: "Acknowledged At", value: g.Text(acknowledgedAtStr)},
	}

	if andon.Severity != model.AndonSeverityInfo {
		attributes = append(attributes,
			attribute{
				label: "Resolved By",
				value: g.Text(resolvedByUsername),
			},
			attribute{
				label: "Resolved At",
				value: g.Text(resolvedAtStr),
			},
		)
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

	return h.Div(
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
	)
}

type andonChangeLogProps struct {
	changeLog []model.AndonChange
}

func andonChangeLog(p *andonChangeLogProps) g.Node {

	var changelogFieldDefs = []components.ChangelogProperty{
		{FieldKey: "Description", Label: g.Text("Description")},
		{FieldKey: "RaisedByUsername", Label: g.Text("Raised By")},
		{FieldKey: "AcknowledgedByUsername", Label: g.Text("Acknowledged By")},
		{FieldKey: "ResolvedByUsername", Label: g.Text("Resolved By")},
		{FieldKey: "CancelledByUsername", Label: g.Text("Cancelled By")},
		{FieldKey: "ReopenedByUsername", Label: g.Text("Reopened By")},
	}

	var changelogEntries []components.ChangelogEntry
	for _, change := range p.changeLog {
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
				"ReopenedByUsername":     change.ReopenedByUsername,
			},
		}
		changelogEntries = append(changelogEntries, entry)
	}

	return components.Changelog(changelogEntries, changelogFieldDefs)

}
