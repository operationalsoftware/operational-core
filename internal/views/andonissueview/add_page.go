package andonissueview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type AddPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	AndonIssueGroups []model.AndonIssueGroup
	Teams            []model.Team
}

func AddPage(p *AddPageProps) g.Node {

	content := g.Group([]g.Node{

		addIssueForm(&addIssueFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			andonIssueGroups: p.AndonIssueGroups,
			teams:            p.Teams,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Add New Andon Issue",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "alert-octagon-outline",
				Title:          "Andon Issues",
				URLPart:        "andon-issues",
			},
			{IconIdentifier: "plus", Title: "Add"},
		},
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonissueview/add_page.css"),
		},
	})
}

type addIssueFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
	andonIssueGroups []model.AndonIssueGroup
	teams            []model.Team
}

func addIssueForm(p *addIssueFormProps) g.Node {

	issueNameLabel := "Issue Name"
	issueNameKey := "IssueName"
	issueNameValue := p.values.Get(issueNameKey)
	issueNameError := ""
	if p.isSubmission || issueNameValue != "" {
		issueNameError = p.validationErrors.GetError(issueNameKey, issueNameLabel)
	}
	issueNameHelperType := components.InputHelperTypeNone
	if issueNameError != "" {
		issueNameHelperType = components.InputHelperTypeError
	}

	parentIDLabel := "Issue Group"
	parentIDKey := "ParentID"
	parentIDValue := p.values.Get(parentIDKey)
	parentIDError := ""
	if p.isSubmission || parentIDValue != "" {
		parentIDError = p.validationErrors.GetError(parentIDKey, parentIDLabel)
	}
	parentIDHelperType := components.InputHelperTypeNone
	if parentIDError != "" {
		parentIDHelperType = components.InputHelperTypeError
	}

	// map andon issues to options for parent select
	parentSelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}
	for _, andonIssue := range p.andonIssueGroups {
		intVal, _ := strconv.Atoi(parentIDValue)
		isSelected := andonIssue.AndonIssueID == intVal

		parentSelectOptions = append(parentSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", andonIssue.AndonIssueID)),
			g.If(isSelected, h.Selected()),
			g.Text(strings.Join(andonIssue.NamePath, " > ")),
		))
	}

	assignedTeamLabel := "Assigned to Team"
	assignedTeamKey := "AssignedTeam"
	assignedTeamValue := p.values.Get(assignedTeamKey)
	assignedTeamError := ""
	if p.isSubmission || assignedTeamValue != "" {
		assignedTeamError = p.validationErrors.GetError(assignedTeamKey, assignedTeamLabel)
	}
	assignedTeamHelperType := components.InputHelperTypeNone
	if assignedTeamError != "" {
		assignedTeamHelperType = components.InputHelperTypeError
	}

	severityLabel := "Severity"
	severityKey := "Severity"
	severityValue := p.values.Get(severityKey)
	severityError := ""
	if p.isSubmission || severityValue != "" {
		severityError = p.validationErrors.GetError(severityKey, severityLabel)
	}
	severityHelperType := components.InputHelperTypeNone
	if severityError != "" {
		severityHelperType = components.InputHelperTypeError
	}

	teamSelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}
	for _, team := range p.teams {
		intVal, _ := strconv.Atoi(assignedTeamValue)
		isSelected := team.TeamID == intVal

		teamSelectOptions = append(teamSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", team.TeamID)),
			g.If(isSelected, h.Selected()),
			g.Text(team.TeamName),
		))
	}

	severitySelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}
	for _, severity := range model.AndonSeverities {
		isSelected := string(severity) == severityValue

		severitySelectOptions = append(severitySelectOptions, h.Option(
			h.Value(string(severity)),
			g.If(isSelected, h.Selected()),
			g.Text(string(severity)),
		))
	}

	return components.Form(
		h.ID("add-andon-issue-form"),
		h.Method("POST"),

		h.Div(
			h.Label(
				g.Text(issueNameLabel),

				h.Input(
					h.Name(issueNameKey),
					h.Placeholder("Enter issue name"),
					g.If(issueNameValue != "", h.Value(issueNameValue)),
					h.AutoComplete("off"),
				),
			),
			g.If(
				issueNameError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: issueNameError,
					Type:  issueNameHelperType,
				})),
		),

		h.Div(
			h.Label(
				g.Text("Child Issue Of"),

				h.Select(
					h.Name(parentIDKey),
					g.Group(parentSelectOptions),
				),
			),
			g.If(parentIDError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: parentIDError,
					Type:  parentIDHelperType,
				})),
		),

		h.Div(
			h.Label(
				g.Text(assignedTeamLabel),

				h.Select(
					h.Name(assignedTeamKey),
					g.Group(teamSelectOptions),
				),
			),
			g.If(assignedTeamError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: assignedTeamError,
					Type:  assignedTeamHelperType,
				})),
		),

		h.Div(
			h.Label(
				g.Text(severityLabel),

				h.Select(
					h.Name(severityKey),
					g.Group(severitySelectOptions),
				),
			),
			g.If(
				severityError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: severityError,
					Type:  severityHelperType,
				}),
			),
		),

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Submit"),
		),
	)
}
