package infra

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	WhatsAppMessagesReceived prometheus.Counter
	WhatsAppMessagesSent     prometheus.Counter
	LLMRequestsTotal         *prometheus.CounterVec
	EventsCreatedTotal       prometheus.Counter
	RemindersSentTotal       prometheus.Counter
	HTTPRequestDuration      *prometheus.HistogramVec
	HTTPRequestsTotal        *prometheus.CounterVec
	ActiveEvents             prometheus.Gauge
	DatabaseConnections      prometheus.Gauge
}

func NewMetrics() *Metrics {
	return &Metrics{
		WhatsAppMessagesReceived: promauto.NewCounter(prometheus.CounterOpts{
			Name: "whatsapp_messages_received_total",
			Help: "Total number of WhatsApp messages received",
		}),
		WhatsAppMessagesSent: promauto.NewCounter(prometheus.CounterOpts{
			Name: "whatsapp_messages_sent_total",
			Help: "Total number of WhatsApp messages sent",
		}),
		LLMRequestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "llm_requests_total",
			Help: "Total number of LLM requests",
		}, []string{"provider", "model", "intent", "status"}),
		EventsCreatedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "events_created_total",
			Help: "Total number of events created",
		}),
		RemindersSentTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "reminders_sent_total",
			Help: "Total number of reminders sent",
		}),
		HTTPRequestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		}, []string{"method", "path", "status_code"}),
		HTTPRequestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		}, []string{"method", "path", "status_code"}),
		ActiveEvents: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "active_events_total",
			Help: "Number of active (scheduled/confirmed) events",
		}),
		DatabaseConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		}),
	}
}
