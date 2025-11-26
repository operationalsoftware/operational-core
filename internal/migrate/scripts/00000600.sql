-- 00000600.sql: Rename metric records to resource recordings and refresh dependent views

DROP VIEW IF EXISTS resource_service_metric_status_view;
DROP VIEW IF EXISTS resource_service_current_metric_view;

ALTER TABLE resource_usage_record
    RENAME TO resource_metric_recording;

ALTER TABLE resource_metric_recording
    RENAME COLUMN resource_usage_record_id TO resource_metric_recording_id;

CREATE OR REPLACE VIEW resource_service_current_metric_view AS
SELECT
    ss.resource_id,
    m.resource_service_metric_id,
    m.name AS metric_name,
    CASE
        WHEN m.is_cumulative THEN
            COALESCE(SUM(rmr.value), 0)
        ELSE
            COALESCE((
                SELECT rmr2.value
                FROM resource_metric_recording rmr2
                WHERE rmr2.resource_id = ss.resource_id
                  AND rmr2.resource_service_metric_id = m.resource_service_metric_id
                ORDER BY rmr2.recorded_at DESC
                LIMIT 1
            ), 0)
    END AS current_value
FROM resource_service_schedule ss
JOIN resource_service_metric m
  ON m.resource_service_metric_id = ss.resource_service_metric_id
LEFT JOIN resource_metric_recording rmr
  ON rmr.resource_id = ss.resource_id
  AND rmr.resource_service_metric_id = m.resource_service_metric_id
  AND rmr.closed_by_resource_service_id IS NULL
GROUP BY ss.resource_id, m.resource_service_metric_id, m.name, m.is_cumulative;

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
        SELECT MAX(rmr2.recorded_at)
        FROM resource_metric_recording rmr2
        WHERE rmr2.resource_id = r.resource_id
          AND (ss.resource_service_metric_id IS NULL OR rmr2.resource_service_metric_id = ss.resource_service_metric_id)
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
