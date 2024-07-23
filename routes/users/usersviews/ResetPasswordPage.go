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

type ResetPasswordPageProps struct {
	User             usermodel.User
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validation.ValidationErrors
	IsSubmission     bool
}

func ResetPasswordPage(p *ResetPasswordPageProps) g.Node {

	resetPasswordContent := g.Group([]g.Node{
		resetPasswordForm(&resetPasswordFormProps{
			userID:           p.User.UserID,
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Title:   "Reset Password: " + p.User.Username,
		Content: resetPasswordContent,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/routes/users/usersviews/resetPassword.css"),
		},
	})
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
