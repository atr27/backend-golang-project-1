package patient

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

// Service provides patient management services
type Service struct {
	db           *gorm.DB
	natsClient   *messaging.NATSClient
}

// NewService creates a new patient service
func NewService(db *gorm.DB, natsClient *messaging.NATSClient) *Service {
	return &Service{
		db:         db,
		natsClient: natsClient,
	}
}

// CreatePatientRequest represents create patient request
type CreatePatientRequest struct {
	FirstName        string                   `json:"first_name" binding:"required"`
	LastName         string                   `json:"last_name" binding:"required"`
	MiddleName       string                   `json:"middle_name"`
	DateOfBirth      time.Time                `json:"date_of_birth" binding:"required"`
	Gender           models.Gender            `json:"gender" binding:"required"`
	BloodType        string                   `json:"blood_type"`
	MaritalStatus    models.MaritalStatus     `json:"marital_status"`
	Nationality      string                   `json:"nationality"`
	Religion         string                   `json:"religion"`
	SSN              string                   `json:"ssn"`
	PassportNumber   string                   `json:"passport_number"`
	Email            string                   `json:"email"`
	PhoneNumber      string                   `json:"phone_number"`
	MobileNumber     string                   `json:"mobile_number"`
	Address          string                   `json:"address"`
	City             string                   `json:"city"`
	State            string                   `json:"state"`
	ZipCode          string                   `json:"zip_code"`
	Country          string                   `json:"country"`
	EmergencyContact models.EmergencyContact  `json:"emergency_contact"`
	Insurance        models.Insurance         `json:"insurance"`
	Language         string                   `json:"language"`
	Occupation       string                   `json:"occupation"`
}

// CreatePatient creates a new patient
func (s *Service) CreatePatient(ctx context.Context, req *CreatePatientRequest, createdBy uuid.UUID) (*models.Patient, error) {
	// Generate MRN (Medical Record Number)
	mrn := s.generateMRN()

	patient := &models.Patient{
		MRN:              mrn,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		MiddleName:       req.MiddleName,
		DateOfBirth:      req.DateOfBirth,
		Gender:           req.Gender,
		BloodType:        req.BloodType,
		MaritalStatus:    req.MaritalStatus,
		Nationality:      req.Nationality,
		Religion:         req.Religion,
		SSN:              req.SSN,
		PassportNumber:   req.PassportNumber,
		Email:            req.Email,
		PhoneNumber:      req.PhoneNumber,
		MobileNumber:     req.MobileNumber,
		Address:          req.Address,
		City:             req.City,
		State:            req.State,
		ZipCode:          req.ZipCode,
		Country:          req.Country,
		EmergencyContact: req.EmergencyContact,
		Insurance:        req.Insurance,
		Language:         req.Language,
		Occupation:       req.Occupation,
		Status:           models.PatientStatusActive,
	}
	patient.CreatedBy = createdBy
	patient.UpdatedBy = createdBy

	if err := s.db.WithContext(ctx).Create(patient).Error; err != nil {
		return nil, errors.ErrDatabaseError.WithDetails(err.Error())
	}

	// Publish event
	s.natsClient.Publish(messaging.SubjectPatientCreated, map[string]interface{}{
		"patient_id": patient.ID,
		"mrn":        patient.MRN,
		"created_by": createdBy,
	})

	return patient, nil
}

// GetPatient retrieves a patient by ID
func (s *Service) GetPatient(ctx context.Context, id uuid.UUID) (*models.Patient, error) {
	var patient models.Patient
	if err := s.db.WithContext(ctx).
		Preload("Allergies").
		Preload("Medications").
		Where("id = ?", id).
		First(&patient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrPatientNotFound(id.String())
		}
		return nil, errors.ErrDatabaseError
	}

	return &patient, nil
}

