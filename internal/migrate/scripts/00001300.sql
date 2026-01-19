-- 00001300.sql: Add require_acknowledgement to andon_issue and update andon status logic

ALTER TABLE andon_issue
ADD COLUMN require_acknowledgement BOOLEAN NOT NULL DEFAULT TRUE;

DROP VIEW IF EXISTS andon_view;
DROP VIEW IF EXISTS andon_issue_group_view;
DROP VIEW IF EXISTS andon_issue_view;
DROP VIEW IF EXISTS andon_issue_tree_view;

CREATE VIEW andon_issue_tree_view AS
WITH RECURSIVE andon_issue_tree AS (
    SELECT
        ai.andon_issue_id,
        ai.issue_name,
        ai.parent_id,
        ARRAY[ai.issue_name] AS name_path,
        1 AS depth,
        ai.is_group,
        ai.is_archived,
        (
            SELECT COUNT(*)
            FROM andon_issue c
            WHERE c.parent_id = ai.andon_issue_id
        ) AS children_count
    FROM andon_issue ai
    WHERE ai.parent_id IS NULL

    UNION ALL

    SELECT
        child.andon_issue_id,
        child.issue_name,
        child.parent_id,
        parent.name_path || child.issue_name,
        parent.depth + 1,
        child.is_group,
        child.is_archived,
        (
            SELECT COUNT(*)
            FROM andon_issue c
            WHERE c.parent_id = child.andon_issue_id
        ) AS children_count
    FROM andon_issue child
    JOIN andon_issue_tree parent ON child.parent_id = parent.andon_issue_id
    WHERE parent.is_group = TRUE
),
 down_depths AS (
    SELECT
        g.andon_issue_id,

        -- Downward depth
        (
            SELECT COALESCE(MAX(depth) - 1, 0)
            FROM (
                WITH RECURSIVE downward AS (
                    SELECT andon_issue_id, parent_id, 1 AS depth
                    FROM andon_issue
                    WHERE andon_issue_id = g.andon_issue_id

                    UNION ALL

                    SELECT ai.andon_issue_id, ai.parent_id, d.depth + 1
                    FROM andon_issue ai
                    JOIN downward d ON ai.parent_id = d.andon_issue_id
                )
                SELECT * FROM downward
            ) AS down_sub
        ) AS down_depth
    FROM andon_issue g
 )
SELECT
    ait.andon_issue_id,
    ait.issue_name,
    ait.parent_id,
    ait.name_path,
    ait.depth,
    ait.is_group,
    ait.is_archived,
    ait.children_count,
    ai.severity,
    ai.require_acknowledgement,
    ai.assigned_team,
    t.team_name AS assigned_team_name,
    ai.created_at,
    ai.created_by,
    cu.username AS created_by_username,
    ai.updated_at,
    ai.updated_by,
    uu.username AS updated_by_username,
    COALESCE(dd.down_depth, 0) + 1 AS down_depth
FROM
    andon_issue_tree ait
    INNER JOIN andon_issue ai USING(andon_issue_id)
    LEFT JOIN team t ON t.team_id = ai.assigned_team
    INNER JOIN app_user cu ON cu.user_id = ai.created_by
    LEFT JOIN app_user uu ON uu.user_id = ai.updated_by
    LEFT JOIN down_depths dd ON dd.andon_issue_id = ait.andon_issue_id;

CREATE VIEW andon_issue_view AS
SELECT
    andon_issue_id,
    issue_name,
    parent_id,
    name_path,
    depth,
    is_archived,
    severity,
    require_acknowledgement,
    assigned_team,
    assigned_team_name,
    created_at,
    created_by,
    created_by_username,
    updated_at,
    updated_by,
    updated_by_username
FROM andon_issue_tree_view
WHERE is_group = false;

CREATE VIEW andon_issue_group_view AS
SELECT
    andon_issue_id,
    issue_name,
    parent_id,
    name_path,
    depth,
    children_count,
    is_archived,
    is_group,
    down_depth
FROM andon_issue_tree_view
WHERE is_group = true;

CREATE VIEW andon_view AS
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
    aiv.require_acknowledgement,
    (a.acknowledged_at IS NOT NULL) AS is_acknowledged,
    (a.resolved_at IS NOT NULL) AS is_resolved,
    (a.cancelled_at IS NOT NULL) AS is_cancelled,
    CASE
      WHEN a.cancelled_at IS NOT NULL THEN false
      WHEN aiv.require_acknowledgement = false THEN
        CASE
          WHEN aiv.severity = 'Info' THEN false
          WHEN a.resolved_at IS NOT NULL THEN false
          ELSE true
        END
      WHEN aiv.severity = 'Info' AND a.acknowledged_at IS NOT NULL THEN false
      WHEN aiv.severity IN ('Self-resolvable', 'Requires Intervention')
           AND a.acknowledged_at IS NOT NULL
           AND a.resolved_at IS NOT NULL
      THEN false
      ELSE true
    END AS is_open,
    CASE
      WHEN a.cancelled_at IS NOT NULL THEN 'Cancelled'
      WHEN aiv.require_acknowledgement = false THEN
        CASE
          WHEN aiv.severity = 'Info' THEN 'Closed'
          WHEN a.resolved_at IS NOT NULL THEN 'Closed'
          WHEN aiv.severity = 'Self-resolvable' THEN 'Work In Progress'
          WHEN aiv.severity = 'Requires Intervention' THEN 'Outstanding'
          ELSE 'Invalid Status'
        END

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
      WHEN aiv.require_acknowledgement = false THEN
        CASE
          WHEN aiv.severity = 'Info' THEN a.raised_at
          WHEN a.resolved_at IS NOT NULL THEN a.resolved_at
          ELSE NULL
        END
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
  CASE
    WHEN base.resolved_at IS NULL OR base.severity = 'Info' THEN NULL
    ELSE EXTRACT(EPOCH FROM (base.resolved_at - base.raised_at))::bigint
  END AS downtime_duration_seconds,
  EXTRACT(
    EPOCH FROM (COALESCE(base.closed_at, NOW()) - base.raised_at)
  )::bigint AS open_duration_seconds
FROM base;
