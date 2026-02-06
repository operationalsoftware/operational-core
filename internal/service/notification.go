package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationService struct {
	db                *pgxpool.Pool
	notificationRepo  *repository.NotificationRepository
	defaultPageSize   int
	defaultFilterName string
}

func NewNotificationService(
	db *pgxpool.Pool,
	notificationRepository *repository.NotificationRepository,
) *NotificationService {
	return &NotificationService{
		db:                db,
		notificationRepo:  notificationRepository,
		defaultPageSize:   50,
		defaultFilterName: "unread",
	}
}

func (s *NotificationService) CreateNotification(
	ctx context.Context,
	notification model.NewNotification,
) (int, error) {
	notification.Category = strings.TrimSpace(strings.ToLower(notification.Category))
	if notification.Category == "" {
		notification.Category = "general"
	}
	notification.ReasonType = model.NormalizeNotificationReasonType(notification.ReasonType)
	newID, err := s.notificationRepo.CreateNotification(ctx, s.db, notification)
	if err != nil {
		return 0, fmt.Errorf("error creating notification: %v", err)
	}

	return newID, nil
}

func (s *NotificationService) ListNotifications(
	ctx context.Context,
	userID int,
	q model.ListNotificationsQuery,
) ([]model.Notification, model.NotificationCounts, model.ListNotificationsQuery, error) {
	q = s.normalizeQuery(q)

	notifications, err := s.notificationRepo.ListNotifications(ctx, s.db, userID, q)
	if err != nil {
		return nil, model.NotificationCounts{}, q, err
	}
	for i := range notifications {
		notifications[i].ReasonType = model.NormalizeNotificationReasonType(notifications[i].ReasonType)
	}

	counts, err := s.notificationRepo.Count(ctx, s.db, userID)
	if err != nil {
		return nil, model.NotificationCounts{}, q, err
	}

	return notifications, counts, q, nil
}

func (s *NotificationService) GetNotificationURL(
	ctx context.Context,
	userID int,
	notificationID int,
) (string, bool, error) {
	url, found, err := s.notificationRepo.GetNotificationURL(ctx, s.db, userID, notificationID)
	if err != nil {
		return "", false, fmt.Errorf("error getting notification URL: %v", err)
	}
	return url, found, nil
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID int) error {
	err := s.notificationRepo.MarkAllRead(ctx, s.db, userID)
	if err != nil {
		return fmt.Errorf("error marking all notifications as read: %v", err)
	}
	return nil
}

func (s *NotificationService) MarkRead(ctx context.Context, userID int, notificationID int) error {
	err := s.notificationRepo.MarkRead(ctx, s.db, userID, notificationID)
	if err != nil {
		return fmt.Errorf("error marking notification as read: %v", err)
	}
	return nil
}

func (s *NotificationService) MarkUnread(ctx context.Context, userID int, notificationID int) error {
	err := s.notificationRepo.MarkUnread(ctx, s.db, userID, notificationID)
	if err != nil {
		return fmt.Errorf("error marking notification as unread: %w", err)
	}
	return nil
}

func (s *NotificationService) SavePushSubscription(
	ctx context.Context,
	userID int,
	subscription model.PushSubscription,
) error {
	err := s.notificationRepo.UpsertPushSubscription(ctx, s.db, userID, subscription)
	if err != nil {
		return fmt.Errorf("error saving push subscription: %v", err)
	}
	return nil
}

func (s *NotificationService) DeletePushSubscription(
	ctx context.Context,
	userID int,
	endpoint string,
) error {
	endpoint = strings.TrimSpace(endpoint)
	if userID == 0 || endpoint == "" {
		return nil
	}
	err := s.notificationRepo.DeletePushSubscription(ctx, s.db, userID, endpoint)
	if err != nil {
		return fmt.Errorf("error deleting push subscription: %v", err)
	}
	return nil
}

func (s *NotificationService) SendPushNotification(
	ctx context.Context,
	userID int,
	payload model.PushNotificationPayload,
	excludeEndpoint string,
) error {
	vapidPublicKey := strings.TrimSpace(os.Getenv("VAPID_PUBLIC_KEY"))
	vapidPrivateKey := strings.TrimSpace(os.Getenv("VAPID_PRIVATE_KEY"))
	if vapidPublicKey == "" || vapidPrivateKey == "" {
		return errors.New("missing VAPID keys")
	}

	subject := strings.TrimSpace(os.Getenv("VAPID_SUBJECT"))
	if strings.HasPrefix(strings.ToLower(subject), "mailto:") {
		subject = strings.TrimSpace(strings.TrimPrefix(subject, "mailto:"))
	}
	if subject == "" {
		subject = "notifications@localhost"
	}

	subscriptions, err := s.notificationRepo.ListPushSubscriptions(ctx, s.db, userID)
	if err != nil {
		return err
	}
	if len(subscriptions) == 0 {
		return nil
	}

	message, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	excludeEndpoint = strings.TrimSpace(excludeEndpoint)

	var firstErr error
	for _, subscription := range subscriptions {
		if excludeEndpoint != "" && subscription.Endpoint == excludeEndpoint {
			continue
		}
		resp, err := webpush.SendNotification(
			message,
			&webpush.Subscription{
				Endpoint: subscription.Endpoint,
				Keys: webpush.Keys{
					P256dh: subscription.Keys.P256dh,
					Auth:   subscription.Keys.Auth,
				},
			},
			&webpush.Options{
				Subscriber:      subject,
				VAPIDPublicKey:  vapidPublicKey,
				VAPIDPrivateKey: vapidPrivateKey,
				TTL:             300,
			},
		)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if resp != nil {
			if resp.StatusCode == http.StatusGone || resp.StatusCode == http.StatusNotFound {
				if deleteErr := s.notificationRepo.DeletePushSubscription(ctx, s.db, userID, subscription.Endpoint); deleteErr != nil {
					log.Println("failed to delete push subscription:", deleteErr)
				}
			}
			_ = resp.Body.Close()
		}
	}

	return firstErr
}

func (s *NotificationService) normalizeQuery(q model.ListNotificationsQuery) model.ListNotificationsQuery {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = s.defaultPageSize
	}
	filter := strings.ToLower(strings.TrimSpace(q.Filter))
	if filter != "read" && filter != "unread" {
		q.Filter = s.defaultFilterName
	} else {
		q.Filter = filter
	}
	return q
}
