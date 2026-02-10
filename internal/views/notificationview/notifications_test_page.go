package notificationview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"net/url"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type NotificationTestPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	SuccessText      string
	ErrorText        string
}

func NotificationTestPage(p NotificationTestPageProps) g.Node {
	content := h.Div(
		h.Class("notification-test"),
		g.If(p.SuccessText != "",
			components.Alert(&components.AlertProps{
				AlertType: components.AlertSuccess,
				Message:   p.SuccessText,
			}),
		),
		g.If(p.ErrorText != "",
			components.Alert(&components.AlertProps{
				AlertType: components.AlertError,
				Message:   p.ErrorText,
			}),
		),
		notificationTestForm(&notificationTestFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
		}),
	)

	return layout.Page(layout.PageProps{
		Title:   "Test Notification",
		Content: content,
		Ctx:     p.Ctx,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "bell-outline",
				Title:          "Notifications",
				URLPart:        "notifications",
			},
			{
				Title: "Test notification",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/notificationview/notifications_test_page.css"),
		},
	})
}

type notificationTestFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
}

func notificationTestForm(p *notificationTestFormProps) g.Node {
	messageKey := "Message"
	messageValue := strings.TrimSpace(p.values.Get(messageKey))

	messageError := ""
	helperType := components.InputHelperTypeNone
	if p.validationErrors.HasErrors() || messageValue != "" {
		messageError = p.validationErrors.GetError(messageKey, "Message")
		if messageError != "" {
			helperType = components.InputHelperTypeError
		}
	}

	return h.Form(
		h.Method("POST"),
		h.Class("form"),
		h.Action("/notifications/test"),

		h.Div(
			h.Label(
				g.Text("Message"),
				h.Input(
					h.Name(messageKey),
					h.Placeholder("Type a test message"),
					h.Value(messageValue),
					h.AutoComplete("off"),
				),
			),
			g.If(
				messageError != "",
				components.InputHelper(&components.InputHelperProps{
					Label: messageError,
					Type:  helperType,
				}),
			),
		),

		h.Button(
			h.Class("button primary"),
			h.Type("submit"),
			g.Text("Send test notification"),
		),
	)
}
