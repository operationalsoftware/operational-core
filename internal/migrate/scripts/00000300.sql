-- Migration 00000300
-- Goal:
--   * Introduce dedicated comment threads decoupled from (entity, entity_id) pairs
--   * Persist a comment_thread_id on domain tables: andon, stock_item
--   * Backfill existing comments (entity IN ('Andon','StockItem')) into threads
--   * Remove legacy coupling: drop entity/entity_id from comment & comment_thread
--   * Recreate comment_view using only comment_thread_id
-- Assumptions:
--   * Existing comments only reference entities: 'Andon' or 'StockItem'.
--   * If other entities exist, migration aborts (to avoid silent data loss).
-- Safety:
--   * All work done in a transaction.
--   * Validation blocks ensure no null thread ids remain.

BEGIN;

-- 1. Drop dependent views so we can modify underlying tables
DROP VIEW IF EXISTS comment_view;
DROP VIEW IF EXISTS andon_view; -- will be recreated to expose comment_thread_id

-- 2. Create comment_thread table (temporarily includes entity, entity_id for backfill mapping)
CREATE TABLE comment_thread (
  comment_thread_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  entity TEXT NOT NULL,
  entity_id INT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  UNIQUE(entity, entity_id)
);

-- 3. Validation: ensure only supported entities appear in existing comments
DO $$
DECLARE bad_entity TEXT;
BEGIN
  SELECT c.entity INTO bad_entity
  FROM comment c
  WHERE c.entity NOT IN ('Andon','StockItem')
  LIMIT 1;
  IF bad_entity IS NOT NULL THEN
  RAISE EXCEPTION 'Unsupported entity in comment table: %', bad_entity
    USING HINT = 'Migration aborted to avoid data loss.';
  END IF;
END $$;

-- 4. Insert one thread per existing Andon referenced by comments
INSERT INTO comment_thread (entity, entity_id, created_at)
SELECT 'Andon', a.andon_id, COALESCE(MIN(c.commented_at), NOW())
FROM andon a
LEFT JOIN comment c ON c.entity = 'Andon' AND c.entity_id = a.andon_id
GROUP BY a.andon_id
ON CONFLICT (entity, entity_id) DO NOTHING;

-- 5. Insert one thread per existing StockItem referenced by comments
INSERT INTO comment_thread (entity, entity_id, created_at)
SELECT 'StockItem', s.stock_item_id, COALESCE(MIN(c.commented_at), NOW())
FROM stock_item s
LEFT JOIN comment c ON c.entity = 'StockItem' AND c.entity_id = s.stock_item_id
GROUP BY s.stock_item_id
ON CONFLICT (entity, entity_id) DO NOTHING;

-- 6. Add thread id columns to domain tables (nullable initially)
ALTER TABLE andon ADD COLUMN IF NOT EXISTS comment_thread_id INTEGER;
ALTER TABLE stock_item ADD COLUMN IF NOT EXISTS comment_thread_id INTEGER;

-- 7. Populate domain tables with their thread ids
UPDATE andon a
SET comment_thread_id = ct.comment_thread_id
FROM comment_thread ct
WHERE ct.entity = 'Andon' AND ct.entity_id = a.andon_id AND a.comment_thread_id IS NULL;

UPDATE stock_item s
SET comment_thread_id = ct.comment_thread_id
FROM comment_thread ct
WHERE ct.entity = 'StockItem' AND ct.entity_id = s.stock_item_id AND s.comment_thread_id IS NULL;

-- 8. Add new nullable column to comment table for FK
ALTER TABLE comment
  ADD COLUMN IF NOT EXISTS comment_thread_id INTEGER REFERENCES comment_thread(comment_thread_id) ON DELETE CASCADE;

-- 9. Backfill comments (Andon)
UPDATE comment c
SET comment_thread_id = a.comment_thread_id
FROM andon a
WHERE c.entity = 'Andon' AND c.entity_id = a.andon_id AND c.comment_thread_id IS NULL;

-- 10. Backfill comments (StockItem)
UPDATE comment c
SET comment_thread_id = s.comment_thread_id
FROM stock_item s
WHERE c.entity = 'StockItem' AND c.entity_id = s.stock_item_id AND c.comment_thread_id IS NULL;

-- 11. Validation: ensure all comments have a thread id now
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM comment WHERE comment_thread_id IS NULL) THEN
    RAISE EXCEPTION 'Backfill of comment_thread_id failed for some rows';
  END IF;
END $$;

-- 12. Drop old comment entity index & columns
DROP INDEX IF EXISTS idx_comment_entity;
ALTER TABLE comment
  DROP COLUMN IF EXISTS entity,
  DROP COLUMN IF EXISTS entity_id;

-- 13. Set NOT NULL constraints & add FKs on domain tables
ALTER TABLE comment
  ALTER COLUMN comment_thread_id SET NOT NULL;
