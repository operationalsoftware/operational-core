package users

import (
	"app/components"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type confirmPasswordInputProps struct {
	ValidationError string
	Value           string
}

func confirmPasswordInput(p *confirmPasswordInputProps) g.Node {
	inputProps := &components.InputProps{
		Label:       "Confirm password",
		Name:        "ConfirmPassword",
		InputType:   "password",
		Placeholder: "Confirm password",
		InputProps: []g.Node{
			hx.Post("/users/validate/confirm-password"),
			h.Value(p.Value),
		},
	}

	if p.ValidationError != "" {
		inputProps.HelperText = p.ValidationError
		inputProps.HelperType = components.InputHelperTypeError
	}

	return components.Input(inputProps,
		hx.Target("this"),
		hx.Swap("outerHTML"),
	)
}

type emailInputProps struct {
	ValidationError string
	Value           string
}

func emailInput(p *emailInputProps) g.Node {
	inputProps := &components.InputProps{
		Label:       "Email",
		Name:        "Email",
		InputType:   "email",
		Placeholder: "Enter email",
		InputProps: []g.Node{
			hx.Post("/users/validate/email"),
			h.Value(p.Value),
		},
	}

	if p.ValidationError != "" {
		inputProps.HelperText = p.ValidationError
		inputProps.HelperType = components.InputHelperTypeError
	}

	return components.Input(inputProps,
		hx.Target("this"),
		hx.Swap("outerHTML"),
	)
}

type firstNameInputProps struct {
	ValidationError string
	Value           string
}

func firstNameInput(p *firstNameInputProps) g.Node {
	inputProps := &components.InputProps{
		Label:       "First Name",
		Name:        "FirstName",
		Placeholder: "Enter first name",
		InputProps: []g.Node{
			hx.Post("/users/validate/first-name"),
			h.Value(p.Value),
		},
	}

	if p.ValidationError != "" {
		inputProps.HelperText = p.ValidationError
		inputProps.HelperType = components.InputHelperTypeError
	}

	return components.Input(inputProps,
		hx.Target("this"),
		hx.Swap("outerHTML"),
	)
}

type lastNameInputProps struct {
	ValidationError string
	Value           string
}

func lastNameInput(p *lastNameInputProps) g.Node {
	inputProps := &components.InputProps{
		Label:       "Last Name",
		Name:        "LastName",
		Placeholder: "Enter last name",
		InputProps: []g.Node{
			hx.Post("/users/validate/last-name"),
			h.Value(p.Value),
		},
	}

	if p.ValidationError != "" {
		inputProps.HelperText = p.ValidationError
		inputProps.HelperType = components.InputHelperTypeError
	}

	return components.Input(inputProps,
		hx.Target("this"),
		hx.Swap("outerHTML"),
	)
}

type passwordInputProps struct {
	ValidationError string
	Value           string
}

func passwordInput(p *passwordInputProps) g.Node {
	inputProps := &components.InputProps{
		Label:       "Password",
		Name:        "Password",
		InputType:   "password",
		Placeholder: "Enter password",
		InputProps: []g.Node{
			hx.Post("/users/validate/password"),
			h.Value(p.Value),
		},
	}

	if p.ValidationError != "" {
		inputProps.HelperText = p.ValidationError
		inputProps.HelperType = components.InputHelperTypeError
	}

	return components.Input(inputProps,
		hx.Target("this"),
		hx.Swap("outerHTML"),
	)
}

type usernameInputProps struct {
	ValidationError string
	Value           string
}

func usernameInput(p *usernameInputProps) g.Node {
	inputProps := &components.InputProps{
		Label:       "Username",
		Name:        "Username",
		Placeholder: "Enter username",
		InputProps: []g.Node{
			hx.Post("/users/validate/username"),
			h.Value(p.Value),
		},
	}

	if p.ValidationError != "" {
		inputProps.HelperText = p.ValidationError
		inputProps.HelperType = components.InputHelperTypeError
	}

	return components.Input(inputProps,
		hx.Target("this"),
		hx.Swap("outerHTML"),
	)
}

// type rolesInputProps struct {
// 	Roles []string
// }

// func rolesCheckBoxes(p *rolesInputProps) g.Node {
// 	checkboxes := []g.Node{}
// 	for _, role := range p.Roles {
// 		checkboxes = append(checkboxes, components.Checkbox(&components.CheckboxProps{
// 			Name:    role,
// 			Label:   role,
// 			Value:   role,
// 			Checked: strings.Contains(strings.Join(p.Roles, ","), role),
// 		}))
// 	}

// 	return h.Div(checkboxes...)
// }
