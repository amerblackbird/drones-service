-- Seed users with different types
-- Admin users
INSERT INTO users ( name, type, email, phone, country, locale, active) VALUES
    ( 'admin', 'admin', 'admin@drones.com', '+1234567890', 'US', 'en', true),
    ('superadmin', 'admin', 'superadmin@drones.com', '+1234567891', 'US', 'en', true);

-- End users (customers)
INSERT INTO users (name, type, email, phone, country, locale, active, bio) VALUES
    ('john_doe', 'enduser', 'john.doe@example.com', '+1234567892', 'US', 'en', true, 'Regular customer'),
    ( 'jane_smith', 'enduser', 'jane.smith@example.com', '+1234567893', 'US', 'en', true, 'Premium customer'),
    ('bob_wilson', 'enduser', 'bob.wilson@example.com', '+1234567894', 'CA', 'en', true, 'Frequent flyer'),
    ( 'alice_brown', 'enduser', 'alice.brown@example.com', '+1234567895', 'UK', 'en', true, 'New customer');

-- Drone users (drone operators/pilots)
INSERT INTO users ( name, type, email, phone, country, locale, active, bio) VALUES
    ('pilot_alpha', 'drone', 'pilot.alpha@drones.com', '+1234567896', 'US', 'en', true, 'Experienced drone operator'),
    ('pilot_beta', 'drone', 'pilot.beta@drones.com', '+1234567897', 'US', 'en', true, 'Senior drone pilot'),
    ('pilot_gamma', 'drone', 'pilot.gamma@drones.com', '+1234567898', 'CA', 'en', true, 'Night operations specialist'),
    ('pilot_delta', 'drone', 'pilot.delta@drones.com', '+1234567899', 'US', 'en', false, 'On leave');



-- Seed drones for each drone user
INSERT INTO drones (
    user_id,
    model,
    serial_number,
    manufacturer,
    max_weight_kg,
    max_speed_kmh,
    max_range_km,
    battery_capacity_mah,
    status,
    battery_level_percent,
    active
) VALUES
    -- Drones for pilot_alpha
    (
        (SELECT id FROM users WHERE email = 'pilot.alpha@drones.com'),
        'DJI Matrice 300',
        'SN-ALPHA-001',
        'DJI',
        2.70,
        82.00,
        15.00,
        5880,
        'idle',
        100.00,
        true
    ),
    -- Drones for pilot_beta
    (
        (SELECT id FROM users WHERE email = 'pilot.beta@drones.com'),
        'Autel EVO Max 4T',
        'SN-BETA-001',
        'Autel Robotics',
        2.00,
        72.00,
        12.00,
        8200,
        'idle',
        85.00,
        true
    ),
    -- Drones for pilot_gamma
    (
        (SELECT id FROM users WHERE email = 'pilot.gamma@drones.com'),
        'DJI M30T',
        'SN-GAMMA-001',
        'DJI',
        3.20,
        82.00,
        18.00,
        6200,
        'idle',
        95.00,
        true
    ),
    -- Drones for pilot_delta (inactive pilot)
    (
        (SELECT id FROM users WHERE email = 'pilot.delta@drones.com'),
        'DJI Matrice 300',
        'SN-DELTA-001',
        'DJI',
        2.70,
        82.00,
        15.00,
        5880,
        'offline',
        60.00,
        false
    );
