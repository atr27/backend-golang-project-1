package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Timestamp     time.Time  `gorm:"not null;index" json:"timestamp"`
	UserID        *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Username      string     `json:"username"`
	Action        string     `gorm:"not null;index" json:"action"` // CREATE, READ, UPDATE, DELETE, LOGIN, LOGOUT
	Resource      string     `gorm:"not null;index" json:"resource"` // patient, encounter, order, etc.
	ResourceID    *uuid.UUID `gorm:"type:uuid;index" json:"resource_id"`
	Description   string     `json:"description"`
	IPAddress     string     `json:"ip_address"`
	UserAgent     string     `json:"user_agent"`
	RequestMethod string     `json:"request_method"`
	RequestPath   string     `json:"request_path"`
	StatusCode    int        `json:"status_code"`
	ChangesOld    string     `gorm:"type:jsonb" json:"changes_old"`
	ChangesNew    string     `gorm:"type:jsonb" json:"changes_new"`
	Metadata      string     `gorm:"type:jsonb" json:"metadata"`
	Severity      AuditSeverity `gorm:"type:varchar(20)" json:"severity"`
}

// AuditSeverity represents audit log severity
type AuditSeverity string

const (
	AuditSeverityInfo     AuditSeverity = "info"
	AuditSeverityWarning  AuditSeverity = "warning"
	AuditSeverityError    AuditSeverity = "error"
	AuditSeverityCritical AuditSeverity = "critical"
)

// TableName specifies table name
func (AuditLog) TableName() string { return "audit_logs" }

// BeforeCreate sets timestamp
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now().UTC()
	}
	return nil
}
