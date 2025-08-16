package infra

import (
	"time"

	"github.com/alarm-agent/internal/ports"
)

type RealTimeProvider struct{}

func NewRealTimeProvider() ports.TimeProvider {
	return &RealTimeProvider{}
}

func (p *RealTimeProvider) Now() time.Time {
	return time.Now()
}

func (p *RealTimeProvider) Sleep(duration time.Duration) {
	time.Sleep(duration)
}
