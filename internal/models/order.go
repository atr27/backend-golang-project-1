package models

import (
	"time"

	"github.com/google/uuid"
)

// Order represents a clinical order (lab, radiology, prescription)
type Order struct {
	AuditableModel
	OrderNumber   string      `gorm:"uniqueIndex;not null" json:"order_number"`
	EncounterID   uuid.UUID   `gorm:"type:uuid;not null;index" json:"encounter_id"`
	Encounter     Encounter   `gorm:"foreignKey:EncounterID" json:"encounter,omitempty"`
	PatientID     uuid.UUID   `gorm:"type:uuid;not null;index" json:"patient_id"`
	Patient       Patient     `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	OrderType     OrderType   `gorm:"type:varchar(50);not null" json:"order_type"`
	Status        OrderStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Priority      Priority    `gorm:"type:varchar(20)" json:"priority"`
	OrderedBy     uuid.UUID   `gorm:"type:uuid;not null" json:"ordered_by"`
	OrderedAt     time.Time   `gorm:"not null" json:"ordered_at"`
	ScheduledFor  *time.Time  `json:"scheduled_for"`
	CompletedAt   *time.Time  `json:"completed_at"`
	CancelledAt   *time.Time  `json:"cancelled_at"`
	CancelReason  string      `json:"cancel_reason"`
	Instructions  string      `json:"instructions"`
	ClinicalNotes string      `json:"clinical_notes"`
	
	// Lab Order specific fields
	LabTests      []LabTest   `gorm:"foreignKey:OrderID" json:"lab_tests,omitempty"`
	
	// Radiology Order specific fields
	RadiologyExams []RadiologyExam `gorm:"foreignKey:OrderID" json:"radiology_exams,omitempty"`
	
	// Prescription specific fields
	Prescriptions []Prescription `gorm:"foreignKey:OrderID" json:"prescriptions,omitempty"`
}

// OrderType represents type of order
type OrderType string

const (
	OrderTypeLab        OrderType = "lab"
	OrderTypeRadiology  OrderType = "radiology"
	OrderTypePrescription OrderType = "prescription"
	OrderTypeProcedure  OrderType = "procedure"
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusScheduled  OrderStatus = "scheduled"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusOnHold     OrderStatus = "on_hold"
)

// LabTest represents a laboratory test order
type LabTest struct {
	AuditableModel
	OrderID       uuid.UUID   `gorm:"type:uuid;not null;index" json:"order_id"`
	Order         Order       `gorm:"foreignKey:OrderID" json:"-"`
	TestCode      string      `gorm:"not null" json:"test_code"` // LOINC code
	TestName      string      `gorm:"not null" json:"test_name"`
	Category      string      `json:"category"`
	Status        OrderStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	SampleType    string      `json:"sample_type"`
	SampleCollectedAt *time.Time `json:"sample_collected_at"`
	ResultsAvailableAt *time.Time `json:"results_available_at"`
	Results       []LabResult `gorm:"foreignKey:LabTestID" json:"results,omitempty"`
}

// LabResult represents laboratory test results
type LabResult struct {
	AuditableModel
	LabTestID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"lab_test_id"`
	LabTest       LabTest    `gorm:"foreignKey:LabTestID" json:"-"`
	ParameterName string     `gorm:"not null" json:"parameter_name"`
	Value         string     `gorm:"not null" json:"value"`
	Unit          string     `json:"unit"`
	ReferenceRange string    `json:"reference_range"`
	Flag          ResultFlag `gorm:"type:varchar(20)" json:"flag"`
	Notes         string     `json:"notes"`
	VerifiedBy    *uuid.UUID `gorm:"type:uuid" json:"verified_by"`
	VerifiedAt    *time.Time `json:"verified_at"`
}

// ResultFlag represents result flag (normal, abnormal, critical)
type ResultFlag string

const (
	ResultFlagNormal   ResultFlag = "normal"
	ResultFlagHigh     ResultFlag = "high"
	ResultFlagLow      ResultFlag = "low"
	ResultFlagCritical ResultFlag = "critical"
	ResultFlagAbnormal ResultFlag = "abnormal"
)

// RadiologyExam represents a radiology examination order
type RadiologyExam struct {
	AuditableModel
	OrderID        uuid.UUID   `gorm:"type:uuid;not null;index" json:"order_id"`
	Order          Order       `gorm:"foreignKey:OrderID" json:"-"`
	ExamCode       string      `gorm:"not null" json:"exam_code"` // CPT code
	ExamName       string      `gorm:"not null" json:"exam_name"`
	Modality       string      `json:"modality"` // X-Ray, CT, MRI, Ultrasound
	BodyPart       string      `json:"body_part"`
	Status         OrderStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ScheduledAt    *time.Time  `json:"scheduled_at"`
	PerformedAt    *time.Time  `json:"performed_at"`
	ReportedAt     *time.Time  `json:"reported_at"`
	Findings       string      `gorm:"type:text" json:"findings"`
	Impression     string      `gorm:"type:text" json:"impression"`
	Radiologist    *uuid.UUID  `gorm:"type:uuid" json:"radiologist"`
	DICOMStudyUID  string      `json:"dicom_study_uid"`
	ImageURL       string      `json:"image_url"`
}

// Prescription represents a medication prescription
type Prescription struct {
	AuditableModel
	OrderID         uuid.UUID          `gorm:"type:uuid;not null;index" json:"order_id"`
	Order           Order              `gorm:"foreignKey:OrderID" json:"-"`
	MedicationName  string             `gorm:"not null" json:"medication_name"`
	GenericName     string             `json:"generic_name"`
	DrugCode        string             `json:"drug_code"` // RxNorm code
	Dosage          string             `gorm:"not null" json:"dosage"`
	Unit            string             `json:"unit"`
	Route           string             `json:"route"`
	Frequency       string             `json:"frequency"`
	Duration        string             `json:"duration"`
	Quantity        int                `json:"quantity"`
	Refills         int                `json:"refills"`
	Instructions    string             `json:"instructions"`
	Status          PrescriptionStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	PrescribedAt    time.Time          `gorm:"not null" json:"prescribed_at"`
	StartDate       *time.Time         `json:"start_date"`
	EndDate         *time.Time         `json:"end_date"`
	PharmacyID      *uuid.UUID         `gorm:"type:uuid" json:"pharmacy_id"`
	DispensedAt     *time.Time         `json:"dispensed_at"`
	DispensedBy     *uuid.UUID         `gorm:"type:uuid" json:"dispensed_by"`
}

// PrescriptionStatus represents prescription status
type PrescriptionStatus string

const (
	PrescriptionStatusPending   PrescriptionStatus = "pending"
	PrescriptionStatusActive    PrescriptionStatus = "active"
	PrescriptionStatusDispensed PrescriptionStatus = "dispensed"
	PrescriptionStatusCompleted PrescriptionStatus = "completed"
	PrescriptionStatusCancelled PrescriptionStatus = "cancelled"
	PrescriptionStatusExpired   PrescriptionStatus = "expired"
)

// TableName specifies table names
func (Order) TableName() string          { return "orders" }
func (LabTest) TableName() string        { return "lab_tests" }
func (LabResult) TableName() string      { return "lab_results" }
func (RadiologyExam) TableName() string  { return "radiology_exams" }
func (Prescription) TableName() string   { return "prescriptions" }
