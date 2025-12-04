package encounter

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hospital-emr/backend/internal/common/errors"
	"github.com/hospital-emr/backend/internal/models"
	"github.com/hospital-emr/backend/pkg/messaging"
	"gorm.io/gorm"
)

// Service provides encounter management services
type Service struct {
	db         *gorm.DB
	natsClient *messaging.NATSClient
}

// NewService creates a new encounter service
func NewService(db *gorm.DB, natsClient *messaging.NATSClient) *Service {
	return &Service{
		db:         db,
		natsClient: natsClient,
	}
}

// CreateEncounterRequest represents create encounter request
type CreateEncounterRequest struct {
	PatientID      uuid.UUID              `json:"patient_id" binding:"required"`
	ProviderID     uuid.UUID              `json:"provider_id" binding:"required"`
	EncounterType  models.EncounterType   `json:"encounter_type" binding:"required"`
	Priority       models.Priority        `json:"priority"`
	Department     string                 `json:"department"`
	Location       string                 `json:"location"`
	AdmissionDate  time.Time              `json:"admission_date" binding:"required"`
	ChiefComplaint string                 `json:"chief_complaint"`
	ReasonForVisit string                 `json:"reason_for_visit"`
}

// CreateEncounter creates a new encounter
func (s *Service) CreateEncounter(ctx context.Context, req *CreateEncounterRequest, createdBy uuid.UUID) (*models.Encounter, error) {
	// Verify patient exists
	var patient models.Patient
	if err := s.db.WithContext(ctx).Where("id = ?", req.PatientID).First(&patient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrPatientNotFound(req.PatientID.String())
		}
		return nil, errors.ErrDatabaseError
	}

	// Verify provider exists
	var provider models.User
	if err := s.db.WithContext(ctx).Where("id = ?", req.ProviderID).First(&provider).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound(req.ProviderID.String())
		}
		return nil, errors.ErrDatabaseError
	}

	// Generate encounter number
	encounterNumber := s.generateEncounterNumber()

	encounter := &models.Encounter{
		EncounterNumber: encounterNumber,
		PatientID:       req.PatientID,
		ProviderID:      req.ProviderID,
		EncounterType:   req.EncounterType,
		Status:          models.EncounterStatusScheduled,
		Priority:        req.Priority,
		Department:      req.Department,
		Location:        req.Location,
		AdmissionDate:   req.AdmissionDate,
		ChiefComplaint:  req.ChiefComplaint,
		ReasonForVisit:  req.ReasonForVisit,
	}
	encounter.CreatedBy = createdBy
	encounter.UpdatedBy = createdBy

	if err := s.db.WithContext(ctx).Create(encounter).Error; err != nil {
		return nil, errors.ErrDatabaseError.WithDetails(err.Error())
	}

	// Publish event
	s.natsClient.Publish(messaging.SubjectEncounterCreated, map[string]interface{}{
		"encounter_id":     encounter.ID,
		"encounter_number": encounter.EncounterNumber,
		"patient_id":       encounter.PatientID,
		"provider_id":      encounter.ProviderID,
		"created_by":       createdBy,
	})

	return encounter, nil
}

// GetEncounter retrieves an encounter by ID
func (s *Service) GetEncounter(ctx context.Context, id uuid.UUID) (*models.Encounter, error) {
	var encounter models.Encounter
	if err := s.db.WithContext(ctx).
		Preload("Patient").
		Preload("Provider").
		Preload("ClinicalNotes.Author").
		Preload("Diagnoses").
		Preload("Procedures").
		Preload("VitalSigns").
		Preload("Orders").
		Where("id = ?", id).
		First(&encounter).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrEncounterNotFound(id.String())
		}
		return nil, errors.ErrDatabaseError
	}

	return &encounter, nil
}

