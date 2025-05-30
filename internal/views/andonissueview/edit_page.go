package andonissueview

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

type EditPageProps struct {
	Ctx              reqcontext.ReqContext
	AndonIssue       model.AndonIssue
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func EditPage(p *EditPageProps) g.Node {

	andonIssue := p.AndonIssue

	content := g.Group([]g.Node{
		editForm(&editFormProps{
			andonIssue:       p.AndonIssue,
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Title: fmt.Sprintf("Edit: %s", andonIssue.IssueName),
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			andonIssuesBreadCrumb,
			{
				IconIdentifier: "account",
				Title:          andonIssue.IssueName,
				URLPart:        fmt.Sprintf("%d", andonIssue.AndonIssueID),
			},
			{Title: "Edit"},
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
}

// same as addUserForm, but no password fields
func editForm(p *editFormProps) g.Node {

	andonIssue := p.andonIssue

	issueNameLabel := "Issue Name"
	issueNameKey := "TeamName"
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

	return components.Form(
		h.ID("edit-andon-issue-form"),
		h.Method("POST"),

		components.Input(&components.InputProps{
			Label:       issueNameLabel,
			Name:        issueNameKey,
			Placeholder: "Enter issue name",
			HelperText:  issueNameError,
			HelperType:  issueNameHelperType,
			InputProps: []g.Node{
				h.Value(issueNameValue),
				h.AutoComplete("off"),
			},
		}),

		components.Checkbox(
			&components.CheckboxProps{
				Name:    "IsArchived",
				Label:   "Is archived?",
				Checked: andonIssue.IsArchived,
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
