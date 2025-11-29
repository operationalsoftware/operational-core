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

type AddSchedulePageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	ServiceMetrics   []model.ServiceMetric
}

func AddSchedulePage(p *AddSchedulePageProps) g.Node {

	content := g.Group([]g.Node{

		scheduleForm(&scheduleFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			serviceMetrics:   p.ServiceMetrics,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add Service Schedule",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "account-wrench",
				Title:          "Services",
				URLPart:        "services",
			},
			{
				IconIdentifier: "clock",
				Title:          "Schedules",
				URLPart:        "schedules",
			},
			{
				IconIdentifier: "plus",
				Title:          "New",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/serviceview/add_schedule_page.css"),
		},
	})
}

type scheduleFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
	serviceMetrics   []model.ServiceMetric
}

func scheduleForm(p *scheduleFormProps) g.Node {

	nameLabel := "Schedule Name"
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

	serviceMetricLabel := "Service Metric"
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

	metricOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}

	selectedMetricID, _ := strconv.Atoi(serviceMetricValue)
	for _, metric := range p.serviceMetrics {
		metricOptions = append(metricOptions,
			h.Option(
				h.Value(fmt.Sprintf("%d", metric.ServiceMetricID)),
				g.If(metric.ServiceMetricID == selectedMetricID, h.Selected()),
				g.Text(metric.Name),
			),
		)
	}

	return h.Form(
		h.Method("POST"),
		h.Class("form"),

		h.Div(
			h.Label(
				g.Text(nameLabel),

				h.Input(
					h.Name(nameKey),
					h.Placeholder("Enter schedule name"),
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
				g.Text(serviceMetricLabel),

				h.Select(
					h.Name(serviceMetricKey),
					g.Group(metricOptions),
				),
			),
			g.If(
				serviceMetricError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: serviceMetricError,
					Type:  serviceMetricHelperType,
				}),
			),
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
					g.Attr("step", "any"),
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
