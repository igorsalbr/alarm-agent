package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type UserRepository struct {
	db QueryExecutor
}

func NewUserRepository(db QueryExecutor) ports.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByWANumber(ctx context.Context, waNumber string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, wa_number, name, timezone, default_remind_before_minutes, 
		       default_remind_frequency_minutes, default_require_confirmation, 
		       llm_provider, llm_model, rate_limit_per_minute, is_active,
		       created_at, updated_at
		FROM users 
		WHERE wa_number = $1`

	err := r.db.GetContext(ctx, &user, query, waNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID int) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, wa_number, name, timezone, default_remind_before_minutes, 
		       default_remind_frequency_minutes, default_require_confirmation, 
		       llm_provider, llm_model, rate_limit_per_minute, is_active,
		       created_at, updated_at
		FROM users 
		WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (wa_number, name, timezone, default_remind_before_minutes, 
		                   default_remind_frequency_minutes, default_require_confirmation,
		                   llm_provider, llm_model, rate_limit_per_minute, is_active)
		VALUES (:wa_number, :name, :timezone, :default_remind_before_minutes, 
		        :default_remind_frequency_minutes, :default_require_confirmation,
		        :llm_provider, :llm_model, :rate_limit_per_minute, :is_active)
		RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	id, err := rows.LastInsertId()
	if err == nil && id > 0 {
		user.ID = int(id)
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET name = :name, timezone = :timezone, 
		    default_remind_before_minutes = :default_remind_before_minutes,
		    default_remind_frequency_minutes = :default_remind_frequency_minutes,
		    default_require_confirmation = :default_require_confirmation,
		    llm_provider = :llm_provider, llm_model = :llm_model,
		    rate_limit_per_minute = :rate_limit_per_minute, is_active = :is_active,
		    updated_at = NOW()
		WHERE id = :id`

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *UserRepository) UpdateConfig(ctx context.Context, userID int, config *domain.UserConfig) error {
	query := `
		UPDATE users 
		SET name = $2, timezone = $3, 
		    default_remind_before_minutes = $4,
		    default_remind_frequency_minutes = $5,
		    default_require_confirmation = $6,
		    llm_provider = $7, llm_model = $8,
		    rate_limit_per_minute = $9, is_active = $10,
		    updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, userID, config.Name, config.Timezone,
		config.DefaultRemindBeforeMinutes, config.DefaultRemindFrequencyMinutes,
		config.DefaultRequireConfirmation, config.LLMProvider, config.LLMModel,
		config.RateLimitPerMinute, config.IsActive)
	return err
}
