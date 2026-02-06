package model

import (
	"strings"
	"time"
)

const (
	NotificationReasonInfo    = "info"
	NotificationReasonSuccess = "success"
	NotificationReasonWarning = "warning"
	NotificationReasonDanger  = "danger"
)

func NormalizeNotificationReasonType(reasonType string) string {
	trimmed := strings.TrimSpace(strings.ToLower(reasonType))
	switch trimmed {
	case NotificationReasonInfo,
		NotificationReasonSuccess,
		NotificationReasonWarning,
		NotificationReasonDanger:
		return trimmed
	default:
		return NotificationReasonInfo
	}
}

type Notification struct {
	NotificationID int
	UserID         int
	ActorUserID    *int
	ActorUsername  string
	Category       string
	Title          string
	Summary        string
	URL            string
	Reason         string
	ReasonType     string
	IsRead         bool
	ReadAt         *time.Time
	CreatedAt      time.Time
}

type NewNotification struct {
	UserID      int
	ActorUserID *int
	Category    string
	Title       string
	Summary     string
	URL         string
	Reason      string
	ReasonType  string
}

type ListNotificationsQuery struct {
	Filter   string
	Page     int
	PageSize int
}

type NotificationCounts struct {
	Total  int
	Unread int
}

type NotificationFilter struct {
	ID    string
	Label string
	Count int
	URL   string
}

type NotificationGroup struct {
	Title string
	Items []NotificationItem
}

type NotificationItem struct {
	NotificationID int
	ActorUsername  string
	Category       string
	Title          string
	URL            string
	Summary        string
	Time           string
	Reason         string
	ReasonType     string
	Unread         bool
}

type PushSubscription struct {
	Endpoint       string               `json:"endpoint"`
	ExpirationTime *int64               `json:"expirationTime"`
	Keys           PushSubscriptionKeys `json:"keys"`
}

type PushSubscriptionKeys struct {
	P256dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

type PushNotificationPayload struct {
	Type  string `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	URL   string `json:"url,omitempty"`
}
