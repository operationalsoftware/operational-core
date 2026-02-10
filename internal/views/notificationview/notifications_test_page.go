package notificationview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type NotificationTestPageProps struct {
	Ctx reqcontext.ReqContext
}

func NotificationTestPage(p NotificationTestPageProps) g.Node {
	content := h.Div(
		h.Class("notification-test"),
		h.H2(g.Text("Test notification Page")),
		h.P(g.Text("This is a test notification page. Notifications can take you to any page on the app which is set at during the creation of the notification.")),
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
