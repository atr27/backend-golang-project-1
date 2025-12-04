package models

import (
	"time"

	"github.com/google/uuid"
)

// Appointment represents a scheduled appointment
type Appointment struct {
	AuditableModel
	AppointmentNumber string            `gorm:"uniqueIndex;not null" json:"appointment_number"`
	PatientID         uuid.UUID         `gorm:"type:uuid;not null;index" json:"patient_id"`
	Patient           Patient           `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	ProviderID        uuid.UUID         `gorm:"type:uuid;not null;index" json:"provider_id"`
	Provider          User              `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	AppointmentType   AppointmentType   `gorm:"type:varchar(50);not null" json:"appointment_type"`
	Status            AppointmentStatus `gorm:"type:varchar(20);not null;default:'scheduled'" json:"status"`
	StartTime         time.Time         `gorm:"not null;index" json:"start_time"`
	EndTime           time.Time         `gorm:"not null" json:"end_time"`
	Duration          int               `json:"duration"` // in minutes
	Department        string            `json:"department"`
	Location          string            `json:"location"`
	Room              string            `json:"room"`
	ReasonForVisit    string            `json:"reason_for_visit"`
	Notes             string            `json:"notes"`
	ReminderSent      bool              `gorm:"default:false" json:"reminder_sent"`
	ReminderSentAt    *time.Time        `json:"reminder_sent_at"`
	CheckedInAt       *time.Time        `json:"checked_in_at"`
	CancelledAt       *time.Time        `json:"cancelled_at"`
	CancellationReason string           `json:"cancellation_reason"`
}

// AppointmentType represents type of appointment
type AppointmentType string

const (
	AppointmentTypeConsultation AppointmentType = "consultation"
	AppointmentTypeFollowUp     AppointmentType = "follow_up"
	AppointmentTypeWellness     AppointmentType = "wellness"
	AppointmentTypeProcedure    AppointmentType = "procedure"
	AppointmentTypeEmergency    AppointmentType = "emergency"
	AppointmentTypeTelehealth   AppointmentType = "telehealth"
)

// AppointmentStatus represents appointment status
type AppointmentStatus string

const (
	AppointmentStatusScheduled  AppointmentStatus = "scheduled"
	AppointmentStatusConfirmed  AppointmentStatus = "confirmed"
	AppointmentStatusCheckedIn  AppointmentStatus = "checked_in"
	AppointmentStatusInProgress AppointmentStatus = "in_progress"
	AppointmentStatusCompleted  AppointmentStatus = "completed"
	AppointmentStatusCancelled  AppointmentStatus = "cancelled"
	AppointmentStatusNoShow     AppointmentStatus = "no_show"
)

// TableName specifies table name
func (Appointment) TableName() string { return "appointments" }
