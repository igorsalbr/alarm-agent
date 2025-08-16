package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/alarm-agent/internal/adapters/http/dto"
	"github.com/alarm-agent/internal/adapters/http/middleware"
	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type UserConfigHandler struct {
	repos ports.Repositories
}

func NewUserConfigHandler(repos ports.Repositories) *UserConfigHandler {
	return &UserConfigHandler{
		repos: repos,
	}
}

// GetUserConfig retrieves the current user's configuration
// GET /api/v1/user/config
func (h *UserConfigHandler) GetUserConfig(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	user, err := h.repos.User().GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve user configuration",
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "user_not_found",
			Message: "User not found",
		})
		return
	}

	response := dto.UserConfigResponse{
		UserID:                        user.ID,
		Name:                          user.Name,
		Timezone:                      user.Timezone,
		DefaultRemindBeforeMinutes:    user.DefaultRemindBeforeMinutes,
		DefaultRemindFrequencyMinutes: user.DefaultRemindFrequencyMinutes,
		DefaultRequireConfirmation:    user.DefaultRequireConfirmation,
		LLMProvider:                   user.LLMProvider,
		LLMModel:                      user.LLMModel,
		RateLimitPerMinute:            user.RateLimitPerMinute,
		IsActive:                      user.IsActive,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUserConfig updates the current user's configuration
// PUT /api/v1/user/config
func (h *UserConfigHandler) UpdateUserConfig(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req dto.UpdateUserConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Get current user to fill in unchanged fields
	user, err := h.repos.User().GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve current configuration",
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "user_not_found",
			Message: "User not found",
		})
		return
	}

	// Apply updates (only change non-nil fields)
	config := &domain.UserConfig{
		UserID:                        userID,
		Name:                          user.Name,
		Timezone:                      user.Timezone,
		DefaultRemindBeforeMinutes:    user.DefaultRemindBeforeMinutes,
		DefaultRemindFrequencyMinutes: user.DefaultRemindFrequencyMinutes,
		DefaultRequireConfirmation:    user.DefaultRequireConfirmation,
		LLMProvider:                   user.LLMProvider,
		LLMModel:                      user.LLMModel,
		RateLimitPerMinute:            user.RateLimitPerMinute,
		IsActive:                      user.IsActive,
	}

	if req.Name != nil {
		config.Name = req.Name
	}
	if req.Timezone != nil {
		config.Timezone = *req.Timezone
	}
	if req.DefaultRemindBeforeMinutes != nil {
		config.DefaultRemindBeforeMinutes = *req.DefaultRemindBeforeMinutes
	}
	if req.DefaultRemindFrequencyMinutes != nil {
		config.DefaultRemindFrequencyMinutes = *req.DefaultRemindFrequencyMinutes
	}
	if req.DefaultRequireConfirmation != nil {
		config.DefaultRequireConfirmation = *req.DefaultRequireConfirmation
	}
	if req.LLMProvider != nil {
		config.LLMProvider = req.LLMProvider
	}
	if req.LLMModel != nil {
		config.LLMModel = req.LLMModel
	}
	if req.RateLimitPerMinute != nil {
		config.RateLimitPerMinute = *req.RateLimitPerMinute
	}
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	}

	if err := h.repos.User().UpdateConfig(c.Request.Context(), userID, config); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to update configuration",
		})
		return
	}

	response := dto.UserConfigResponse{
		UserID:                        config.UserID,
		Name:                          config.Name,
		Timezone:                      config.Timezone,
		DefaultRemindBeforeMinutes:    config.DefaultRemindBeforeMinutes,
		DefaultRemindFrequencyMinutes: config.DefaultRemindFrequencyMinutes,
		DefaultRequireConfirmation:    config.DefaultRequireConfirmation,
		LLMProvider:                   config.LLMProvider,
		LLMModel:                      config.LLMModel,
		RateLimitPerMinute:            config.RateLimitPerMinute,
		IsActive:                      config.IsActive,
	}

	c.JSON(http.StatusOK, response)
}

// ListAllowedContacts lists the current user's allowed contacts
// GET /api/v1/user/allowed-contacts
func (h *UserConfigHandler) ListAllowedContacts(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	contacts, err := h.repos.UserAllowedContact().List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve allowed contacts",
		})
		return
	}

	response := make([]dto.AllowedContactResponse, len(contacts))
	for i, contact := range contacts {
		response[i] = dto.AllowedContactResponse{
			ID:            contact.ID,
			ContactNumber: contact.ContactNumber,
			Note:          contact.Note,
			CreatedAt:     contact.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     contact.UpdatedAt.Format(time.RFC3339),
		}
	}

	c.JSON(http.StatusOK, response)
}

// AddAllowedContact adds a new allowed contact for the current user
// POST /api/v1/user/allowed-contacts
func (h *UserConfigHandler) AddAllowedContact(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req dto.AddAllowedContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	contact := &domain.UserAllowedContact{
		UserID:        userID,
		ContactNumber: req.ContactNumber,
		Note:          req.Note,
	}

	if err := h.repos.UserAllowedContact().Add(c.Request.Context(), contact); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to add allowed contact",
		})
		return
	}

	// Retrieve the added contact to get the full details
	addedContact, err := h.repos.UserAllowedContact().GetByUserAndNumber(c.Request.Context(), userID, req.ContactNumber)
	if err != nil || addedContact == nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Contact added but failed to retrieve details",
		})
		return
	}

	response := dto.AllowedContactResponse{
		ID:            addedContact.ID,
		ContactNumber: addedContact.ContactNumber,
		Note:          addedContact.Note,
		CreatedAt:     addedContact.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     addedContact.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, response)
}

// RemoveAllowedContact removes an allowed contact for the current user
// DELETE /api/v1/user/allowed-contacts/:contactNumber
func (h *UserConfigHandler) RemoveAllowedContact(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	contactNumber := c.Param("contactNumber")
	if contactNumber == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Contact number is required",
		})
		return
	}

	// Check if contact exists for this user
	existing, err := h.repos.UserAllowedContact().GetByUserAndNumber(c.Request.Context(), userID, contactNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to check contact existence",
		})
		return
	}

	if existing == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "contact_not_found",
			Message: "Allowed contact not found",
		})
		return
	}

	if err := h.repos.UserAllowedContact().Remove(c.Request.Context(), userID, contactNumber); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to remove allowed contact",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Allowed contact removed successfully"})
}
