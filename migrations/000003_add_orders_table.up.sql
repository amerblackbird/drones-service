-- Create sequence for order numbers
CREATE SEQUENCE orders_seq START 1;

-- Create function to generate order number with date prefix
CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS VARCHAR AS $$
DECLARE
    date_prefix VARCHAR(6);
    seq_num VARCHAR(6);
    order_num VARCHAR(100);
BEGIN
    -- Get current date in YYMMDD format
    date_prefix := TO_CHAR(CURRENT_DATE, 'YYMMDD');
    
    -- Get next sequence value and pad with zeros to 6 digits
    seq_num := LPAD(NEXTVAL('orders_seq')::TEXT, 6, '0');
    
    -- Combine to create order number: ORD + YYMMDD + sequence
    order_num := 'ORD' || date_prefix || seq_num;
    
    RETURN order_num;
END;
$$ LANGUAGE plpgsql;

-- Create the orders table
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_number VARCHAR(100) UNIQUE NOT NULL DEFAULT generate_order_number(),

    -- Customer & Receiver
    user_id UUID NOT NULL,
    receiver_name VARCHAR(150),
    receiver_phone VARCHAR(20),
    delivery_note TEXT,

    -- Package
    package_weight_kg NUMERIC(10,2),

    -- Origin details
    origin_address TEXT NOT NULL,
    origin_lat DOUBLE PRECISION NOT NULL,
    origin_lon DOUBLE PRECISION NOT NULL,

    destination_address TEXT NOT NULL,
    destination_lat DOUBLE PRECISION NOT NULL,
    destination_lon DOUBLE PRECISION NOT NULL,

    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'in_progress', 'delivered', 'failed', 'in_flight', 'cancelled'
    
    -- Timing
    scheduled_at TIMESTAMP,
    picked_up_at TIMESTAMP,
    delivered_at TIMESTAMP,
    failed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    withdrawn_at TIMESTAMP,

    -- Delivery details
    drone_id UUID,
    delivered_by_drone_id UUID,
    
    -- Real-time tracking
    current_lat DOUBLE PRECISION,
    current_lon DOUBLE PRECISION,
    current_altitude DOUBLE PRECISION,
    last_location_update_at TIMESTAMP,
    estimated_arrival_at TIMESTAMP,

     -- Base fields
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by_id UUID,
    updated_by_id UUID,
    
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (drone_id) REFERENCES drones(id),
    FOREIGN KEY (delivered_by_drone_id) REFERENCES drones(id)
);

CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_active ON orders(active);
CREATE INDEX idx_orders_drone_id ON orders(drone_id);
CREATE INDEX idx_orders_delivered_by_drone_id ON orders(delivered_by_drone_id);
CREATE INDEX idx_orders_status_active ON orders(status, active) WHERE active = TRUE;

CREATE TRIGGER trg_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();