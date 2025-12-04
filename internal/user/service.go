package user

import (
	"context"

	"github.com/hospital-emr/backend/internal/common/errors"
	"github.com/hospital-emr/backend/internal/models"
	"gorm.io/gorm"
)

// Service provides user management services
type Service struct {
	db *gorm.DB
}

// NewService creates a new user service
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// ListUsers lists users with optional role filtering
func (s *Service) ListUsers(ctx context.Context, roleFilter string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.WithContext(ctx).Model(&models.User{}).Preload("Roles")

	// Apply role filter if provided
	if roleFilter != "" {
		query = query.Joins("JOIN user_roles ON user_roles.user_id = users.id").
			Joins("JOIN roles ON roles.id = user_roles.role_id").
			Where("roles.code = ?", roleFilter)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.ErrDatabaseError.WithDetails(err.Error())
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := query.
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, errors.ErrDatabaseError.WithDetails(err.Error())
	}

	return users, total, nil
}
