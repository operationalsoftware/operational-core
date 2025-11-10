package serviceview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"net/url"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AddServiceMetricPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func AddServiceMetricPage(p *AddServiceMetricPageProps) g.Node {

	content := g.Group([]g.Node{

		addResourceMetricForm(&addResourceMetricFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add Service Metric",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "account-wrench",
				Title:          "Services",
				URLPart:        "services",
			},
			{
				IconIdentifier: "speedometer",
				Title:          "Metrics",
				URLPart:        "metrics",
			},
			{
				IconIdentifier: "plus",
				Title:          "Metric",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/serviceview/add_metric_page.css"),
		},
	})
}

type addResourceMetricFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

func addResourceMetricForm(p *addResourceMetricFormProps) g.Node {

	nameLabel := "Name"
	nameKey := "Name"
	nameValue := p.values.Get(nameKey)
	nameError := ""
	if p.isSubmission || nameValue != "" {
		nameError = p.validationErrors.GetError(nameKey, nameLabel)
	}
	nameHelperType := components.InputHelperTypeNone
	if nameError != "" {
		nameHelperType = components.InputHelperTypeError
	}

	descriptionLabel := "Description"
	descriptionKey := "Description"
	descriptionValue := p.values.Get(descriptionKey)
	descriptionError := ""
	if p.isSubmission || descriptionValue != "" {
		descriptionError = p.validationErrors.GetError(descriptionKey, descriptionLabel)
	}
	descriptionHelperType := components.InputHelperTypeNone
	if descriptionError != "" {
		descriptionHelperType = components.InputHelperTypeError
	}

	isCumulativeLabel := "Is Cumulative?"
	isCumulativeKey := "IsCumulative"
	isCumulativeValue := false
	if p.values.Get(isCumulativeKey) == "true" {
		isCumulativeValue = true
	}
	isCumulativeError := ""
	if p.isSubmission {
		isCumulativeError = p.validationErrors.GetError(isCumulativeKey, isCumulativeLabel)
	}
	isCumulativeHelperType := components.InputHelperTypeNone
	if isCumulativeError != "" {
		isCumulativeHelperType = components.InputHelperTypeError
	}

	return h.Form(
		h.Method("POST"),
		h.Class("form"),

		h.Div(
			h.Label(
				g.Text(nameLabel),

				h.Input(
					h.Name(nameKey),
					h.Placeholder("Enter name"),
					h.Value(nameValue),
					h.AutoComplete("off"),
				),
			),
			g.If(
				nameError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: nameError,
					Type:  nameHelperType,
				}),
			),
		),

		h.Div(
			h.Label(
				g.Text(descriptionLabel),

				h.Input(
					h.Name(descriptionKey),
					h.Placeholder("Enter description"),
					h.Value(descriptionValue),
					h.AutoComplete("off"),
				),
			),
			g.If(
				descriptionError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: descriptionError,
					Type:  descriptionHelperType,
				}),
			),
		),

		h.Div(
			h.Label(
				g.Text(isCumulativeLabel),

				h.Input(
					h.Type("checkbox"),
					h.Name(isCumulativeKey),
					g.If(isCumulativeValue, h.Checked()),
					h.Value("true"),
				),
			),
			g.If(isCumulativeError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: isCumulativeError,
					Type:  isCumulativeHelperType,
				}),
			),
			h.P(
				h.Class("note"),
				g.Textf("Note: Cumulative metrics sum all usage records to calculate current usage. Non-cumulative metrics only use the latest usage record."),
			),
		),

		h.Button(
			h.Class("button primary"),
			h.Type("submit"),
			g.Text("Add Metric"),
		),
	)
}
