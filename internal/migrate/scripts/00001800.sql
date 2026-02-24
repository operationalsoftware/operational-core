-- 00001800.sql: add target_url to comment_thread and backfill existing threads

ALTER TABLE comment_thread
    ADD COLUMN IF NOT EXISTS target_url TEXT;

UPDATE comment_thread ct
SET target_url = COALESCE(
    (
        SELECT '/andons/' || a.andon_id::text
        FROM andon a
        WHERE a.comment_thread_id = ct.comment_thread_id
        ORDER BY a.andon_id DESC
        LIMIT 1
    ),
    (
        SELECT '/stock-items/' || si.stock_item_id::text
        FROM stock_item si
        WHERE si.comment_thread_id = ct.comment_thread_id
        ORDER BY si.stock_item_id DESC
        LIMIT 1
    ),
    (
        SELECT '/services/' || rs.resource_service_id::text
        FROM resource_service rs
        WHERE rs.comment_thread_id = ct.comment_thread_id
        ORDER BY rs.resource_service_id DESC
        LIMIT 1
    )
)
WHERE ct.target_url IS NULL OR btrim(ct.target_url) = '';

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM comment_thread
        WHERE target_url IS NULL OR btrim(target_url) = ''
    ) THEN
        RAISE EXCEPTION 'comment_thread.target_url backfill failed for one or more threads';
    END IF;
END;
$$;
