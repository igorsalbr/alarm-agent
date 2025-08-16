package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type WhitelistRepository struct {
	db QueryExecutor
}

func NewWhitelistRepository(db QueryExecutor) ports.WhitelistRepository {
	return &WhitelistRepository{db: db}
}

func (r *WhitelistRepository) IsWhitelisted(ctx context.Context, number string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM whitelist_numbers WHERE number = $1)"

	err := r.db.GetContext(ctx, &exists, query, number)
	return exists, err
}

func (r *WhitelistRepository) Add(ctx context.Context, whitelist *domain.WhitelistNumber) error {
	query := `
		INSERT INTO whitelist_numbers (number, note)
		VALUES (:number, :note)
		ON CONFLICT (number) DO UPDATE SET note = :note`

	_, err := r.db.NamedExecContext(ctx, query, whitelist)
	return err
}

func (r *WhitelistRepository) Remove(ctx context.Context, number string) error {
	query := "DELETE FROM whitelist_numbers WHERE number = $1"
	_, err := r.db.ExecContext(ctx, query, number)
	return err
}

func (r *WhitelistRepository) List(ctx context.Context) ([]domain.WhitelistNumber, error) {
	var whitelist []domain.WhitelistNumber
	query := "SELECT number, note, created_at FROM whitelist_numbers ORDER BY created_at DESC"

	err := r.db.SelectContext(ctx, &whitelist, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.WhitelistNumber{}, nil
		}
		return nil, err
	}

	return whitelist, nil
}
