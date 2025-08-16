package repo

import (
	"context"

	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type InboundMessageRepository struct {
	db QueryExecutor
}

func NewInboundMessageRepository(db QueryExecutor) ports.InboundMessageRepository {
	return &InboundMessageRepository{db: db}
}

func (r *InboundMessageRepository) Create(ctx context.Context, message *domain.InboundMessage) error {
	query := `
		INSERT INTO inbound_messages (provider_message_id, from_number, raw_payload)
		VALUES (:provider_message_id, :from_number, :raw_payload)
		RETURNING id, processed_at, created_at`
	
	rows, err := r.db.NamedExecContext(ctx, query, message)
	if err != nil {
		return err
	}
	
	id, err := rows.LastInsertId()
	if err == nil && id > 0 {
		message.ID = int(id)
	}
	
	return nil
}

func (r *InboundMessageRepository) Exists(ctx context.Context, providerMessageID string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM inbound_messages WHERE provider_message_id = $1)"
	
	err := r.db.GetContext(ctx, &exists, query, providerMessageID)
	return exists, err
}