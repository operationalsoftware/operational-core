package users

import (
	"app/components"
	userModel "app/src/users/model"
	"app/validation"
	"fmt"
	"net/url"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
	"golang.org/x/exp/slices"
)

type addUserFormProps struct {
	values           url.Values
	validationErrors validation.ValidationErrors
	isSubmission     bool
}

func addUserForm(p *addUserFormProps) g.Node {

	firstNameLabel := "First Name"
	firstNameKey := "FirstName"
	firstNameValue := p.values.Get(firstNameKey)
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
	lastNameValue := p.values.Get(lastNameKey)
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
	usernameValue := p.values.Get(usernameKey)
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
	emailValue := p.values.Get(emailKey)
	emailError := ""
	if p.isSubmission || emailValue != "" {
		emailError = p.validationErrors.GetError(emailKey, emailLabel)
	}
	emailHelperType := components.InputHelperTypeNone
	if emailError != "" {
		emailHelperType = components.InputHelperTypeError
	}

	passwordLabel := "Password"
	passwordKey := "Password"
	passwordValue := p.values.Get(passwordKey)
	passwordError := ""
	if p.isSubmission || passwordValue != "" {
		passwordError = p.validationErrors.GetError(passwordKey, passwordLabel)
	}
	passwordHelperType := components.InputHelperTypeNone
	if passwordError != "" {
		passwordHelperType = components.InputHelperTypeError
	}

	confirmPasswordLabel := "Confirm Password"
	confirmPasswordKey := "ConfirmPassword"
	confirmPasswordValue := p.values.Get(confirmPasswordKey)
	confirmPasswordError := ""
	if p.isSubmission || confirmPasswordValue != "" {
		confirmPasswordError = p.validationErrors.GetError(confirmPasswordKey, confirmPasswordLabel)
	}
	confirmPasswordHelperType := components.InputHelperTypeNone
	if confirmPasswordError != "" {
		confirmPasswordHelperType = components.InputHelperTypeError
	}

	commonHtmx := g.Group([]g.Node{
		hx.Post("/users/add/validate"),
		hx.Target("closest form"),
		hx.Select("form > *"),
	})

	return components.Form(
		h.ID("add-user-form"),
		hx.Post("/users/add"),

		components.Input(&components.InputProps{
			Label:       firstNameLabel,
			Name:        firstNameKey,
			Placeholder: "Enter first name",
			HelperText:  firstNameError,
			HelperType:  firstNameHelperType,
			InputProps: []g.Node{
				h.Value(firstNameValue),
				commonHtmx,
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
				commonHtmx,
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
				commonHtmx,
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
				commonHtmx,
			},
		}),

		components.Input(&components.InputProps{
			Label:       passwordLabel,
			Name:        passwordKey,
			InputType:   "password",
			Placeholder: "Enter password",
			HelperText:  passwordError,
			HelperType:  passwordHelperType,
			InputProps: []g.Node{
				h.Value(passwordValue),
				commonHtmx,
			},
		}),

		components.Input(&components.InputProps{
			Label:       confirmPasswordLabel,
			Name:        confirmPasswordKey,
			InputType:   "password",
			Placeholder: "Confirm password",
			HelperText:  confirmPasswordError,
			HelperType:  confirmPasswordHelperType,
			InputProps: []g.Node{
				h.Value(confirmPasswordValue),
				commonHtmx,
			},
		}),

		components.Button(
			&components.ButtonProps{
				Disabled: p.validationErrors.HasErrors(),
			},
			h.Type("submit"),
			g.Text("Submit"),
		),
	)

}

type editUserFormProps struct {
	user             userModel.User
	values           url.Values
	validationErrors validation.ValidationErrors
	isSubmission     bool
}

// same as addUserForm, but no password fields
func editUserForm(p *editUserFormProps) g.Node {

	var userRoleOptions = []components.CheckboxOption{
		{
			Value: "User Admin",
			Label: "User Admin",
		},
	}

	// update user role options, checked if user has the role
	for idx, option := range userRoleOptions {
		// check if user.Roles contains option.Value
		if slices.Contains(p.user.Roles, option.Value) {
			userRoleOptions[idx].Checked = true
		}
	}

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

	commonHtmx := g.Group([]g.Node{
		hx.Post(fmt.Sprintf("/users/%d/edit/validate", p.user.UserID)),
		hx.Target("closest form"),
		hx.Select("form > *"),
	})

	return components.Form(
		h.ID("edit-user-form"),
		hx.Post(fmt.Sprintf("/users/%d/edit", p.user.UserID)),

		components.Input(&components.InputProps{
			Label:       firstNameLabel,
			Name:        firstNameKey,
			Placeholder: "Enter first name",
			HelperText:  firstNameError,
			HelperType:  firstNameHelperType,
			InputProps: []g.Node{
				h.Value(firstNameValue),
				commonHtmx,
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
				commonHtmx,
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
				commonHtmx,
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
				commonHtmx,
			},
		}),

		components.CheckboxGroup(&components.CheckboxGroupProps{
			Name:    "Roles",
			Label:   "Roles",
			Options: userRoleOptions,
		}),

		components.Button(
			&components.ButtonProps{
				Disabled: p.validationErrors.HasErrors(),
			},
			h.Type("submit"),
			g.Text("Submit"),
		),
	)
}
