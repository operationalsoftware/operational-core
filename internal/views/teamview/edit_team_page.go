package teamview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/url"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type EditTeamPageProps struct {
	Ctx              reqcontext.ReqContext
	Team             model.Team
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func EditTeamPage(p *EditTeamPageProps) g.Node {

	team := p.Team

	content := g.Group([]g.Node{
		editTeamForm(&editTeamFormProps{
			team:             p.Team,
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Title: fmt.Sprintf("Edit: %s", team.TeamName),
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			teamsBreadCrumb,
			{
				IconIdentifier: "account",
				Title:          team.TeamName,
				URLPart:        fmt.Sprintf("%d", team.TeamID),
			},
			{Title: "Edit"},
		},
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/teamview/edit_team_page.css"),
		},
	})
}

type editTeamFormProps struct {
	team             model.Team
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

// same as addUserForm, but no password fields
func editTeamForm(p *editTeamFormProps) g.Node {

	team := p.team

	teamNameLabel := "Team Name"
	teamNameKey := "TeamName"
	var teamNameValue string
	if p.values.Get(teamNameKey) != "" {
		teamNameValue = p.values.Get(teamNameKey)
	} else {
		teamNameValue = team.TeamName
	}
	teamNameError := ""
	if p.isSubmission || teamNameValue != "" {
		teamNameError = p.validationErrors.GetError(teamNameKey, teamNameLabel)
	}
	teamNameHelperType := components.InputHelperTypeNone
	if teamNameError != "" {
		teamNameHelperType = components.InputHelperTypeError
	}

	return components.Form(
		h.ID("edit-team-form"),
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

		components.Checkbox(
			&components.CheckboxProps{
				Name:    "IsArchived",
				Label:   "Is archived?",
				Checked: team.IsArchived,
				Value:   "true",
			},
		),

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Submit"),
		),
	)
}
