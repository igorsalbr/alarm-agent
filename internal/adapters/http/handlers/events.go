package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/alarm-agent/internal/adapters/http/dto"
	"github.com/alarm-agent/internal/adapters/http/middleware"
	"github.com/alarm-agent/internal/domain"
	"github.com/alarm-agent/internal/usecase"
)

type EventsHandler struct {
	eventUseCase *usecase.EventUseCase
}

func NewEventsHandler(eventUseCase *usecase.EventUseCase) *EventsHandler {
	return &EventsHandler{
		eventUseCase: eventUseCase,
	}
}

// CreateEvent creates a new event for the authenticated user
// POST /api/v1/events
func (h *EventsHandler) CreateEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Create event entities from request
	entities := &domain.EventEntities{
		Title:                  req.Title,
		Location:               req.Location,
		StartsAt:               &req.StartsAt,
		RemindBeforeMinutes:    req.RemindBeforeMinutes,
		RemindFrequencyMinutes: req.RemindFrequencyMinutes,
		RequireConfirmation:    req.RequireConfirmation,
		MaxNotifications:       req.MaxNotifications,
	}

	event, err := h.eventUseCase.CreateEvent(c.Request.Context(), userID, entities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "create_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Event created successfully",
		Data:    dto.EventToResponse(event),
	})
}

// GetEvent retrieves a specific event by ID
// GET /api/v1/events/:id
func (h *EventsHandler) GetEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_event_id",
			Message: "Invalid event ID format",
		})
		return
	}

	event, err := h.eventUseCase.GetEventByID(c.Request.Context(), userID, eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "event_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.EventToResponse(event))
}

// UpdateEvent updates an existing event
// PUT /api/v1/events/:id
func (h *EventsHandler) UpdateEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_event_id",
			Message: "Invalid event ID format",
		})
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Create event entities with ID for update
	entities := &domain.EventEntities{
		Identifier: &domain.EventIdentifier{
			ID: &eventID,
		},
		Title:                  req.Title,
		Location:               req.Location,
		StartsAt:               req.StartsAt,
		RemindBeforeMinutes:    req.RemindBeforeMinutes,
		RemindFrequencyMinutes: req.RemindFrequencyMinutes,
		RequireConfirmation:    req.RequireConfirmation,
		MaxNotifications:       req.MaxNotifications,
	}

	event, err := h.eventUseCase.UpdateEvent(c.Request.Context(), userID, entities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "update_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Event updated successfully",
		Data:    dto.EventToResponse(event),
	})
}

// DeleteEvent deletes an event
// DELETE /api/v1/events/:id
func (h *EventsHandler) DeleteEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_event_id",
			Message: "Invalid event ID format",
		})
		return
	}

	identifier := &domain.EventIdentifier{
		ID: &eventID,
	}

	event, err := h.eventUseCase.CancelEvent(c.Request.Context(), userID, identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "delete_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Event deleted successfully",
		Data:    dto.EventToResponse(event),
	})
}

// ListEvents retrieves events for the authenticated user
// GET /api/v1/events
func (h *EventsHandler) ListEvents(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var query dto.ListEventsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_query",
			Message: err.Error(),
		})
		return
	}

	// Set defaults
	limit := 50
	offset := 0
	if query.Limit != nil {
		limit = *query.Limit
	}
	if query.Offset != nil {
		offset = *query.Offset
	}

	events, err := h.eventUseCase.ListEvents(c.Request.Context(), userID, query.StartDate, query.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "list_failed",
			Message: err.Error(),
		})
		return
	}

	// Filter by status if provided
	if query.Status != nil {
		filtered := make([]domain.Event, 0)
		for _, event := range events {
			if string(event.Status) == *query.Status {
				filtered = append(filtered, event)
			}
		}
		events = filtered
	}

	// Apply pagination
	totalCount := len(events)
	start := offset
	end := offset + limit
	if start > len(events) {
		start = len(events)
	}
	if end > len(events) {
		end = len(events)
	}
	paginatedEvents := events[start:end]

	// Convert to response DTOs
	eventResponses := make([]dto.EventResponse, len(paginatedEvents))
	for i, event := range paginatedEvents {
		eventResponses[i] = dto.EventToResponse(&event)
	}

	response := dto.ListEventsResponse{
		Events:     eventResponses,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}

	c.JSON(http.StatusOK, response)
}

// ConfirmEvent confirms an event
// POST /api/v1/events/:id/confirm
func (h *EventsHandler) ConfirmEvent(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_event_id",
			Message: "Invalid event ID format",
		})
		return
	}

	identifier := &domain.EventIdentifier{
		ID: &eventID,
	}

	event, err := h.eventUseCase.ConfirmEvent(c.Request.Context(), userID, identifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "confirm_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Event confirmed successfully",
		Data:    dto.EventToResponse(event),
	})
}
