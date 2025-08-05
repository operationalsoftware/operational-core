package andonissueview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type EditPageProps struct {
	Ctx              reqcontext.ReqContext
	AndonIssue       model.AndonIssue
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	AndonIssues      []model.AndonIssueNode
	AndonIssueGroups []model.AndonIssueGroup
	Teams            []model.Team
}

func EditPage(p *EditPageProps) g.Node {

	andonIssue := p.AndonIssue

	content := g.Group([]g.Node{
		editForm(&editFormProps{
			andonIssue:       p.AndonIssue,
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			andonIssues:      p.AndonIssues,
			andonIssueGroups: p.AndonIssueGroups,
			teams:            p.Teams,
		}),
	})

	return layout.Page(layout.PageProps{
		Title: fmt.Sprintf("Edit: %s", andonIssue.IssueName),
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "alert-octagon-outline",
				Title:          "Andons",
				URLPart:        "andons",
			},
			{
				Title:   "Andon Issues",
				URLPart: "andon-issues",
			},
			{
				Title:   andonIssue.IssueName,
				URLPart: fmt.Sprintf("%d", andonIssue.AndonIssueID),
			},
			{
				IconIdentifier: "pencil",
				Title:          "Edit",
			},
		},
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonissueview/edit_page.css"),
		},
	})
}

type editFormProps struct {
	andonIssue       model.AndonIssue
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
	andonIssues      []model.AndonIssueNode
	andonIssueGroups []model.AndonIssueGroup
	teams            []model.Team
}

// same as addUserForm, but no password fields
func editForm(p *editFormProps) g.Node {

	andonIssue := p.andonIssue

	issueNameLabel := "Issue Name"
	issueNameKey := "IssueName"
	var issueNameValue string
	if p.values.Get(issueNameKey) != "" {
		issueNameValue = p.values.Get(issueNameKey)
	} else {
		issueNameValue = andonIssue.IssueName
	}
	issueNameError := ""
	if p.isSubmission || issueNameValue != "" {
		issueNameError = p.validationErrors.GetError(issueNameKey, issueNameLabel)
	}
	issueNameHelperType := components.InputHelperTypeNone
	if issueNameError != "" {
		issueNameHelperType = components.InputHelperTypeError
	}

	parentIDLabel := "Child Issue Of"
	parentIDKey := "ParentID"
	parentIDValue := ""
	if p.values.Get(parentIDKey) != "" {
		parentIDValue = p.values.Get(parentIDKey)
	} else if andonIssue.ParentID != nil {
		parentIDValue = fmt.Sprintf("%d", *andonIssue.ParentID)
	}
	parentIDError := ""
	if p.isSubmission || parentIDValue != "" {
		parentIDError = p.validationErrors.GetError(parentIDKey, parentIDLabel)
	}
	parentIDHelperType := components.InputHelperTypeNone
	if parentIDError != "" {
		parentIDHelperType = components.InputHelperTypeError
	}

	assignedToTeamLabel := "Assigned to Team"
	assignedToTeamKey := "AssignedToTeam"
	assignedToTeamValue := p.values.Get(assignedToTeamKey)
	assignedToTeamError := ""
	if p.isSubmission || assignedToTeamValue != "" {
		assignedToTeamError = p.validationErrors.GetError(assignedToTeamKey, assignedToTeamLabel)
	}
	assignedToTeamHelperType := components.InputHelperTypeNone
	if assignedToTeamError != "" {
		assignedToTeamHelperType = components.InputHelperTypeError
	}

	parentSelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}
	for _, andonIssue := range p.andonIssueGroups {
		if p.andonIssue.AndonIssueID == andonIssue.AndonIssueID {
			continue
		}

		intVal, _ := strconv.Atoi(parentIDValue)
		isSelected := andonIssue.AndonIssueID == intVal

		parentSelectOptions = append(parentSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", andonIssue.AndonIssueID)),
			g.If(isSelected, h.Selected()),
			g.Text(strings.Join(andonIssue.NamePath, " > ")),
		))
	}

	teamSelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}
	for _, team := range p.teams {
		intVal := p.andonIssue.AssignedToTeam
		isSelected := team.TeamID == nilsafe.Int(intVal)

		teamSelectOptions = append(teamSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", team.TeamID)),
			g.If(isSelected, h.Selected()),
			g.Text(team.TeamName),
		))
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

	isArchivedLabel := "Is Archived?"
	isArchivedKey := "IsArchived"
	isArchivedValue := false
	if p.values.Get(isArchivedKey) == "true" {
		isArchivedValue = true
	} else if andonIssue.IsArchived {
		isArchivedValue = true
	}
	isArchivedError := ""
	if p.isSubmission {
		isArchivedError = p.validationErrors.GetError(isArchivedKey, isArchivedLabel)
	}
	isArchivedHelperType := components.InputHelperTypeNone
	if isArchivedError != "" {
		isArchivedHelperType = components.InputHelperTypeError
	}

	severitySelectOptions := []g.Node{}
	for _, severity := range model.AndonSeverities {
		isSelected := string(severity) == severityValue

		severitySelectOptions = append(severitySelectOptions, h.Option(
			h.Value(string(severity)),
			g.If(isSelected, h.Selected()),
			g.Text(string(severity)),
		))
	}

	return components.Form(
		h.ID("edit-andon-issue-form"),
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
				}),
			),
		),

		h.Div(
			h.Label(
				g.Text(parentIDLabel),

				h.Select(
					h.Name(parentIDKey),
					g.Group(parentSelectOptions),
				),
			),
			g.If(parentIDError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: parentIDError,
					Type:  parentIDHelperType,
				}),
			),
		),

		h.Div(
			h.Label(
				g.Text(assignedToTeamLabel),

				h.Select(
					h.Name(assignedToTeamKey),
					g.Group(teamSelectOptions),
				),
			),
			g.If(assignedToTeamError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: assignedToTeamError,
					Type:  assignedToTeamHelperType,
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
			g.If(isArchivedError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: isArchivedError,
					Type:  isArchivedHelperType,
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
