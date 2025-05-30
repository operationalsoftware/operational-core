package andonissueview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"net/url"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type AddPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func AddPage(p *AddPageProps) g.Node {

	content := g.Group([]g.Node{

		addForm(&addTeamFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
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

	return components.Form(
		h.ID("add-andon-issue-form"),
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

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Submit"),
		),
	)
}
