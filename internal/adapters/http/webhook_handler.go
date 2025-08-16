package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alarm-agent/internal/adapters/whatsapp"
	"github.com/alarm-agent/internal/ports"
	"github.com/alarm-agent/internal/usecase"
)

type WebhookHandler struct {
	messageUseCase *usecase.MessageUseCase
	verifier       ports.WhatsAppWebhookVerifier
	logger         *zap.Logger
}

func NewWebhookHandler(
	messageUseCase *usecase.MessageUseCase,
	verifier ports.WhatsAppWebhookVerifier,
	logger *zap.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		messageUseCase: messageUseCase,
		verifier:       verifier,
		logger:         logger,
	}
}

func (h *WebhookHandler) HandleWhatsAppWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error("Failed to read request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	signature := c.GetHeader("X-Signature-256")
	if !h.verifier.VerifySignature(body, signature) {
		h.logger.Warn("Invalid webhook signature", 
			zap.String("signature", signature),
			zap.String("from", c.ClientIP()),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	var webhookRequest whatsapp.InfobipWebhookRequest
	if err := c.ShouldBindJSON(&webhookRequest); err != nil {
		h.logger.Error("Failed to parse webhook payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	messages := webhookRequest.ExtractMessages()
	h.logger.Info("Received WhatsApp messages", zap.Int("count", len(messages)))

	for _, message := range messages {
		go func(msg whatsapp.ParsedMessage) {
			if err := h.messageUseCase.ProcessInboundMessage(c.Request.Context(), msg); err != nil {
				h.logger.Error("Failed to process inbound message",
					zap.Error(err),
					zap.String("message_id", msg.ID),
					zap.String("from", msg.From),
				)
			}
		}(message)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}