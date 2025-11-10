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

type AddResourcePageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	Teams            []model.Team
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func AddResourcePage(p *AddResourcePageProps) g.Node {

	content := g.Group([]g.Node{

		addResourceForm(&addResourceFormProps{
			values:           p.Values,
			teams:            p.Teams,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Add Resource",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "cube-scan",
				Title:          "Resources",
				URLPart:        "resources",
			},
			{
				IconIdentifier: "plus",
				Title:          "Add",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/resourceview/add_resource_page.css"),
		},
	})
}

type addResourceFormProps struct {
	values           url.Values
	teams            []model.Team
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

func addResourceForm(p *addResourceFormProps) g.Node {

	typeLabel := "Type"
	typeKey := "Type"
	typeValue := p.values.Get(typeKey)
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
	referenceValue := p.values.Get(referenceKey)
	referenceError := ""
	if p.isSubmission || referenceValue != "" {
		referenceError = p.validationErrors.GetError(referenceKey, referenceLabel)
	}
	referenceHelperType := components.InputHelperTypeNone
	if referenceError != "" {
		referenceHelperType = components.InputHelperTypeError
	}

	teamLabel := "Service Ownership Team"
	teamKey := "ServiceOwnershipTeamID"
	teamValue := p.values.Get(teamKey)

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

		h.Button(
			h.Class("button primary"),
			h.Type("submit"),
			g.Text("Add Resource"),
		),
	)
}

func serviceTeamOptions(teams []model.Team, selectedValue string) g.Node {
	options := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("Unassigned"),
			g.If(selectedValue == "", h.Selected()),
		),
	}

	for _, team := range teams {
		value := strconv.Itoa(team.TeamID)
		isSelected := selectedValue == value
		options = append(options, h.Option(
			h.Value(value),
			g.Text(team.TeamName),
			g.If(isSelected, h.Selected()),
		))
	}

	return g.Group(options)
}
