-- Drop indecies
DROP INDEX IF EXISTS idx_drones_identifier;
DROP INDEX IF EXISTS idx_drones_status;
DROP INDEX IF EXISTS idx_drones_active;
DROP INDEX IF EXISTS idx_drones_user_id;
DROP INDEX IF EXISTS idx_drones_serial_number;
DROP INDEX IF EXISTS idx_drones_status_active;
DROP INDEX IF EXISTS idx_drones_location;

-- Drop the drones table and related objects
DROP TABLE IF EXISTS drones CASCADE;
DROP FUNCTION IF EXISTS generate_drone_identifier() CASCADE;
DROP SEQUENCE IF EXISTS drones_seq CASCADE;
