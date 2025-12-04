package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hospital-emr/backend/internal/common/config"
	"github.com/hospital-emr/backend/internal/common/errors"
	"github.com/hospital-emr/backend/internal/models"
	"github.com/hospital-emr/backend/pkg/encryption"
	"github.com/hospital-emr/backend/pkg/jwt"
	"gorm.io/gorm"
)

// Service provides authentication services
type Service struct {
	db     *gorm.DB
	config *config.Config
}

// NewService creates a new auth service
func NewService(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{
		db:     db,
		config: cfg,
	}
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	MFACode  string `json:"mfa_code,omitempty"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	User         *models.User `json:"user"`
	MFARequired  bool        `json:"mfa_required,omitempty"`
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*LoginResponse, error) {
	// Find user by email
	var user models.User
	if err := s.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("email = ? AND status = ?", req.Email, models.UserStatusActive).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrInvalidCredentials
		}
		return nil, errors.ErrDatabaseError
	}

	// Verify password
	if !encryption.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.ErrInvalidCredentials
	}

	// Check if MFA is enabled
	if user.MFAEnabled {
		if req.MFACode == "" {
			return &LoginResponse{
				MFARequired: true,
			}, nil
		}

		// TODO: Verify MFA code
		// if !verifyMFACode(user.MFASecret, req.MFACode) {
		// 	return nil, errors.ErrMFAInvalid
		// }
	}

	// Extract roles
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Code
	}

	// Generate tokens
	accessToken, err := jwt.GenerateToken(
		user.ID,
		user.Email,
		roles,
		s.config.JWT.Secret,
		s.config.GetJWTExpiration(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.GenerateToken(
		user.ID,
		user.Email,
		roles,
		s.config.JWT.Secret,
		s.config.GetJWTRefreshExpiration(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session
	session := models.Session{
		UserID:       user.ID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.config.GetJWTExpiration()),
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		IsActive:     true,
	}

	if err := s.db.WithContext(ctx).Create(&session).Error; err != nil {
		return nil, errors.ErrDatabaseError
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = ipAddress
	s.db.WithContext(ctx).Save(&user)

	// Clear sensitive data
	user.PasswordHash = ""
	user.MFASecret = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.config.GetJWTExpiration().Seconds()),
		User:         &user,
		MFARequired:  false,
	}, nil
}

// Logout logs out a user
func (s *Service) Logout(ctx context.Context, userID uuid.UUID, token string) error {
	now := time.Now()
	return s.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("user_id = ? AND token = ?", userID, token).
		Updates(map[string]interface{}{
			"is_active":  false,
			"revoked_at": now,
		}).Error
}

// RefreshToken refreshes an access token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	// Validate refresh token
	claims, err := jwt.ValidateToken(refreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.ErrTokenInvalid
	}

	// Check if session exists and is active
	var session models.Session
	if err := s.db.WithContext(ctx).
		Where("user_id = ? AND refresh_token = ? AND is_active = ?", claims.UserID, refreshToken, true).
		First(&session).Error; err != nil {
		return nil, errors.ErrTokenInvalid
	}

	// Get user
	var user models.User
	if err := s.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("id = ? AND status = ?", claims.UserID, models.UserStatusActive).
		First(&user).Error; err != nil {
		return nil, errors.ErrUserNotFound(claims.UserID.String())
	}

	// Generate new tokens
	accessToken, err := jwt.GenerateToken(
		user.ID,
		user.Email,
		claims.Roles,
		s.config.JWT.Secret,
		s.config.GetJWTExpiration(),
	)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := jwt.GenerateToken(
		user.ID,
		user.Email,
		claims.Roles,
		s.config.JWT.Secret,
		s.config.GetJWTRefreshExpiration(),
	)
	if err != nil {
		return nil, err
	}

	// Update session
	session.Token = accessToken
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(s.config.GetJWTExpiration())
	s.db.WithContext(ctx).Save(&session)

	// Clear sensitive data
	user.PasswordHash = ""
	user.MFASecret = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(s.config.GetJWTExpiration().Seconds()),
		User:         &user,
	}, nil
}

// VerifyToken verifies if a token is valid
func (s *Service) VerifyToken(ctx context.Context, token string) (*models.User, error) {
	claims, err := jwt.ValidateToken(token, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.ErrTokenInvalid
	}

	var user models.User
	if err := s.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("id = ? AND status = ?", claims.UserID, models.UserStatusActive).
		First(&user).Error; err != nil {
		return nil, errors.ErrUserNotFound(claims.UserID.String())
	}

	user.PasswordHash = ""
	user.MFASecret = ""

	return &user, nil
}
