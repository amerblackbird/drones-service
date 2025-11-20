-- Drop indexes
DROP INDEX IF EXISTS idx_activity_logs_performed_at;
DROP INDEX IF EXISTS idx_activity_logs_action;
DROP INDEX IF EXISTS idx_activity_logs_actor;

-- Drop triggers
DROP TRIGGER IF EXISTS trg_activity_logs_updated_at ON activity_logs;

-- Drop table
DROP TABLE IF EXISTS activity_logs;
