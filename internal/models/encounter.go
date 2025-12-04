package models

import (
	"time"

	"github.com/google/uuid"
)

// Encounter represents a clinical encounter (visit)
type Encounter struct {
	AuditableModel
	EncounterNumber string          `gorm:"uniqueIndex;not null" json:"encounter_number"`
	PatientID       uuid.UUID       `gorm:"type:uuid;not null;index" json:"patient_id"`
	Patient         Patient         `gorm:"foreignKey:PatientID;references:ID" json:"patient,omitempty"`
	ProviderID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"provider_id"`
	Provider        User            `gorm:"foreignKey:ProviderID;references:ID" json:"provider,omitempty"`
	EncounterType   EncounterType   `gorm:"type:varchar(50);not null" json:"encounter_type"`
	Status          EncounterStatus `gorm:"type:varchar(20);not null;default:'scheduled'" json:"status"`
	Priority        Priority        `gorm:"type:varchar(20)" json:"priority"`
	Department      string          `json:"department"`
	Location        string          `json:"location"`
	AdmissionDate   time.Time       `gorm:"not null" json:"admission_date"`
	DischargeDate   *time.Time      `json:"discharge_date"`
	ChiefComplaint  string          `json:"chief_complaint"`
	ReasonForVisit  string          `json:"reason_for_visit"`
	ClinicalNotes   []ClinicalNote  `gorm:"foreignKey:EncounterID" json:"clinical_notes,omitempty"`
	Diagnoses       []Diagnosis     `gorm:"foreignKey:EncounterID" json:"diagnoses,omitempty"`
	Procedures      []Procedure     `gorm:"foreignKey:EncounterID" json:"procedures,omitempty"`
	Orders          []Order         `gorm:"foreignKey:EncounterID" json:"orders,omitempty"`
	VitalSigns      []VitalSign     `gorm:"foreignKey:EncounterID" json:"vital_signs,omitempty"`
}

// EncounterType represents type of encounter
type EncounterType string

const (
	EncounterTypeOutpatient EncounterType = "outpatient"
	EncounterTypeInpatient  EncounterType = "inpatient"
	EncounterTypeEmergency  EncounterType = "emergency"
	EncounterTypeWellness   EncounterType = "wellness"
	EncounterTypeTelehealth EncounterType = "telehealth"
)

// EncounterStatus represents encounter status
type EncounterStatus string

const (
	EncounterStatusScheduled  EncounterStatus = "scheduled"
	EncounterStatusInProgress EncounterStatus = "in_progress"
	EncounterStatusCompleted  EncounterStatus = "completed"
	EncounterStatusCancelled  EncounterStatus = "cancelled"
)

// Priority represents encounter priority
type Priority string

const (
	PriorityRoutine  Priority = "routine"
	PriorityUrgent   Priority = "urgent"
	PriorityEmergent Priority = "emergent"
)

// ClinicalNote represents SOAP notes and other clinical documentation
type ClinicalNote struct {
	AuditableModel
	EncounterID uuid.UUID    `gorm:"type:uuid;not null;index" json:"encounter_id"`
	Encounter   Encounter    `gorm:"foreignKey:EncounterID" json:"-"`
	NoteType    NoteType     `gorm:"type:varchar(50);not null" json:"note_type"`
	Subjective  string       `gorm:"type:text" json:"subjective"` // S - Subjective
	Objective   string       `gorm:"type:text" json:"objective"`  // O - Objective
	Assessment  string       `gorm:"type:text" json:"assessment"` // A - Assessment
	Plan        string       `gorm:"type:text" json:"plan"`       // P - Plan
	Content     string       `gorm:"type:text" json:"content"`    // General content
	AuthorID    uuid.UUID    `gorm:"type:uuid;not null" json:"author_id"`
	Author      User         `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	SignedAt    *time.Time   `json:"signed_at"`
	SignedBy    *uuid.UUID   `gorm:"type:uuid" json:"signed_by"`
}

