-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    wa_number VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255),
    timezone VARCHAR(50) DEFAULT 'America/Sao_Paulo',
    default_remind_before_minutes INTEGER DEFAULT 30,
    default_remind_frequency_minutes INTEGER DEFAULT 15,
    default_require_confirmation BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create whitelist_numbers table
CREATE TABLE whitelist_numbers (
    number VARCHAR(20) PRIMARY KEY,
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create events table
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    location VARCHAR(500),
    starts_at TIMESTAMP WITH TIME ZONE NOT NULL,
    remind_before_minutes INTEGER DEFAULT 30,
    remind_frequency_minutes INTEGER DEFAULT 15,
    require_confirmation BOOLEAN DEFAULT true,
    max_notifications INTEGER DEFAULT 3,
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'confirmed', 'canceled', 'completed')),
    notifications_sent INTEGER DEFAULT 0,
    last_notified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create inbound_messages table for idempotency
CREATE TABLE inbound_messages (
    id SERIAL PRIMARY KEY,
    provider_message_id VARCHAR(255) UNIQUE NOT NULL,
    from_number VARCHAR(20) NOT NULL,
    raw_payload JSONB NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_starts_at ON events(starts_at);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_reminders ON events(starts_at, status, notifications_sent, max_notifications) WHERE status IN ('scheduled', 'confirmed');
CREATE INDEX idx_users_wa_number ON users(wa_number);
CREATE INDEX idx_inbound_messages_from_number ON inbound_messages(from_number);
CREATE INDEX idx_inbound_messages_created_at ON inbound_messages(created_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updating updated_at columns
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();