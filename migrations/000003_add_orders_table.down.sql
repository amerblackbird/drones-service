-- Drop indexes
DROP INDEX IF EXISTS idx_orders_order_number;
DROP INDEX IF EXISTS idx_orders_active;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_drone_id;
DROP INDEX IF EXISTS idx_orders_delivered_by_drone_id;
DROP INDEX IF EXISTS idx_orders_status_active;

-- Drop triggers
DROP TRIGGER IF EXISTS trg_orders_updated_at ON orders;

-- Drop table
DROP TABLE IF EXISTS orders;

-- Drop function
DROP FUNCTION IF EXISTS generate_order_number();

-- Drop sequence
DROP SEQUENCE IF EXISTS orders_seq;

