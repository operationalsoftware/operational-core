package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (r *NotificationRepository) CreateNotification(
	ctx context.Context,
	exec db.PGExecutor,
	notification model.NewNotification,
) (int, error) {
	query := `
INSERT INTO notification (
	user_id,
	actor_user_id,
	category,
	title,
	summary,
	url,
	reason,
	reason_type
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING notification_id
`

	var id int
	err := exec.QueryRow(
		ctx,
		query,
		notification.UserID,
		notification.ActorUserID,
		notification.Category,
		notification.Title,
		notification.Summary,
		notification.URL,
		notification.Reason,
		notification.ReasonType,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *NotificationRepository) ListNotifications(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
	q model.ListNotificationsQuery,
) ([]model.Notification, error) {
	limit := q.PageSize
	if limit <= 0 {
		limit = 50
	}
	page := q.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	whereClause, args := generateNotificationWhereClause(userID, q.Filter)
	args = append(args, limit, offset)

	query := `
SELECT
	notification_id,
	user_id,
	actor_user_id,
	actor_username,
	category,
	title,
	summary,
	url,
	reason,
	reason_type,
	is_read,
	read_at,
	created_at
FROM
	notification_view
` + whereClause + `
ORDER BY created_at DESC
LIMIT $` + strconv.Itoa(len(args)-1) + ` OFFSET $` + strconv.Itoa(len(args)) + `
`

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		var actorUserID *int
		var actorUsername *string
		var readAt *time.Time

		err = rows.Scan(
			&notification.NotificationID,
			&notification.UserID,
			&actorUserID,
			&actorUsername,
			&notification.Category,
			&notification.Title,
			&notification.Summary,
			&notification.URL,
			&notification.Reason,
			&notification.ReasonType,
			&notification.IsRead,
			&readAt,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		notification.ActorUserID = actorUserID
		if actorUsername != nil {
			notification.ActorUsername = *actorUsername
		}
		notification.ReadAt = readAt
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *NotificationRepository) Count(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
) (model.NotificationCounts, error) {
	query := `
SELECT
	COUNT(*) AS total_count,
	COUNT(*) FILTER (WHERE is_read = FALSE) AS unread_count
FROM
	notification
WHERE
	user_id = $1
`

	var counts model.NotificationCounts
	err := exec.QueryRow(ctx, query, userID).Scan(&counts.Total, &counts.Unread)
	if err != nil {
		return model.NotificationCounts{}, err
	}

	return counts, nil
}

func (r *NotificationRepository) GetNotificationURL(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
	notificationID int,
) (string, error) {
	query := `
SELECT
	url
FROM
	notification
WHERE
	user_id = $1
	AND notification_id = $2
`

	var url string
	err := exec.QueryRow(ctx, query, userID, notificationID).Scan(&url)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return url, nil
}

func (r *NotificationRepository) MarkAllRead(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
) error {
	query := `
UPDATE notification
SET
	is_read = TRUE,
	read_at = NOW()
WHERE
	user_id = $1
	AND is_read = FALSE
`

	_, err := exec.Exec(ctx, query, userID)
	return err
}

func (r *NotificationRepository) MarkRead(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
	notificationID int,
) error {
	query := `
UPDATE notification
SET
	is_read = TRUE,
	read_at = NOW()
WHERE
	user_id = $1
	AND notification_id = $2
`

	_, err := exec.Exec(ctx, query, userID, notificationID)
	return err
}

func (r *NotificationRepository) MarkUnread(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
	notificationID int,
) error {
	query := `
UPDATE notification
SET
	is_read = FALSE,
	read_at = NULL
WHERE
	user_id = $1
	AND notification_id = $2
`

	_, err := exec.Exec(ctx, query, userID, notificationID)
	return err
}

func (r *NotificationRepository) UpsertPushSubscription(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
	subscription model.PushSubscription,
) error {
	query := `
INSERT INTO notification_subscription (
	user_id,
	endpoint,
	p256dh,
	auth,
	created_at,
	updated_at
)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (user_id, endpoint)
DO UPDATE SET
	p256dh = EXCLUDED.p256dh,
	auth = EXCLUDED.auth,
	updated_at = NOW()
`

	_, err := exec.Exec(
		ctx,
		query,
		userID,
		subscription.Endpoint,
		subscription.Keys.P256dh,
		subscription.Keys.Auth,
	)
	return err
}

func (r *NotificationRepository) ListPushSubscriptions(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
) ([]model.PushSubscription, error) {
	query := `
SELECT
	endpoint,
	p256dh,
	auth
FROM
	notification_subscription
WHERE
	user_id = $1
`

	rows, err := exec.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []model.PushSubscription
	for rows.Next() {
		var sub model.PushSubscription
		if err := rows.Scan(&sub.Endpoint, &sub.Keys.P256dh, &sub.Keys.Auth); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *NotificationRepository) DeletePushSubscription(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
	endpoint string,
) error {
	query := `
DELETE FROM notification_subscription
WHERE
	user_id = $1
	AND endpoint = $2
`

	_, err := exec.Exec(ctx, query, userID, endpoint)
	return err
}

func generateNotificationWhereClause(userID int, filter string) (string, []any) {
	where := []string{"notification_view.user_id = $1"}
	args := []any{userID}

	switch strings.ToLower(filter) {
	case "unread":
		where = append(where, "notification_view.is_read = FALSE")
	case "read":
		where = append(where, "notification_view.is_read = TRUE")
	}

	return "WHERE " + strings.Join(where, " AND "), args
}
