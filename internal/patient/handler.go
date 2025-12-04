package patient

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hospital-emr/backend/internal/common/errors"
	_ "github.com/hospital-emr/backend/internal/models"
)

// Handler handles patient HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new patient handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreatePatient godoc
// @Summary Create a new patient
// @Description Register a new patient in the system
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePatientRequest true "Patient information"
// @Success 201 {object} models.Patient
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /api/v1/patients [post]
func (h *Handler) CreatePatient(c *gin.Context) {
	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	userIDValue, _ := c.Get("user_id")
	createdBy, _ := userIDValue.(uuid.UUID)

	patient, err := h.service.CreatePatient(c.Request.Context(), &req, createdBy)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// GetPatient godoc
// @Summary Get patient by ID
// @Description Get detailed information about a patient
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Patient ID"
// @Success 200 {object} models.Patient
// @Failure 404 {object} errors.AppError
// @Router /api/v1/patients/{id} [get]
func (h *Handler) GetPatient(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid patient ID"))
		return
	}

	patient, err := h.service.GetPatient(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, patient)
}

// ListPatients godoc
// @Summary List patients
// @Description Get a paginated list of patients
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search term"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/patients [get]
func (h *Handler) ListPatients(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	patients, total, err := h.service.ListPatients(c.Request.Context(), page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       patients,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// UpdatePatient godoc
// @Summary Update patient
// @Description Update patient information
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Patient ID"
// @Param request body CreatePatientRequest true "Patient information"
// @Success 200 {object} models.Patient
// @Failure 404 {object} errors.AppError
// @Router /api/v1/patients/{id} [put]
func (h *Handler) UpdatePatient(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid patient ID"))
		return
	}

	var req CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	userIDValue, _ := c.Get("user_id")
	updatedBy, _ := userIDValue.(uuid.UUID)

	patient, err := h.service.UpdatePatient(c.Request.Context(), id, &req, updatedBy)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, patient)
}

// DeletePatient godoc
// @Summary Delete patient
// @Description Soft delete a patient
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Patient ID"
// @Success 204
// @Failure 404 {object} errors.AppError
// @Router /api/v1/patients/{id} [delete]
func (h *Handler) DeletePatient(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid patient ID"))
		return
	}

	if err := h.service.DeletePatient(c.Request.Context(), id); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPatientTimeline godoc
// @Summary Get patient timeline
// @Description Get complete medical timeline for a patient
// @Tags patients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Patient ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} errors.AppError
// @Router /api/v1/patients/{id}/timeline [get]
func (h *Handler) GetPatientTimeline(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid patient ID"))
		return
	}

	timeline, err := h.service.GetPatientTimeline(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, timeline)
}