// GetPatientByMRN retrieves a patient by MRN
func (s *Service) GetPatientByMRN(ctx context.Context, mrn string) (*models.Patient, error) {
	var patient models.Patient
	if err := s.db.WithContext(ctx).
		Preload("Allergies").
		Preload("Medications").
		Where("mrn = ?", mrn).
		First(&patient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrPatientNotFound(mrn)
		}
		return nil, errors.ErrDatabaseError
	}

	return &patient, nil
}

// ListPatients lists patients with pagination
func (s *Service) ListPatients(ctx context.Context, page, pageSize int, search string) ([]models.Patient, int64, error) {
	var patients []models.Patient
	var total int64

	query := s.db.WithContext(ctx).Model(&models.Patient{})

	// Apply search filter
	if search != "" {
		query = query.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR mrn ILIKE ? OR email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := query.
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&patients).Error; err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	return patients, total, nil
}

// UpdatePatient updates a patient
func (s *Service) UpdatePatient(ctx context.Context, id uuid.UUID, req *CreatePatientRequest, updatedBy uuid.UUID) (*models.Patient, error) {
	var patient models.Patient
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&patient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrPatientNotFound(id.String())
		}
		return nil, errors.ErrDatabaseError
	}

	// Update fields
	patient.FirstName = req.FirstName
	patient.LastName = req.LastName
	patient.MiddleName = req.MiddleName
	patient.DateOfBirth = req.DateOfBirth
	patient.Gender = req.Gender
	patient.BloodType = req.BloodType
	patient.MaritalStatus = req.MaritalStatus
	patient.Nationality = req.Nationality
	patient.Religion = req.Religion
	patient.SSN = req.SSN
	patient.PassportNumber = req.PassportNumber
	patient.Email = req.Email
	patient.PhoneNumber = req.PhoneNumber
	patient.MobileNumber = req.MobileNumber
	patient.Address = req.Address
	patient.City = req.City
	patient.State = req.State
	patient.ZipCode = req.ZipCode
	patient.Country = req.Country
	patient.EmergencyContact = req.EmergencyContact
	patient.Insurance = req.Insurance
	patient.Language = req.Language
	patient.Occupation = req.Occupation
	patient.UpdatedBy = updatedBy

	if err := s.db.WithContext(ctx).Save(&patient).Error; err != nil {
		return nil, errors.ErrDatabaseError
	}

	// Publish event
	s.natsClient.Publish(messaging.SubjectPatientUpdated, map[string]interface{}{
		"patient_id": patient.ID,
		"mrn":        patient.MRN,
		"updated_by": updatedBy,
	})

	return &patient, nil
}

// DeletePatient soft deletes a patient
func (s *Service) DeletePatient(ctx context.Context, id uuid.UUID) error {
	result := s.db.WithContext(ctx).Delete(&models.Patient{}, id)
	if result.Error != nil {
		return errors.ErrDatabaseError
	}
	if result.RowsAffected == 0 {
		return errors.ErrPatientNotFound(id.String())
	}
	return nil
}

// generateMRN generates a unique Medical Record Number
func (s *Service) generateMRN() string {
	// Simple implementation - in production, use a more sophisticated approach
	return fmt.Sprintf("MRN%d", time.Now().UnixNano()%1000000000)
}

// GetPatientTimeline gets patient medical timeline
func (s *Service) GetPatientTimeline(ctx context.Context, patientID uuid.UUID) (map[string]interface{}, error) {
	var patient models.Patient
	if err := s.db.WithContext(ctx).
		Preload("Encounters.Provider").
		Preload("Encounters.Diagnoses").
		Preload("Appointments.Provider").
		Preload("Allergies").
		Preload("Medications").
		Where("id = ?", patientID).
		First(&patient).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrPatientNotFound(patientID.String())
		}
		return nil, errors.ErrDatabaseError
	}

	timeline := map[string]interface{}{
		"patient":      patient,
		"encounters":   patient.Encounters,
		"appointments": patient.Appointments,
		"allergies":    patient.Allergies,
		"medications":  patient.Medications,
	}

	return timeline, nil
}
