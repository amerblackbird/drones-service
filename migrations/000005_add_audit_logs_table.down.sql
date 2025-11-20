-- Drop indexes
DROP INDEX IF EXISTS idx_audit_logs_performed_at;
DROP INDEX IF EXISTS idx_audit_logs_performed_by;
DROP INDEX IF EXISTS idx_audit_logs_resource;


-- Drop triggers
DROP TRIGGER IF EXISTS trg_audit_logs_updated_at ON audit_logs;

-- Drop table
DROP TABLE IF EXISTS audit_logs;