ALTER TABLE andon
  ALTER COLUMN comment_thread_id SET NOT NULL,
  ADD CONSTRAINT fk_andon_comment_thread
    FOREIGN KEY (comment_thread_id) REFERENCES comment_thread(comment_thread_id) ON DELETE RESTRICT;
ALTER TABLE stock_item
  ALTER COLUMN comment_thread_id SET NOT NULL,
  ADD CONSTRAINT fk_stock_item_comment_thread
    FOREIGN KEY (comment_thread_id) REFERENCES comment_thread(comment_thread_id) ON DELETE RESTRICT;

-- 14. Helpful indexes
CREATE INDEX IF NOT EXISTS idx_comment_thread ON comment(comment_thread_id);
CREATE INDEX IF NOT EXISTS idx_andon_comment_thread ON andon(comment_thread_id);
CREATE INDEX IF NOT EXISTS idx_stock_item_comment_thread ON stock_item(comment_thread_id);

-- 15. Remove entity/entity_id from comment_thread (decouple)
ALTER TABLE comment_thread
  DROP COLUMN IF EXISTS entity,
  DROP COLUMN IF EXISTS entity_id;

-- 16. Recreate view
CREATE VIEW comment_view AS
SELECT
  c.comment_id,
  c.comment_thread_id,
  c.comment,
  u.username as commented_by_username,
  c.commented_at,
  COALESCE(json_agg(json_build_object(
    'file_id', f.file_id,
    'filename', f.filename,
    'content_type', f.content_type,
    'size_bytes', f.size_bytes,
    'status', f.status,
    'user_id', f.user_id
  ) ORDER BY f.created_at ASC
  ) FILTER (WHERE f.file_id IS NOT NULL), '[]') AS attachments
FROM comment c
LEFT JOIN app_user u ON c.commented_by = u.user_id
LEFT JOIN file f ON f.entity = 'Comment' AND f.entity_id = c.comment_id
GROUP BY c.comment_id, c.comment_thread_id, c.comment, u.username, c.commented_at
ORDER BY c.commented_at ASC;

-- 17. Recreate andon_view adding the new a.comment_thread_id column (after gallery_id)
CREATE VIEW andon_view AS
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
  (acknowledged_at IS NOT NULL) AS is_acknowledged,
  (resolved_at IS NOT NULL) AS is_resolved,
  (cancelled_at IS NOT NULL) AS is_cancelled,
  CASE
    WHEN cancelled_at IS NOT NULL THEN false
    WHEN severity = 'Info' AND acknowledged_at IS NOT NULL THEN false
    WHEN severity IN ('Self-resolvable', 'Requires Intervention')
       AND acknowledged_at IS NOT NULL
       AND resolved_at IS NOT NULL
    THEN false
    ELSE true
  END AS is_open,
  CASE
    WHEN cancelled_at IS NOT NULL THEN 'Cancelled'

    -- Info
    WHEN severity = 'Info'
      AND acknowledged_at IS NOT NULL THEN 'Closed'
    WHEN severity = 'Info'
      AND acknowledged_at IS NULL THEN 'Requires Acknowledgement'

    -- Self-resolvable
    WHEN severity ='Self-resolvable'
      AND acknowledged_at IS NOT NULL
      AND resolved_at IS NOT NULL THEN 'Closed'
    -- Self-resolvable andons are considered WIP immediately upon creation
    WHEN severity = 'Self-resolvable'
      AND resolved_at IS NULL THEN 'Work In Progress'
    WHEN severity = 'Self-resolvable'
      AND acknowledged_at IS NULL THEN 'Requires Acknowledgement'

    -- Requires Intervention
    WHEN severity = 'Requires Intervention'
      AND acknowledged_at IS NOT NULL
      AND resolved_at IS NOT NULL THEN 'Closed'
    WHEN severity = 'Requires Intervention'
      AND acknowledged_at IS NULL
      AND resolved_at IS NULL THEN 'Outstanding'
    WHEN severity = 'Requires Intervention'
      AND acknowledged_at IS NOT NULL THEN 'Work In Progress'
    WHEN severity = 'Requires Intervention'
      AND resolved_at IS NOT NULL THEN 'Requires Acknowledgement'

     -- should never see this
    ELSE 'Invalid Status'
  END AS status
FROM
  andon a
  INNER JOIN app_user u ON a.raised_by = u.user_id
  LEFT JOIN app_user acku ON a.acknowledged_by = acku.user_id
  LEFT JOIN app_user ru ON a.resolved_by = ru.user_id
  LEFT JOIN app_user cu ON a.cancelled_by = cu.user_id
  INNER JOIN andon_issue_view aiv ON a.andon_issue_id = aiv.andon_issue_id;

COMMIT;