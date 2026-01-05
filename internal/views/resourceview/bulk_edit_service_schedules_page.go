package resourceview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"fmt"
	"net/url"
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type BulkEditServiceSchedulesPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ResourceIDs      []int
	ServiceSchedules []model.ServiceSchedule
}

func BulkEditServiceSchedulesPage(p *BulkEditServiceSchedulesPageProps) g.Node {
	selectedCount := len(p.ResourceIDs)

	resourceLabel := "resources"
	if selectedCount == 1 {
		resourceLabel = "resource"
	}

	content := g.Group([]g.Node{
		h.Div(
			h.Class("bulk-edit-header"),
			h.H3(g.Text("Bulk Edit Service Schedules")),
			h.P(
				h.Class("bulk-edit-count"),
				g.Text(fmt.Sprintf("%d %s selected", selectedCount, resourceLabel)),
			),
			g.If(
				selectedCount == 0,
				h.P(
					h.Class("bulk-edit-warning"),
					g.Text("No resources selected. Return to Resources to pick resources for bulk edit."),
				),
			),
		),
		bulkEditServiceSchedulesForm(&bulkEditServiceSchedulesFormProps{
			resourceIDs:      p.ResourceIDs,
			serviceSchedules: p.ServiceSchedules,
			assignSelected:   p.Values["AssignServiceScheduleIDs"],
			unassignSelected: p.Values["UnassignServiceScheduleIDs"],
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Bulk Edit Service Schedules",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "cube-scan",
				Title:          "Resources",
				URLPart:        "resources",
			},
			{
				Title: "Bulk Edit Schedules",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/resourceview/bulk_edit_service_schedules_page.css"),
		},
	})
}

type bulkEditServiceSchedulesFormProps struct {
	resourceIDs      []int
	serviceSchedules []model.ServiceSchedule
	assignSelected   []string
	unassignSelected []string
}

func bulkEditServiceSchedulesForm(p *bulkEditServiceSchedulesFormProps) g.Node {
	disabled := len(p.resourceIDs) == 0

	assignOptions := mapServiceSchedulesToOptions(p.serviceSchedules, p.assignSelected)
	unassignOptions := mapServiceSchedulesToOptions(p.serviceSchedules, p.unassignSelected)

	var resourceInputs []g.Node
	for _, resourceID := range p.resourceIDs {
		resourceInputs = append(resourceInputs, h.Input(
			h.Type("hidden"),
			h.Name("ResourceIDs"),
			h.Value(strconv.Itoa(resourceID)),
		))
	}

	return h.Form(
		h.Method("POST"),
		h.Class("form bulk-edit-form"),
		g.Group(resourceInputs),

		h.Div(
			h.Label(
				g.Text("Assign service schedules"),
			),
			components.SearchSelect(&components.SearchSelectProps{
				Name:        "AssignServiceScheduleIDs",
				Placeholder: "Select schedules to assign",
				Mode:        "multi",
				Options:     assignOptions,
			}),
		),

		h.Div(
			h.Label(
				g.Text("Unassign service schedules"),
			),
			components.SearchSelect(&components.SearchSelectProps{
				Name:        "UnassignServiceScheduleIDs",
				Placeholder: "Select schedules to unassign",
				Mode:        "multi",
				Options:     unassignOptions,
			}),
		),

		components.Button(
			&components.ButtonProps{
				Disabled: disabled,
			},
			h.Type("submit"),
			g.Text("Apply changes"),
		),
	)
}

func mapServiceSchedulesToOptions(
	schedules []model.ServiceSchedule,
	selectedValues []string,
) []components.SearchSelectOption {
	selectedSet := make(map[string]struct{}, len(selectedValues))
	for _, val := range selectedValues {
		selectedSet[val] = struct{}{}
	}

	options := make([]components.SearchSelectOption, 0, len(schedules))
	for _, schedule := range schedules {
		value := strconv.Itoa(schedule.ServiceScheduleID)
		label := schedule.Name
		_, isSelected := selectedSet[value]
		options = append(options, components.SearchSelectOption{
			Text:     label,
			Value:    value,
			Selected: isSelected,
		})
	}

	return options
}
