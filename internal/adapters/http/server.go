package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/alarm-agent/internal/adapters/http/handlers"
	"github.com/alarm-agent/internal/adapters/http/middleware"
	"github.com/alarm-agent/internal/config"
	"github.com/alarm-agent/internal/ports"
	"github.com/alarm-agent/internal/usecase"
)

type Server struct {
	config       *config.Config
	repos        ports.Repositories
	eventUseCase *usecase.EventUseCase
	logger       *zap.Logger
	router       *gin.Engine
	server       *http.Server
}

func NewServer(
	cfg *config.Config,
	repos ports.Repositories,
	messageUseCase *usecase.MessageUseCase,
	eventUseCase *usecase.EventUseCase,
	verifier ports.WhatsAppWebhookVerifier,
	logger *zap.Logger,
) *Server {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(LoggingMiddleware(logger))
	router.Use(CORSMiddleware())
	router.Use(RateLimitMiddleware(cfg.Security.RateLimitPerMinute))
	router.Use(TimeoutMiddleware(30 * time.Second))

	server := &Server{
		config:       cfg,
		repos:        repos,
		eventUseCase: eventUseCase,
		logger:       logger,
		router:       router,
		server: &http.Server{
			Addr:    ":" + cfg.App.Port,
			Handler: router,
		},
	}

	server.setupRoutes(messageUseCase, verifier)
	return server
}

func (s *Server) setupRoutes(messageUseCase *usecase.MessageUseCase, verifier ports.WhatsAppWebhookVerifier) {
	webhookHandler := NewWebhookHandler(messageUseCase, verifier, s.logger)
	healthHandler := NewHealthHandler(s.repos, s.logger)

	s.router.GET("/health", healthHandler.Health)
	s.router.GET("/ready", healthHandler.Ready)

	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	webhookGroup := s.router.Group("/webhook")
	{
		webhookGroup.POST("/whatsapp", webhookHandler.HandleWhatsAppWebhook)
	}

	// Setup API routes
	s.setupAPIRoutes()

	if s.config.IsDevelopment() {
		s.setupDevelopmentRoutes()
	}
}

func (s *Server) setupAPIRoutes() {
	// Initialize handlers
	eventsHandler := handlers.NewEventsHandler(s.eventUseCase)
	usersHandler := handlers.NewUsersHandler(s.repos.User())
	llmHandler := handlers.NewLLMHandler(s.repos.LLMConfig())

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(s.repos.User(), s.repos.Whitelist())

	// Public routes (no authentication required)
	publicAPI := s.router.Group("/api/v1")
	{
		publicAPI.POST("/auth", usersHandler.AuthenticateUser)
		publicAPI.GET("/llm/providers", llmHandler.GetProviders)
		publicAPI.GET("/llm/providers/:providerId/models", llmHandler.GetModels)
		publicAPI.GET("/llm/default", llmHandler.GetDefaultModel)
	}

	// Protected routes (authentication required)
	protectedAPI := s.router.Group("/api/v1")
	protectedAPI.Use(authMiddleware.AuthenticateByWANumber())
	{
		// User profile routes
		protectedAPI.GET("/profile", usersHandler.GetProfile)
		protectedAPI.PUT("/profile", usersHandler.UpdateProfile)

		// Events routes
		protectedAPI.POST("/events", eventsHandler.CreateEvent)
		protectedAPI.GET("/events", eventsHandler.ListEvents)
		protectedAPI.GET("/events/:id", eventsHandler.GetEvent)
		protectedAPI.PUT("/events/:id", eventsHandler.UpdateEvent)
		protectedAPI.DELETE("/events/:id", eventsHandler.DeleteEvent)
		protectedAPI.POST("/events/:id/confirm", eventsHandler.ConfirmEvent)
	}
}

func (s *Server) setupDevelopmentRoutes() {
	devGroup := s.router.Group("/dev")
	{
		devGroup.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message":     "Development server is running",
				"environment": s.config.App.Environment,
				"version":     "1.0.0",
			})
		})
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server",
		zap.String("address", s.server.Addr),
		zap.String("environment", s.config.App.Environment),
	)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")
	return s.server.Shutdown(ctx)
}
