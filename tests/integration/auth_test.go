// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hospital-emr/backend/internal/auth"
	"github.com/hospital-emr/backend/internal/common/config"
	"github.com/hospital-emr/backend/internal/common/database"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *database.DB) {
	gin.SetMode(gin.TestMode)
	
	cfg, _ := config.Load()
	db, _ := database.New(cfg)
	
	authService := auth.NewService(db.DB, cfg)
	authHandler := auth.NewHandler(authService)
	
	router := gin.New()
	router.POST("/api/v1/auth/login", authHandler.Login)
	
	return router, db
}

func TestIntegrationLogin(t *testing.T) {
	router, db := setupTestRouter()
	defer db.Close()
	
	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Valid Login",
			requestBody: map[string]string{
				"email":    "admin@hospital-emr.com",
				"password": "admin123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response["access_token"])
				assert.NotEmpty(t, response["refresh_token"])
			},
		},
		{
			name: "Invalid Password",
			requestBody: map[string]string{
				"email":    "admin@hospital-emr.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "INVALID_CREDENTIALS", response["code"])
			},
		},
		{
			name: "User Not Found",
			requestBody: map[string]string{
				"email":    "nonexistent@hospital-emr.com",
				"password": "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse:  func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name: "Missing Email",
			requestBody: map[string]string{
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}
