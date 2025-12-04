package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hospital-emr/backend/internal/common/config"
	"github.com/hospital-emr/backend/internal/common/database"
	"github.com/hospital-emr/backend/internal/models"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test querying patients
	var patients []models.Patient
	result := db.DB.WithContext(context.Background()).
		Limit(5).
		Find(&patients)

	if result.Error != nil {
		log.Fatalf("Error querying patients: %v", result.Error)
	}

	fmt.Printf("Successfully queried %d patients\n", len(patients))
	for _, p := range patients {
		fmt.Printf("Patient: %s %s (MRN: %s)\n", p.FirstName, p.LastName, p.MRN)
	}
}
