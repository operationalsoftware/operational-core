package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"strings"

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
	return s.notificationRepo.CreateNotification(ctx, s.db, notification)
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

	counts, err := s.notificationRepo.CountNotifications(ctx, s.db, userID)
	if err != nil {
		return nil, model.NotificationCounts{}, q, err
	}

	return notifications, counts, q, nil
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID int) error {
	return s.notificationRepo.MarkAllRead(ctx, s.db, userID)
}

func (s *NotificationService) MarkRead(ctx context.Context, userID int, notificationID int) error {
	return s.notificationRepo.MarkRead(ctx, s.db, userID, notificationID)
}

func (s *NotificationService) MarkUnread(ctx context.Context, userID int, notificationID int) error {
	return s.notificationRepo.MarkUnread(ctx, s.db, userID, notificationID)
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
