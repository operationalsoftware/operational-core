package router

import (
	"app/internal/handler"
	"app/internal/service"
	"app/pkg/env"
	"net/http"
)

func addNotificationRoutes(
	mux *http.ServeMux,
	notificationService service.NotificationService,
) {
	if env.IsProduction() {
		return
	}

	notificationHandler := handler.NewNotificationHandler(notificationService)

	mux.HandleFunc("GET /notifications", notificationHandler.NotificationsPage)
	mux.HandleFunc("GET /notifications/tray", notificationHandler.NotificationsTray)
	mux.HandleFunc("POST /notifications/subscriptions", notificationHandler.SavePushSubscription)
	mux.HandleFunc("POST /notifications/subscriptions/delete", notificationHandler.DeletePushSubscription)
	mux.HandleFunc("POST /notifications/mark-all-read", notificationHandler.MarkAllRead)
	mux.HandleFunc("POST /notifications/{id}/read", notificationHandler.MarkRead)
	mux.HandleFunc("POST /notifications/{id}/unread", notificationHandler.MarkUnread)
}
