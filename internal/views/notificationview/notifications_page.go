package notificationview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type NotificationPageProps struct {
	Ctx            reqcontext.ReqContext
	Filters        []model.NotificationFilter
	ActiveFilter   string
	Groups         []model.NotificationGroup
	UnreadCount    int
	Page           int
	PageSize       int
	TotalRecords   int
	VAPIDPublicKey string
}

func NotificationPage(p NotificationPageProps) g.Node {
	filters := p.Filters
	if len(filters) == 0 {
		filters = defaultNotificationFilters(p.UnreadCount)
	}

	activeFilter := p.ActiveFilter
	if activeFilter == "" {
		activeFilter = "unread"
	}
	p.ActiveFilter = activeFilter

	content := h.Div(
		h.Class("notifications-page"),

		notificationsHeader(p.UnreadCount, activeFilter, p.VAPIDPublicKey),

		h.Div(
			h.Class("notifications-layout"),
			notificationSidebar(filters, activeFilter),
			notificationsFeedSection(&p),
		),
	)

	return layout.Page(layout.PageProps{
		Title:   "Notifications",
		Content: content,
		Ctx:     p.Ctx,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "bell-outline",
				Title:          "Notifications",
				URLPart:        "notifications",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/notificationview/notifications_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/components/Table.js"),
			components.InlineScript("/internal/views/notificationview/notifications_page.js"),
		},
	})
}

func notificationsHeader(unreadCount int, activeFilter string, vapidPublicKey string) g.Node {
	subtitle := "Activity across the platform."
	showMarkAll := strings.ToLower(activeFilter) != "read"
	showPush := strings.TrimSpace(vapidPublicKey) != ""

	return h.Div(
		h.Class("notifications-header"),
		h.Div(
			h.Class("notifications-title"),
			h.H1(g.Text("Notifications")),
			h.P(g.Text(subtitle)),
		),
		g.If(
			showMarkAll || showPush,
			h.Div(
				h.Class("notifications-actions"),
				g.If(showPush, notificationsPushAction(vapidPublicKey)),
				g.If(showMarkAll,
					h.Form(
						h.Method("POST"),
						h.Action("/notifications/mark-all-read"),
						h.Class("notifications-mark-all"),
						components.Button(&components.ButtonProps{
							ButtonType: components.ButtonPrimary,
							Size:       components.ButtonSm,
						},
							components.Icon(&components.IconProps{Identifier: "text-box-check-outline"}),
							g.Text("Mark all as read"),
						),
					),
				),
			),
		),
	)
}

func notificationsPushAction(vapidPublicKey string) g.Node {
	vapidPublicKey = strings.TrimSpace(vapidPublicKey)
	if vapidPublicKey == "" {
		return g.Group(nil)
	}

	return h.Div(
		h.Class("notifications-push-action"),
		h.Button(
			h.Type("button"),
			h.Class("button small notifications-push-button"),
			h.Data("push-toggle", "true"),
			h.Data("vapid-public-key", vapidPublicKey),
			h.Span(
				h.Class("notifications-push-icon"),
				components.Icon(&components.IconProps{
					Identifier: "bell-outline",
					Classes:    c.Classes{"icon-off": true},
				}),
				components.Icon(&components.IconProps{
					Identifier: "bell",
					Classes:    c.Classes{"icon-on": true},
				}),
			),
			h.Span(
				h.Class("notifications-push-button-text"),
				g.Text("Notifications: Off"),
			),
		),
	)
}

func notificationsFeedSection(p *NotificationPageProps) g.Node {
	return h.Div(
		h.Class("notifications-feed-section"),
		notificationFeed(p.Groups, p.UnreadCount, p),
		notificationsPaginationForm(p),
	)
}

func notificationSidebar(filters []model.NotificationFilter, activeFilter string) g.Node {
	return h.Aside(
		h.Class("notifications-sidebar"),
		components.Card(
			h.Div(
				h.Class("sidebar-section"),
				h.H3(h.Class("sidebar-title"), g.Text("Filters")),
				notificationFilters(filters, activeFilter),
			),
		),
	)
}

func notificationFilters(filters []model.NotificationFilter, activeFilter string) g.Node {
	var items []g.Node
	for _, filter := range filters {
		items = append(items,
			h.Li(
				h.A(
					c.Classes{
						"notification-filter": true,
						"active":              filter.ID == activeFilter,
					},
					h.Href(filter.URL),
					h.Span(g.Text(filter.Label)),
					h.Span(
						h.Class("notification-filter-count"),
						g.Text(strconv.Itoa(filter.Count)),
					),
				),
			),
		)
	}

	return h.Ul(
		h.Class("notification-filter-list"),
		g.Group(items),
	)
}

