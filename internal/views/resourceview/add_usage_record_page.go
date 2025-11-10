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

type AddResourceUsagePageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	Resource         model.Resource
	ServiceMetrics   []model.ServiceMetric
}

func AddResourceUsagePage(p *AddResourceUsagePageProps) g.Node {

	content := g.Group([]g.Node{

		addResourceUsageRecordForm(&addResourceUsageRecordFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			serviceMetrics:   p.ServiceMetrics,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add Resource Usage Record",
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
				Title:          "Usage Record",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/resourceview/add_usage_record_page.css"),
		},
	})
}

type addResourceUsageRecordFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
	serviceMetrics   []model.ServiceMetric
}

func addResourceUsageRecordForm(p *addResourceUsageRecordFormProps) g.Node {

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

	usageLabel := "Value"
	usagePlaceholder := "Enter value"
	usageKey := "Value"
	usageValue := p.values.Get(usageKey)
	usageError := ""
	if p.isSubmission || usageValue != "" {
		usageError = p.validationErrors.GetError(usageKey, usageLabel)
	}
	usageHelperType := components.InputHelperTypeNone
	if usageError != "" {
		usageHelperType = components.InputHelperTypeError
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

		if isSelected {
			if metric.IsCumulative {
				usageLabel = "Usage Since Last Service"
			} else {
				usageLabel = "Current Reading"
			}
		}

		metricSelectOptions = append(metricSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", metric.ServiceMetricID)),
			h.Data("is-cumulative", fmt.Sprintf("%t", metric.IsCumulative)),
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
					g.Attr("onchange", "handleMetricChange(event)"),
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
				g.Text(usageLabel),

				h.Input(
					h.Name(usageKey),
					h.Type("number"),
					h.Placeholder(usagePlaceholder),
					h.Value(usageValue),
					h.AutoComplete("off"),
				),
			),
			g.If(
				usageError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: usageError,
					Type:  usageHelperType,
				}),
			),
		),

		h.Button(
			h.Class("button primary"),
			h.Type("submit"),
			g.Text("Add Usage Record"),
		),
	)
}
