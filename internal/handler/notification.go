package handler

import (
	"app/internal/components"
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/notificationview"
	"app/pkg/appurl"
	"app/pkg/cookie"
	"app/pkg/env"
	"app/pkg/reqcontext"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) NotificationsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var uv notificationsPageURLVals
	if err := appurl.Unmarshal(r.URL.Query(), &uv); err != nil {
		log.Println("error decoding url values:", err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}
	if uv.Filter == "" {
		uv.Filter = r.URL.Query().Get("filter")
	}

	notifications, counts, normalizedQuery, err := h.notificationService.ListNotifications(r.Context(), ctx.User.UserID, model.ListNotificationsQuery{
		Filter:   uv.Filter,
		Page:     uv.Page,
		PageSize: uv.PageSize,
	})
	if err != nil {
		log.Println("error loading notifications:", err)
		http.Error(w, "Error loading notifications", http.StatusInternalServerError)
		return
	}

	_ = notificationview.NotificationPage(notificationview.NotificationPageProps{
		Ctx:            ctx,
		Filters:        notificationFilters(counts),
		ActiveFilter:   normalizedQuery.Filter,
		Groups:         groupNotifications(notifications),
		UnreadCount:    counts.Unread,
		Page:           normalizedQuery.Page,
		PageSize:       normalizedQuery.PageSize,
		TotalRecords:   notificationTotalRecords(counts, normalizedQuery.Filter),
		VAPIDPublicKey: vapidPublicKeyForEnv(),
		ShowTestSent:   r.URL.Query().Get("TestSent") == "1",
	}).Render(w)
}

func (h *NotificationHandler) NotificationsTray(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	const trayPageSize = 6
	notifications, counts, _, err := h.notificationService.ListNotifications(
		r.Context(),
		ctx.User.UserID,
		model.ListNotificationsQuery{
			Filter:   "unread",
			Page:     1,
			PageSize: trayPageSize,
		},
	)
	if err != nil {
		log.Println("error loading notifications tray:", err)
		http.Error(w, "Error loading notifications", http.StatusInternalServerError)
		return
	}

	items := make([]model.NotificationItem, 0, len(notifications))
	for _, notification := range notifications {
		items = append(items, notificationItem(notification))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = components.NotificationsTray(components.NotificationsTrayProps{
		Items:       items,
		UnreadCount: counts.Unread,
	}).Render(w)
}

func (h *NotificationHandler) NotificationTestPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	_ = notificationview.NotificationTestPage(notificationview.NotificationTestPageProps{
		Ctx: ctx,
	}).Render(w)
}

func (h *NotificationHandler) SendTestNotification(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	message := "This is a test notification."

	notificationID, err := h.notificationService.CreateNotification(r.Context(), model.NewNotification{
		UserID:     ctx.User.UserID,
		Category:   "test",
		Title:      "Test notification",
		Summary:    message,
		URL:        "/notifications/test",
		Reason:     "Test",
		ReasonType: model.NotificationReasonInfo,
	})
	if err != nil {
		log.Println("error creating test notification:", err)
		_ = notificationview.NotificationTestPage(notificationview.NotificationTestPageProps{
			Ctx: ctx,
		}).Render(w)
		return
	}

	targetURL := "/notifications/test"
	pushURL := targetURL
	if notificationID > 0 {
		query := url.Values{}
		query.Set("Redirect", targetURL)
		pushURL = fmt.Sprintf("/notifications/%d?%s", notificationID, query.Encode())
	}

	payload := model.PushNotificationPayload{
		Title:          "Test notification",
		Body:           message,
		URL:            pushURL,
		NotificationID: notificationID,
	}

	if err := h.notificationService.SendPushNotification(r.Context(), ctx.User.UserID, payload, ""); err != nil {
		log.Println("error sending test push notification:", err)
		_ = notificationview.NotificationTestPage(notificationview.NotificationTestPageProps{
			Ctx: ctx,
		}).Render(w)
		return
	}

	http.Redirect(w, r, "/notifications?TestSent=1", http.StatusSeeOther)
}