func notificationFeed(groups []model.NotificationGroup, unreadCount int, p *NotificationPageProps) g.Node {
	var body []g.Node

	body = append(body, notificationsFeedHeader(unreadCount, p.ActiveFilter))

	hasItems := false
	for _, group := range groups {
		if len(group.Items) == 0 {
			continue
		}
		hasItems = true
		body = append(body, notificationGroup(group, p))
	}

	if !hasItems {
		body = append(body, notificationsEmptyState(p.ActiveFilter))
	}

	return h.Section(
		h.Class("notifications-feed"),
		components.Card(
			g.Group(body),
		),
	)
}

func notificationsPaginationForm(p *NotificationPageProps) g.Node {
	if p.TotalRecords <= 0 {
		return g.Group(nil)
	}

	return h.Form(
		h.ID("notifications-pagination-form"),
		h.Method("GET"),
		h.Input(
			h.Type("hidden"),
			h.Name("Filter"),
			h.Value(p.ActiveFilter),
		),
		h.Input(
			h.Type("radio"),
			h.Checked(),
			h.Name("Page"),
			h.Value(fmt.Sprintf("%d", p.Page)),
			h.Style("display: none"),
		),
		h.Div(
			h.Class("notifications-pagination"),
			components.TablePagination(&components.TablePaginationProps{
				TotalRecords:        p.TotalRecords,
				CurrentPage:         p.Page,
				CurrentPageQueryKey: "Page",
				PageSize:            p.PageSize,
				PageSizeQueryKey:    "PageSize",
			}),
		),
	)
}

func notificationsFeedHeader(unreadCount int, activeFilter string) g.Node {
	filter := strings.ToLower(strings.TrimSpace(activeFilter))
	title := "Unread"
	showBadge := true

	if filter == "read" {
		title = "Read"
		showBadge = false
	}

	var badge g.Node
	if showBadge {
		if unreadCount > 0 {
			badge = components.Badge(&components.BadgeProps{
				Type: components.BadgePrimary,
				Size: components.BadgeSm,
			}, g.Text(fmt.Sprintf("%d unread", unreadCount)))
		} else {
			badge = components.Badge(&components.BadgeProps{
				Type: components.BadgeSecondary,
				Size: components.BadgeSm,
			}, g.Text("All read"))
		}
	}

	return h.Div(
		h.Class("notifications-feed-header"),
		h.Div(
			h.Class("notifications-feed-title"),
			h.H2(g.Text(title)),
			g.If(showBadge, badge),
		),
	)
}

func notificationGroup(group model.NotificationGroup, p *NotificationPageProps) g.Node {
	return h.Div(
		h.Class("notification-group"),
		h.Div(
			h.Class("notification-group-header"),
			h.Span(g.Text(group.Title)),
			h.Span(g.Text(fmt.Sprintf("%d updates", len(group.Items)))),
		),
		h.Ul(
			h.Class("notification-list"),
			g.Group(g.Map(group.Items, func(item model.NotificationItem) g.Node {
				return notificationItem(item, p)
			})),
		),
	)
}

func notificationItem(item model.NotificationItem, p *NotificationPageProps) g.Node {
	iconIdentifier := components.NotificationIconIdentifier(item)
	linkURL := notificationOpenActionURL(item, p)
	titleClasses := c.Classes{
		"notification-title": true,
	}
	titleNode := h.Span(
		titleClasses,
		g.Text(item.Title),
	)
	if linkURL != "" && item.NotificationID > 0 {
		titleNode = h.Form(
			h.Method("POST"),
			h.Action(linkURL),
			h.Class("notification-title-form"),
			h.Button(
				h.Type("submit"),
				h.Class("notification-title notification-title-button"),
				g.Text(item.Title),
			),
		)
	} else if linkURL != "" {
		titleNode = h.A(
			titleClasses,
			h.Href(linkURL),
			g.Text(item.Title),
		)
	}

	return h.Li(
		c.Classes{
			"notification-item": true,
			"unread":            item.Unread,
		},
		h.Span(h.Class("notification-unread")),
		h.Div(
			components.NotificationIconClasses(item),
			components.Icon(&components.IconProps{Identifier: iconIdentifier}),
		),
		h.Div(
			h.Class("notification-content"),
			h.Div(
				h.Class("notification-title-row"),
				titleNode,
				notificationBadge(item),
			),
			g.If(item.Summary != "",
				h.P(
					h.Class("notification-summary"),
					g.Text(item.Summary),
				),
			),
		),
		notificationActions(item, p),
	)
}

func notificationOpenURL(item model.NotificationItem) string {
	url := strings.TrimSpace(item.URL)
	if url == "" {
		return "/notifications"
	}
	return url
}

