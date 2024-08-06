package usersviews

import (
	"app/components"
	"app/internal/reqcontext"
	"app/internal/validate"
	"app/layout"
	"app/models/usermodel"
	"fmt"
	"net/url"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type EditUserPageProps struct {
	Ctx              reqcontext.ReqContext
	User             usermodel.User
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func EditUserPage(p *EditUserPageProps) g.Node {

	content := g.Group([]g.Node{
		editUserForm(&editUserFormProps{
			user: p.User,
		}),
	})

	return layout.Page(layout.PageProps{
		Title:   fmt.Sprintf("Edit User: %s", p.User.Username),
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/routes/users/usersviews/editUser.css"),
		},
	})
}

type editUserFormProps struct {
	user             usermodel.User
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

// same as addUserForm, but no password fields
func editUserForm(p *editUserFormProps) g.Node {

	firstNameLabel := "First Name"
	firstNameKey := "FirstName"
	var firstNameValue string
	if p.values.Get(firstNameKey) != "" {
		firstNameValue = p.values.Get(firstNameKey)
	} else {
		firstNameValue = p.user.FirstName.String
	}
	firstNameError := ""
	if p.isSubmission || firstNameValue != "" {
		firstNameError = p.validationErrors.GetError(firstNameKey, firstNameLabel)
	}
	firstNameHelperType := components.InputHelperTypeNone
	if firstNameError != "" {
		firstNameHelperType = components.InputHelperTypeError
	}

	lastNameLabel := "Last Name"
	lastNameKey := "LastName"
	var lastNameValue string
	if p.values.Get(lastNameKey) != "" {
		lastNameValue = p.values.Get(lastNameKey)
	} else {
		lastNameValue = p.user.LastName.String
	}
	lastNameError := ""
	if p.isSubmission || lastNameValue != "" {
		lastNameError = p.validationErrors.GetError(lastNameKey, lastNameLabel)
	}
	lastNameHelperType := components.InputHelperTypeNone
	if lastNameError != "" {
		lastNameHelperType = components.InputHelperTypeError
	}

	usernameLabel := "Username"
	usernameKey := "Username"
	var usernameValue string
	if p.values.Get(usernameKey) != "" {
		usernameValue = p.values.Get(usernameKey)
	} else {
		usernameValue = p.user.Username
	}
	usernameError := ""
	if p.isSubmission || usernameValue != "" {
		usernameError = p.validationErrors.GetError(usernameKey, usernameLabel)
	}
	usernameHelperType := components.InputHelperTypeNone
	if usernameError != "" {
		usernameHelperType = components.InputHelperTypeError
	}

	emailLabel := "Email"
	emailKey := "Email"
	var emailValue string
	if p.values.Get(emailKey) != "" {
		emailValue = p.values.Get(emailKey)
	} else {
		emailValue = p.user.Email.String
	}
	emailError := ""
	if p.isSubmission || emailValue != "" {
		emailError = p.validationErrors.GetError(emailKey, emailLabel)
	}
	emailHelperType := components.InputHelperTypeNone
	if emailError != "" {
		emailHelperType = components.InputHelperTypeError
	}

	return components.Form(
		h.ID("edit-user-form"),
		h.Method("POST"),

		components.Input(&components.InputProps{
			Label:       firstNameLabel,
			Name:        firstNameKey,
			Placeholder: "Enter first name",
			HelperText:  firstNameError,
			HelperType:  firstNameHelperType,
			InputProps: []g.Node{
				h.Value(firstNameValue),
				h.AutoComplete("off"),
			},
		}),

		components.Input(&components.InputProps{
			Label:       lastNameLabel,
			Name:        lastNameKey,
			Placeholder: "Enter last name",
			HelperText:  lastNameError,
			HelperType:  lastNameHelperType,
			InputProps: []g.Node{
				h.Value(lastNameValue),
				h.AutoComplete("off"),
			},
		}),

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

		components.Input(&components.InputProps{
			Label:       emailLabel,
			Name:        emailKey,
			Placeholder: "Enter email",
			HelperText:  emailError,
			HelperType:  emailHelperType,
			InputProps: []g.Node{
				h.Value(emailValue),
				h.AutoComplete("off"),
			},
		}),

		permissionsCheckboxes(p.user.Permissions),

		components.Button(
			&components.ButtonProps{
				Disabled: p.validationErrors.HasErrors(),
			},
			h.Type("submit"),
			g.Text("Submit"),
		),
	)
}
