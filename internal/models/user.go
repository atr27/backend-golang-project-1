package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a system user (doctor, nurse, admin, etc.)
type User struct {
	BaseModel
	Email           string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash    string     `gorm:"not null" json:"-"`
	FirstName       string     `gorm:"not null" json:"first_name"`
	LastName        string     `gorm:"not null" json:"last_name"`
	PhoneNumber     string     `json:"phone_number"`
	Status          UserStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	MFAEnabled      bool       `gorm:"default:false" json:"mfa_enabled"`
	MFASecret       string     `json:"-"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	LastLoginIP     string     `json:"last_login_ip"`
	PasswordExpiry  *time.Time `json:"password_expiry"`
	LicenseNumber   string     `json:"license_number"`
	Specialty       string     `json:"specialty"`
	Department      string     `json:"department"`
	Roles           []Role     `gorm:"many2many:user_roles;" json:"roles"`
	Sessions        []Session  `gorm:"foreignKey:UserID" json:"-"`
}

// UserStatus represents user account status
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusLocked    UserStatus = "locked"
)

// Role represents a user role
type Role struct {
	BaseModel
	Name        string       `gorm:"uniqueIndex;not null" json:"name"`
	Code        string       `gorm:"uniqueIndex;not null" json:"code"`
	Description string       `json:"description"`
	IsActive    bool         `gorm:"default:true" json:"is_active"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
	Users       []User       `gorm:"many2many:user_roles;" json:"-"`
}

// Permission represents a system permission
type Permission struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Code        string `gorm:"uniqueIndex;not null" json:"code"`
	Resource    string `gorm:"not null" json:"resource"`
	Action      string `gorm:"not null" json:"action"`
	Description string `json:"description"`
	Roles       []Role `gorm:"many2many:role_permissions;" json:"-"`
}

// Session represents a user session
type Session struct {
	BaseModel
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"-"`
	Token        string    `gorm:"uniqueIndex;not null" json:"-"`
	RefreshToken string    `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	RevokedAt    *time.Time `json:"revoked_at"`
}

// TableName specifies table names
func (User) TableName() string       { return "users" }
func (Role) TableName() string       { return "roles" }
func (Permission) TableName() string { return "permissions" }
func (Session) TableName() string    { return "sessions" }

// Common roles
const (
	RoleAdmin       = "admin"
	RoleDoctor      = "doctor"
	RoleNurse       = "nurse"
	RoleReceptionist = "receptionist"
	RolePharmacist  = "pharmacist"
	RoleLabTech     = "lab_technician"
	RoleRadiologist = "radiologist"
	RolePatient     = "patient"
)

// Common permissions
const (
	PermissionViewPatients   = "view_patients"
	PermissionCreatePatients = "create_patients"
	PermissionUpdatePatients = "update_patients"
	PermissionDeletePatients = "delete_patients"
	
	PermissionViewEncounters   = "view_encounters"
	PermissionCreateEncounters = "create_encounters"
	PermissionUpdateEncounters = "update_encounters"
	
	PermissionViewOrders   = "view_orders"
	PermissionCreateOrders = "create_orders"
	PermissionUpdateOrders = "update_orders"
	
	PermissionViewResults   = "view_results"
	PermissionUpdateResults = "update_results"
	
	PermissionManageUsers = "manage_users"
	PermissionManageRoles = "manage_roles"
	PermissionViewAuditLog = "view_audit_log"
)
