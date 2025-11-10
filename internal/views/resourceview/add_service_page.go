package resourceview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/url"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AddResourceServicePageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	Resource         model.Resource
}

func AddResourceServicePage(p *AddResourceServicePageProps) g.Node {

	content := g.Group([]g.Node{

		addServiceForm(&addServiceFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			resource:         p.Resource,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add Resource Service Metric",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "account-wrench",
				Title:          "Services",
				URLPart:        "services",
			},
			{
				IconIdentifier: "plus",
				Title:          "Service",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/resourceview/add_service_page.css"),
		},
	})
}

type addServiceFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
	resource         model.Resource
}

func addServiceForm(p *addServiceFormProps) g.Node {

	notesLabel := "Notes"
	notesKey := "Notes"
	notesValue := p.values.Get(notesKey)
	notesError := ""
	if p.isSubmission || notesValue != "" {
		notesError = p.validationErrors.GetError(notesKey, notesLabel)
	}
	notesHelperType := components.InputHelperTypeNone
	if notesError != "" {
		notesHelperType = components.InputHelperTypeError
	}

	return components.Form(
		h.Method("POST"),

		h.H3(
			g.Text(fmt.Sprintf("Add Service for %s", p.resource.Reference)),
		),

		h.Div(
			h.Label(
				g.Text(notesLabel),

				h.Textarea(
					h.Name(notesKey),
					h.Placeholder("Leave some notes about the service "+
						"(you can change these later)"),
					h.Value(notesValue),
					h.Rows("7"),
				),
			),
			g.If(
				notesError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: notesError,
					Type:  notesHelperType,
				}),
			),
			h.P(
				h.Class("note"),
				g.Textf("* Add images in the next step."),
			),
		),

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Start Service"),
		),
	)
}
