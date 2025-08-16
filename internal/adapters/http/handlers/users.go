package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alarm-agent/internal/adapters/http/dto"
	"github.com/alarm-agent/internal/adapters/http/middleware"
	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/ports"
)

type UsersHandler struct {
	userRepo ports.UserRepository
}

func NewUsersHandler(userRepo ports.UserRepository) *UsersHandler {
	return &UsersHandler{
		userRepo: userRepo,
	}
}

// GetProfile retrieves the current user's profile
// GET /api/v1/profile
func (h *UsersHandler) GetProfile(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	domainUser := user.(*domain.User)
	c.JSON(http.StatusOK, dto.UserToResponse(domainUser))
}

// UpdateProfile updates the current user's profile
// PUT /api/v1/profile
func (h *UsersHandler) UpdateProfile(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	domainUser := user.(*domain.User)

	var req dto.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Update user fields if provided
	if req.Name != nil {
		domainUser.Name = req.Name
	}
	if req.Timezone != nil {
		domainUser.Timezone = *req.Timezone
	}
	if req.DefaultRemindBeforeMinutes != nil {
		domainUser.DefaultRemindBeforeMinutes = *req.DefaultRemindBeforeMinutes
	}
	if req.DefaultRemindFrequencyMinutes != nil {
		domainUser.DefaultRemindFrequencyMinutes = *req.DefaultRemindFrequencyMinutes
	}
	if req.DefaultRequireConfirmation != nil {
		domainUser.DefaultRequireConfirmation = *req.DefaultRequireConfirmation
	}
	if req.LLMProvider != nil {
		domainUser.LLMProvider = req.LLMProvider
	}
	if req.LLMModel != nil {
		domainUser.LLMModel = req.LLMModel
	}

	if err := h.userRepo.Update(c.Request.Context(), domainUser); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Profile updated successfully",
		Data:    dto.UserToResponse(domainUser),
	})
}

// AuthenticateUser creates or authenticates a user by WhatsApp number
// POST /api/v1/auth
func (h *UsersHandler) AuthenticateUser(c *gin.Context) {
	var req dto.AuthenticateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid_request",
			Message: err.Error(),
		})
		return
	}

	user, err := h.userRepo.GetByWANumber(c.Request.Context(), req.WANumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "lookup_failed",
			Message: "Failed to lookup user",
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "user_not_found",
			Message: "User not found. Please send a WhatsApp message first to create your account.",
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "User authenticated successfully",
		Data:    dto.UserToResponse(user),
	})
}