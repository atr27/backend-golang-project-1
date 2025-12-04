package main

import (
	"fmt"
	"os"

	"github.com/hospital-emr/backend/internal/common/config"
	"github.com/hospital-emr/backend/internal/common/database"
	"github.com/hospital-emr/backend/internal/common/logger"
	"github.com/hospital-emr/backend/internal/models"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate [up|down|create]")
		os.Exit(1)
	}

	command := os.Args[1]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(logger.Config{
		Level:  "info",
		Format: "console",
	})

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	switch command {
	case "up":
		migrateUp(db)
	case "down":
		migrateDown(db)
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Usage: migrate create [migration_name]")
			os.Exit(1)
		}
		createMigration(os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func migrateUp(db *database.DB) {
	logger.Info("Running migrations...")

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
			logger.Fatalf("Failed to migrate %T: %v", model, err)
		}
		logger.Infof("Migrated: %T", model)
	}

	logger.Info("All migrations completed successfully")
}

func migrateDown(db *database.DB) {
	logger.Info("Rolling back migrations...")

	// Drop join tables first
	if err := db.Exec("DROP TABLE IF EXISTS user_roles CASCADE").Error; err != nil {
		logger.Errorf("Failed to drop user_roles: %v", err)
	}
	if err := db.Exec("DROP TABLE IF EXISTS role_permissions CASCADE").Error; err != nil {
		logger.Errorf("Failed to drop role_permissions: %v", err)
	}

	models := []interface{}{
		&models.AuditLog{},
		&models.Prescription{},
		&models.RadiologyExam{},
		&models.LabResult{},
		&models.LabTest{},
		&models.Order{},
		&models.Appointment{},
		&models.VitalSign{},
		&models.Procedure{},
		&models.Diagnosis{},
		&models.ClinicalNote{},
		&models.Encounter{},
		&models.Medication{},
		&models.Allergy{},
		&models.Patient{},
		&models.Session{},
		&models.Permission{},
		&models.Role{},
		&models.User{},
	}

	for _, model := range models {
		if err := db.Migrator().DropTable(model); err != nil {
			logger.Errorf("Failed to drop table for %T: %v", model, err)
		} else {
			logger.Infof("Dropped table: %T", model)
		}
	}

	logger.Info("All migrations rolled back")
}

func createMigration(name string) {
	logger.Infof("Creating migration: %s", name)
	// TODO: Implement migration file creation
	logger.Info("Migration file created (not implemented yet)")
}
