package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type LLMConfigRepository struct {
	db QueryExecutor
}

func NewLLMConfigRepository(db QueryExecutor) ports.LLMConfigRepository {
	return &LLMConfigRepository{db: db}
}

func (r *LLMConfigRepository) GetDefaultModel(ctx context.Context) (*domain.LLMModel, error) {
	query := `
		SELECT m.id, m.provider_id, m.name, m.display_name, m.description, 
			   m.is_active, m.is_default, m.config, m.created_at, m.updated_at,
			   p.id as "provider.id", p.name as "provider.name", p.display_name as "provider.display_name", 
			   p.description as "provider.description", p.is_active as "provider.is_active", 
			   p.created_at as "provider.created_at", p.updated_at as "provider.updated_at"
		FROM llm_models m
		JOIN llm_providers p ON m.provider_id = p.id
		WHERE m.is_default = true AND m.is_active = true AND p.is_active = true
		LIMIT 1`

	var result struct {
		domain.LLMModel
		Provider domain.LLMProvider `db:"provider"`
	}

	err := r.db.GetContext(ctx, &result, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no default LLM model found")
		}
		return nil, fmt.Errorf("failed to get default model: %w", err)
	}

	result.LLMModel.Provider = &result.Provider
	return &result.LLMModel, nil
}

func (r *LLMConfigRepository) GetModelByProviderAndName(ctx context.Context, providerName, modelName string) (*domain.LLMModel, error) {
	query := `
		SELECT m.id, m.provider_id, m.name, m.display_name, m.description, 
			   m.is_active, m.is_default, m.config, m.created_at, m.updated_at,
			   p.id as "provider.id", p.name as "provider.name", p.display_name as "provider.display_name",
			   p.description as "provider.description", p.is_active as "provider.is_active", 
			   p.created_at as "provider.created_at", p.updated_at as "provider.updated_at"
		FROM llm_models m
		JOIN llm_providers p ON m.provider_id = p.id
		WHERE p.name = $1 AND m.name = $2 AND m.is_active = true AND p.is_active = true`

	var result struct {
		domain.LLMModel
		Provider domain.LLMProvider `db:"provider"`
	}

	err := r.db.GetContext(ctx, &result, query, providerName, modelName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("model %s not found for provider %s", modelName, providerName)
		}
		return nil, fmt.Errorf("failed to get model: %w", err)
	}

	result.LLMModel.Provider = &result.Provider
	return &result.LLMModel, nil
}

func (r *LLMConfigRepository) GetActiveProviders(ctx context.Context) ([]domain.LLMProvider, error) {
	query := `
		SELECT id, name, display_name, description, is_active, created_at, updated_at
		FROM llm_providers
		WHERE is_active = true
		ORDER BY name`

	var providers []domain.LLMProvider
	err := r.db.SelectContext(ctx, &providers, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active providers: %w", err)
	}

	return providers, nil
}

func (r *LLMConfigRepository) GetActiveModelsByProvider(ctx context.Context, providerID int) ([]domain.LLMModel, error) {
	query := `
		SELECT m.id, m.provider_id, m.name, m.display_name, m.description, 
			   m.is_active, m.is_default, m.config, m.created_at, m.updated_at
		FROM llm_models m
		WHERE m.provider_id = $1 AND m.is_active = true
		ORDER BY m.is_default DESC, m.name`

	var models []domain.LLMModel
	err := r.db.SelectContext(ctx, &models, query, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active models for provider %d: %w", providerID, err)
	}

	return models, nil
}

func (r *LLMConfigRepository) GetUserLLMConfig(ctx context.Context, userID int) (*domain.LLMModel, error) {
	query := `
		SELECT m.id, m.provider_id, m.name, m.display_name, m.description, 
			   m.is_active, m.is_default, m.config, m.created_at, m.updated_at,
			   p.id as "provider.id", p.name as "provider.name", p.display_name as "provider.display_name",
			   p.description as "provider.description", p.is_active as "provider.is_active", 
			   p.created_at as "provider.created_at", p.updated_at as "provider.updated_at"
		FROM users u
		LEFT JOIN llm_providers p ON u.llm_provider = p.name
		LEFT JOIN llm_models m ON u.llm_model = m.name AND m.provider_id = p.id
		WHERE u.id = $1 AND (m.is_active IS NULL OR m.is_active = true) AND (p.is_active IS NULL OR p.is_active = true)`

	var result struct {
		domain.LLMModel
		Provider domain.LLMProvider `db:"provider"`
	}

	err := r.db.GetContext(ctx, &result, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// User doesn't have custom LLM config, use default
			return r.GetDefaultModel(ctx)
		}
		return nil, fmt.Errorf("failed to get user LLM config: %w", err)
	}

	// If user has no custom config (null values), use default
	if result.ID == 0 {
		return r.GetDefaultModel(ctx)
	}

	result.LLMModel.Provider = &result.Provider
	return &result.LLMModel, nil
}
