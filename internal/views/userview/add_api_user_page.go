package userview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"net/url"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AddAPIUserPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func AddAPIUserPage(p *AddAPIUserPageProps) g.Node {
	content := g.Group([]g.Node{

		addApiUserForm(&addApiUserFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Add New API User",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			usersBreadCrumb,
			{IconIdentifier: "add", Title: "Add API User"},
		},
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/userview/add_api_user_page.css"),
		},
	})
}

type addApiUserFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

func addApiUserForm(p *addApiUserFormProps) g.Node {
	usernameLabel := "Username"
	usernameKey := "Username"
	usernameValue := p.values.Get(usernameKey)
	usernameError := ""
	if p.isSubmission || usernameValue != "" {
		usernameError = p.validationErrors.GetError(usernameKey, usernameLabel)
	}
	usernameHelperType := components.InputHelperTypeNone
	if usernameError != "" {
		usernameHelperType = components.InputHelperTypeError
	}

	return components.Form(
		h.ID("add-api-user-form"),
		h.Method("POST"),

		components.Input(&components.InputProps{
			Label:       usernameLabel,
			Name:        usernameKey,
			Placeholder: "Enter username",
			HelperText:  usernameError,
			HelperType:  usernameHelperType,
			InputProps: []g.Node{
				h.Value(usernameValue),
				h.AutoComplete("off"),
			},
		}),

		permissionsCheckboxesPartial(model.UserPermissions{}),

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Submit"),
		),

		h.Div(
			h.ID("result"),
		),
	)
}