// ListEncounters lists encounters with pagination
func (s *Service) ListEncounters(ctx context.Context, page, pageSize int, patientID *uuid.UUID, providerID *uuid.UUID, status *models.EncounterStatus) ([]models.Encounter, int64, error) {
	var encounters []models.Encounter
	var total int64

	query := s.db.WithContext(ctx).Model(&models.Encounter{})

	// Apply filters
	if patientID != nil {
		query = query.Where("patient_id = ?", *patientID)
	}
	if providerID != nil {
		query = query.Where("provider_id = ?", *providerID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := query.
		Preload("Patient").
		Preload("Provider").
		Offset(offset).
		Limit(pageSize).
		Order("admission_date DESC").
		Find(&encounters).Error; err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	return encounters, total, nil
}

// UpdateEncounter updates an encounter
func (s *Service) UpdateEncounter(ctx context.Context, id uuid.UUID, status models.EncounterStatus, updatedBy uuid.UUID) (*models.Encounter, error) {
	var encounter models.Encounter
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&encounter).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrEncounterNotFound(id.String())
		}
		return nil, errors.ErrDatabaseError
	}

	encounter.Status = status
	encounter.UpdatedBy = updatedBy

	// If completing encounter, set discharge date
	if status == models.EncounterStatusCompleted {
		now := time.Now()
		encounter.DischargeDate = &now
	}

	if err := s.db.WithContext(ctx).Save(&encounter).Error; err != nil {
		return nil, errors.ErrDatabaseError
	}

	// Publish event
	s.natsClient.Publish(messaging.SubjectEncounterUpdated, map[string]interface{}{
		"encounter_id": encounter.ID,
		"status":       encounter.Status,
		"updated_by":   updatedBy,
	})

	return &encounter, nil
}

// AddClinicalNote adds a clinical note to an encounter
func (s *Service) AddClinicalNote(ctx context.Context, encounterID uuid.UUID, req *AddClinicalNoteRequest, authorID uuid.UUID) (*models.ClinicalNote, error) {
	// Verify encounter exists
	var encounter models.Encounter
	if err := s.db.WithContext(ctx).Where("id = ?", encounterID).First(&encounter).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrEncounterNotFound(encounterID.String())
		}
		return nil, errors.ErrDatabaseError
	}

	// Verify author (user) exists
	var author models.User
	if err := s.db.WithContext(ctx).Where("id = ?", authorID).First(&author).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound(authorID.String())
		}
		return nil, errors.ErrDatabaseError.WithDetails("Failed to verify author")
	}

	note := &models.ClinicalNote{
		EncounterID: encounterID,
		NoteType:    req.NoteType,
		Subjective:  req.Subjective,
		Objective:   req.Objective,
		Assessment:  req.Assessment,
		Plan:        req.Plan,
		Content:     req.Content,
		AuthorID:    authorID,
	}
	note.CreatedBy = authorID
	note.UpdatedBy = authorID

	if err := s.db.WithContext(ctx).Create(note).Error; err != nil {
		return nil, errors.ErrDatabaseError.WithDetails(err.Error())
	}

	return note, nil
}


// AddDiagnosis adds a diagnosis to an encounter
func (s *Service) AddDiagnosis(ctx context.Context, encounterID uuid.UUID, req *AddDiagnosisRequest, diagnosedBy uuid.UUID) (*models.Diagnosis, error) {
	// Verify encounter exists
	var encounter models.Encounter
	if err := s.db.WithContext(ctx).Where("id = ?", encounterID).First(&encounter).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrEncounterNotFound(encounterID.String())
		}
		return nil, errors.ErrDatabaseError
	}

	diagnosis := &models.Diagnosis{
		EncounterID:   encounterID,
		ICD10Code:     req.ICD10Code,
		Description:   req.Description,
		DiagnosisType: req.DiagnosisType,
		Status:        "active",
		OnsetDate:     req.OnsetDate,
		Severity:      req.Severity,
		Notes:         req.Notes,
		DiagnosedBy:   diagnosedBy,
	}
	diagnosis.CreatedBy = diagnosedBy
	diagnosis.UpdatedBy = diagnosedBy

	if err := s.db.WithContext(ctx).Create(diagnosis).Error; err != nil {
		return nil, errors.ErrDatabaseError.WithDetails(err.Error())
	}

	return diagnosis, nil
}

