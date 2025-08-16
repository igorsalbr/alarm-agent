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

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (wa_number, name, timezone, default_remind_before_minutes, 
		                   default_remind_frequency_minutes, default_require_confirmation)
		VALUES (:wa_number, :name, :timezone, :default_remind_before_minutes, 
		        :default_remind_frequency_minutes, :default_require_confirmation)
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
		    updated_at = NOW()
		WHERE id = :id`

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}
