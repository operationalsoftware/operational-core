package userview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/cookie"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/url"
	"strconv"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type EditUserPageProps struct {
	Ctx              reqcontext.ReqContext
	User             model.User
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func EditUserPage(p *EditUserPageProps) g.Node {

	content := g.Group([]g.Node{
		editUserForm(&editUserFormProps{
			user:             p.User,
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Title: fmt.Sprintf("Edit: %s", p.User.Username),
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			usersBreadCrumb,
			{
				IconIdentifier: "account",
				Title:          p.User.Username,
				URLPart:        fmt.Sprintf("%d", p.User.UserID),
			},
			{Title: "Edit"},
		},
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/userview/edit_user_page.css"),
		},
	})
}

type editUserFormProps struct {
	user             model.User
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

// same as addUserForm, but no password fields
func editUserForm(p *editUserFormProps) g.Node {

	isAPIUser := p.user.IsAPIUser

	firstNameLabel := "First Name"
	firstNameKey := "FirstName"
	var firstNameValue string
	if p.values.Get(firstNameKey) != "" {
		firstNameValue = p.values.Get(firstNameKey)
	} else {
		firstNameValue = nilsafe.Str(p.user.FirstName)
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
		lastNameValue = nilsafe.Str(p.user.LastName)
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
		emailValue = nilsafe.Str(p.user.Email)
	}
	emailError := ""
	if p.isSubmission || emailValue != "" {
		emailError = p.validationErrors.GetError(emailKey, emailLabel)
	}
	emailHelperType := components.InputHelperTypeNone
	if emailError != "" {
		emailHelperType = components.InputHelperTypeError
	}

	sessdurationLabel := "Session Duration In Minutes"
	sessdurationKey := "SessionDurationMinutes"
	var sessdurationValue string
	if p.values.Get(sessdurationKey) != "" {
		sessdurationValue = p.values.Get(sessdurationKey)
	} else {
		if p.user.SessionDurationMinutes != nil {
			sessdurationValue = strconv.Itoa(*p.user.SessionDurationMinutes)
		}
	}
	sessdurationError := ""
	if p.isSubmission || sessdurationValue != "" {
		sessdurationError = p.validationErrors.GetError(sessdurationKey, sessdurationLabel)
	}
	sessdurationHelperType := components.InputHelperTypeNone
	if sessdurationError != "" {
		sessdurationHelperType = components.InputHelperTypeError
	}

	return components.Form(
		h.ID("edit-user-form"),
		h.Method("POST"),

		g.If(
			!isAPIUser,
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
		),

		g.If(
			!isAPIUser,
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
		),

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

		g.If(
			!isAPIUser,
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
		),

		g.If(
			!isAPIUser,
			h.Div(
				components.Input(&components.InputProps{
					Label:       sessdurationLabel,
					Name:        sessdurationKey,
					Placeholder: "Enter session duration in minutes",
					HelperText:  sessdurationError,
					HelperType:  sessdurationHelperType,
					InputType:   "number",
					InputProps: []g.Node{
						h.Value(sessdurationValue),
						h.AutoComplete("off"),
						h.Min("1"),
						h.Max("525600"),
					},
				}),

				h.P(
					h.Class("session-helper"),

					g.Text(fmt.Sprintf("* Leave unset to use organization default: %d", int(cookie.DefaultSessionDurationMinutes.Minutes()))),
				),
			),
		),

		permissionsCheckboxesPartial(p.user.Permissions),

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Submit"),
		),
	)
}
