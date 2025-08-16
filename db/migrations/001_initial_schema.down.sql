-- Drop triggers
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_inbound_messages_created_at;
DROP INDEX IF EXISTS idx_inbound_messages_from_number;
DROP INDEX IF EXISTS idx_users_wa_number;
DROP INDEX IF EXISTS idx_events_reminders;
DROP INDEX IF EXISTS idx_events_status;
DROP INDEX IF EXISTS idx_events_starts_at;
DROP INDEX IF EXISTS idx_events_user_id;

-- Drop tables in reverse order
DROP TABLE IF EXISTS inbound_messages;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS whitelist_numbers;
DROP TABLE IF EXISTS users;