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
	ServiceSchedules []model.ServiceSchedule
}

func AddServiceSchedulePage(p *AddResourceServiceSchedulePageProps) g.Node {

	content := g.Group([]g.Node{

		addServiceScheduleForm(&addServiceScheduleFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			serviceSchedules: p.ServiceSchedules,
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
	serviceSchedules []model.ServiceSchedule
}

func addServiceScheduleForm(p *addServiceScheduleFormProps) g.Node {

	serviceScheduleLabel := "Select service schedule"
	serviceScheduleKey := "ServiceScheduleID"
	serviceScheduleValue := p.values.Get(serviceScheduleKey)
	serviceScheduleError := ""
	if p.isSubmission || serviceScheduleValue != "" {
		serviceScheduleError = p.validationErrors.GetError(serviceScheduleKey, serviceScheduleLabel)
	}
	serviceScheduleHelperType := components.InputHelperTypeNone
	if serviceScheduleError != "" {
		serviceScheduleHelperType = components.InputHelperTypeError
	}

	teamSelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}
	for _, schedule := range p.serviceSchedules {
		intVal, _ := strconv.Atoi(serviceScheduleValue)
		isSelected := schedule.ServiceScheduleID == intVal

		teamSelectOptions = append(teamSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", schedule.ServiceScheduleID)),
			g.If(isSelected, h.Selected()),
			g.Text(schedule.Name),
		))
	}

	return h.Form(
		h.Method("POST"),
		h.Class("form"),

		h.Div(
			h.Label(
				g.Text(serviceScheduleLabel),

				h.Select(
					h.Name(serviceScheduleKey),
					g.Group(teamSelectOptions),
				),
			),
			g.If(serviceScheduleError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: serviceScheduleError,
					Type:  serviceScheduleHelperType,
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