func notificationOpenActionURL(item model.NotificationItem, p *NotificationPageProps) string {
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

func notificationBadge(item model.NotificationItem) g.Node {
	if isMentionNotification(item) {
		return notificationBadgeWithLabel(item, mentionBadgeLabel(item))
	}
	if strings.TrimSpace(item.Reason) == "" {
		return g.Group(nil)
	}
	return notificationBadgeWithLabel(item, item.Reason)
}

func notificationBadgeWithLabel(item model.NotificationItem, label string) g.Node {
	label = strings.TrimSpace(label)
	if label == "" {
		return g.Group(nil)
	}

	badgeType := components.BadgeSecondary
	switch model.NormalizeNotificationReasonType(item.ReasonType) {
	case model.NotificationReasonInfo:
		badgeType = components.BadgeType("info")
	case model.NotificationReasonSuccess:
		badgeType = components.BadgeSuccess
	case model.NotificationReasonWarning:
		badgeType = components.BadgeWarning
	case model.NotificationReasonDanger:
		badgeType = components.BadgeDanger
	}

	return components.Badge(&components.BadgeProps{
		Type: badgeType,
		Size: components.BadgeSm,
	}, g.Text(label))
}

func isMentionNotification(item model.NotificationItem) bool {
	switch strings.ToLower(strings.TrimSpace(item.Category)) {
	case "mention", "mentions":
		return true
	default:
		return false
	}
}

func mentionBadgeLabel(item model.NotificationItem) string {
	if strings.TrimSpace(item.ActorUsername) != "" {
		return "@" + strings.TrimSpace(item.ActorUsername)
	}
	trimmed := strings.TrimSpace(item.Reason)
	if trimmed == "" {
		return "@ Mention"
	}
	if strings.HasPrefix(trimmed, "@") {
		return trimmed
	}
	return "@" + trimmed
}

func notificationActions(item model.NotificationItem, p *NotificationPageProps) g.Node {
	if item.Time == "" && item.NotificationID == 0 {
		return g.Group(nil)
	}

	return h.Div(
		h.Class("notification-actions"),
		g.If(item.Time != "",
			h.Div(
				h.Class("notification-time"),
				g.Text(item.Time),
			),
		),
		notificationToggleForm(item, p),
	)
}

func notificationToggleForm(item model.NotificationItem, p *NotificationPageProps) g.Node {
	if item.NotificationID == 0 {
		return g.Group(nil)
	}

	action := "read"
	label := "Mark as read"
	if !item.Unread {
		action = "unread"
		label = "Mark as unread"
	}

	return h.Form(
		h.Method("POST"),
		h.Action(notificationToggleActionURL(item.NotificationID, action, p)),
		h.Button(
			h.Type("submit"),
			h.Class("notification-toggle"),
			h.Aria("label", label),
			components.Icon(&components.IconProps{
				Identifier: notificationToggleIcon(item),
			}),
		),
	)
}

func notificationToggleIcon(item model.NotificationItem) string {
	if item.Unread {
		return "email-open-outline"
	}
	return "email-outline"
}

func notificationToggleActionURL(notificationID int, action string, p *NotificationPageProps) string {
	redirectURL := notificationsListURL(p)
	query := url.Values{}
	query.Set("Redirect", redirectURL)
	return fmt.Sprintf("/notifications/%d/%s?%s", notificationID, action, query.Encode())
}

func notificationsListURL(p *NotificationPageProps) string {
	query := url.Values{}
	if p.ActiveFilter != "" {
		query.Set("Filter", p.ActiveFilter)
	}
	if p.Page > 0 {
		query.Set("Page", strconv.Itoa(p.Page))
	}
	if p.PageSize > 0 {
		query.Set("PageSize", strconv.Itoa(p.PageSize))
	}
	if len(query) == 0 {
		return "/notifications"
	}
	return "/notifications?" + query.Encode()
}

func notificationsEmptyState(activeFilter string) g.Node {
	filter := strings.ToLower(strings.TrimSpace(activeFilter))
	title := "No unread notifications"
	subtitle := "We'll let you know when something needs your attention."
	if filter == "read" {
		title = "No read notifications"
		subtitle = "Unread updates are still in your inbox."
	}
	return h.Div(
		h.Class("notifications-empty"),
		components.Icon(&components.IconProps{Identifier: "inbox-outline"}),
		h.H3(g.Text(title)),
		h.P(g.Text(subtitle)),
	)
}

func defaultNotificationFilters(unreadCount int) []model.NotificationFilter {
	return []model.NotificationFilter{
		{
			ID:    "unread",
			Label: "Unread",
			Count: unreadCount,
			URL:   "/notifications?Filter=unread",
		},
		{
			ID:    "read",
			Label: "Read",
			Count: 0,
			URL:   "/notifications?Filter=read",
		},
	}
}
