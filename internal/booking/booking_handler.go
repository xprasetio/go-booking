package booking

import (
	"booking/internal/booking/model"
	"booking/pkg/logger"
	"booking/pkg/response"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type BookingHandler struct {
	service BookingServiceInterface
	logger  logger.Logger
}

func NewBookingHandler(service BookingServiceInterface, logger logger.Logger) *BookingHandler {
	return &BookingHandler{
		service: service,
		logger:  logger,
	}
}

type CreateBookingRequest struct {
	SpaceID   uuid.UUID `json:"space_id"`
	StartDate string    `json:"start_date"` // Format: "2006-01-02"
	EndDate   string    `json:"end_date"`   // Format: "2006-01-02"
}

func (h *BookingHandler) Create(c echo.Context) error {
	var req CreateBookingRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request body", err)
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid start_date format. Use YYYY-MM-DD", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid end_date format. Use YYYY-MM-DD", err)
	}

	// Set time to check-in (14:00) and check-out (12:00)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 14, 0, 0, 0, time.Local)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 12, 0, 0, 0, time.Local)

	// Get user ID from context
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "invalid user id", err)
	}

	input := model.CreateBookingInput{
		UserID:    userID,
		SpaceID:   req.SpaceID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Validate input
	if input.SpaceID == uuid.Nil {
		return response.Error(c, http.StatusBadRequest, "space_id is required", nil)
	}

	if input.StartDate.IsZero() {
		return response.Error(c, http.StatusBadRequest, "start_date is required", nil)
	}

	if input.EndDate.IsZero() {
		return response.Error(c, http.StatusBadRequest, "end_date is required", nil)
	}

	if input.StartDate.After(input.EndDate) {
		return response.Error(c, http.StatusBadRequest, "start_date must be before end_date", nil)
	}

	booking, err := h.service.Create(c.Request().Context(), input)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"user_id":  userID,
			"space_id": input.SpaceID,
			"error":    err.Error(),
		}).Error(c.Request().Context(), "failed to create booking")
		return response.Error(c, http.StatusInternalServerError, err.Error(), err)
	}

	return response.Success(c, http.StatusCreated, "booking created successfully", booking.ToResponse())
}

func (h *BookingHandler) GetByID(c echo.Context) error {
	// Get user ID from context
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "invalid user id", err)
	}

	bookingID := c.Param("id")

	booking, err := h.service.GetByID(c.Request().Context(), bookingID)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"booking_id": bookingID,
			"error":      err.Error(),
		}).Error(c.Request().Context(), "failed to get booking")
		return response.Error(c, http.StatusNotFound, "booking not found", err)
	}

	// Check if user is authorized to view this booking
	if booking.UserID != userID {
		return response.Error(c, http.StatusForbidden, "you are not authorized to view this booking", nil)
	}

	return response.Success(c, http.StatusOK, "booking retrieved successfully", booking.ToResponse())
}

func (h *BookingHandler) GetAll(c echo.Context) error {
	// Get user ID from context
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "invalid user id", err)
	}

	bookings, err := h.service.GetAll(c.Request().Context(), userID)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error(c.Request().Context(), "failed to get bookings")
		return response.Error(c, http.StatusInternalServerError, "failed to get bookings", err)
	}

	// Convert bookings to response format
	var responseBookings []model.BookingResponse
	for _, booking := range bookings {
		responseBookings = append(responseBookings, booking.ToResponse())
	}

	return response.Success(c, http.StatusOK, "bookings retrieved successfully", responseBookings)
}

func (h *BookingHandler) Cancel(c echo.Context) error {
	// Get user ID from context
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "invalid user id", err)
	}

	bookingID := c.Param("id")

	err = h.service.Cancel(c.Request().Context(), bookingID, userID)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"booking_id": bookingID,
			"user_id":    userID,
			"error":      err.Error(),
		}).Error(c.Request().Context(), "failed to cancel booking")
		return response.Error(c, http.StatusInternalServerError, err.Error(), err)
	}

	return response.Success(c, http.StatusOK, "booking cancelled successfully", nil)
}
