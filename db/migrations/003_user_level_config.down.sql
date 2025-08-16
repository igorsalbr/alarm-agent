-- Drop user_allowed_contacts table
DROP TABLE IF EXISTS user_allowed_contacts;

-- Remove user-level configuration fields
ALTER TABLE users DROP COLUMN IF EXISTS rate_limit_per_minute;
ALTER TABLE users DROP COLUMN IF EXISTS is_active;