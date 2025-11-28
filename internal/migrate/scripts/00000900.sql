-- 00000900.sql: Add service metric lifetime total view for cumulative, non-archived metrics

DROP VIEW IF EXISTS service_metric_lifetime_total_view;

CREATE OR REPLACE VIEW service_metric_lifetime_total_view AS
SELECT
    r.resource_id,
    r.type AS resource_type,
    r.reference,
    m.name AS metric_name,
    COALESCE((
        SELECT SUM(rmr.value)
        FROM resource_metric_recording rmr
        WHERE rmr.resource_id = r.resource_id
          AND rmr.resource_service_metric_id = m.resource_service_metric_id
    ), 0) AS lifetime_total
FROM resource r
JOIN resource_service_schedule ss
  ON ss.resource_id = r.resource_id
JOIN resource_service_metric m
  ON m.resource_service_metric_id = ss.resource_service_metric_id
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
