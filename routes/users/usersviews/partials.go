package usersviews

import (
	"app/components"
	"app/models/usermodel"
	"app/validation"
	"fmt"
	"net/url"
	"reflect"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func getPermissionDescription(module, permission string) string {
	// Use reflection to get the description of
	// Create an instance of UserPermissions
	permissions := usermodel.UserPermissions{}
	permissionsType := reflect.TypeOf(permissions)
	moduleField, found := permissionsType.FieldByName(module)
	if found {
		permissionField, found := moduleField.Type.FieldByName(permission)
		if found {
			description := permissionField.Tag.Get("description")
			if description != "" {
				return description
			}
		}
	}

	return ""
}

func permissionsCheckboxes(userPermissions usermodel.UserPermissions) g.Node {

	return components.Fieldset(
		h.Label(g.Text("Permissions")),
		h.H4(h.Class("module-title"), g.Text("User Admin")),
		components.Checkbox(
			&components.CheckboxProps{
				Classes: c.Classes{
					"permission-checkbox": true,
				},
				Name:    "Permissions.UserAdmin.Access",
				Label:   getPermissionDescription("UserAdmin", "Access"),
				Checked: userPermissions.UserAdmin.Access,
				Value:   "true",
			},
		),

		components.InlineStyle("/src/users/permissions.css"),
	)

}

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

	return components.Form(
		h.ID("add-user-form"),
		hx.Post("/users/add"),
		hx.Select("form > *"),

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

		components.Input(&components.InputProps{
			Label:       passwordLabel,
			Name:        passwordKey,
			InputType:   "password",
			Placeholder: "Enter password",
			HelperText:  passwordError,
			HelperType:  passwordHelperType,
			InputProps: []g.Node{
				h.Value(passwordValue),
				h.AutoComplete("off"),
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
	)

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
		hx.Post("/users/add-api-user"),
		// Default behaviour here assumes invalid form. Overridden with headers from the server
		hx.Select("form > *"),
		hx.Target("closest form"),

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

type apiUserCredentialsProps struct {
	Username string
	Password string
}

func apiUserCredentials(p *apiUserCredentialsProps) g.Node {
	return components.Card(
		h.Div(
			h.Class("api-user-credentials"),
			h.Div(
				h.Class("content"),
				h.H2(
					g.Text("API User Credentials"),
				),
				h.Span(g.Text("Username: ")),
				h.Span(g.Text(p.Username)),
				h.Br(),
				h.Span(g.Text("Password: ")),
				h.Span(g.Text(p.Password)),
			),
			components.Button(&components.ButtonProps{
				ButtonType: components.ButtonPrimary,
				Link:       "/users",
				Classes: c.Classes{
					"users-btn": true,
				},
			},
				g.Text("Back to Users"),
			),
		),
	)
}

type editUserFormProps struct {
	user             usermodel.User
	values           url.Values
	validationErrors validation.ValidationErrors
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
		hx.Post(fmt.Sprintf("/users/%d/edit", p.user.UserID)),
		hx.Select("form > *"),

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

type resetPasswordFormProps struct {
	userID           int
	values           url.Values
	validationErrors validation.ValidationErrors
	isSubmission     bool
}

func resetPasswordForm(p *resetPasswordFormProps) g.Node {

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

	return components.Form(
		h.ID("reset-password-form"),
		hx.Post(fmt.Sprintf("/users/%d/reset-password", p.userID)),
		hx.Select("form > *"),

		components.Input(&components.InputProps{
			Label:       passwordLabel,
			Name:        passwordKey,
			InputType:   "password",
			Placeholder: "Enter password",
			HelperText:  passwordError,
			HelperType:  passwordHelperType,
			InputProps: []g.Node{
				h.Value(passwordValue),
				h.AutoComplete("off"),
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
				h.AutoComplete("off"),
			},
		}),

		components.Button(&components.ButtonProps{}, h.Type("submit"), g.Text("Submit")),
	)
}
