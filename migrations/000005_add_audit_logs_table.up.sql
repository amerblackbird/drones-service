--- Audit Logs Table
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resource_name VARCHAR(100) NOT NULL,   -- e.g., 'orders', 'drones', 'assignments'
    resource_id VARCHAR(100) NOT NULL,     -- id of the resource
    action VARCHAR(50) NOT NULL,         -- 'CREATE', 'UPDATE', 'DELETE'
    old_data JSONB,                      -- previous row data
    new_data JSONB,                      -- new row data
    performed_by_id UUID,                   -- user or system who did the action
    performed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    performed_user_type VARCHAR(50),  -- 'admin', 'enduser', 'system'
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by_id UUID,
    updated_by_id UUID
    
);

CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_name, resource_id);
CREATE INDEX idx_audit_logs_performed_by ON audit_logs(performed_by_id);
CREATE INDEX idx_audit_logs_performed_at ON audit_logs(performed_at);


CREATE TRIGGER trg_audit_logs_updated_at
BEFORE UPDATE ON audit_logs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();