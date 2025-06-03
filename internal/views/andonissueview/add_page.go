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
	AndonIssues      []model.AndonIssue
	Teams            []model.Team
}

func AddPage(p *AddPageProps) g.Node {

	content := g.Group([]g.Node{

		addForm(&addTeamFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			andonIssues:      p.AndonIssues,
			teams:            p.Teams,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Add New Andon Issue",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			andonIssuesBreadCrumb,
			{IconIdentifier: "plus", Title: "Add"},
		},
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonissueview/add_page.css"),
		},
	})
}

type addTeamFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
	andonIssues      []model.AndonIssue
	teams            []model.Team
}

func addForm(p *addTeamFormProps) g.Node {

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

	parentIDLabel := "Child Issue Of"
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
	for _, andonIssue := range p.andonIssues {
		intVal, _ := strconv.Atoi(parentIDValue)
		isSelected := andonIssue.AndonIssueID == intVal

		parentSelectOptions = append(parentSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", andonIssue.AndonIssueID)),
			g.If(isSelected, h.Selected()),
			g.Text(strings.Join(andonIssue.NamePath, " > ")),
		))
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

	teamSelectOptions := []g.Node{}
	for _, team := range p.teams {
		intVal, _ := strconv.Atoi(assignedToTeamValue)
		isSelected := team.TeamID == intVal

		teamSelectOptions = append(teamSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", team.TeamID)),
			g.If(isSelected, h.Selected()),
			g.Text(team.TeamName),
		))
	}

	resolvableByRaiserLabel := "Resolvable by Raiser?"
	resolvableByRaiserKey := "ResolvableByRaiser"
	resolvableByRaiserValue := false
	if p.values.Get(resolvableByRaiserKey) == "true" {
		resolvableByRaiserValue = true
	}

	willStopProcessLabel := "Will Stop Process?"
	willStopProcessKey := "WillStopProcess"
	willStopProcessValue := true
	if p.values.Get(willStopProcessKey) == "false" {
		willStopProcessValue = false
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
				})),
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
				g.Text(resolvableByRaiserLabel),

				h.Input(
					h.Type("checkbox"),
					h.Name(resolvableByRaiserKey),
					g.If(resolvableByRaiserValue, h.Checked()),
					h.Value("true"),
				),
			),
		),

		h.Div(
			h.Label(
				g.Text(willStopProcessLabel),

				h.Input(
					h.Type("checkbox"),
					h.Name(willStopProcessKey),
					g.If(willStopProcessValue, h.Checked()),
					h.Value("true"),
				),
			),
		),

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Submit"),
		),
	)
}
