package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hospital-emr/backend/internal/common/errors"
	_ "github.com/hospital-emr/backend/internal/models"
)

// Handler handles authentication HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new auth handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	resp, err := h.service.Login(
		c.Request.Context(),
		&req,
		c.ClientIP(),
		c.Request.UserAgent(),
	)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate session
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} errors.AppError
// @Router /api/v1/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	token := c.GetHeader("Authorization")
	if len(token) > 7 {
		token = token[7:] // Remove "Bearer " prefix
	}

	if err := h.service.Logout(c.Request.Context(), userID.(uuid.UUID), token); err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} errors.AppError
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrBadRequest.WithDetails(err.Error()))
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.ErrInternal)
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// VerifyToken godoc
// @Summary Verify JWT token
// @Description Verify if the provided token is valid
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} errors.AppError
// @Router /api/v1/auth/verify [get]
func (h *Handler) VerifyToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) > 7 {
		token = token[7:] // Remove "Bearer " prefix
	}

	user, err := h.service.VerifyToken(c.Request.Context(), token)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.StatusCode, appErr)
		} else {
			c.JSON(http.StatusUnauthorized, errors.ErrTokenInvalid)
		}
		return
	}

	c.JSON(http.StatusOK, user)
}
