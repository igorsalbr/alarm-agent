package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/alarm-agent/internal/adapters/whatsapp"
	"github.com/alarm-agent/internal/usecase"
)

type MockMessageUseCase struct {
	mock.Mock
}

func (m *MockMessageUseCase) ProcessInboundMessage(ctx context.Context, message whatsapp.ParsedMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

type MockWebhookVerifier struct {
	mock.Mock
}

func (m *MockWebhookVerifier) VerifySignature(payload []byte, signature string) bool {
	args := m.Called(payload, signature)
	return args.Bool(0)
}

func TestWebhookHandler_HandleWhatsAppWebhook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockMessageUseCase := &MockMessageUseCase{}
	mockVerifier := &MockWebhookVerifier{}
	logger, _ := zap.NewDevelopment()

	handler := NewWebhookHandler(mockMessageUseCase, mockVerifier, logger)

	payload := `{
		"results": [{
			"messageId": "test-123",
			"from": "5511999999999",
			"to": "5511888888888",
			"receivedAt": "2024-01-01T10:00:00Z",
			"message": {
				"type": "TEXT",
				"text": "Test message"
			}
		}]
	}`

	mockVerifier.On("VerifySignature", mock.AnythingOfType("[]uint8"), "valid-signature").Return(true)
	mockMessageUseCase.On("ProcessInboundMessage", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("whatsapp.ParsedMessage")).Return(nil)

	router := gin.New()
	router.POST("/webhook/whatsapp", handler.HandleWhatsAppWebhook)

	req := httptest.NewRequest("POST", "/webhook/whatsapp", bytes.NewString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature-256", "valid-signature")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockVerifier.AssertExpectations(t)
}

func TestWebhookHandler_HandleWhatsAppWebhook_InvalidSignature(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockMessageUseCase := &MockMessageUseCase{}
	mockVerifier := &MockWebhookVerifier{}
	logger, _ := zap.NewDevelopment()

	handler := NewWebhookHandler(mockMessageUseCase, mockVerifier, logger)

	payload := `{"results": []}`

	mockVerifier.On("VerifySignature", mock.AnythingOfType("[]uint8"), "invalid-signature").Return(false)

	router := gin.New()
	router.POST("/webhook/whatsapp", handler.HandleWhatsAppWebhook)

	req := httptest.NewRequest("POST", "/webhook/whatsapp", bytes.NewString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature-256", "invalid-signature")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockVerifier.AssertExpectations(t)
	mockMessageUseCase.AssertNotCalled(t, "ProcessInboundMessage")
}

func TestWebhookHandler_HandleWhatsAppWebhook_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockMessageUseCase := &MockMessageUseCase{}
	mockVerifier := &MockWebhookVerifier{}
	logger, _ := zap.NewDevelopment()

	handler := NewWebhookHandler(mockMessageUseCase, mockVerifier, logger)

	payload := `invalid json`

	mockVerifier.On("VerifySignature", mock.AnythingOfType("[]uint8"), "valid-signature").Return(true)

	router := gin.New()
	router.POST("/webhook/whatsapp", handler.HandleWhatsAppWebhook)

	req := httptest.NewRequest("POST", "/webhook/whatsapp", bytes.NewString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature-256", "valid-signature")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockVerifier.AssertExpectations(t)
	mockMessageUseCase.AssertNotCalled(t, "ProcessInboundMessage")
}
