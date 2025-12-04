package scheduling

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

// Service provides appointment scheduling services
type Service struct {
	db         *gorm.DB
	natsClient *messaging.NATSClient
}

// NewService creates a new scheduling service
func NewService(db *gorm.DB, natsClient *messaging.NATSClient) *Service {
	return &Service{
		db:         db,
		natsClient: natsClient,
	}
}

// CreateAppointmentRequest represents create appointment request
type CreateAppointmentRequest struct {
	PatientID       uuid.UUID                `json:"patient_id" binding:"required"`
	ProviderID      uuid.UUID                `json:"provider_id" binding:"required"`
	AppointmentType models.AppointmentType   `json:"appointment_type" binding:"required"`
	StartTime       time.Time                `json:"start_time" binding:"required"`
	Duration        int                      `json:"duration" binding:"required"` // in minutes
	Department      string                   `json:"department"`
	Location        string                   `json:"location"`
	Room            string                   `json:"room"`
	ReasonForVisit  string                   `json:"reason_for_visit"`
	Notes           string                   `json:"notes"`
}

// CreateAppointment creates a new appointment
func (s *Service) CreateAppointment(ctx context.Context, req *CreateAppointmentRequest, createdBy uuid.UUID) (*models.Appointment, error) {
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

	// Check provider availability
	endTime := req.StartTime.Add(time.Duration(req.Duration) * time.Minute)
	if !s.isProviderAvailable(ctx, req.ProviderID, req.StartTime, endTime) {
		return nil, errors.ErrAppointmentConflict()
	}

	// Generate appointment number
	appointmentNumber := s.generateAppointmentNumber()

	appointment := &models.Appointment{
		AppointmentNumber: appointmentNumber,
		PatientID:         req.PatientID,
		ProviderID:        req.ProviderID,
		AppointmentType:   req.AppointmentType,
		Status:            models.AppointmentStatusScheduled,
		StartTime:         req.StartTime,
		EndTime:           endTime,
		Duration:          req.Duration,
		Department:        req.Department,
		Location:          req.Location,
		Room:              req.Room,
		ReasonForVisit:    req.ReasonForVisit,
		Notes:             req.Notes,
	}
	appointment.CreatedBy = createdBy
	appointment.UpdatedBy = createdBy

	if err := s.db.WithContext(ctx).Create(appointment).Error; err != nil {
		return nil, errors.ErrDatabaseError.WithDetails(err.Error())
	}

	// Publish event
	s.natsClient.Publish(messaging.SubjectAppointmentBooked, map[string]interface{}{
		"appointment_id":     appointment.ID,
		"appointment_number": appointment.AppointmentNumber,
		"patient_id":         appointment.PatientID,
		"provider_id":        appointment.ProviderID,
		"start_time":         appointment.StartTime,
		"created_by":         createdBy,
	})

	return appointment, nil
}

// GetAppointment retrieves an appointment by ID
func (s *Service) GetAppointment(ctx context.Context, id uuid.UUID) (*models.Appointment, error) {
	var appointment models.Appointment
	if err := s.db.WithContext(ctx).
		Preload("Patient").
		Preload("Provider").
		Where("id = ?", id).
		First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrAppointmentNotFound(id.String())
		}
		return nil, errors.ErrDatabaseError
	}

	return &appointment, nil
}

// ListAppointments lists appointments with pagination and filters
func (s *Service) ListAppointments(ctx context.Context, page, pageSize int, patientID *uuid.UUID, providerID *uuid.UUID, status *models.AppointmentStatus, date *time.Time) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	query := s.db.WithContext(ctx).Model(&models.Appointment{})

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
	if date != nil {
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		query = query.Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay)
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
		Order("start_time ASC").
		Find(&appointments).Error; err != nil {
		return nil, 0, errors.ErrDatabaseError
	}

	return appointments, total, nil
}

// UpdateAppointment updates an appointment
func (s *Service) UpdateAppointment(ctx context.Context, id uuid.UUID, req *CreateAppointmentRequest, updatedBy uuid.UUID) (*models.Appointment, error) {
	var appointment models.Appointment
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrAppointmentNotFound(id.String())
		}
		return nil, errors.ErrDatabaseError
	}

	// Check provider availability (excluding current appointment)
	endTime := req.StartTime.Add(time.Duration(req.Duration) * time.Minute)
	if !s.isProviderAvailableExcept(ctx, req.ProviderID, req.StartTime, endTime, id) {
		return nil, errors.ErrAppointmentConflict()
	}

	// Update fields
	appointment.PatientID = req.PatientID
	appointment.ProviderID = req.ProviderID
	appointment.AppointmentType = req.AppointmentType
	appointment.StartTime = req.StartTime
	appointment.EndTime = endTime
	appointment.Duration = req.Duration
	appointment.Department = req.Department
	appointment.Location = req.Location
	appointment.Room = req.Room
	appointment.ReasonForVisit = req.ReasonForVisit
	appointment.Notes = req.Notes
	appointment.UpdatedBy = updatedBy

	if err := s.db.WithContext(ctx).Save(&appointment).Error; err != nil {
		return nil, errors.ErrDatabaseError
	}

	return &appointment, nil
}

