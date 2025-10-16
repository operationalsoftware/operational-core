-- 00000300.sql: Add closed_at and open_duration_seconds to andon_view
-- This migration centralizes the computation of when an Andon is considered closed
-- and exposes a reusable open_duration_seconds derived from raised_at -> closed_at (or now if open).

CREATE OR REPLACE VIEW andon_view AS
WITH base AS (
  SELECT
    a.andon_id,
    a.description,
    a.andon_issue_id,
    a.gallery_id,
    a.comment_thread_id,
    a.source,
    a.location,
    a.raised_at,
    a.raised_by,
    a.acknowledged_at,
    a.resolved_at,
    a.cancelled_at,
    a.last_updated,
    aiv.issue_name,
    aiv.assigned_team,
    aiv.assigned_team_name,
    aiv.name_path,
    u.username AS raised_by_username,
    acku.username AS acknowledged_by_username,
    ru.username AS resolved_by_username,
    cu.username AS cancelled_by_username,
    aiv.severity,
    (a.acknowledged_at IS NOT NULL) AS is_acknowledged,
    (a.resolved_at IS NOT NULL) AS is_resolved,
    (a.cancelled_at IS NOT NULL) AS is_cancelled,
    CASE
      WHEN a.cancelled_at IS NOT NULL THEN false
      WHEN aiv.severity = 'Info' AND a.acknowledged_at IS NOT NULL THEN false
      WHEN aiv.severity IN ('Self-resolvable', 'Requires Intervention')
           AND a.acknowledged_at IS NOT NULL
           AND a.resolved_at IS NOT NULL
      THEN false
      ELSE true
    END AS is_open,
    CASE
      WHEN a.cancelled_at IS NOT NULL THEN 'Cancelled'

      -- Info
      WHEN aiv.severity = 'Info' AND a.acknowledged_at IS NOT NULL THEN 'Closed'
      WHEN aiv.severity = 'Info' AND a.acknowledged_at IS NULL THEN 'Requires Acknowledgement'

      -- Self-resolvable
      WHEN aiv.severity = 'Self-resolvable' AND a.acknowledged_at IS NOT NULL AND a.resolved_at IS NOT NULL THEN 'Closed'
      -- Self-resolvable andons are considered WIP immediately upon creation
      WHEN aiv.severity = 'Self-resolvable' AND a.resolved_at IS NULL THEN 'Work In Progress'
      WHEN aiv.severity = 'Self-resolvable' AND a.acknowledged_at IS NULL THEN 'Requires Acknowledgement'

      -- Requires Intervention
      WHEN aiv.severity = 'Requires Intervention' AND a.acknowledged_at IS NOT NULL AND a.resolved_at IS NOT NULL THEN 'Closed'
      WHEN aiv.severity = 'Requires Intervention' AND a.acknowledged_at IS NULL AND a.resolved_at IS NULL THEN 'Outstanding'
      WHEN aiv.severity = 'Requires Intervention' AND a.acknowledged_at IS NOT NULL THEN 'Work In Progress'
      WHEN aiv.severity = 'Requires Intervention' AND a.resolved_at IS NOT NULL THEN 'Requires Acknowledgement'

      ELSE 'Invalid Status'
    END AS status,
    -- closed_at follows our severity-driven close rules
    CASE
      WHEN a.cancelled_at IS NOT NULL THEN a.cancelled_at
      WHEN aiv.severity = 'Info' AND a.acknowledged_at IS NOT NULL THEN a.acknowledged_at
      WHEN aiv.severity IN ('Self-resolvable', 'Requires Intervention')
           AND a.acknowledged_at IS NOT NULL AND a.resolved_at IS NOT NULL
      THEN GREATEST(a.acknowledged_at, a.resolved_at)
      ELSE NULL
    END AS closed_at
  FROM andon a
  INNER JOIN app_user u ON a.raised_by = u.user_id
  LEFT JOIN app_user acku ON a.acknowledged_by = acku.user_id
  LEFT JOIN app_user ru ON a.resolved_by = ru.user_id
  LEFT JOIN app_user cu ON a.cancelled_by = cu.user_id
  INNER JOIN andon_issue_view aiv ON a.andon_issue_id = aiv.andon_issue_id
)
SELECT
  base.*,
  GREATEST(
    EXTRACT(EPOCH FROM (COALESCE(base.closed_at, NOW()) - base.raised_at)),
    0
  )::bigint AS open_duration_seconds
FROM base;
