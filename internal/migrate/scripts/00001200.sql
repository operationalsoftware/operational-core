CREATE TABLE notification (
    notification_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL REFERENCES app_user(user_id) ON DELETE CASCADE,
    actor_user_id INT REFERENCES app_user(user_id),
    category TEXT NOT NULL DEFAULT 'general',
    title TEXT NOT NULL,
    summary TEXT,
    url TEXT,
    reason TEXT,
    reason_type TEXT,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX notification_user_created_idx
    ON notification(user_id, created_at DESC);

CREATE INDEX notification_user_read_idx
    ON notification(user_id, is_read);

CREATE INDEX notification_user_category_idx
    ON notification(user_id, category);

CREATE VIEW notification_view AS
SELECT
    notification.notification_id,
    notification.user_id,
    notification.actor_user_id,
    actor.username AS actor_username,
    notification.category,
    notification.title,
    notification.summary,
    notification.url,
    notification.reason,
    notification.reason_type,
    notification.is_read,
    notification.read_at,
    notification.created_at
FROM
    notification
LEFT JOIN app_user actor ON actor.user_id = notification.actor_user_id;
