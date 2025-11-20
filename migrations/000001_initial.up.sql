-- Enable the uuid-ossp extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create the users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- Basic Info
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'admin', 'enduser', 'drone'
    email VARCHAR(255),
    phone VARCHAR(20),

    country VARCHAR(100),
    locale VARCHAR(10),
    device_id VARCHAR(255),
    notification_token VARCHAR(255),
    avatar_url TEXT,
    bio TEXT,
    
    -- Base fields
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by_id UUID,
    updated_by_id UUID,

    -- Name and type unique constraint
    UNIQUE (name, type)
);

CREATE INDEX idx_users_active ON users(name, type);
CREATE INDEX idx_users_type ON users(type);
CREATE INDEX idx_users_email ON users(email);

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


-- -- Create the drones table
-- CREATE TABLE drones (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

--     -- Identification
--     serial_number VARCHAR(100) UNIQUE NOT NULL,
--     model VARCHAR(100),
--     manufacturer VARCHAR(100),
    
--     -- Specifications
--     battery_capacity INT,
--     payload_capacity FLOAT,

--     -- Battery
--     last_charged_at TIMESTAMP,
--     is_charging BOOLEAN DEFAULT FALSE,

--     -- Telemetry
--     last_known_lat DOUBLE PRECISION,
--     last_known_lng DOUBLE PRECISION,
--     last_altitude_m DOUBLE PRECISION,
--     last_speed_kmh DOUBLE PRECISION,
--     status VARCHAR(50) NOT NULL DEFAULT 'available', -- 'available', 'assigned', 'broken'
    
--     -- Operational
--     current_order_id UUID,
--     crashes_count INT DEFAULT 0,
--     maintenance_required BOOLEAN DEFAULT FALSE,

--     -- Maintenance
--     last_maintenance_at TIMESTAMP,
--     next_maintenance_at TIMESTAMP,

--     -- Base fields
--     created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
--     updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
--     active BOOLEAN NOT NULL DEFAULT TRUE,
--     created_by_id UUID,
--     updated_by_id UUID
    
-- );

-- CREATE INDEX idx_drones_serial_number ON drones(serial_number);
-- CREATE INDEX idx_drones_status ON drones(status);
-- CREATE INDEX idx_drones_active ON drones(active);

-- CREATE TRIGGER trg_drones_updated_at
-- BEFORE UPDATE ON drones
-- FOR EACH ROW
-- EXECUTE FUNCTION update_updated_at_column();





