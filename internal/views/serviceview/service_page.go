package serviceview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/url"
	"strconv"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type ResourceServicePageProps struct {
	Ctx                     reqcontext.ReqContext
	Values                  url.Values
	ValidationErrors        validate.ValidationErrors
	IsSubmission            bool
	ResourceService         model.ResourceService
	LastResourceService     *model.ResourceService
	ResourceServiceComments []model.Comment
	ServiceChangelog        []model.ResourceServiceChange
	CommentHMACEnvelope     string
	ReturnTo                string
	GalleryImageURLs        []string
}

func ResourceServicePage(p *ResourceServicePageProps) g.Node {

	service := p.ResourceService
	isWIPService := p.ResourceService.Status == model.ServiceStatusWorkInProgress
	canUserEdit := p.Ctx.User.Permissions.UserAdmin.Access

	galleryButtonText := "View Service Images"
	if len(p.GalleryImageURLs) == 0 {
		galleryButtonText = "Add Service Images"
	}

	content := g.Group([]g.Node{

		h.Div(
			h.Class("header"),

			h.Div(
				h.Class("title"),
				h.H3(g.Textf("Service of %s", service.ResourceReference)),
			),

			h.Div(
				h.Class("actions"),

				g.If(
					!isWIPService,

					h.Button(
						c.Classes{"button": true},
						g.Attr("onclick", "updateService(event)"),
						h.Title("Reopen"),
						h.Data("id", strconv.Itoa(p.ResourceService.ResourceID)),
						h.Data("service-id", strconv.Itoa(p.ResourceService.ResourceServiceID)),
						h.Data("action", "reopen"),

						components.Icon(&components.IconProps{
							Identifier: "restore",
						}),

						g.Text("Reopen"),
					),
				),

				g.If(
					isWIPService,

					h.Button(
						c.Classes{"button": true, "resolve": true},
						g.Attr("onclick", "updateService(event)"),
						h.Title("Complete"),
						h.Data("id", strconv.Itoa(p.ResourceService.ResourceID)),
						h.Data("service-id", strconv.Itoa(p.ResourceService.ResourceServiceID)),
						h.Data("action", "complete"),

						components.Icon(&components.IconProps{
							Identifier: "check",
						}),

						g.Text("Complete"),
					),
				),
				g.If(
					isWIPService,
					h.Button(
						c.Classes{"button": true, "danger": true},
						g.Attr("onclick", "updateService(event)"),
						h.Title("Cancel"),
						h.Data("id", strconv.Itoa(p.ResourceService.ResourceID)),
						h.Data("service-id", strconv.Itoa(p.ResourceService.ResourceServiceID)),
						h.Data("action", "cancel"),

						components.Icon(&components.IconProps{
							Identifier: "cancel",
						}),

						g.Text("Cancel"),
					),
				),
			),
		),

		h.Div(
			h.Class("two-column-flex"),

			h.Div(
				serviceAttributesList(&serviceAttributesListProps{
					service:             p.ResourceService,
					lastResourceService: p.LastResourceService,
				}),

				h.Form(
					h.Method("POST"),

					h.H3(g.Text("Notes")),
					h.Div(
						h.Class("description"),
						h.Textarea(
							h.Rows("5"),
							h.Name("Notes"),
							h.Placeholder("Enter notes"),
							g.Text(service.Notes),
						),
					),

					h.Button(
						h.Class("button primary notes-btn"),
						h.Type("submit"),
						g.Text("Save"),
					),
				),
			),

			h.Div(
				h.Class("gallery-container"),

				g.If(
					len(p.GalleryImageURLs) > 0,
					components.Gallery(p.GalleryImageURLs),
				),

				h.A(
					h.Class("button primary"),
					h.Href(p.ResourceService.GalleryURL),

					g.Text(galleryButtonText),

					components.Icon(&components.IconProps{
						Identifier: "arrow-right-thin",
					}),
				),
			),
		),

		h.Div(
			h.Class("two-column-flex"),
			components.CommentsThread(&components.CommentsThreadProps{
				Comments:        p.ResourceServiceComments,
				CommentThreadID: p.ResourceService.CommentThreadID,
				HMACEnvelope:    p.CommentHMACEnvelope,
				CanAddComment:   canUserEdit,
			}),
			serviceChangeLog(&serviceChangeLogProps{
				changeLog: p.ServiceChangelog,
			}),
		),
	})

	pageTitle := fmt.Sprintf(
		"Service Details • %s (%s)",
		service.ResourceReference,
		service.ResourceType,
	)

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: pageTitle,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "cube-scan",
				Title:          "Services",
				URLPart:        "services",
			},
			{Title: fmt.Sprintf("Service of %s", service.ResourceReference)},
		},
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/serviceview/service_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/serviceview/service_page.js"),
		},
	})
}

type serviceAttributesListProps struct {
	service             model.ResourceService
	lastResourceService *model.ResourceService
}

func serviceAttributesList(p *serviceAttributesListProps) g.Node {

	service := p.service
	lastService := p.lastResourceService

	raisedByUsername := "\u2013"
	if service.StartedByUsername != "" {
		raisedByUsername = service.StartedByUsername
	}

	type attribute struct {
		label string
		value g.Node
	}
	attributes := []attribute{
		{
			label: "Resource",
			value: h.A(
				h.Href(fmt.Sprintf("/resources/%d", service.ResourceID)),
				h.Class("resource-link"),
				g.Text(service.ResourceReference),
			),
		},
		{label: "Resource Type", value: g.Text(service.ResourceType)},
		{label: "Raised By", value: g.Text(raisedByUsername)},
		{label: "Raised At", value: g.Text(service.StartedAt.Format("2006-01-02 15:04:05"))},
	}

	lastServiceValue := g.Text("\u2013")
	if lastService != nil {
		lastServiceValue = h.A(
			h.Href(fmt.Sprintf("/services/%d", lastService.ResourceServiceID)),
			h.Class("resource-link"),
			g.Textf(
				"%s (%s)",
				lastService.StartedAt.Format("2006-01-02 15:04:05"),
				lastService.Status,
			),
		)
	}

	attributes = append(attributes, attribute{
		label: "Last Service",
		value: lastServiceValue,
	})

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

type serviceChangeLogProps struct {
	changeLog []model.ResourceServiceChange
}

func serviceChangeLog(p *serviceChangeLogProps) g.Node {

	var changelogFieldDefs = []components.ChangelogProperty{
		{FieldKey: "Notes", Label: g.Text("Notes")},
		{FieldKey: "RaisedByUsername", Label: g.Text("Raised By")},
		{FieldKey: "CompletedByUsername", Label: g.Text("Completed By")},
		{FieldKey: "ReopenedByUsername", Label: g.Text("Reopened By")},
		{FieldKey: "CancelledByUsername", Label: g.Text("Cancelled By")},
	}

	var changelogEntries []components.ChangelogEntry
	for _, change := range p.changeLog {

		entry := components.ChangelogEntry{
			ChangedAt:        change.ChangeAt,
			ChangeByUsername: change.ChangeByUsername,
			IsCreation:       change.IsCreation,
			Changes: map[string]any{
				"Notes":               change.Notes,
				"RaisedByUsername":    change.StartedByUsername,
				"CompletedByUsername": change.CompletedByUsername,
				"ReopenedByUsername":  change.ReopenedByUsername,
				"CancelledByUsername": change.CancelledByUsername,
			},
		}
		changelogEntries = append(changelogEntries, entry)
	}

	return components.Changelog(changelogEntries, changelogFieldDefs)
}