// CancelAppointment cancels an appointment
func (s *Service) CancelAppointment(ctx context.Context, id uuid.UUID, reason string, cancelledBy uuid.UUID) error {
	var appointment models.Appointment
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrAppointmentNotFound(id.String())
		}
		return errors.ErrDatabaseError
	}

	now := time.Now()
	appointment.Status = models.AppointmentStatusCancelled
	appointment.CancelledAt = &now
	appointment.CancellationReason = reason
	appointment.UpdatedBy = cancelledBy

	if err := s.db.WithContext(ctx).Save(&appointment).Error; err != nil {
		return errors.ErrDatabaseError
	}

	// Publish event
	s.natsClient.Publish(messaging.SubjectAppointmentCancelled, map[string]interface{}{
		"appointment_id": appointment.ID,
		"patient_id":     appointment.PatientID,
		"provider_id":    appointment.ProviderID,
		"reason":         reason,
		"cancelled_by":   cancelledBy,
	})

	return nil
}

// CheckInAppointment marks appointment as checked in
func (s *Service) CheckInAppointment(ctx context.Context, id uuid.UUID) error {
	var appointment models.Appointment
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrAppointmentNotFound(id.String())
		}
		return errors.ErrDatabaseError
	}

	now := time.Now()
	appointment.Status = models.AppointmentStatusCheckedIn
	appointment.CheckedInAt = &now

	return s.db.WithContext(ctx).Save(&appointment).Error
}

// GetProviderAvailability gets available time slots for a provider
func (s *Service) GetProviderAvailability(ctx context.Context, providerID uuid.UUID, date time.Time) ([]TimeSlot, error) {
	// Verify provider exists
	var provider models.User
	if err := s.db.WithContext(ctx).Where("id = ?", providerID).First(&provider).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound(providerID.String())
		}
		return nil, errors.ErrDatabaseError
	}

	// Get appointments for the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var appointments []models.Appointment
	if err := s.db.WithContext(ctx).
		Where("provider_id = ? AND start_time >= ? AND start_time < ? AND status NOT IN ?",
			providerID, startOfDay, endOfDay, []string{string(models.AppointmentStatusCancelled)}).
		Order("start_time ASC").
		Find(&appointments).Error; err != nil {
		return nil, errors.ErrDatabaseError
	}

	// Generate available slots (assuming 8 AM - 5 PM, 30-minute slots)
	workStart := time.Date(date.Year(), date.Month(), date.Day(), 8, 0, 0, 0, date.Location())
	workEnd := time.Date(date.Year(), date.Month(), date.Day(), 17, 0, 0, 0, date.Location())
	slotDuration := 30 * time.Minute

	var availableSlots []TimeSlot
	currentSlot := workStart

	for currentSlot.Before(workEnd) {
		slotEnd := currentSlot.Add(slotDuration)
		isAvailable := true

		// Check if slot conflicts with any appointment
		for _, appt := range appointments {
			if (currentSlot.Before(appt.EndTime) && slotEnd.After(appt.StartTime)) {
				isAvailable = false
				break
			}
		}

		availableSlots = append(availableSlots, TimeSlot{
			StartTime: currentSlot,
			EndTime:   slotEnd,
			Available: isAvailable,
		})

		currentSlot = slotEnd
	}

	return availableSlots, nil
}

// Helper functions

func (s *Service) isProviderAvailable(ctx context.Context, providerID uuid.UUID, startTime, endTime time.Time) bool {
	var count int64
	s.db.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("provider_id = ? AND status NOT IN ? AND start_time < ? AND end_time > ?",
			providerID,
			[]string{string(models.AppointmentStatusCancelled)},
			endTime,
			startTime).
		Count(&count)

	return count == 0
}

func (s *Service) isProviderAvailableExcept(ctx context.Context, providerID uuid.UUID, startTime, endTime time.Time, exceptID uuid.UUID) bool {
	var count int64
	s.db.WithContext(ctx).
		Model(&models.Appointment{}).
		Where("provider_id = ? AND id != ? AND status NOT IN ? AND start_time < ? AND end_time > ?",
			providerID,
			exceptID,
			[]string{string(models.AppointmentStatusCancelled)},
			endTime,
			startTime).
		Count(&count)

	return count == 0
}

func (s *Service) generateAppointmentNumber() string {
	return fmt.Sprintf("APT%d", time.Now().UnixNano()%1000000000)
}

// TimeSlot represents an available time slot
type TimeSlot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Available bool      `json:"available"`
}
