package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type EventRepository struct {
	db QueryExecutor
}

func NewEventRepository(db QueryExecutor) ports.EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, event *domain.Event) error {
	query := `
		INSERT INTO events (user_id, title, location, starts_at, remind_before_minutes, 
		                   remind_frequency_minutes, require_confirmation, max_notifications, status)
		VALUES (:user_id, :title, :location, :starts_at, :remind_before_minutes, 
		        :remind_frequency_minutes, :require_confirmation, :max_notifications, :status)
		RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return err
	}

	id, err := rows.LastInsertId()
	if err == nil && id > 0 {
		event.ID = int(id)
	}

	return nil
}

func (r *EventRepository) Update(ctx context.Context, event *domain.Event) error {
	query := `
		UPDATE events 
		SET title = :title, location = :location, starts_at = :starts_at,
		    remind_before_minutes = :remind_before_minutes,
		    remind_frequency_minutes = :remind_frequency_minutes,
		    require_confirmation = :require_confirmation,
		    max_notifications = :max_notifications,
		    status = :status,
		    notifications_sent = :notifications_sent,
		    last_notified_at = :last_notified_at,
		    updated_at = NOW()
		WHERE id = :id`

	_, err := r.db.NamedExecContext(ctx, query, event)
	return err
}

func (r *EventRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM events WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *EventRepository) GetByID(ctx context.Context, id int) (*domain.Event, error) {
	var event domain.Event
	query := `
		SELECT id, user_id, title, location, starts_at, remind_before_minutes,
		       remind_frequency_minutes, require_confirmation, max_notifications,
		       status, notifications_sent, last_notified_at, created_at, updated_at
		FROM events 
		WHERE id = $1`

	err := r.db.GetContext(ctx, &event, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

func (r *EventRepository) GetByUserID(ctx context.Context, userID int) ([]domain.Event, error) {
	var events []domain.Event
	query := `
		SELECT id, user_id, title, location, starts_at, remind_before_minutes,
		       remind_frequency_minutes, require_confirmation, max_notifications,
		       status, notifications_sent, last_notified_at, created_at, updated_at
		FROM events 
		WHERE user_id = $1
		ORDER BY starts_at ASC`

	err := r.db.SelectContext(ctx, &events, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Event{}, nil
		}
		return nil, err
	}

	return events, nil
}

func (r *EventRepository) GetByUserIDAndDateRange(ctx context.Context, userID int, start, end time.Time) ([]domain.Event, error) {
	var events []domain.Event
	query := `
		SELECT id, user_id, title, location, starts_at, remind_before_minutes,
		       remind_frequency_minutes, require_confirmation, max_notifications,
		       status, notifications_sent, last_notified_at, created_at, updated_at
		FROM events 
		WHERE user_id = $1 AND starts_at BETWEEN $2 AND $3
		ORDER BY starts_at ASC`

	err := r.db.SelectContext(ctx, &events, query, userID, start, end)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Event{}, nil
		}
		return nil, err
	}

	return events, nil
}

func (r *EventRepository) GetPendingReminders(ctx context.Context, reminderWindow time.Duration) ([]domain.EventWithUser, error) {
	var eventsWithUsers []domain.EventWithUser
	now := time.Now()
	windowEnd := now.Add(reminderWindow)

	query := `
		SELECT 
		    e.id, e.user_id, e.title, e.location, e.starts_at, e.remind_before_minutes,
		    e.remind_frequency_minutes, e.require_confirmation, e.max_notifications,
		    e.status, e.notifications_sent, e.last_notified_at, e.created_at, e.updated_at,
		    u.id as "user.id", u.wa_number as "user.wa_number", u.name as "user.name", 
		    u.timezone as "user.timezone", u.default_remind_before_minutes as "user.default_remind_before_minutes",
		    u.default_remind_frequency_minutes as "user.default_remind_frequency_minutes",
		    u.default_require_confirmation as "user.default_require_confirmation",
		    u.created_at as "user.created_at", u.updated_at as "user.updated_at"
		FROM events e
		JOIN users u ON e.user_id = u.id
		WHERE e.status IN ('scheduled', 'confirmed')
		  AND e.notifications_sent < e.max_notifications
		  AND (e.starts_at - INTERVAL '1 minute' * e.remind_before_minutes) BETWEEN $1 AND $2
		  AND (e.last_notified_at IS NULL 
		       OR e.last_notified_at <= $1 - INTERVAL '1 minute' * e.remind_frequency_minutes)
		ORDER BY e.starts_at ASC`

	err := r.db.SelectContext(ctx, &eventsWithUsers, query, now, windowEnd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.EventWithUser{}, nil
		}
		return nil, err
	}

	return eventsWithUsers, nil
}

func (r *EventRepository) FindByUserAndIdentifier(ctx context.Context, userID int, identifier *domain.EventIdentifier) ([]domain.Event, error) {
	var events []domain.Event
	var conditions []string
	var args []interface{}

	baseQuery := `
		SELECT id, user_id, title, location, starts_at, remind_before_minutes,
		       remind_frequency_minutes, require_confirmation, max_notifications,
		       status, notifications_sent, last_notified_at, created_at, updated_at
		FROM events 
		WHERE user_id = $1`

	args = append(args, userID)
	argIndex := 2

	if identifier.EventID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *identifier.EventID)
		argIndex++
	}

	if identifier.Title != nil {
		conditions = append(conditions, fmt.Sprintf("LOWER(title) LIKE LOWER($%d)", argIndex))
		args = append(args, "%"+*identifier.Title+"%")
		argIndex++
	}

	if identifier.DateHint != nil {
		conditions = append(conditions, fmt.Sprintf("DATE(starts_at) = $%d", argIndex))
		args = append(args, *identifier.DateHint)
		argIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " AND (" + strings.Join(conditions, " OR ") + ")"
	}

	baseQuery += " ORDER BY starts_at ASC"

	err := r.db.SelectContext(ctx, &events, baseQuery, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Event{}, nil
		}
		return nil, err
	}

	return events, nil
}
