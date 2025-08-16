package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alarm-agent/internal/ports"
)

type HealthHandler struct {
	repos  ports.Repositories
	logger *zap.Logger
}

func NewHealthHandler(repos ports.Repositories, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		repos:  repos,
		logger: logger,
	}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "alarm-agent",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	checks := map[string]string{}
	
	if _, err := h.repos.User().GetByWANumber(ctx, "health-check"); err != nil {
		if err.Error() != "sql: no rows in result set" {
			checks["database"] = "unhealthy"
			h.logger.Error("Database health check failed", zap.Error(err))
		} else {
			checks["database"] = "healthy"
		}
	} else {
		checks["database"] = "healthy"
	}

	overall := "healthy"
	statusCode := http.StatusOK
	
	for _, status := range checks {
		if status != "healthy" {
			overall = "unhealthy"
			statusCode = http.StatusServiceUnavailable
			break
		}
	}

	c.JSON(statusCode, gin.H{
		"status":    overall,
		"timestamp": time.Now().UTC(),
		"checks":    checks,
	})
}