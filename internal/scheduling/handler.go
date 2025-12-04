package scheduling

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hospital-emr/backend/internal/common/errors"
	"github.com/hospital-emr/backend/internal/models"
)

// Handler handles appointment scheduling HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new scheduling handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateAppointment godoc
// @Summary Create a new appointment
// @Description Schedule a new appointment for a patient
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAppointmentRequest true "Appointment information"
// @Success 201 {object} models.Appointment
// @Failure 400 {object} errors.AppError
// @Failure 409 {object} errors.AppError
// @Router /api/v1/appointments [post]
func (h *Handler) CreateAppointment(c *gin.Context) {
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	userIDValue, _ := c.Get("user_id")
	createdBy, _ := userIDValue.(uuid.UUID)

	appointment, err := h.service.CreateAppointment(c.Request.Context(), &req, createdBy)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

// GetAppointment godoc
// @Summary Get appointment by ID
// @Description Get detailed information about an appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Success 200 {object} models.Appointment
// @Failure 404 {object} errors.AppError
// @Router /api/v1/appointments/{id} [get]
func (h *Handler) GetAppointment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid appointment ID"))
		return
	}

	appointment, err := h.service.GetAppointment(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, appointment)
}

// ListAppointments godoc
// @Summary List appointments
// @Description Get a paginated list of appointments
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param patient_id query string false "Filter by patient ID"
// @Param provider_id query string false "Filter by provider ID"
// @Param status query string false "Filter by status"
// @Param date query string false "Filter by date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/appointments [get]
func (h *Handler) ListAppointments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var patientID *uuid.UUID
	if patientIDStr := c.Query("patient_id"); patientIDStr != "" {
		id, err := uuid.Parse(patientIDStr)
		if err == nil {
			patientID = &id
		}
	}

	var providerID *uuid.UUID
	if providerIDStr := c.Query("provider_id"); providerIDStr != "" {
		id, err := uuid.Parse(providerIDStr)
		if err == nil {
			providerID = &id
		}
	}

	var status *models.AppointmentStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s := models.AppointmentStatus(statusStr)
		status = &s
	}

	var date *time.Time
	if dateStr := c.Query("date"); dateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = &parsedDate
		}
	}

	appointments, total, err := h.service.ListAppointments(c.Request.Context(), page, pageSize, patientID, providerID, status, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        appointments,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// UpdateAppointment godoc
// @Summary Update appointment
// @Description Update appointment information
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Param request body CreateAppointmentRequest true "Appointment information"
// @Success 200 {object} models.Appointment
// @Failure 404 {object} errors.AppError
// @Router /api/v1/appointments/{id} [put]
func (h *Handler) UpdateAppointment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid appointment ID"))
		return
	}

	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	userIDValue, _ := c.Get("user_id")
	updatedBy, _ := userIDValue.(uuid.UUID)

	appointment, err := h.service.UpdateAppointment(c.Request.Context(), id, &req, updatedBy)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, appointment)
}

// CancelAppointment godoc
// @Summary Cancel appointment
// @Description Cancel an appointment
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Param request body map[string]string true "Cancellation reason"
// @Success 200 {object} map[string]string
// @Failure 404 {object} errors.AppError
// @Router /api/v1/appointments/{id}/cancel [post]
func (h *Handler) CancelAppointment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid appointment ID"))
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	userIDValue, _ := c.Get("user_id")
	cancelledBy, _ := userIDValue.(uuid.UUID)

	if err := h.service.CancelAppointment(c.Request.Context(), id, req.Reason, cancelledBy); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment cancelled successfully"})
}

// CheckInAppointment godoc
// @Summary Check in appointment
// @Description Mark appointment as checked in
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Appointment ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} errors.AppError
// @Router /api/v1/appointments/{id}/checkin [post]
func (h *Handler) CheckInAppointment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid appointment ID"))
		return
	}

	if err := h.service.CheckInAppointment(c.Request.Context(), id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checked in successfully"})
}

// GetProviderAvailability godoc
// @Summary Get provider availability
// @Description Get available time slots for a provider on a specific date
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Provider ID"
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} []TimeSlot
// @Failure 400 {object} errors.AppError
// @Router /api/v1/providers/{id}/availability [get]
func (h *Handler) GetProviderAvailability(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid provider ID"))
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Date is required"))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid date format (use YYYY-MM-DD)"))
		return
	}

	slots, err := h.service.GetProviderAvailability(c.Request.Context(), providerID, date)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, slots)
}

// GetAvailability godoc
// @Summary Get general availability
// @Description Get available time slots for a provider
// @Tags appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param provider_id query string true "Provider ID"
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} []TimeSlot
// @Failure 400 {object} errors.AppError
// @Router /api/v1/appointments/availability [get]
func (h *Handler) GetAvailability(c *gin.Context) {
	providerIDStr := c.Query("provider_id")
	if providerIDStr == "" {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Provider ID is required"))
		return
	}

	providerID, err := uuid.Parse(providerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid provider ID"))
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Date is required"))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid date format (use YYYY-MM-DD)"))
		return
	}

	slots, err := h.service.GetProviderAvailability(c.Request.Context(), providerID, date)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, slots)
}