// RecordVitalSigns records vital signs for an encounter
func (s *Service) RecordVitalSigns(ctx context.Context, encounterID uuid.UUID, req *RecordVitalSignsRequest, recordedBy uuid.UUID) (*models.VitalSign, error) {
	// Verify encounter exists
	var encounter models.Encounter
	if err := s.db.WithContext(ctx).Where("id = ?", encounterID).First(&encounter).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrEncounterNotFound(encounterID.String())
		}
		return nil, errors.ErrDatabaseError
	}

	// Calculate BMI if height and weight are provided
	var bmi *float64
	if req.Height != nil && req.Weight != nil && *req.Height > 0 {
		heightM := *req.Height / 100 // Convert cm to meters
		calculatedBMI := *req.Weight / (heightM * heightM)
		bmi = &calculatedBMI
	}

	vitalSign := &models.VitalSign{
		EncounterID:            encounterID,
		PatientID:              encounter.PatientID,
		MeasuredAt:             time.Now(),
		Temperature:            req.Temperature,
		TemperatureUnit:        req.TemperatureUnit,
		HeartRate:              req.HeartRate,
		RespiratoryRate:        req.RespiratoryRate,
		BloodPressureSystolic:  req.BloodPressureSystolic,
		BloodPressureDiastolic: req.BloodPressureDiastolic,
		OxygenSaturation:       req.OxygenSaturation,
		Weight:                 req.Weight,
		Height:                 req.Height,
		BMI:                    bmi,
		Pain:                   req.Pain,
		RecordedBy:             recordedBy,
		Notes:                  req.Notes,
	}
	vitalSign.CreatedBy = recordedBy
	vitalSign.UpdatedBy = recordedBy

	if err := s.db.WithContext(ctx).Create(vitalSign).Error; err != nil {
		return nil, errors.ErrDatabaseError.WithDetails(err.Error())
	}

	return vitalSign, nil
}

// generateEncounterNumber generates a unique encounter number
func (s *Service) generateEncounterNumber() string {
	return fmt.Sprintf("ENC%d", time.Now().UnixNano()%1000000000)
}

// Request structures
type AddClinicalNoteRequest struct {
	NoteType   models.NoteType `json:"note_type" binding:"required"`
	Subjective string          `json:"subjective"`
	Objective  string          `json:"objective"`
	Assessment string          `json:"assessment"`
	Plan       string          `json:"plan"`
	Content    string          `json:"content"`
}

type AddDiagnosisRequest struct {
	ICD10Code     string                `json:"icd10_code" binding:"required"`
	Description   string                `json:"description" binding:"required"`
	DiagnosisType models.DiagnosisType  `json:"diagnosis_type" binding:"required"`
	OnsetDate     *time.Time            `json:"onset_date"`
	Severity      string                `json:"severity"`
	Notes         string                `json:"notes"`
}

type RecordVitalSignsRequest struct {
	Temperature            *float64 `json:"temperature"`
	TemperatureUnit        string   `json:"temperature_unit"`
	HeartRate              *int     `json:"heart_rate"`
	RespiratoryRate        *int     `json:"respiratory_rate"`
	BloodPressureSystolic  *int     `json:"blood_pressure_systolic"`
	BloodPressureDiastolic *int     `json:"blood_pressure_diastolic"`
	OxygenSaturation       *float64 `json:"oxygen_saturation"`
	Weight                 *float64 `json:"weight"`
	Height                 *float64 `json:"height"`
	Pain                   *int     `json:"pain"`
	Notes                  string   `json:"notes"`
}
