package userview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/http"
	"net/url"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type ResetPasswordPageProps struct {
	User             model.User
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func ResetPasswordPage(p *ResetPasswordPageProps) g.Node {

	resetPasswordContent := g.Group([]g.Node{
		resetPasswordForm(&resetPasswordFormProps{
			req:              p.Ctx.Req,
			userID:           p.User.UserID,
			userName:         p.Ctx.User.Username,
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Title: "Reset Password: " + p.User.Username,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			usersBreadCrumb,
			{
				IconIdentifier: "account",
				Title:          p.User.Username,
				URLPart:        fmt.Sprintf("%d", p.User.UserID),
			},
			{Title: "Reset Password"},
		},
		Content: resetPasswordContent,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/userview/reset_password_page.css"),
		},
	})
}

type resetPasswordFormProps struct {
	req              *http.Request
	userID           int
	userName         string
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

func resetPasswordForm(p *resetPasswordFormProps) g.Node {
	// Generate encrypted value
	encryptedData := p.req.URL.Query().Get("encryptedData")

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

	formElements := []g.Node{
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
	}

	if encryptedData != "" {
		formElements = append(formElements,
			h.Div(
				h.Class("form-group"),
				h.Div(
					h.Class("form-success"), g.Text("Password set successfully!"),
				),
				h.Label(
					h.Class("form-label"), g.Text("Encrypted Credentials:"),
				),
				h.Code(
					h.ID("encrypted-string"),
					h.Class("encrypted-code"),
					g.Text(encryptedData),
				),
			),
		)
	}

	formElements = append(formElements,
		components.Button(&components.ButtonProps{}, h.Type("submit"), g.Text("Submit")),
	)

	return components.Form(
		h.ID("reset-password-form"),
		h.Method("POST"),
		g.Group(formElements),
	)
}
