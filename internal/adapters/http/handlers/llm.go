package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alarm-agent/internal/adapters/http/dto"
	"github.com/alarm-agent/internal/ports"
)

type LLMHandler struct {
	llmConfigRepo ports.LLMConfigRepository
}

func NewLLMHandler(llmConfigRepo ports.LLMConfigRepository) *LLMHandler {
	return &LLMHandler{
		llmConfigRepo: llmConfigRepo,
	}
}

// GetProviders retrieves all available LLM providers
// GET /api/v1/llm/providers
func (h *LLMHandler) GetProviders(c *gin.Context) {
	providers, err := h.llmConfigRepo.GetActiveProviders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "providers_fetch_failed",
			Message: err.Error(),
		})
		return
	}

	providerResponses := make([]dto.LLMProviderResponse, len(providers))
	for i, provider := range providers {
		providerResponse := dto.LLMProviderToResponse(&provider)

		// Get models for this provider
		models, err := h.llmConfigRepo.GetActiveModelsByProvider(c.Request.Context(), provider.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "models_fetch_failed",
				Message: err.Error(),
			})
			return
		}

		modelResponses := make([]dto.LLMModelResponse, len(models))
		for j, model := range models {
			modelResponses[j] = dto.LLMModelToResponse(&model)
		}
		providerResponse.Models = modelResponses

		providerResponses[i] = providerResponse
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Providers retrieved successfully",
		Data:    providerResponses,
	})
}

// GetModels retrieves models for a specific provider
// GET /api/v1/llm/providers/:providerId/models
func (h *LLMHandler) GetModels(c *gin.Context) {
	providerIDStr := c.Param("providerId")

	// For simplicity, we'll treat providerId as provider name
	// In a real implementation, you might want to convert to int ID
	providers, err := h.llmConfigRepo.GetActiveProviders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "providers_fetch_failed",
			Message: err.Error(),
		})
		return
	}

	var providerID int
	found := false
	for _, provider := range providers {
		if provider.Name == providerIDStr {
			providerID = provider.ID
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "provider_not_found",
			Message: "Provider not found",
		})
		return
	}

	models, err := h.llmConfigRepo.GetActiveModelsByProvider(c.Request.Context(), providerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "models_fetch_failed",
			Message: err.Error(),
		})
		return
	}

	modelResponses := make([]dto.LLMModelResponse, len(models))
	for i, model := range models {
		modelResponses[i] = dto.LLMModelToResponse(&model)
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Models retrieved successfully",
		Data:    modelResponses,
	})
}

// GetDefaultModel retrieves the default LLM model
// GET /api/v1/llm/default
func (h *LLMHandler) GetDefaultModel(c *gin.Context) {
	model, err := h.llmConfigRepo.GetDefaultModel(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "default_model_fetch_failed",
			Message: err.Error(),
		})
		return
	}

	modelResponse := dto.LLMModelToResponse(model)

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Default model retrieved successfully",
		Data:    modelResponse,
	})
}
