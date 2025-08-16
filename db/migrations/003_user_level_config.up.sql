-- Add user-level configuration fields
ALTER TABLE users ADD COLUMN rate_limit_per_minute INTEGER DEFAULT 30;
ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT true;

-- Create user_allowed_contacts table to replace global whitelist
CREATE TABLE user_allowed_contacts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    contact_number VARCHAR(20) NOT NULL,
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, contact_number)
);

-- Create index for performance
CREATE INDEX idx_user_allowed_contacts_user_id ON user_allowed_contacts(user_id);
CREATE INDEX idx_user_allowed_contacts_contact_number ON user_allowed_contacts(contact_number);

-- Create trigger for updating updated_at column
CREATE TRIGGER update_user_allowed_contacts_updated_at BEFORE UPDATE ON user_allowed_contacts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Migrate existing whitelist data to user-level (we'll handle this in the application)
-- For now, we keep both tables to allow gradual migration