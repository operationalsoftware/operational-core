-- 00000700.sql: Track last_active timestamps on users
-- Adds a nullable last_active column and refreshes user_view to expose it.

ALTER TABLE app_user
  ADD COLUMN IF NOT EXISTS last_active TIMESTAMPTZ;

DROP VIEW IF EXISTS user_view;

CREATE VIEW user_view AS
SELECT
    u.user_id,
    u.is_api_user,
    u.username,
    u.email,
    u.first_name,
    u.last_name,
    u.created,
    u.last_login,
    u.last_active,
    u.session_duration_minutes,
    u.permissions,
    COALESCE(json_agg(
        json_build_object(
            'team_id', t.team_id,
            'team_name', t.team_name,
            'role', ut.role
        )
    ) FILTER (WHERE t.team_id IS NOT NULL), '[]') AS teams
FROM
    app_user u
LEFT JOIN user_team ut ON u.user_id = ut.user_id
LEFT JOIN team t ON ut.team_id = t.team_id
GROUP BY u.user_id;
