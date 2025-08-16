-- Drop indexes
DROP INDEX IF EXISTS idx_llm_models_default_per_provider;
DROP INDEX IF EXISTS idx_users_llm_model;
DROP INDEX IF EXISTS idx_users_llm_provider;
DROP INDEX IF EXISTS idx_llm_models_default;
DROP INDEX IF EXISTS idx_llm_models_active;
DROP INDEX IF EXISTS idx_llm_models_provider_id;
DROP INDEX IF EXISTS idx_llm_providers_active;
DROP INDEX IF EXISTS idx_llm_providers_name;

-- Drop triggers
DROP TRIGGER IF EXISTS update_llm_models_updated_at ON llm_models;
DROP TRIGGER IF EXISTS update_llm_providers_updated_at ON llm_providers;

-- Remove columns from users table
ALTER TABLE users DROP COLUMN IF EXISTS llm_model;
ALTER TABLE users DROP COLUMN IF EXISTS llm_provider;

-- Drop tables
DROP TABLE IF EXISTS llm_models;
DROP TABLE IF EXISTS llm_providers;