-- Create llm_providers table for available LLM providers
CREATE TABLE llm_providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create llm_models table for available models
CREATE TABLE llm_models (
    id SERIAL PRIMARY KEY,
    provider_id INTEGER REFERENCES llm_providers(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,
    config JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(provider_id, name)
);

-- Add llm_provider and llm_model columns to users table for user-level configuration
ALTER TABLE users ADD COLUMN llm_provider VARCHAR(50);
ALTER TABLE users ADD COLUMN llm_model VARCHAR(100);

-- Create indexes
CREATE INDEX idx_llm_providers_name ON llm_providers(name);
CREATE INDEX idx_llm_providers_active ON llm_providers(is_active);
CREATE INDEX idx_llm_models_provider_id ON llm_models(provider_id);
CREATE INDEX idx_llm_models_active ON llm_models(is_active);
CREATE INDEX idx_llm_models_default ON llm_models(is_default);
CREATE INDEX idx_users_llm_provider ON users(llm_provider);
CREATE INDEX idx_users_llm_model ON users(llm_model);

-- Create triggers for updating updated_at columns
CREATE TRIGGER update_llm_providers_updated_at BEFORE UPDATE ON llm_providers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_llm_models_updated_at BEFORE UPDATE ON llm_models FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default providers
INSERT INTO llm_providers (name, display_name, description, is_active) VALUES 
    ('anthropic', 'Anthropic', 'Anthropic Claude models', true),
    ('openai', 'OpenAI', 'OpenAI GPT models', true);

-- Insert default models
INSERT INTO llm_models (provider_id, name, display_name, description, is_default, is_active) VALUES 
    ((SELECT id FROM llm_providers WHERE name = 'anthropic'), 'claude-3-haiku-20240307', 'Claude 3 Haiku', 'Fast and efficient Claude model', true, true),
    ((SELECT id FROM llm_providers WHERE name = 'anthropic'), 'claude-3-sonnet-20240229', 'Claude 3 Sonnet', 'Balanced Claude model', false, true),
    ((SELECT id FROM llm_providers WHERE name = 'openai'), 'gpt-3.5-turbo', 'GPT-3.5 Turbo', 'OpenAI GPT-3.5 Turbo', false, true),
    ((SELECT id FROM llm_providers WHERE name = 'openai'), 'gpt-4', 'GPT-4', 'OpenAI GPT-4', false, true);

-- Add constraint to ensure only one default model per provider
CREATE UNIQUE INDEX idx_llm_models_default_per_provider ON llm_models(provider_id) WHERE is_default = true;