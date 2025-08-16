package domain

import (
	"time"
)

type EventStatus string

const (
	EventStatusScheduled EventStatus = "scheduled"
	EventStatusConfirmed EventStatus = "confirmed"
	EventStatusCanceled  EventStatus = "canceled"
	EventStatusCompleted EventStatus = "completed"
)

type Event struct {
	ID                     int         `json:"id" db:"id"`
	UserID                 int         `json:"user_id" db:"user_id"`
	Title                  string      `json:"title" db:"title"`
	Location               *string     `json:"location,omitempty" db:"location"`
	StartsAt               time.Time   `json:"starts_at" db:"starts_at"`
	RemindBeforeMinutes    int         `json:"remind_before_minutes" db:"remind_before_minutes"`
	RemindFrequencyMinutes int         `json:"remind_frequency_minutes" db:"remind_frequency_minutes"`
	RequireConfirmation    bool        `json:"require_confirmation" db:"require_confirmation"`
	MaxNotifications       int         `json:"max_notifications" db:"max_notifications"`
	Status                 EventStatus `json:"status" db:"status"`
	NotificationsSent      int         `json:"notifications_sent" db:"notifications_sent"`
	LastNotifiedAt         *time.Time  `json:"last_notified_at,omitempty" db:"last_notified_at"`
	CreatedAt              time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time   `json:"updated_at" db:"updated_at"`
}

type EventWithUser struct {
	Event
	User User `json:"user"`
}
