-- Migration 00000400
-- Goal:
--  * Allow SSO-only accounts by making hashed_password nullable
--  * Add auth_provider to track authentication source (default 'local')
--  * Add external_id to store provider-specific user ID (e.g., Microsoft Graph id)

ALTER TABLE app_user 
  ALTER COLUMN hashed_password DROP NOT NULL;

ALTER TABLE app_user 
  ADD COLUMN IF NOT EXISTS auth_provider TEXT DEFAULT 'local' NOT NULL;

ALTER TABLE app_user 
  ADD COLUMN IF NOT EXISTS external_id TEXT;

CREATE INDEX IF NOT EXISTS idx_app_user_external_id ON app_user(external_id);
