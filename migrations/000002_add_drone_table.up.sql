-- Create sequence for drone identifiers
CREATE SEQUENCE drones_seq START 1;

-- Create function to generate drone identifier
CREATE OR REPLACE FUNCTION generate_drone_identifier()
RETURNS VARCHAR AS $$
DECLARE
    seq_num VARCHAR(6);
    drone_id VARCHAR(100);
BEGIN
    -- Get next sequence value and pad with zeros to 6 digits
    seq_num := LPAD(NEXTVAL('drones_seq')::TEXT, 6, '0');
    
    -- Create drone identifier: DRN + sequence
    drone_id := 'DRN' || seq_num;
    
    RETURN drone_id;
END;
$$ LANGUAGE plpgsql;

-- Create the drones table
CREATE TABLE drones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    drone_identifier VARCHAR(100) UNIQUE NOT NULL DEFAULT generate_drone_identifier(),
    
    -- Owner/Operator (one-to-one relationship)
    user_id UUID NOT NULL UNIQUE,
    
    -- Drone details
    model VARCHAR(100),
    serial_number VARCHAR(100) UNIQUE,
    manufacturer VARCHAR(100),
    
    -- Capacity
    max_weight_kg NUMERIC(10,2),
    max_speed_kmh NUMERIC(10,2),
    max_range_km NUMERIC(10,2),
    battery_capacity_mah INTEGER,
    
    -- Current status
    status VARCHAR(50) NOT NULL DEFAULT 'idle', -- Workflow: idle -> assigned -> in_flight -> delivering -> returning -> idle (or charging/maintenance/broken)
     battery_level_percent NUMERIC(5,2),
    
    -- Location tracking
    current_lat DOUBLE PRECISION,
    current_lon DOUBLE PRECISION,
    current_altitude DOUBLE PRECISION,
    last_location_update_at TIMESTAMP,
    
    
    -- Operational data
    total_flight_hours NUMERIC(10,2) DEFAULT 0,
    total_deliveries INTEGER DEFAULT 0,
    last_maintenance_at TIMESTAMP,
    next_maintenance_due_at TIMESTAMP,
    
    -- Base fields
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by_id UUID,
    updated_by_id UUID,
    
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Indexes for efficient querying
CREATE INDEX idx_drones_identifier ON drones(drone_identifier);
CREATE INDEX idx_drones_status ON drones(status);
CREATE INDEX idx_drones_active ON drones(active);
CREATE INDEX idx_drones_user_id ON drones(user_id);
CREATE INDEX idx_drones_serial_number ON drones(serial_number);
CREATE INDEX idx_drones_status_active ON drones(status, active) WHERE active = TRUE;
CREATE INDEX idx_drones_location ON drones(current_lat, current_lon) WHERE status = 'in_flight';


CREATE TRIGGER trg_drones_updated_at
BEFORE UPDATE ON drones
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();