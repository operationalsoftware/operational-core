package resourceview

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

type AddResourceMetricRecordPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	Resource         model.Resource
	ServiceMetrics   []model.ServiceMetric
}

func AddResourceMetricRecordPage(p *AddResourceMetricRecordPageProps) g.Node {

	content := g.Group([]g.Node{

		addResourceMetricRecordForm(&addResourceMetricRecordFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			serviceMetrics:   p.ServiceMetrics,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add Resource Recording",
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
				Title:          "Resource Recording",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/resourceview/add_metric_record_page.css"),
		},
	})
}

type addResourceMetricRecordFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
	serviceMetrics   []model.ServiceMetric
}

func addResourceMetricRecordForm(p *addResourceMetricRecordFormProps) g.Node {

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

	valueLabel := "Recording Value"
	valuePlaceholder := "Enter value"
	valueKey := "Value"
	valueValue := p.values.Get(valueKey)
	valueError := ""
	if p.isSubmission || valueValue != "" {
		valueError = p.validationErrors.GetError(valueKey, valueLabel)
	}
	valueHelperType := components.InputHelperTypeNone
	if valueError != "" {
		valueHelperType = components.InputHelperTypeError
	}

	metricSelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}
	for _, metric := range p.serviceMetrics {
		intVal, _ := strconv.Atoi(serviceMetricValue)
		isSelected := metric.ServiceMetricID == intVal

		metricSelectOptions = append(metricSelectOptions, h.Option(
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
					h.ID("service-metric-select"),
					h.Name(serviceMetricKey),
					g.Group(metricSelectOptions),
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
				g.Text(valueLabel),

				h.Input(
					h.Name(valueKey),
					h.Type("number"),
					h.Placeholder(valuePlaceholder),
					h.Value(valueValue),
					h.AutoComplete("off"),
				),
			),
			g.If(
				valueError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: valueError,
					Type:  valueHelperType,
				}),
			),
		),

		h.Button(
			h.Class("button primary"),
			h.Type("submit"),
			g.Text("Add Recording"),
		),
	)
}
