package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/alarm-agent/internal/adapters/http/dto"
	"github.com/alarm-agent/internal/ports"
)

const (
	UserContextKey   = "user"
	UserIDContextKey = "user_id"
)

type AuthMiddleware struct {
	userRepo ports.UserRepository
}

func NewAuthMiddleware(userRepo ports.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: userRepo,
	}
}

// AuthenticateByWANumber validates that the user exists and is active
func (a *AuthMiddleware) AuthenticateByWANumber() gin.HandlerFunc {
	return func(c *gin.Context) {
		waNumber := c.GetHeader("X-WA-Number")
		if waNumber == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "missing_wa_number",
				Message: "WhatsApp number is required in X-WA-Number header",
			})
			c.Abort()
			return
		}

		// Clean the phone number format
		waNumber = strings.TrimSpace(waNumber)

		// Get user by WhatsApp number
		user, err := a.userRepo.GetByWANumber(c.Request.Context(), waNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "user_lookup_failed",
				Message: "Failed to lookup user",
			})
			c.Abort()
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found. Please send a WhatsApp message first to create your account.",
			})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "user_inactive",
				Message: "User account is inactive",
			})
			c.Abort()
			return
		}

		// Store user in context
		c.Set(UserContextKey, user)
		c.Set(UserIDContextKey, user.ID)
		c.Next()
	}
}

// AuthenticateByUserID validates user ID from URL parameter
func (a *AuthMiddleware) AuthenticateByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("userId")
		if userIDStr == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "missing_user_id",
				Message: "User ID is required",
			})
			c.Abort()
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "invalid_user_id",
				Message: "Invalid user ID format",
			})
			c.Abort()
			return
		}

		// Get user by ID
		user, err := a.userRepo.GetByWANumber(c.Request.Context(), "")
		if err != nil {
			// Need to implement GetByID method
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "user_lookup_failed",
				Message: "Failed to lookup user",
			})
			c.Abort()
			return
		}

		if user == nil || user.ID != userID {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found",
			})
			c.Abort()
			return
		}

		// Store user in context
		c.Set(UserContextKey, user)
		c.Set(UserIDContextKey, userID)
		c.Next()
	}
}

// GetCurrentUser helper to extract user from context
func GetCurrentUser(c *gin.Context) interface{} {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return nil
	}
	return user
}

// GetCurrentUserID helper to extract user ID from context
func GetCurrentUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get(UserIDContextKey)
	if !exists {
		return 0, false
	}
	return userID.(int), true
}
