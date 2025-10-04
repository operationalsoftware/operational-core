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

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type EditPageProps struct {
	Ctx              reqcontext.ReqContext
	AndonIssue       model.AndonIssue
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
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
			andonIssueGroups: p.AndonIssueGroups,
			teams:            p.Teams,
		}),
	})

	return layout.Page(layout.PageProps{
		Title: fmt.Sprintf("Edit: %s", andonIssue.IssueName),
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
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
		parentIDValue = fmt.Sprintf("%d", andonIssue.ParentID)
	}
	parentIDError := ""
	if p.isSubmission || parentIDValue != "" {
		parentIDError = p.validationErrors.GetError(parentIDKey, parentIDLabel)
	}
	parentIDHelperType := components.InputHelperTypeNone
	if parentIDError != "" {
		parentIDHelperType = components.InputHelperTypeError
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
		intVal := p.andonIssue.AssignedTeam
		isSelected := team.TeamID == intVal

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

	severitySelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		)}
	for _, severity := range model.AndonSeverities {
		isSelected := string(severity) == string(andonIssue.Severity)

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

		h.Button(
			h.Class("button primary"),
			h.Type("submit"),
			g.Text("Save"),
		),
	)
}
