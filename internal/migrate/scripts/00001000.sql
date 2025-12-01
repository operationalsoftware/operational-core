-- 00001000.sql: Rename resource service schedule assignments, drop archival flag, and refresh dependent views
DROP VIEW IF EXISTS resource_service_metric_status_view;
DROP VIEW IF EXISTS resource_service_current_metric_view;
DROP VIEW IF EXISTS service_metric_lifetime_total_view;

-- Remove archived assignments now that unassignment deletes rows
DELETE FROM resource_service_schedule WHERE is_archived = TRUE;

ALTER TABLE resource_service_schedule
    RENAME TO service_schedule_assignment;

ALTER TABLE service_schedule_assignment
    RENAME COLUMN resource_service_schedule_id TO service_schedule_assignment_id;

ALTER TABLE service_schedule_assignment
    DROP COLUMN IF EXISTS is_archived;

CREATE OR REPLACE VIEW resource_service_current_metric_view AS
SELECT
    r.resource_id,
    m.resource_service_metric_id,
    m.name AS metric_name,
    CASE
        WHEN m.is_cumulative THEN
            COALESCE(SUM(rmr.value), 0)
        ELSE
            COALESCE((
                SELECT rmr2.value
                FROM resource_metric_recording rmr2
                WHERE rmr2.resource_id = r.resource_id
                  AND rmr2.resource_service_metric_id = m.resource_service_metric_id
                ORDER BY rmr2.recorded_at DESC
                LIMIT 1
            ), 0)
    END AS current_value
FROM resource r
JOIN service_schedule_assignment ssa
  ON ssa.resource_id = r.resource_id
JOIN service_schedule ss
  ON ss.service_schedule_id = ssa.service_schedule_id
JOIN resource_service_metric m
  ON m.resource_service_metric_id = ss.resource_service_metric_id
LEFT JOIN resource_metric_recording rmr
  ON rmr.resource_id = r.resource_id
  AND rmr.resource_service_metric_id = m.resource_service_metric_id
  AND rmr.closed_by_resource_service_id IS NULL
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
        SELECT MAX(rmr2.recorded_at)
        FROM resource_metric_recording rmr2
        WHERE rmr2.resource_id = r.resource_id
          AND rmr2.resource_service_metric_id = m.resource_service_metric_id
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
    ssv.is_archived AS schedule_is_archived,
    ssv.metric_is_archived AS metric_is_archived
FROM resource r
JOIN service_schedule_assignment ssa
  ON ssa.resource_id = r.resource_id
JOIN service_schedule_view ssv
  ON ssv.service_schedule_id = ssa.service_schedule_id
JOIN resource_service_metric m
  ON m.resource_service_metric_id = ssv.resource_service_metric_id
JOIN resource_service_current_metric_view cmv
  ON cmv.resource_id = r.resource_id
  AND cmv.resource_service_metric_id = ssv.resource_service_metric_id
LEFT JOIN team t
  ON t.team_id = r.service_ownership_team_id
WHERE
    r.is_archived = FALSE;

CREATE OR REPLACE VIEW service_metric_lifetime_total_view AS
SELECT
  r.resource_id,
  r.type AS resource_type,
  r.reference,
  m.name AS metric_name,
  COALESCE(SUM(rmr.value), 0) AS lifetime_total
FROM resource r
JOIN service_schedule_assignment ssa
  ON ssa.resource_id = r.resource_id
JOIN service_schedule ss
  ON ss.service_schedule_id = ssa.service_schedule_id
JOIN resource_service_metric m
  ON m.resource_service_metric_id = ss.resource_service_metric_id
LEFT JOIN resource_metric_recording rmr
  ON rmr.resource_id = r.resource_id
  AND rmr.resource_service_metric_id = m.resource_service_metric_id
WHERE
  r.is_archived = FALSE
  AND ss.is_archived = FALSE
  AND m.is_archived = FALSE
  AND m.is_cumulative = TRUE
GROUP BY
  r.resource_id,
  r.type,
  r.reference,
  m.resource_service_metric_id,
  m.name;
