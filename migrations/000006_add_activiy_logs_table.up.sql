--- Activity Logs Table
CREATE TABLE activity_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_id VARCHAR(50),              -- drone or user
    actor_type VARCHAR(50) NOT NULL,     -- 'drone', 'user'
    action VARCHAR(100) NOT NULL,        -- 'ORDER_SUBMITTED', 'ORDER_PICKED', 'ORDER_DELIVERED', 'LOCATION_UPDATE', etc.
    resource_name VARCHAR(100),            -- optional: related entity, e.g., 'orders'
    resource_id UUID,                      -- optional: related entity id
    metadata JSONB,                      -- optional extra data
    performed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ip VARCHAR(45),                
    device VARCHAR(255),            -- device info if applicable
    location VARCHAR(255),          -- geo-location if applicable
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by_id UUID,
    updated_by_id UUID
    
);

CREATE INDEX idx_activity_logs_actor ON activity_logs(actor_id, actor_type);
CREATE INDEX idx_activity_logs_action ON activity_logs(action);
CREATE INDEX idx_activity_logs_performed_at ON activity_logs(performed_at);

CREATE TRIGGER trg_activity_logs_updated_at
BEFORE UPDATE ON activity_logs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();