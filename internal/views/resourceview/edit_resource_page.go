package resourceview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"net/url"
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type EditResourcePageProps struct {
	Ctx              reqcontext.ReqContext
	Resource         model.Resource
	Teams            []model.Team
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func EditResourcePage(p *EditResourcePageProps) g.Node {

	content := g.Group([]g.Node{
		editResourceForm(&editResourceFormProps{
			resource:         p.Resource,
			teams:            p.Teams,
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Edit Resource",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "cube-scan",
				Title:          "Resources",
				URLPart:        "resources",
			},
			{
				Title: "Edit",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/resourceview/add_resource_page.css"),
		},
	})
}

type editResourceFormProps struct {
	resource         model.Resource
	teams            []model.Team
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

func editResourceForm(p *editResourceFormProps) g.Node {
	typeLabel := "Type"
	typeKey := "Type"
	typeValue := p.resource.Type
	if v := p.values.Get(typeKey); p.isSubmission || v != "" {
		typeValue = v
	}
	typeError := ""
	if p.isSubmission || typeValue != "" {
		typeError = p.validationErrors.GetError(typeKey, typeLabel)
	}
	typeHelperType := components.InputHelperTypeNone
	if typeError != "" {
		typeHelperType = components.InputHelperTypeError
	}

	referenceLabel := "Reference"
	referenceKey := "Reference"
	referenceValue := p.resource.Reference
	if v := p.values.Get(referenceKey); p.isSubmission || v != "" {
		referenceValue = v
	}
	referenceError := ""
	if p.isSubmission || referenceValue != "" {
		referenceError = p.validationErrors.GetError(referenceKey, referenceLabel)
	}
	referenceHelperType := components.InputHelperTypeNone
	if referenceError != "" {
		referenceHelperType = components.InputHelperTypeError
	}

	archiveLabel := "Is Archived?"
	archiveKey := "IsArchived"
	archiveError := ""
	if p.isSubmission {
		archiveError = p.validationErrors.GetError(archiveKey, archiveLabel)
	}
	archiveHelperType := components.InputHelperTypeNone
	if archiveError != "" {
		archiveHelperType = components.InputHelperTypeError
	}

	checkedArchived := p.resource.IsArchived
	if p.isSubmission {
		checkedArchived = p.values.Get(archiveKey) == "true"
	}

	teamLabel := "Service Ownership Team"
	teamKey := "ServiceOwnershipTeamID"
	teamValue := ""
	if p.resource.ServiceOwnershipTeamID != nil {
		teamValue = strconv.Itoa(*p.resource.ServiceOwnershipTeamID)
	}
	if v := p.values.Get(teamKey); p.isSubmission || v != "" {
		teamValue = v
	}

	return h.Form(
		h.Method("POST"),
		h.Class("form"),

		h.Div(
			h.Label(
				g.Text(typeLabel),

				h.Input(
					h.Name(typeKey),
					h.Placeholder("Enter type"),
					h.Value(typeValue),
					h.AutoComplete("off"),
				),
			),
			g.If(
				typeError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: typeError,
					Type:  typeHelperType,
				}),
			),
		),

		h.Div(
			h.Label(
				g.Text(referenceLabel),

				h.Input(
					h.Name(referenceKey),
					h.Placeholder("Enter reference"),
					h.Value(referenceValue),
					h.AutoComplete("off"),
				),
			),
			g.If(
				referenceError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: referenceError,
					Type:  referenceHelperType,
				}),
			),
		),

		h.Div(
			h.Label(
				g.Text(teamLabel),
				h.Select(
					h.Name(teamKey),
					h.Class("select"),
					serviceTeamOptions(p.teams, teamValue),
				),
			),
		),

		h.Div(
			components.Checkbox(
				&components.CheckboxProps{
					Name:    archiveKey,
					Label:   archiveLabel,
					Checked: checkedArchived,
					Value:   "true",
				},
			),
			g.If(
				archiveError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: archiveError,
					Type:  archiveHelperType,
				}),
			),
		),

		h.Button(
			h.Class("button primary"),
			h.Type("submit"),
			g.Text("Save Resource"),
		),
	)
}
