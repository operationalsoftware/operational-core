package usersviews

import (
	"app/components"
	"app/internal/reqcontext"
	"app/internal/validation"
	"app/layout"
	"app/models/usermodel"
	"net/url"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type AddAPIUserPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validation.ValidationErrors
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
		Ctx:     p.Ctx,
		Title:   "Add New API User",
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/routes/users/usersviews/addApiUser.css"),
		},
	})
}

type addApiUserFormProps struct {
	values           url.Values
	validationErrors validation.ValidationErrors
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

		permissionsCheckboxes(usermodel.UserPermissions{}),

		components.Button(
			&components.ButtonProps{
				Disabled: p.validationErrors.HasErrors(),
			},
			h.Type("submit"),
			g.Text("Submit"),
		),

		h.Div(
			h.ID("result"),
		),
	)
}
