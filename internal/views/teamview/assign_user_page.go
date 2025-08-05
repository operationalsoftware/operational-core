package teamview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"net/url"
	"strconv"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type AssignUserPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	Users            []model.User
}

func AssignUserPage(p *AssignUserPageProps) g.Node {

	content := g.Group([]g.Node{

		assignUserForm(&assignUserFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
			users:            p.Users,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Add New Team",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			teamsBreadCrumb,
			{IconIdentifier: "account-plus", Title: "Assign User"},
		},
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/teamview/assign_user_page.css"),
		},
	})
}

type assignUserFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
	users            []model.User
}

func assignUserForm(p *assignUserFormProps) g.Node {

	usernameLabel := "Username"
	usernameKey := "UserID"
	usernameValue := p.values.Get(usernameKey)
	usernameError := ""
	if p.isSubmission || usernameValue != "" {
		usernameError = p.validationErrors.GetError(usernameKey, usernameLabel)
	}
	usernameHelperType := components.InputHelperTypeNone
	if usernameError != "" {
		usernameHelperType = components.InputHelperTypeError
	}

	roleLabel := "Role"
	roleKey := "Role"
	roleValue := p.values.Get(roleKey)
	roleError := ""
	if p.isSubmission || roleValue != "" {
		roleError = p.validationErrors.GetError(roleKey, roleLabel)
	}
	roleHelperType := components.InputHelperTypeNone
	if roleError != "" {
		roleHelperType = components.InputHelperTypeError
	}

	roleSelectOptions := []g.Node{
		h.Option(
			h.Value(""),
			g.Text("\u2013"),
		),
	}
	for _, role := range []string{"Member", "Manager"} {

		isSelected := role == roleValue

		roleSelectOptions = append(roleSelectOptions, h.Option(
			h.Value(role),
			g.If(isSelected, h.Selected()),
			g.Text(role),
		))
	}

	return components.Form(
		h.ID("add-team-form"),
		h.Method("POST"),

		h.Div(
			h.Label(
				g.Text(usernameLabel),
			),
			components.SearchSelect(&components.SearchSelectProps{
				Name:          usernameKey,
				Placeholder:   "Select username",
				Mode:          "single",
				Options:       MapUsersToOptions(p.users),
				ShowOnlyLabel: true,
				Selected:      usernameValue,
			}),
			g.If(
				usernameError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: usernameError,
					Type:  usernameHelperType,
				},
				),
			),
		),

		h.Div(
			h.Label(
				g.Text(roleLabel),

				h.Select(
					h.Name(roleKey),
					g.Group(roleSelectOptions),
				),
			),
			g.If(roleError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: roleError,
					Type:  roleHelperType,
				}),
			),
		),

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Submit"),
		),
	)
}

func MapUsersToOptions(users []model.User) []components.SearchSelectOption {
	out := make([]components.SearchSelectOption, len(users))
	for i, v := range users {
		out[i] = components.SearchSelectOption{
			Label: v.Username,
			Value: strconv.Itoa(v.UserID),
		}
	}
	return out
}
