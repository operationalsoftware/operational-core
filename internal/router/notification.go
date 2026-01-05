package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addNotificationRoutes(
	mux *http.ServeMux,
	notificationService service.NotificationService,
) {
	notificationHandler := handler.NewNotificationHandler(notificationService)

	mux.HandleFunc("GET /notifications", notificationHandler.NotificationsPage)
	mux.HandleFunc("POST /notifications/{id}/read", notificationHandler.MarkRead)
	mux.HandleFunc("POST /notifications/{id}/unread", notificationHandler.MarkUnread)
	mux.HandleFunc("POST /notifications/mark-all-read", notificationHandler.MarkAllRead)
}
