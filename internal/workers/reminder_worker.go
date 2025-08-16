package workers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type ReminderWorker struct {
	repos          ports.Repositories
	whatsappSender ports.WhatsAppSender
	timeProvider   ports.TimeProvider
	logger         *zap.Logger
	tickInterval   time.Duration
	stopCh         chan struct{}
}

func NewReminderWorker(
	repos ports.Repositories,
	whatsappSender ports.WhatsAppSender,
	timeProvider ports.TimeProvider,
	logger *zap.Logger,
	tickInterval time.Duration,
) *ReminderWorker {
	return &ReminderWorker{
		repos:          repos,
		whatsappSender: whatsappSender,
		timeProvider:   timeProvider,
		logger:         logger,
		tickInterval:   tickInterval,
		stopCh:         make(chan struct{}),
	}
}

func (w *ReminderWorker) Start(ctx context.Context) error {
	w.logger.Info("Starting reminder worker", zap.Duration("tick_interval", w.tickInterval))
	
	ticker := time.NewTicker(w.tickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Reminder worker stopped by context")
			return ctx.Err()
		case <-w.stopCh:
			w.logger.Info("Reminder worker stopped")
			return nil
		case <-ticker.C:
			if err := w.processReminders(ctx); err != nil {
				w.logger.Error("Failed to process reminders", zap.Error(err))
			}
		}
	}
}

func (w *ReminderWorker) Stop() {
	close(w.stopCh)
}

func (w *ReminderWorker) processReminders(ctx context.Context) error {
	reminderWindow := 30 * time.Minute
	eventsWithUsers, err := w.repos.Event().GetPendingReminders(ctx, reminderWindow)
	if err != nil {
		return fmt.Errorf("failed to get pending reminders: %w", err)
	}

	if len(eventsWithUsers) == 0 {
		return nil
	}

	w.logger.Info("Processing reminders", zap.Int("count", len(eventsWithUsers)))

	for _, eventWithUser := range eventsWithUsers {
		if err := w.processEventReminder(ctx, &eventWithUser); err != nil {
			w.logger.Error("Failed to process event reminder",
				zap.Error(err),
				zap.Int("event_id", eventWithUser.Event.ID),
				zap.String("user_number", eventWithUser.User.WANumber),
			)
		}
	}

	return nil
}

func (w *ReminderWorker) processEventReminder(ctx context.Context, eventWithUser *domain.EventWithUser) error {
	event := &eventWithUser.Event
	user := &eventWithUser.User

	now := w.timeProvider.Now()
	reminderTime := event.StartsAt.Add(-time.Duration(event.RemindBeforeMinutes) * time.Minute)

	if now.Before(reminderTime) {
		return nil
	}

	if event.NotificationsSent >= event.MaxNotifications {
		return nil
	}

	var message string
	if event.RequireConfirmation && event.Status == domain.EventStatusScheduled {
		message = w.buildConfirmationMessage(event)
	} else {
		message = w.buildReminderMessage(event)
	}

	if err := w.whatsappSender.SendText(ctx, user.WANumber, message); err != nil {
		return fmt.Errorf("failed to send reminder message: %w", err)
	}

	event.NotificationsSent++
	event.LastNotifiedAt = &now

	if err := w.repos.Event().Update(ctx, event); err != nil {
		return fmt.Errorf("failed to update event after sending reminder: %w", err)
	}

	w.logger.Info("Sent reminder",
		zap.Int("event_id", event.ID),
		zap.String("user_number", user.WANumber),
		zap.String("event_title", event.Title),
		zap.Int("notifications_sent", event.NotificationsSent),
	)

	return nil
}

func (w *ReminderWorker) buildReminderMessage(event *domain.Event) string {
	var parts []string
	parts = append(parts, "⏰ *Lembrete de Compromisso*")
	parts = append(parts, fmt.Sprintf("📅 %s", event.Title))
	parts = append(parts, fmt.Sprintf("🕐 %s", event.StartsAt.Format("02/01/2006 15:04")))

	if event.Location != nil {
		parts = append(parts, fmt.Sprintf("📍 %s", *event.Location))
	}

	timeUntil := time.Until(event.StartsAt)
	if timeUntil > 0 {
		if timeUntil < time.Hour {
			minutes := int(timeUntil.Minutes())
			parts = append(parts, fmt.Sprintf("⏱️ Começa em %d minutos", minutes))
		} else {
			hours := int(timeUntil.Hours())
			parts = append(parts, fmt.Sprintf("⏱️ Começa em %d horas", hours))
		}
	}

	return strings.Join(parts, "\n")
}

func (w *ReminderWorker) buildConfirmationMessage(event *domain.Event) string {
	var parts []string
	parts = append(parts, "❓ *Confirmação de Compromisso*")
	parts = append(parts, fmt.Sprintf("📅 %s", event.Title))
	parts = append(parts, fmt.Sprintf("🕐 %s", event.StartsAt.Format("02/01/2006 15:04")))

	if event.Location != nil {
		parts = append(parts, fmt.Sprintf("📍 %s", *event.Location))
	}

	parts = append(parts, "")
	parts = append(parts, "Por favor, confirme sua presença:")
	parts = append(parts, "✅ Responda 'OK' ou 'Confirmo' para confirmar")
	parts = append(parts, "❌ Responda 'Cancelar' para cancelar")

	return strings.Join(parts, "\n")
}