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

type EditGroupPageProps struct {
	Ctx              reqcontext.ReqContext
	AndonIssueGroup  model.AndonIssueGroup
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	AndonIssueGroups []model.AndonIssueGroup
	ErrorText        string
}

const MaxIssueDepth = 4

func EditGroupPage(p *EditGroupPageProps) g.Node {

	andonIssue := p.AndonIssueGroup

	content := g.Group([]g.Node{
		editGroupForm(&editGroupFormProps{
			andonIssueGroup:  p.AndonIssueGroup,
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			andonIssueGroups: p.AndonIssueGroups,
			errorText:        p.ErrorText,
		}),
	})

	return layout.Page(layout.PageProps{
		Title: fmt.Sprintf("Edit: %s", andonIssue.IssueName),
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "alert-octagon-outline",
				Title:          "Andon Issues",
				URLPart:        "andon-issues",
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
			components.InlineStyle("/internal/views/andonissueview/edit_group_page.css"),
		},
	})
}

type editGroupFormProps struct {
	andonIssueGroup  model.AndonIssueGroup
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
	andonIssueGroups []model.AndonIssueGroup
	errorText        string
}

func editGroupForm(p *editGroupFormProps) g.Node {

	andonIssue := p.andonIssueGroup

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

	parentIDLabel := "Child Group Of"
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

	parentSelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}

	spareDepth := MaxIssueDepth - p.andonIssueGroup.DownDepth

	for _, aig := range p.andonIssueGroups {
		if p.andonIssueGroup.AndonIssueID == aig.AndonIssueID {
			continue
		}

		isDisabled := false
		if aig.Depth > spareDepth {
			isDisabled = true
		}

		intVal, _ := strconv.Atoi(parentIDValue)
		isSelected := aig.AndonIssueID == intVal

		parentSelectOptions = append(parentSelectOptions, h.Option(
			h.Value(fmt.Sprintf("%d", aig.AndonIssueID)),
			g.If(isSelected, h.Selected()),
			g.Text(strings.Join(aig.NamePath, " > ")),
			g.If(isDisabled, h.Disabled()),
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
					h.Placeholder("Enter issue group name"),
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

			h.P(
				h.Class("note"),
				g.Text("* Only two levels of groups are supported"),
			),
		),

		g.If(
			p.errorText != "",
			h.Div(
				h.Class("error-msg"),
				g.Text(p.errorText),
			),
		),

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Save"),
		),
	)
}
