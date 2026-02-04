package components

import (
	"app/internal/model"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type NotificationsTrayProps struct {
	Items       []model.NotificationItem
	UnreadCount int
}

func NotificationsTray(p NotificationsTrayProps) g.Node {
	var body g.Node
	if len(p.Items) == 0 {
		body = notificationsTrayEmpty()
	} else {
		body = h.Ul(
			h.Class("notifications-tray-list"),
			g.Group(g.Map(p.Items, func(item model.NotificationItem) g.Node {
				return notificationTrayItem(item)
			})),
		)
	}

	return h.Div(
		h.Class("notifications-tray"),
		h.Data("unread-count", strconv.Itoa(p.UnreadCount)),
		h.Div(
			h.Class("notifications-tray-header"),
			h.Div(
				h.Class("notifications-tray-title"),
				h.Span(g.Text("Notifications")),
				notificationsTrayBadge(p.UnreadCount),
			),
			h.A(
				h.Class("notifications-tray-view-all"),
				h.Href("/notifications"),
				g.Text("View all"),
			),
		),
		body,
	)
}

func notificationsTrayBadge(unreadCount int) g.Node {
	if unreadCount <= 0 {
		return g.Text("")
	}

	return Badge(&BadgeProps{
		Type: BadgePrimary,
		Size: BadgeSm,
	}, g.Text(fmt.Sprintf("%d unread", unreadCount)))
}

func notificationTrayItem(item model.NotificationItem) g.Node {
	itemURL := notificationOpenURL(item)
	actionURL := notificationOpenActionURL(item)

	title := strings.TrimSpace(item.Title)
	if title == "" {
		title = "Notification"
	}

	linkContent := g.Group([]g.Node{
		h.Span(h.Class("notifications-tray-dot")),
		h.Span(
			NotificationIconClasses(item),
			Icon(&IconProps{
				Identifier: NotificationIconIdentifier(item),
			}),
		),
		h.Span(
			h.Class("notifications-tray-content"),
			h.Span(
				h.Class("notifications-tray-item-title"),
				g.Text(title),
			),
			g.If(item.Summary != "",
				h.Span(
					h.Class("notifications-tray-summary"),
					g.Text(item.Summary),
				),
			),
			g.If(item.Time != "",
				h.Span(
					h.Class("notifications-tray-time"),
					g.Text(item.Time),
				),
			),
		),
	})

	linkNode := h.A(
		h.Class("notifications-tray-link"),
		h.Href(itemURL),
		linkContent,
	)

	if item.NotificationID > 0 {
		linkNode = h.Form(
			h.Method("POST"),
			h.Action(actionURL),
			h.Class("notifications-tray-form"),
			h.Button(
				h.Type("submit"),
				h.Class("notifications-tray-link notifications-tray-button"),
				linkContent,
			),
		)
	}

	return h.Li(
		c.Classes{
			"notifications-tray-item": true,
			"unread":                  item.Unread,
		},
		linkNode,
	)
}

func notificationsTrayEmpty() g.Node {
	return h.Div(
		h.Class("notifications-tray-empty"),
		Icon(&IconProps{Identifier: "inbox-outline"}),
		h.P(g.Text("You're all caught up.")),
	)
}

func notificationOpenURL(item model.NotificationItem) string {
	itemURL := strings.TrimSpace(item.URL)
	if itemURL == "" {
		return "/notifications"
	}
	return itemURL
}

func notificationOpenActionURL(item model.NotificationItem) string {
	if item.NotificationID == 0 {
		return notificationOpenURL(item)
	}
	target := notificationOpenURL(item)
	query := url.Values{}
	if target != "" {
		query.Set("Redirect", target)
	}
	return fmt.Sprintf("/notifications/%d/read?%s", item.NotificationID, query.Encode())
}

func NotificationIconIdentifier(item model.NotificationItem) string {
	switch model.NormalizeNotificationReasonType(item.ReasonType) {
	case model.NotificationReasonDanger:
		return "close"
	case model.NotificationReasonWarning:
		return "exclamation"
	case model.NotificationReasonSuccess:
		return "check"
	case model.NotificationReasonInfo:
		return "comment-text-outline"
	}
	return "comment-text-outline"
}

func NotificationIconClasses(item model.NotificationItem) c.Classes {
	classes := c.Classes{
		"notification-icon": true,
	}

	switch model.NormalizeNotificationReasonType(item.ReasonType) {
	case model.NotificationReasonInfo:
		classes["reason-info"] = true
	case model.NotificationReasonSuccess:
		classes["reason-success"] = true
	case model.NotificationReasonWarning:
		classes["reason-warning"] = true
	case model.NotificationReasonDanger:
		classes["reason-danger"] = true
	}

	return classes
}
