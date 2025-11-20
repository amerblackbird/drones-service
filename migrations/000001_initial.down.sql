DROP INDEX IF EXISTS idx_users_type;
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_email;

-- Drop triggers
-- DROP TRIGGER IF EXISTS trg_drones_updated_at ON drones;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;

-- Drop tables in correct order (foreign key dependencies)
-- DROP TABLE IF EXISTS drones;
DROP TABLE IF EXISTS users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";