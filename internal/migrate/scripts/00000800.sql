-- Introduce shared service schedule templates and migrate existing resource schedules
DROP VIEW IF EXISTS resource_service_metric_status_view;
DROP VIEW IF EXISTS resource_service_current_metric_view;
DROP VIEW IF EXISTS service_schedule_view;

-- Preserve existing data
ALTER TABLE resource_service_schedule RENAME TO resource_service_schedule_old;

-- New schedule template table
CREATE TABLE service_schedule (
    service_schedule_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    resource_service_metric_id INT NOT NULL REFERENCES resource_service_metric(resource_service_metric_id),
    threshold DECIMAL(16, 4) NOT NULL,
    name TEXT NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT service_schedule_metric_threshold_key UNIQUE (resource_service_metric_id, threshold)
);

CREATE OR REPLACE VIEW service_schedule_view AS
SELECT
    ss.service_schedule_id,
    ss.name,
    ss.resource_service_metric_id,
    m.name AS metric_name,
    ss.threshold,
    ss.is_archived,
    m.is_archived AS metric_is_archived
FROM service_schedule ss
JOIN resource_service_metric m ON m.resource_service_metric_id = ss.resource_service_metric_id;

-- Seed templates from existing rows (dedupe by metric/threshold)
INSERT INTO service_schedule (
    resource_service_metric_id,
    threshold,
    name,
    is_archived,
    created_at
)
SELECT DISTINCT ON (rss.resource_service_metric_id, rss.threshold)
    rss.resource_service_metric_id,
    rss.threshold,
    COALESCE(m.name, CONCAT('Schedule ', rss.resource_service_metric_id, '-', rss.threshold)) AS name,
    COALESCE(rss.is_archived, FALSE),
    rss.created_at
FROM resource_service_schedule_old rss
JOIN resource_service_metric m ON m.resource_service_metric_id = rss.resource_service_metric_id
ORDER BY rss.resource_service_metric_id, rss.threshold, rss.created_at;

-- Map old schedule IDs to new template IDs
CREATE TEMP TABLE service_schedule_map AS
SELECT
    rss.resource_service_schedule_id AS old_service_schedule_id,
    ss.service_schedule_id AS new_service_schedule_id
FROM resource_service_schedule_old rss
JOIN service_schedule ss
  ON ss.resource_service_metric_id = rss.resource_service_metric_id
 AND ss.threshold = rss.threshold;

-- New resource-to-schedule assignments
CREATE TABLE resource_service_schedule (
    resource_service_schedule_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    resource_id INT NOT NULL REFERENCES resource(resource_id) ON DELETE CASCADE,
    service_schedule_id INT NOT NULL REFERENCES service_schedule(service_schedule_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    is_archived BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO resource_service_schedule (
    resource_id,
    service_schedule_id,
    created_at,
    is_archived
)
SELECT
    rss.resource_id,
    map.new_service_schedule_id,
    rss.created_at,
    COALESCE(rss.is_archived, FALSE)
FROM resource_service_schedule_old rss
JOIN service_schedule_map map
  ON map.old_service_schedule_id = rss.resource_service_schedule_id;

DROP TABLE resource_service_schedule_old;

-- Recreate views to reflect new schema
CREATE OR REPLACE VIEW resource_service_current_metric_view AS
SELECT
    r.resource_id,
    m.resource_service_metric_id,
    m.name AS metric_name,
    CASE
        WHEN m.is_cumulative THEN
            COALESCE(SUM(ru.value), 0)
        ELSE
            COALESCE((
                SELECT ru2.value
                FROM resource_metric_recording ru2
                WHERE ru2.resource_id = r.resource_id
                  AND ru2.resource_service_metric_id = m.resource_service_metric_id
                ORDER BY ru2.recorded_at DESC
                LIMIT 1
            ), 0)
    END AS current_value
FROM resource r
JOIN resource_service_schedule rss
  ON rss.resource_id = r.resource_id
JOIN service_schedule ss
  ON ss.service_schedule_id = rss.service_schedule_id
JOIN resource_service_metric m
  ON m.resource_service_metric_id = ss.resource_service_metric_id
LEFT JOIN resource_metric_recording ru
  ON ru.resource_id = r.resource_id
  AND ru.resource_service_metric_id = m.resource_service_metric_id
  AND ru.closed_by_resource_service_id IS NULL
WHERE
    rss.is_archived = FALSE
GROUP BY
    r.resource_id,
    m.resource_service_metric_id,
    m.name,
    m.is_cumulative;


CREATE OR REPLACE VIEW resource_service_metric_status_view AS
SELECT
    ssv.service_schedule_id,
    ssv.name AS service_schedule_name,
    r.resource_id,
    r.type,
    r.reference,
    r.service_ownership_team_id,
    t.team_name AS service_ownership_team_name,
    ssv.resource_service_metric_id,
    ssv.metric_name,
    COALESCE(cmv.current_value, 0) AS current_value,
    ssv.threshold,
    CASE
        WHEN ssv.threshold > 0 THEN ROUND(COALESCE(cmv.current_value, 0) / ssv.threshold, 2)
        ELSE 0
    END AS normalised_value,
    CASE
        WHEN ssv.threshold > 0 THEN ROUND((COALESCE(cmv.current_value, 0) / ssv.threshold) * 100, 0)
        ELSE 0
    END AS normalised_percentage,
    CASE
        WHEN ssv.threshold > 0
             AND (COALESCE(cmv.current_value, 0) / ssv.threshold) >= 1
        THEN TRUE
        ELSE FALSE
    END AS is_due,
    (
        SELECT MAX(ru2.recorded_at)
        FROM resource_metric_recording ru2
        WHERE ru2.resource_id = r.resource_id
          AND ru2.resource_service_metric_id = m.resource_service_metric_id
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
            AND rs2.status = 'Work In Progress'
        ORDER BY rs2.started_at DESC
        LIMIT 1
    ) AS wip_service_id,
    EXISTS (
        SELECT 1
        FROM resource_service rs3
        WHERE
            rs3.resource_id = r.resource_id
            AND rs3.status = 'Work In Progress'
    ) AS has_wip_service,
    (rss.is_archived OR ssv.is_archived) AS schedule_is_archived,
    ssv.metric_is_archived AS metric_is_archived
FROM resource r
JOIN resource_service_schedule rss
  ON rss.resource_id = r.resource_id
JOIN service_schedule_view ssv
  ON ssv.service_schedule_id = rss.service_schedule_id
JOIN resource_service_metric m
  ON m.resource_service_metric_id = ssv.resource_service_metric_id
JOIN resource_service_current_metric_view cmv
  ON cmv.resource_id = r.resource_id
  AND cmv.resource_service_metric_id = ssv.resource_service_metric_id
LEFT JOIN team t
  ON t.team_id = r.service_ownership_team_id
WHERE
    r.is_archived = FALSE
    AND rss.is_archived = FALSE
    AND ssv.is_archived = FALSE
    AND ssv.metric_is_archived = FALSE;
