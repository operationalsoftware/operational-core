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

type AddGroupPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	AndonIssueGroups []model.AndonIssueGroup
	Teams            []model.Team
}

func AddGroupPage(p *AddGroupPageProps) g.Node {

	content := g.Group([]g.Node{

		addGroupForm(&addGroupFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			andonIssueGroups: p.AndonIssueGroups,
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
			{IconIdentifier: "plus", Title: "Add Group"},
		},
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/andonissueview/add_group_page.css"),
		},
	})
}

type addGroupFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
	andonIssueGroups []model.AndonIssueGroup
}

func addGroupForm(p *addGroupFormProps) g.Node {

	issueNameLabel := "Group Name"
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

	parentIDLabel := "Child Group Of"
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
			g.Text("Select child group"),
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

				h.P(
					h.Class("note"),
					g.Text("* Only two levels of groups are supported"),
				),
			),
			g.If(parentIDError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: parentIDError,
					Type:  parentIDHelperType,
				})),
		),

		components.Button(
			&components.ButtonProps{},

			h.Type("submit"),
			g.Text("Add Issue Group"),
		),
	)
}