// NoteType represents type of clinical note
type NoteType string

const (
	NoteTypeSOAP       NoteType = "soap"
	NoteTypeProgress   NoteType = "progress"
	NoteTypeConsult    NoteType = "consult"
	NoteTypeDischarge  NoteType = "discharge"
	NoteTypeProcedure  NoteType = "procedure"
)

// Diagnosis represents a clinical diagnosis
type Diagnosis struct {
	AuditableModel
	EncounterID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"encounter_id"`
	Encounter       Encounter      `gorm:"foreignKey:EncounterID" json:"-"`
	ICD10Code       string         `gorm:"not null" json:"icd10_code"`
	Description     string         `gorm:"not null" json:"description"`
	DiagnosisType   DiagnosisType  `gorm:"type:varchar(20);not null" json:"diagnosis_type"`
	Status          string         `gorm:"default:'active'" json:"status"`
	OnsetDate       *time.Time     `json:"onset_date"`
	ResolvedDate    *time.Time     `json:"resolved_date"`
	Severity        string         `json:"severity"`
	Notes           string         `json:"notes"`
	DiagnosedBy     uuid.UUID      `gorm:"type:uuid;not null" json:"diagnosed_by"`
}

// DiagnosisType represents type of diagnosis
type DiagnosisType string

const (
	DiagnosisTypePrimary   DiagnosisType = "primary"
	DiagnosisTypeSecondary DiagnosisType = "secondary"
	DiagnosisTypeDifferential DiagnosisType = "differential"
)

// Procedure represents a medical procedure
type Procedure struct {
	AuditableModel
	EncounterID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"encounter_id"`
	Encounter       Encounter  `gorm:"foreignKey:EncounterID" json:"-"`
	ProcedureCode   string     `json:"procedure_code"` // CPT code
	ProcedureName   string     `gorm:"not null" json:"procedure_name"`
	Description     string     `json:"description"`
	PerformedAt     time.Time  `gorm:"not null" json:"performed_at"`
	PerformedBy     uuid.UUID  `gorm:"type:uuid;not null" json:"performed_by"`
	Location        string     `json:"location"`
	Duration        int        `json:"duration"` // in minutes
	Status          string     `gorm:"default:'completed'" json:"status"`
	Complications   string     `json:"complications"`
	Notes           string     `json:"notes"`
}

// VitalSign represents patient vital signs
type VitalSign struct {
	AuditableModel
	EncounterID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"encounter_id"`
	Encounter         Encounter  `gorm:"foreignKey:EncounterID" json:"-"`
	PatientID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"patient_id"`
	Patient           Patient    `gorm:"foreignKey:PatientID" json:"-"`
	MeasuredAt        time.Time  `gorm:"not null" json:"measured_at"`
	Temperature       *float64   `json:"temperature"`        // Celsius
	TemperatureUnit   string     `json:"temperature_unit"`
	HeartRate         *int       `json:"heart_rate"`         // bpm
	RespiratoryRate   *int       `json:"respiratory_rate"`   // breaths/min
	BloodPressureSystolic  *int  `json:"blood_pressure_systolic"`  // mmHg
	BloodPressureDiastolic *int  `json:"blood_pressure_diastolic"` // mmHg
	OxygenSaturation  *float64   `json:"oxygen_saturation"`  // %
	Weight            *float64   `json:"weight"`             // kg
	Height            *float64   `json:"height"`             // cm
	BMI               *float64   `json:"bmi"`
	Pain              *int       `json:"pain"`               // 0-10 scale
	RecordedBy        uuid.UUID  `gorm:"type:uuid;not null" json:"recorded_by"`
	Notes             string     `json:"notes"`
}

// TableName specifies table names
func (Encounter) TableName() string     { return "encounters" }
func (ClinicalNote) TableName() string  { return "clinical_notes" }
func (Diagnosis) TableName() string     { return "diagnoses" }
func (Procedure) TableName() string     { return "procedures" }
func (VitalSign) TableName() string     { return "vital_signs" }
