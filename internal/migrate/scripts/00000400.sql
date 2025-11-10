CREATE TABLE resource (
    resource_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    type TEXT NOT NULL,
    reference TEXT UNIQUE NOT NULL,
    service_ownership_team_id INT REFERENCES team(team_id),
    is_archived BOOLEAN NOT NULL DEFAULT FALSE
);


CREATE TABLE resource_service_metric (
    resource_service_metric_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    is_cumulative BOOLEAN NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE
);


CREATE TABLE resource_service (
    resource_service_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    resource_id INT NOT NULL REFERENCES resource(resource_id),
    status TEXT NOT NULL,
    started_at TIMESTAMPTZ DEFAULT NOW(),
    started_by INT REFERENCES app_user(user_id),
    completed_at TIMESTAMPTZ,
    completed_by INT REFERENCES app_user(user_id),
    cancelled_at TIMESTAMPTZ,
    cancelled_by INT REFERENCES app_user(user_id),
    notes TEXT NOT NULL DEFAULT '',
    gallery_id INTEGER REFERENCES gallery(gallery_id) NOT NULL,
    comment_thread_id INTEGER REFERENCES comment_thread(comment_thread_id) NOT NULL
);


CREATE TABLE resource_service_change (
    resource_service_change_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    resource_service_id INT NOT NULL REFERENCES resource_service(resource_service_id) ON DELETE CASCADE,
    change_by INT NOT NULL REFERENCES app_user(user_id),
    change_at TIMESTAMPTZ DEFAULT NOW(),

    notes TEXT,
    started_by INT REFERENCES app_user(user_id),
    completed_by INT REFERENCES app_user(user_id),
    cancelled_by INT REFERENCES app_user(user_id),
    reopened_by INT REFERENCES app_user(user_id)
);


CREATE TABLE resource_service_schedule (
    resource_service_schedule_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    resource_id INT NOT NULL REFERENCES resource(resource_id) ON DELETE CASCADE,
    resource_service_metric_id INT NOT NULL REFERENCES resource_service_metric(resource_service_metric_id),
    threshold NUMERIC NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (resource_id, resource_service_metric_id)
);


CREATE TABLE resource_usage_record (
    resource_usage_record_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    resource_id INT NOT NULL REFERENCES resource(resource_id) ON DELETE CASCADE,
    resource_service_metric_id INT NOT NULL REFERENCES resource_service_metric(resource_service_metric_id),
    value NUMERIC NOT NULL,
    recorded_at TIMESTAMPTZ DEFAULT NOW(),
    closed_by_resource_service_id INT REFERENCES resource_service(resource_service_id)
);

CREATE OR REPLACE VIEW resource_service_current_metric_view AS
SELECT
    ss.resource_id,
    m.resource_service_metric_id,
    m.name AS metric_name,
    CASE
        WHEN m.is_cumulative THEN
            COALESCE(SUM(ru.value), 0)
        ELSE
            COALESCE((
                SELECT ru2.value
                FROM resource_usage_record ru2
                WHERE ru2.resource_id = ss.resource_id
                  AND ru2.resource_service_metric_id = m.resource_service_metric_id
                ORDER BY ru2.recorded_at DESC
                LIMIT 1
            ), 0)
    END AS current_value
FROM resource_service_schedule ss
JOIN resource_service_metric m
  ON m.resource_service_metric_id = ss.resource_service_metric_id
LEFT JOIN resource_usage_record ru
  ON ru.resource_id = ss.resource_id
  AND ru.resource_service_metric_id = m.resource_service_metric_id
  AND ru.closed_by_resource_service_id IS NULL
GROUP BY ss.resource_id, m.resource_service_metric_id, m.name, m.is_cumulative;


CREATE OR REPLACE VIEW resource_view AS
SELECT
    r.resource_id,
    r.type,
    r.reference,
    r.is_archived,
    r.service_ownership_team_id,
    t.team_name AS service_ownership_team_name,
    COALESCE((
        SELECT MAX(completed_at)
        FROM resource_service rs
        WHERE rs.resource_id = r.resource_id
    ), NULL) AS last_serviced_at
FROM resource r
LEFT JOIN team t ON t.team_id = r.service_ownership_team_id;


CREATE VIEW resource_service_view AS
SELECT
    rs.resource_service_id,
    rs.resource_id,
    CASE
        WHEN rs.cancelled_at IS NOT NULL THEN 'Cancelled'
        WHEN rs.completed_at IS NOT NULL THEN 'Completed'
        WHEN rs.started_at IS NOT NULL THEN 'Work In Progress'
        ELSE 'Unknown'
    END AS status,
    rs.started_at,
    rs.started_by,
    su.username AS started_by_username,
    rs.completed_at,
    cu.username AS completed_by_username,
    rs.cancelled_at,
    xu.username AS cancelled_by_username,
    rs.notes,
    rs.gallery_id,
    rs.comment_thread_id,
    r.type,
    r.reference
