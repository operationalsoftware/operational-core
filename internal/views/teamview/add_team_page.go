package teamview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"net/url"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type AddTeamPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func AddTeamPage(p *AddTeamPageProps) g.Node {

	content := g.Group([]g.Node{

		addTeamForm(&addTeamFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Add New Team",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			teamsBreadCrumb,
			{IconIdentifier: "account-plus", Title: "Add"},
		},
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/teamview/add_team_page.css"),
		},
	})
}

type addTeamFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

func addTeamForm(p *addTeamFormProps) g.Node {

	teamNameLabel := "Team Name"
	teamNameKey := "TeamName"
	teamNameValue := p.values.Get(teamNameKey)
	teamNameError := ""
	if p.isSubmission || teamNameValue != "" {
		teamNameError = p.validationErrors.GetError(teamNameKey, teamNameLabel)
	}
	teamNameHelperType := components.InputHelperTypeNone
	if teamNameError != "" {
		teamNameHelperType = components.InputHelperTypeError
	}

	return components.Form(
		h.ID("add-team-form"),
		h.Method("POST"),

		components.Input(&components.InputProps{
			Label:       teamNameLabel,
			Name:        teamNameKey,
			Placeholder: "Enter team name",
			HelperText:  teamNameError,
			HelperType:  teamNameHelperType,
			InputProps: []g.Node{
				h.Value(teamNameValue),
				h.AutoComplete("off"),
			},
		}),

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Submit"),
		),
	)
}
