package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type UserAllowedContactRepository struct {
	db QueryExecutor
}

func NewUserAllowedContactRepository(db QueryExecutor) ports.UserAllowedContactRepository {
	return &UserAllowedContactRepository{db: db}
}

func (r *UserAllowedContactRepository) IsAllowed(ctx context.Context, userID int, contactNumber string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM user_allowed_contacts WHERE user_id = $1 AND contact_number = $2)"
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, userID, contactNumber)
	return exists, err
}

func (r *UserAllowedContactRepository) Add(ctx context.Context, contact *domain.UserAllowedContact) error {
	query := `
		INSERT INTO user_allowed_contacts (user_id, contact_number, note)
		VALUES (:user_id, :contact_number, :note)
		ON CONFLICT (user_id, contact_number) 
		DO UPDATE SET note = EXCLUDED.note, updated_at = NOW()
	`
	_, err := r.db.NamedExecContext(ctx, query, contact)
	return err
}

func (r *UserAllowedContactRepository) Remove(ctx context.Context, userID int, contactNumber string) error {
	query := "DELETE FROM user_allowed_contacts WHERE user_id = $1 AND contact_number = $2"
	_, err := r.db.ExecContext(ctx, query, userID, contactNumber)
	return err
}

func (r *UserAllowedContactRepository) List(ctx context.Context, userID int) ([]domain.UserAllowedContact, error) {
	var contacts []domain.UserAllowedContact
	query := "SELECT id, user_id, contact_number, note, created_at, updated_at FROM user_allowed_contacts WHERE user_id = $1 ORDER BY created_at DESC"
	
	err := r.db.SelectContext(ctx, &contacts, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []domain.UserAllowedContact{}, nil
		}
		return nil, fmt.Errorf("failed to list user allowed contacts: %w", err)
	}
	
	return contacts, nil
}

func (r *UserAllowedContactRepository) GetByUserAndNumber(ctx context.Context, userID int, contactNumber string) (*domain.UserAllowedContact, error) {
	var contact domain.UserAllowedContact
	query := "SELECT id, user_id, contact_number, note, created_at, updated_at FROM user_allowed_contacts WHERE user_id = $1 AND contact_number = $2"
	
	err := r.db.GetContext(ctx, &contact, query, userID, contactNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user allowed contact: %w", err)
	}
	
	return &contact, nil
}