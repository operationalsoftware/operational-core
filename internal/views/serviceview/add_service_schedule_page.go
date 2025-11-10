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
	h "maragu.dev/gomponents/html"
)

type AddResourceServiceSchedulePageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	Resource         model.Resource
	ServiceMetrics   []model.ServiceMetric
}

func AddServiceSchedulePage(p *AddResourceServiceSchedulePageProps) g.Node {

	content := g.Group([]g.Node{

		addServiceScheduleForm(&addServiceScheduleFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			serviceMetrics:   p.ServiceMetrics,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add Resource Service Schedule",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "cube-scan",
				Title:          "Resources",
				URLPart:        "resources",
			},
			{
				Title:   p.Resource.Reference,
				URLPart: fmt.Sprintf("%d", p.Resource.ResourceID),
			},
			{
				IconIdentifier: "plus",
				Title:          "Service Schedule",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/serviceview/add_service_schedule_page.css"),
		},
	})
}

type addServiceScheduleFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
	serviceMetrics   []model.ServiceMetric
}

func addServiceScheduleForm(p *addServiceScheduleFormProps) g.Node {

	serviceMetricLabel := "Select service metric"
	serviceMetricKey := "ServiceMetricID"
	serviceMetricValue := p.values.Get(serviceMetricKey)
	serviceMetricError := ""
	if p.isSubmission || serviceMetricValue != "" {
		serviceMetricError = p.validationErrors.GetError(serviceMetricKey, serviceMetricLabel)
	}
	serviceMetricHelperType := components.InputHelperTypeNone
	if serviceMetricError != "" {
		serviceMetricHelperType = components.InputHelperTypeError
	}

	thresholdLabel := "Threshold"
	thresholdKey := "Threshold"
	thresholdValue := p.values.Get(thresholdKey)
	thresholdError := ""
	if p.isSubmission || thresholdValue != "" {
		thresholdError = p.validationErrors.GetError(thresholdKey, thresholdLabel)
	}
	thresholdHelperType := components.InputHelperTypeNone
	if thresholdError != "" {
		thresholdHelperType = components.InputHelperTypeError
	}

	teamSelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}
	for _, metric := range p.serviceMetrics {
		intVal, _ := strconv.Atoi(serviceMetricValue)
		isSelected := metric.ServiceMetricID == intVal

		teamSelectOptions = append(teamSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", metric.ServiceMetricID)),
			g.If(isSelected, h.Selected()),
			g.Text(metric.Name),
		))
	}

	return h.Form(
		h.Method("POST"),
		h.Class("form"),

		h.Div(
			h.Label(
				g.Text(serviceMetricLabel),

				h.Select(
					h.Name(serviceMetricKey),
					g.Group(teamSelectOptions),
				),
			),
			g.If(serviceMetricError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: serviceMetricError,
					Type:  serviceMetricHelperType,
				})),
		),

		h.Div(
			h.Label(
				g.Text(thresholdLabel),

				h.Input(
					h.Name(thresholdKey),
					h.Type("number"),
					h.Placeholder("Enter threshold"),
					h.Value(thresholdValue),
					h.AutoComplete("off"),
				),
			),
			g.If(
				thresholdError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: thresholdError,
					Type:  thresholdHelperType,
				}),
			),
		),

		h.Button(
			h.Class("button primary"),
			h.Type("submit"),
			g.Text("Submit"),
		),
	)
}
