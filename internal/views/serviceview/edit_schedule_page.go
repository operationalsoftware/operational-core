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

type EditSchedulePageProps struct {
	Ctx              reqcontext.ReqContext
	Schedule         model.ServiceSchedule
	ServiceMetrics   []model.ServiceMetric
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func EditSchedulePage(p *EditSchedulePageProps) g.Node {

	content := g.Group([]g.Node{

		editScheduleForm(&editScheduleFormProps{
			schedule:         p.Schedule,
			serviceMetrics:   p.ServiceMetrics,
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   fmt.Sprintf("Edit Schedule: %s", p.Schedule.Name),
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
				Title: p.Schedule.Name,
			},
			{
				IconIdentifier: "pencil",
				Title:          "Edit",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/serviceview/edit_schedule_page.css"),
		},
	})
}

type editScheduleFormProps struct {
	schedule         model.ServiceSchedule
	serviceMetrics   []model.ServiceMetric
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

func editScheduleForm(p *editScheduleFormProps) g.Node {

	schedule := p.schedule

	nameLabel := "Schedule Name"
	nameKey := "Name"
	nameValue := schedule.Name
	if p.isSubmission {
		nameValue = p.values.Get(nameKey)
	}
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
	serviceMetricValue := fmt.Sprintf("%d", schedule.ResourceServiceMetricID)
	if p.isSubmission {
		serviceMetricValue = p.values.Get(serviceMetricKey)
	}
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
	thresholdValue := schedule.Threshold.String()
	if p.isSubmission {
		thresholdValue = p.values.Get(thresholdKey)
	}
	thresholdError := ""
	if p.isSubmission || thresholdValue != "" {
		thresholdError = p.validationErrors.GetError(thresholdKey, thresholdLabel)
	}
	thresholdHelperType := components.InputHelperTypeNone
	if thresholdError != "" {
		thresholdHelperType = components.InputHelperTypeError
	}

	isArchivedLabel := "Is Archived?"
	isArchivedKey := "IsArchived"
	isArchivedValue := schedule.IsArchived
	if p.isSubmission {
		isArchivedValue = p.values.Get(isArchivedKey) == "true"
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
					h.Step("any"),
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
		),

		h.Button(
			h.Class("button primary"),
			h.Type("submit"),
			g.Text("Submit"),
		),
	)
}
