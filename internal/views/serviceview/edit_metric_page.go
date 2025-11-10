package serviceview

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

type EditServiceMetricPageProps struct {
	Ctx              reqcontext.ReqContext
	Metric           model.ServiceMetric
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func EditServiceMetricPage(p *EditServiceMetricPageProps) g.Node {

	content := g.Group([]g.Node{

		editServiceMetricForm(&editServiceMetricFormProps{
			metric:           p.Metric,
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   fmt.Sprintf("Edit Metric: %s", p.Metric.Name),
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
				Title:   p.Metric.Name,
				URLPart: fmt.Sprintf("%d", p.Metric.ServiceMetricID),
			},
			{
				IconIdentifier: "pencil",
				Title:          "Edit",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/serviceview/add_metric_page.css"),
		},
	})
}

type editServiceMetricFormProps struct {
	metric           model.ServiceMetric
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

func editServiceMetricForm(p *editServiceMetricFormProps) g.Node {

	metric := p.metric

	nameLabel := "Name"
	nameKey := "Name"
	nameValue := metric.Name
	if p.isSubmission {
		if v := p.values.Get(nameKey); v != "" {
			nameValue = v
		} else {
			nameValue = ""
		}
	}
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
	descriptionValue := metric.Description
	if p.isSubmission {
		if v := p.values.Get(descriptionKey); v != "" {
			descriptionValue = v
		} else {
			descriptionValue = ""
		}
	}
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
	isCumulativeValue := metric.IsCumulative
	if p.isSubmission {
		isCumulativeValue = p.values.Get(isCumulativeKey) == "true"
	}

	isArchivedLabel := "Is Archived?"
	isArchivedKey := "IsArchived"
	isArchivedValue := metric.IsArchived
	if p.isSubmission {
		isArchivedValue = p.values.Get(isArchivedKey) == "true"
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
			h.P(
				h.Class("note"),
				g.Textf("Note: Cumulative metrics sum all usage records to calculate current usage. Non-cumulative metrics only use the latest usage record."),
			),
		),

		h.Div(
			h.Label(
				g.Text(isArchivedLabel),

				h.Input(
					h.Type("checkbox"),
					h.Name(isArchivedKey),
					g.If(isArchivedValue, h.Checked()),
					h.Value("true"),
				),
			),
			h.P(
				h.Class("note"),
				g.Textf("Archived metrics are hidden from scheduling and usage flows but can be restored at any time."),
			),
		),

		h.Button(
			h.Class("button primary"),
			h.Type("submit"),
			g.Text("Save Changes"),
		),
	)
}
