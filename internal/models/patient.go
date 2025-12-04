package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Patient represents a patient in the system
type Patient struct {
	AuditableModel
	MRN             string          `gorm:"uniqueIndex;not null" json:"mrn"` // Medical Record Number
	FirstName       string          `gorm:"not null" json:"first_name"`
	LastName        string          `gorm:"not null" json:"last_name"`
	MiddleName      string          `json:"middle_name"`
	DateOfBirth     time.Time       `gorm:"not null" json:"date_of_birth"`
	Gender          Gender          `gorm:"type:varchar(20);not null" json:"gender"`
	BloodType       string          `json:"blood_type"`
	MaritalStatus   MaritalStatus   `gorm:"type:varchar(20)" json:"marital_status"`
	Nationality     string          `json:"nationality"`
	Religion        string          `json:"religion"`
	SSN             string          `gorm:"index" json:"ssn"` // Social Security Number / National ID
	PassportNumber  string          `json:"passport_number"`
	Email           string          `json:"email"`
	PhoneNumber     string          `json:"phone_number"`
	MobileNumber    string          `json:"mobile_number"`
	Address         string          `json:"address"`
	City            string          `json:"city"`
	State           string          `json:"state"`
	ZipCode         string          `json:"zip_code"`
	Country         string          `json:"country"`
	EmergencyContact EmergencyContact `gorm:"type:jsonb" json:"emergency_contact"`
	Insurance       Insurance       `gorm:"type:jsonb" json:"insurance"`
	Status          PatientStatus   `gorm:"type:varchar(20);default:'active'" json:"status"`
	ProfilePhoto    string          `json:"profile_photo"`
	Language        string          `json:"language"`
	Occupation      string          `json:"occupation"`
	Encounters      []Encounter     `gorm:"foreignKey:PatientID" json:"encounters,omitempty"`
	Appointments    []Appointment   `gorm:"foreignKey:PatientID" json:"appointments,omitempty"`
	Allergies       []Allergy       `gorm:"foreignKey:PatientID" json:"allergies,omitempty"`
	Medications     []Medication    `gorm:"foreignKey:PatientID" json:"medications,omitempty"`
}

// Gender represents patient gender
type Gender string

const (
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderOther   Gender = "other"
	GenderUnknown Gender = "unknown"
)

// MaritalStatus represents marital status
type MaritalStatus string

const (
	MaritalStatusSingle   MaritalStatus = "single"
	MaritalStatusMarried  MaritalStatus = "married"
	MaritalStatusDivorced MaritalStatus = "divorced"
	MaritalStatusWidowed  MaritalStatus = "widowed"
)

// PatientStatus represents patient status
type PatientStatus string

const (
	PatientStatusActive   PatientStatus = "active"
	PatientStatusInactive PatientStatus = "inactive"
	PatientStatusDeceased PatientStatus = "deceased"
)

// EmergencyContact represents emergency contact information
type EmergencyContact struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Address      string `json:"address"`
}

// Scan implements sql.Scanner interface for JSONB
func (ec *EmergencyContact) Scan(value interface{}) error {
	if value == nil {
		*ec = EmergencyContact{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	
	return json.Unmarshal(bytes, ec)
}

// Value implements driver.Valuer interface for JSONB
func (ec EmergencyContact) Value() (driver.Value, error) {
	if ec == (EmergencyContact{}) {
		return nil, nil
	}
	return json.Marshal(ec)
}

// Insurance represents insurance information
type Insurance struct {
	Provider       string `json:"provider"`
	PolicyNumber   string `json:"policy_number"`
	GroupNumber    string `json:"group_number"`
	ExpiryDate     string `json:"expiry_date"`
	CoverageType   string `json:"coverage_type"`
	PrimaryInsured string `json:"primary_insured"`
}

// Scan implements sql.Scanner interface for JSONB
func (ins *Insurance) Scan(value interface{}) error {
	if value == nil {
		*ins = Insurance{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	
	return json.Unmarshal(bytes, ins)
}

// Value implements driver.Valuer interface for JSONB
func (ins Insurance) Value() (driver.Value, error) {
	if ins == (Insurance{}) {
		return nil, nil
	}
	return json.Marshal(ins)
}

// Allergy represents patient allergies
type Allergy struct {
	AuditableModel
	PatientID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"patient_id"`
	Patient     Patient        `gorm:"foreignKey:PatientID" json:"-"`
	AllergyType AllergyType    `gorm:"type:varchar(50);not null" json:"allergy_type"`
	Allergen    string         `gorm:"not null" json:"allergen"`
	Reaction    string         `json:"reaction"`
	Severity    AllergySeverity `gorm:"type:varchar(20)" json:"severity"`
	OnsetDate   *time.Time     `json:"onset_date"`
	Notes       string         `json:"notes"`
	Status      string         `gorm:"default:'active'" json:"status"`
}

// AllergyType represents type of allergy
type AllergyType string

const (
	AllergyTypeDrug        AllergyType = "drug"
	AllergyTypeFood        AllergyType = "food"
	AllergyTypeEnvironment AllergyType = "environment"
	AllergyTypeOther       AllergyType = "other"
)

// AllergySeverity represents allergy severity
type AllergySeverity string

const (
	AllergySeverityMild     AllergySeverity = "mild"
	AllergySeverityModerate AllergySeverity = "moderate"
	AllergySeveritySevere   AllergySeverity = "severe"
	AllergySeverityFatal    AllergySeverity = "fatal"
)

// Medication represents patient medication history
type Medication struct {
	AuditableModel
	PatientID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"patient_id"`
	Patient         Patient    `gorm:"foreignKey:PatientID" json:"-"`
	MedicationName  string     `gorm:"not null" json:"medication_name"`
	GenericName     string     `json:"generic_name"`
	Dosage          string     `json:"dosage"`
	Frequency       string     `json:"frequency"`
	Route           string     `json:"route"`
	StartDate       time.Time  `gorm:"not null" json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	PrescribedBy    uuid.UUID  `gorm:"type:uuid" json:"prescribed_by"`
	Reason          string     `json:"reason"`
	Instructions    string     `json:"instructions"`
	Status          string     `gorm:"default:'active'" json:"status"`
	RefillsRemaining int       `json:"refills_remaining"`
}

// TableName specifies table names
func (Patient) TableName() string    { return "patients" }
func (Allergy) TableName() string    { return "allergies" }
func (Medication) TableName() string { return "medications" }