FROM
    resource_service rs
INNER JOIN resource r ON r.resource_id = rs.resource_id
LEFT JOIN app_user su ON rs.started_by = su.user_id
LEFT JOIN app_user cu ON rs.completed_by = cu.user_id
LEFT JOIN app_user xu ON rs.cancelled_by = xu.user_id;


CREATE OR REPLACE VIEW resource_service_change_view AS
SELECT
    rsc.resource_service_change_id,
	rsc.resource_service_id,
    rsc.change_by,
	change_user.username AS change_by_username,
	rsc.change_at,
    CASE
        WHEN rsc.change_at = MIN(rsc.change_at) OVER (PARTITION BY rsc.resource_service_id)
        THEN true
        ELSE false
    END AS is_creation,
	rsc.notes,
    rsc.started_by,
	su.username AS started_by_username,
    rsc.completed_by,
	cmu.username AS completed_by_username,
    rsc.reopened_by,
	rou.username AS reopened_by_username,
    rsc.cancelled_by,
	cu.username AS cancelled_by_username
FROM
    resource_service_change AS rsc
    INNER JOIN
        app_user AS change_user ON rsc.change_by = change_user.user_id
    LEFT JOIN
        app_user AS su ON rsc.started_by = su.user_id
    LEFT JOIN
        app_user AS cmu ON rsc.completed_by = cmu.user_id
    LEFT JOIN
        app_user AS rou ON rsc.reopened_by = rou.user_id
    LEFT JOIN
        app_user AS cu ON rsc.cancelled_by = cu.user_id;


CREATE OR REPLACE VIEW resource_service_metric_status_view AS
SELECT
    ss.resource_service_schedule_id,
    r.resource_id,
    r.type,
    r.reference,
    r.service_ownership_team_id,
    t.team_name AS service_ownership_team_name,
    m.resource_service_metric_id,
    m.name AS metric_name,
    COALESCE(cmv.current_value, 0) AS current_value,
    ss.threshold,
    CASE
        WHEN ss.threshold > 0 THEN ROUND(COALESCE(cmv.current_value, 0) / ss.threshold, 2)
        ELSE 0
    END AS normalised_value,
    CASE
        WHEN ss.threshold > 0 THEN ROUND((COALESCE(cmv.current_value, 0) / ss.threshold) * 100, 0)
        ELSE 0
    END AS normalised_percentage,
    CASE
        WHEN ss.threshold > 0 AND (COALESCE(cmv.current_value, 0) / ss.threshold) >= 1 THEN TRUE
        ELSE FALSE
    END AS is_due,
    (
        SELECT MAX(ru2.recorded_at)
        FROM resource_usage_record ru2
        WHERE ru2.resource_id = r.resource_id
          AND (ss.resource_service_metric_id IS NULL OR ru2.resource_service_metric_id = ss.resource_service_metric_id)
    ) AS last_recorded_at,
    (
        SELECT MAX(COALESCE(rs.completed_at, rs.started_at))
        FROM resource_service rs
        WHERE rs.resource_id = r.resource_id
    ) AS last_serviced_at,
    (
        SELECT rs2.resource_service_id
        FROM resource_service rs2
        WHERE 
            rs2.resource_id = r.resource_id
            AND
            rs2.status = 'Work In Progress'
        ORDER BY rs2.started_at DESC
        LIMIT 1
    ) AS wip_service_id,
    EXISTS (
        SELECT 1
        FROM resource_service rs3
        WHERE
            rs3.resource_id = r.resource_id
            AND
            rs3.status = 'Work In Progress'
    ) AS has_wip_service,
    ss.is_archived AS schedule_is_archived,
    m.is_archived AS metric_is_archived
FROM resource r
INNER JOIN resource_service_schedule ss
  ON ss.resource_id = r.resource_id
INNER JOIN resource_service_metric m
  ON m.resource_service_metric_id = ss.resource_service_metric_id
INNER JOIN resource_service_current_metric_view cmv
  ON cmv.resource_id = r.resource_id
  AND cmv.resource_service_metric_id = ss.resource_service_metric_id
LEFT JOIN team t ON t.team_id = r.service_ownership_team_id
WHERE
	r.is_archived = FALSE
	AND ss.is_archived = FALSE
	AND m.is_archived = FALSE;
