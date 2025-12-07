package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hospital-emr/backend/internal/auth"
	"github.com/hospital-emr/backend/internal/common/config"
	"github.com/hospital-emr/backend/internal/common/database"
	"github.com/hospital-emr/backend/internal/common/logger"
	"github.com/hospital-emr/backend/internal/common/middleware"
	"github.com/hospital-emr/backend/internal/encounter"
	"github.com/hospital-emr/backend/internal/models"
	"github.com/hospital-emr/backend/internal/patient"
	"github.com/hospital-emr/backend/internal/scheduling"
	"github.com/hospital-emr/backend/internal/user"
	"github.com/hospital-emr/backend/pkg/messaging"
	_ "github.com/hospital-emr/backend/api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Hospital EMR System API
// @version 1.0
// @description Comprehensive Electronic Medical Record system for hospitals
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@hospital-emr.com

// @license.name Proprietary
// @license.url http://www.hospital-emr.com/license

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(logger.Config{
		Level:  cfg.App.LogLevel,
		Format: cfg.App.LogFormat,
	})

	logger.Infof("Starting %s v%s in %s mode", cfg.App.Name, cfg.App.Version, cfg.App.Environment)

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Auto-migrate database models
	if err := autoMigrate(db); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}

	// Connect to NATS
	var natsClient *messaging.NATSClient
	if cfg.NATS.URL != "" {
		natsClient, err = messaging.NewNATSClient(cfg.NATS.URL)
		if err != nil {
			logger.Warnf("Failed to connect to NATS: %v", err)
			// Continue without NATS for now
		} else {
			defer natsClient.Close()
			logger.Info("Connected to NATS successfully")
		}
	} else {
		logger.Info("NATS URL not provided, skipping NATS connection")
	}

	// Initialize services
	authService := auth.NewService(db.DB, cfg)
	patientService := patient.NewService(db.DB, natsClient)
	encounterService := encounter.NewService(db.DB, natsClient)
	schedulingService := scheduling.NewService(db.DB, natsClient)
	userService := user.NewService(db.DB)

	// Initialize handlers
	authHandler := auth.NewHandler(authService)
	patientHandler := patient.NewHandler(patientService)
	encounterHandler := encounter.NewHandler(encounterService)
	schedulingHandler := scheduling.NewHandler(schedulingService)
	userHandler := user.NewHandler(userService)

	// Setup router
	router := setupRouter(cfg, authHandler, patientHandler, encounterHandler, schedulingHandler, userHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:           ":" + cfg.App.Port,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server in goroutine
	go func() {
		logger.Infof("Server starting on port %s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

func setupRouter(cfg *config.Config, authHandler *auth.Handler, patientHandler *patient.Handler, encounterHandler *encounter.Handler, schedulingHandler *scheduling.Handler, userHandler *user.Handler) *gin.Engine {
	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))
	router.Use(middleware.RateLimiter(cfg.Security.RateLimitPerMinute))

	// Health check endpoints
	router.GET("/health", healthCheck)
	router.GET("/ready", readyCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/otentikasi")
		{
			auth.POST("/masuk", authHandler.Login)
			auth.POST("/segarkan", authHandler.RefreshToken)
		}

		// Protected routes (authentication required)
		authenticated := v1.Group("")
		authenticated.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		authenticated.Use(middleware.AuditLog())
		{
			// Auth routes
			authRoutes := authenticated.Group("/otentikasi")
			{
				authRoutes.POST("/keluar", authHandler.Logout)
				authRoutes.GET("/verifikasi", authHandler.VerifyToken)
			}

			// Patient routes
			patients := authenticated.Group("/pasien")
			{
				patients.GET("", patientHandler.ListPatients)
				patients.POST("", patientHandler.CreatePatient)
				patients.GET("/:id", patientHandler.GetPatient)
				patients.PUT("/:id", patientHandler.UpdatePatient)
				// Only admins can delete patients
				patients.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), patientHandler.DeletePatient)
				patients.GET("/:id/riwayat", patientHandler.GetPatientTimeline)
			}

			// Encounter routes
			encounters := authenticated.Group("/kunjungan")
			{
				encounters.GET("", encounterHandler.ListEncounters)
				encounters.POST("", encounterHandler.CreateEncounter)
				encounters.GET("/:id", encounterHandler.GetEncounter)
				encounters.PUT("/:id/status", encounterHandler.UpdateEncounterStatus)
				encounters.POST("/:id/selesai", encounterHandler.CompleteEncounter)
				encounters.POST("/:id/catatan", encounterHandler.AddClinicalNote)
				encounters.POST("/:id/diagnosis", encounterHandler.AddDiagnosis)
				encounters.POST("/:id/tanda-vital", encounterHandler.RecordVitalSigns)
			}

			// Appointment/Scheduling routes
			appointments := authenticated.Group("/janji-temu")
			{
				appointments.GET("", schedulingHandler.ListAppointments)
				appointments.POST("", schedulingHandler.CreateAppointment)
				appointments.GET("/:id", schedulingHandler.GetAppointment)
				appointments.PUT("/:id", schedulingHandler.UpdateAppointment)
				appointments.POST("/:id/check-in", schedulingHandler.CheckInAppointment)
				appointments.POST("/:id/batal", schedulingHandler.CancelAppointment)
				appointments.GET("/ketersediaan", schedulingHandler.GetAvailability)
			}

			// User routes
			users := authenticated.Group("/pengguna")
			{
				users.GET("", userHandler.ListUsers)
			}
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve static files from frontend/dist
	router.Static("/assets", "./frontend/dist/assets")

	// SPA fallback - serve index.html for all non-API routes
	router.NoRoute(func(c *gin.Context) {
		// Don't serve index.html for API routes
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}
		// Serve index.html for all other routes (SPA client-side routing)
		c.File("./frontend/dist/index.html")
	})

	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Server is running",
		"time":    time.Now().UTC(),
	})
}

func readyCheck(c *gin.Context) {
	// TODO: Check database and other dependencies
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": gin.H{
			"database": "ok",
			"nats":     "ok",
		},
	})
}

func autoMigrate(db *database.DB) error {
	logger.Info("Running database migrations...")

	models := []interface{}{
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Session{},
		&models.Patient{},
		&models.Allergy{},
		&models.Medication{},
		&models.Encounter{},
		&models.ClinicalNote{},
		&models.Diagnosis{},
		&models.Procedure{},
		&models.VitalSign{},
		&models.Appointment{},
		&models.Order{},
		&models.LabTest{},
		&models.LabResult{},
		&models.RadiologyExam{},
		&models.Prescription{},
		&models.AuditLog{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

