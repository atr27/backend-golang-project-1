package encounter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hospital-emr/backend/internal/common/errors"
	"github.com/hospital-emr/backend/internal/models"
)

// Handler handles encounter HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new encounter handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateEncounter godoc
// @Summary Create a new encounter
// @Description Create a new clinical encounter for a patient
// @Tags encounters
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateEncounterRequest true "Encounter information"
// @Success 201 {object} models.Encounter
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /api/v1/encounters [post]
func (h *Handler) CreateEncounter(c *gin.Context) {
	var req CreateEncounterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	userIDValue, _ := c.Get("user_id")
	createdBy, _ := userIDValue.(uuid.UUID)

	encounter, err := h.service.CreateEncounter(c.Request.Context(), &req, createdBy)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusCreated, encounter)
}

// GetEncounter godoc
// @Summary Get encounter by ID
// @Description Get detailed information about an encounter
// @Tags encounters
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Encounter ID"
// @Success 200 {object} models.Encounter
// @Failure 404 {object} errors.AppError
// @Router /api/v1/encounters/{id} [get]
func (h *Handler) GetEncounter(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid encounter ID"))
		return
	}

	encounter, err := h.service.GetEncounter(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, encounter)
}

// ListEncounters godoc
// @Summary List encounters
// @Description Get a paginated list of encounters
// @Tags encounters
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param patient_id query string false "Filter by patient ID"
// @Param provider_id query string false "Filter by provider ID"
// @Param status query string false "Filter by status"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/encounters [get]
func (h *Handler) ListEncounters(c *gin.Context) {
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

	var status *models.EncounterStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s := models.EncounterStatus(statusStr)
		status = &s
	}

	encounters, total, err := h.service.ListEncounters(c.Request.Context(), page, pageSize, patientID, providerID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        encounters,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// UpdateEncounterStatus godoc
// @Summary Update encounter status
// @Description Update the status of an encounter
// @Tags encounters
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Encounter ID"
// @Param request body map[string]string true "Status update"
// @Success 200 {object} models.Encounter
// @Failure 404 {object} errors.AppError
// @Router /api/v1/encounters/{id}/status [put]
func (h *Handler) UpdateEncounterStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid encounter ID"))
		return
	}

	var req struct {
		Status models.EncounterStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	userIDValue, _ := c.Get("user_id")
	updatedBy, _ := userIDValue.(uuid.UUID)

	encounter, err := h.service.UpdateEncounter(c.Request.Context(), id, req.Status, updatedBy)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, encounter)
}

// AddClinicalNote godoc
// @Summary Add clinical note
// @Description Add a clinical note to an encounter
// @Tags encounters
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Encounter ID"
// @Param request body AddClinicalNoteRequest true "Clinical note"
// @Success 201 {object} models.ClinicalNote
// @Failure 400 {object} errors.AppError
// @Router /api/v1/encounters/{id}/notes [post]
func (h *Handler) AddClinicalNote(c *gin.Context) {
	encounterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid encounter ID"))
		return
	}

	var req AddClinicalNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	userIDValue, _ := c.Get("user_id")
	authorID, _ := userIDValue.(uuid.UUID)

	note, err := h.service.AddClinicalNote(c.Request.Context(), encounterID, &req, authorID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusCreated, note)
}

// AddDiagnosis godoc
// @Summary Add diagnosis
// @Description Add a diagnosis to an encounter
// @Tags encounters
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Encounter ID"
// @Param request body AddDiagnosisRequest true "Diagnosis"
// @Success 201 {object} models.Diagnosis
// @Failure 400 {object} errors.AppError
// @Router /api/v1/encounters/{id}/diagnoses [post]
func (h *Handler) AddDiagnosis(c *gin.Context) {
	encounterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid encounter ID"))
		return
	}

	var req AddDiagnosisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	userIDValue, _ := c.Get("user_id")
	diagnosedBy, _ := userIDValue.(uuid.UUID)

	diagnosis, err := h.service.AddDiagnosis(c.Request.Context(), encounterID, &req, diagnosedBy)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusCreated, diagnosis)
}

// RecordVitalSigns godoc
// @Summary Record vital signs
// @Description Record vital signs for an encounter
// @Tags encounters
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Encounter ID"
// @Param request body RecordVitalSignsRequest true "Vital signs"
// @Success 201 {object} models.VitalSign
// @Failure 400 {object} errors.AppError
// @Router /api/v1/encounters/{id}/vitals [post]
func (h *Handler) RecordVitalSigns(c *gin.Context) {
	encounterID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid encounter ID"))
		return
	}

	var req RecordVitalSignsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	userIDValue, _ := c.Get("user_id")
	recordedBy, _ := userIDValue.(uuid.UUID)

	vitalSign, err := h.service.RecordVitalSigns(c.Request.Context(), encounterID, &req, recordedBy)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusCreated, vitalSign)
}

// CompleteEncounter godoc
// @Summary Complete encounter
// @Description Mark an encounter as completed
// @Tags encounters
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Encounter ID"
// @Success 200 {object} models.Encounter
// @Failure 404 {object} errors.AppError
// @Router /api/v1/encounters/{id}/complete [post]
func (h *Handler) CompleteEncounter(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails("Invalid encounter ID"))
		return
	}

	userIDValue, _ := c.Get("user_id")
	updatedBy, _ := userIDValue.(uuid.UUID)

	encounter, err := h.service.UpdateEncounter(c.Request.Context(), id, models.EncounterStatusCompleted, updatedBy)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, encounter)
}
