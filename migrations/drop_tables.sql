-- migrations/drop_tables.sql


DROP TRIGGER IF EXISTS trigger_update_tickets_sold ON tickets;


DROP FUNCTION IF EXISTS update_tickets_sold();

DROP INDEX IF EXISTS idx_events_owner;
DROP INDEX IF EXISTS idx_tickets_event;
DROP INDEX IF EXISTS idx_tickets_user;
DROP INDEX IF EXISTS idx_orders_user;
DROP INDEX IF EXISTS idx_orders_event;
DROP INDEX IF EXISTS idx_payments_order;
DROP INDEX IF EXISTS idx_events_date;
DROP INDEX IF EXISTS idx_events_status;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_midtrans_transaction;
DROP INDEX IF EXISTS idx_user_profiles_user_id;

DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS tickets CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS user_profiles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS email_verifications CASCADE;