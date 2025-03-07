-- migrations/schema.sql

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    location VARCHAR(200),
    event_date TIMESTAMP NOT NULL,
    max_capacity INTEGER NOT NULL,
    tickets_sold INTEGER DEFAULT 0,
    price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    event_id INTEGER REFERENCES events(id),
    user_id INTEGER REFERENCES users(id),
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    event_id INTEGER REFERENCES events(id),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    midtrans_transaction_id VARCHAR(100) UNIQUE,
    payment_type VARCHAR(50),
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    midtrans_status_code VARCHAR(10),
    midtrans_status_message TEXT,
    payment_time TIMESTAMP,
    expiry_time TIMESTAMP,
    snap_token VARCHAR(200),
    redirect_url TEXT,
    callback_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_events_owner ON events(owner_id);
CREATE INDEX idx_tickets_event ON tickets(event_id);
CREATE INDEX idx_tickets_user ON tickets(user_id);
CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_event ON orders(event_id);
CREATE INDEX idx_payments_order ON payments(order_id);

CREATE INDEX idx_events_date ON events(event_date);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_midtrans_transaction ON payments(midtrans_transaction_id);

ALTER TABLE events ADD CONSTRAINT check_capacity CHECK (tickets_sold <= max_capacity);

ALTER TABLE events ADD CONSTRAINT check_price CHECK (price >= 0);

ALTER TABLE events ADD CONSTRAINT check_capacity_positive CHECK (max_capacity > 0);

CREATE OR REPLACE FUNCTION update_tickets_sold() RETURNS TRIGGER AS $update_tickets_sold$
BEGIN
    IF NEW.status = 'active' AND (TG_OP = 'INSERT' OR OLD.status != 'active') THEN
        UPDATE events SET tickets_sold = tickets_sold + 1 WHERE id = NEW.event_id;
    ELSIF OLD.status = 'active' AND NEW.status != 'active' THEN
        UPDATE events SET tickets_sold = tickets_sold - 1 WHERE id = NEW.event_id;
    END IF;
    RETURN NEW;
END;
$update_tickets_sold$ LANGUAGE plpgsql;


CREATE TRIGGER trigger_update_tickets_sold
AFTER INSERT OR UPDATE ON tickets
FOR EACH ROW
EXECUTE FUNCTION update_tickets_sold();