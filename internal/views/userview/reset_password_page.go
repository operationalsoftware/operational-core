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

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
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
		AppendBody: []g.Node{
			components.InlineScript("/static/js/app_nfc.js"),
			components.InlineScript("/internal/views/userview/reset_password_page.js"),
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
	encryptedCredentials := p.req.URL.Query().Get("EncryptedCredentials")

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

		h.Div(
			h.Class("form-wrapper"),
			h.Div(
				h.Class("form-inputs"),
				h.Div(
					h.Class("password-input-wrapper"),
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

					h.Button(
						h.Type("button"),
						h.Class("toggle-password-btn"),
						g.Attr("data-target", passwordKey),
						components.Icon(&components.IconProps{
							Identifier: "eye-off-outline",
							Classes: c.Classes{
								"eye-icon":      true,
								"eye-open-icon": true,
							},
						}),
						components.Icon(&components.IconProps{
							Identifier: "eye-outline",
							Classes: c.Classes{
								"eye-icon":        true,
								"eye-closed-icon": true,
								"hidden":          true,
							},
						}),
					),
				),
				h.Div(
					h.Class("password-input-wrapper"),
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

					h.Button(
						h.Type("button"),
						h.Class("toggle-password-btn"),
						g.Attr("data-target", confirmPasswordKey),
						components.Icon(&components.IconProps{
							Identifier: "eye-off-outline",
							Classes: c.Classes{
								"eye-icon":      true,
								"eye-open-icon": true,
							},
						}),
						components.Icon(&components.IconProps{
							Identifier: "eye-outline",
							Classes: c.Classes{
								"eye-icon":        true,
								"eye-closed-icon": true,
								"hidden":          true,
							},
						}),
					),
				),
			),

			h.Button(
				h.Class("generate-password-btn"),
				h.Type("button"),
				components.Icon(&components.IconProps{
					Identifier: "lock-plus-outline",
				}),
			),
		),
	}

	formElements = append(formElements,
		components.Button(&components.ButtonProps{}, h.Class("set-password-btn"), h.Type("submit"), g.Text("Set Password")),
	)

	if encryptedCredentials != "" {
		formElements = append(formElements,
			h.Div(
				h.Class("form-group"),
				h.Div(
					h.Class("form-success"), g.Text("Password set successfully!"),
				),

				components.Divider(g.Text("")),

				h.Label(
					h.Class("form-label"), g.Text("Encrypted Credentials:"),
				),
				h.Code(
					h.ID("encrypted-string"),
					h.Class("encrypted-code"),
					g.Text(encryptedCredentials),
				),
			),
		)

		formElements = append(formElements,
			h.Div(
				h.Class("nfc-btns"),

				h.Button(
					h.Class("button write-nfc-btn"),
					h.Type("button"),
					g.Text("Write NFC"),
				),
				h.Button(
					h.Class("button nfc-readonly-btn"),
					h.Type("button"),
					g.Text("Make NFC read-only"),
				),
			),
		)
	}

	return components.Form(
		h.ID("reset-password-form"),
		h.Method("POST"),
		g.Group(formElements),
	)
}