func (h *NotificationHandler) SavePushSubscription(w http.ResponseWriter, r *http.Request) {
	if env.IsProduction() {
		http.NotFound(w, r)
		return
	}

	ctx := reqcontext.GetContext(r)

	var subscription model.PushSubscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	subscription.Endpoint = strings.TrimSpace(subscription.Endpoint)
	subscription.Keys.P256dh = strings.TrimSpace(subscription.Keys.P256dh)
	subscription.Keys.Auth = strings.TrimSpace(subscription.Keys.Auth)

	if subscription.Endpoint == "" || subscription.Keys.P256dh == "" || subscription.Keys.Auth == "" {
		http.Error(w, "Missing subscription fields", http.StatusBadRequest)
		return
	}

	if err := h.notificationService.SavePushSubscription(r.Context(), ctx.User.UserID, subscription); err != nil {
		log.Println("error saving push subscription:", err)
		http.Error(w, "Error saving subscription", http.StatusInternalServerError)
		return
	}

	if err := cookie.SetPushSubscriptionCookie(w, subscription.Endpoint); err != nil {
		log.Println("error setting push subscription cookie:", err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) DeletePushSubscription(w http.ResponseWriter, r *http.Request) {
	if env.IsProduction() {
		http.NotFound(w, r)
		return
	}

	ctx := reqcontext.GetContext(r)

	endpoint, err := cookie.GetPushSubscriptionEndpoint(r)
	if err == nil && endpoint != "" {
		if err := h.notificationService.DeletePushSubscription(r.Context(), ctx.User.UserID, endpoint); err != nil {
			log.Println("error deleting push subscription:", err)
			http.Error(w, "Error deleting subscription", http.StatusInternalServerError)
			return
		}
	}

	cookie.ClearPushSubscriptionCookie(w)

	w.WriteHeader(http.StatusNoContent)
}

func vapidPublicKeyForEnv() string {
	if env.IsProduction() {
		return ""
	}
	return strings.TrimSpace(os.Getenv("VAPID_PUBLIC_KEY"))
}

type notificationsPageURLVals struct {
	Filter   string
	Page     int
	PageSize int
}

func notificationFilters(counts model.NotificationCounts) []model.NotificationFilter {
	readCount := max(counts.Total-counts.Unread, 0)

	return []model.NotificationFilter{
		{
			ID:    "unread",
			Label: "Unread",
			Count: counts.Unread,
			URL:   "/notifications?Filter=unread",
		},
		{
			ID:    "read",
			Label: "Read",
			Count: readCount,
			URL:   "/notifications?Filter=read",
		},
	}
}

func notificationTotalRecords(counts model.NotificationCounts, filter string) int {
	readCount := max(counts.Total-counts.Unread, 0)

	switch strings.ToLower(filter) {
	case "read":
		return readCount
	default:
		return counts.Unread
	}
}

func groupNotifications(notifications []model.Notification) []model.NotificationGroup {
	grouped := map[string][]model.NotificationItem{}

	for _, notification := range notifications {
		label := notificationGroupLabel(notification.CreatedAt)
		grouped[label] = append(grouped[label], notificationItem(notification))
	}

	var groups []model.NotificationGroup
	for _, label := range []string{"Today", "Yesterday", "Earlier"} {
		items := grouped[label]
		if len(items) == 0 {
			continue
		}
		groups = append(groups, model.NotificationGroup{
			Title: label,
			Items: items,
		})
	}

	return groups
}

func notificationGroupLabel(createdAt time.Time) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	if createdAt.After(today) {
		return "Today"
	}
	if createdAt.After(yesterday) {
		return "Yesterday"
	}
	return "Earlier"
}

func notificationItem(notification model.Notification) model.NotificationItem {
	return model.NotificationItem{
		NotificationID: notification.NotificationID,
		ActorUsername:  notification.ActorUsername,
		Category:       notification.Category,
		Title:          notification.Title,
		URL:            notification.URL,
		Summary:        notification.Summary,
		Time:           formatNotificationTime(notification.CreatedAt),
		Reason:         notification.Reason,
		ReasonType:     model.NormalizeNotificationReasonType(notification.ReasonType),
		Unread:         !notification.IsRead,
	}
}

func formatNotificationTime(createdAt time.Time) string {
	diff := time.Since(createdAt)
	if diff < time.Minute {
		return "Just now"
	}
	if diff < time.Hour {
		return fmt.Sprintf("%dm ago", int(diff.Minutes()))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	}
	if diff < 7*24*time.Hour {
		return fmt.Sprintf("%dd ago", int(diff.Hours()/24))
	}
	return createdAt.Format("Jan 2, 2006")
}

func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if err := h.notificationService.MarkAllRead(r.Context(), ctx.User.UserID); err != nil {
		log.Println("error marking notifications as read:", err)
		http.Error(w, "Error marking notifications as read", http.StatusInternalServerError)
		return
	}
	if pushErr := h.notificationService.SendPushNotification(
		r.Context(),
		ctx.User.UserID,
		model.PushNotificationPayload{Type: "notification_read"},
		"",
	); pushErr != nil {
		log.Println("error refreshing notification tray:", pushErr)
	}
	http.Redirect(w, r, "/notifications?Filter=unread", http.StatusSeeOther)
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	notificationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || notificationID <= 0 {
		log.Println("invalid notification id:", err)
		http.Error(w, "Invalid notification id", http.StatusBadRequest)
		return
	}

	err = h.notificationService.MarkRead(r.Context(), ctx.User.UserID, notificationID)
	if err != nil {
		log.Println("error updating notification:", err)
		http.Error(w, "Error updating notification", http.StatusInternalServerError)
		return
	}

	if pushErr := h.notificationService.SendPushNotification(
		r.Context(),
		ctx.User.UserID,
		model.PushNotificationPayload{
			Type:           "notification_read",
			NotificationID: notificationID,
		},
		"",
	); pushErr != nil {
		log.Println("error refreshing notification tray:", pushErr)
	}

	redirectURL := notificationRedirect(r.URL.Query().Get("Redirect"))
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *NotificationHandler) MarkUnread(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	notificationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || notificationID <= 0 {
		log.Println("invalid notification id:", err)
		http.Error(w, "Invalid notification id", http.StatusBadRequest)
		return
	}

	err = h.notificationService.MarkUnread(r.Context(), ctx.User.UserID, notificationID)
	if err != nil {
		log.Println("error updating notification:", err)
		http.Error(w, "Error updating notification", http.StatusInternalServerError)
		return
	}

	redirectURL := notificationRedirect(r.URL.Query().Get("Redirect"))
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *NotificationHandler) OpenNotification(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	notificationID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || notificationID <= 0 {
		log.Println("invalid notification id:", err)
		http.Error(w, "Invalid notification id", http.StatusBadRequest)
		return
	}

	storedURL, found, err := h.notificationService.GetNotificationURL(r.Context(), ctx.User.UserID, notificationID)
	if err != nil {
		log.Println("error fetching notification url:", err)
		http.Error(w, "Error fetching notification", http.StatusInternalServerError)
		return
	}
	if !found {
		http.NotFound(w, r)
		return
	}

	if err := h.notificationService.MarkRead(r.Context(), ctx.User.UserID, notificationID); err != nil {
		log.Println("error marking notification read:", err)
	}

	if pushErr := h.notificationService.SendPushNotification(
		r.Context(),
		ctx.User.UserID,
		model.PushNotificationPayload{
			Type:           "notification_read",
			NotificationID: notificationID,
		},
		"",
	); pushErr != nil {
		log.Println("error refreshing notification tray:", pushErr)
	}

	redirectParam := strings.TrimSpace(r.URL.Query().Get("Redirect"))
	target := strings.TrimSpace(storedURL)
	if redirectParam != "" && (target == "" || redirectParam == target) {
		target = redirectParam
	}
	if target == "" {
		target = "/notifications"
	}

	http.Redirect(w, r, notificationRedirect(target), http.StatusSeeOther)
}

func notificationRedirect(redirect string) string {
	if redirect == "" {
		return "/notifications"
	}
	if strings.HasPrefix(redirect, "/") && !strings.HasPrefix(redirect, "//") {
		return redirect
	}
	return "/notifications"
}
