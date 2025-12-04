package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hospital-emr/backend/internal/common/config"
	"github.com/hospital-emr/backend/internal/common/database"
	"github.com/hospital-emr/backend/internal/common/logger"
	"github.com/hospital-emr/backend/internal/models"
	"github.com/hospital-emr/backend/pkg/encryption"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		return
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

	logger.Info("Starting database seeding...")

	// Seed roles and permissions
	seedRolesAndPermissions(db)

	// Seed default admin user
	seedUsers(db)

	// Seed appointments
	seedAppointments(db)

	// Seed encounters
	seedEncounters(db)

	logger.Info("Database seeding completed successfully")
}

func seedRolesAndPermissions(db *database.DB) {
	logger.Info("Seeding roles and permissions...")

	// Create permissions
	permissions := []models.Permission{
		{Name: "View Patients", Code: models.PermissionViewPatients, Resource: "patient", Action: "view"},
		{Name: "Create Patients", Code: models.PermissionCreatePatients, Resource: "patient", Action: "create"},
		{Name: "Update Patients", Code: models.PermissionUpdatePatients, Resource: "patient", Action: "update"},
		{Name: "Delete Patients", Code: models.PermissionDeletePatients, Resource: "patient", Action: "delete"},
		{Name: "View Encounters", Code: models.PermissionViewEncounters, Resource: "encounter", Action: "view"},
		{Name: "Create Encounters", Code: models.PermissionCreateEncounters, Resource: "encounter", Action: "create"},
		{Name: "Update Encounters", Code: models.PermissionUpdateEncounters, Resource: "encounter", Action: "update"},
		{Name: "View Orders", Code: models.PermissionViewOrders, Resource: "order", Action: "view"},
		{Name: "Create Orders", Code: models.PermissionCreateOrders, Resource: "order", Action: "create"},
		{Name: "Update Orders", Code: models.PermissionUpdateOrders, Resource: "order", Action: "update"},
		{Name: "View Results", Code: models.PermissionViewResults, Resource: "result", Action: "view"},
		{Name: "Update Results", Code: models.PermissionUpdateResults, Resource: "result", Action: "update"},
		{Name: "Manage Users", Code: models.PermissionManageUsers, Resource: "user", Action: "manage"},
		{Name: "Manage Roles", Code: models.PermissionManageRoles, Resource: "role", Action: "manage"},
		{Name: "View Audit Log", Code: models.PermissionViewAuditLog, Resource: "audit", Action: "view"},
	}

	for i := range permissions {
		db.FirstOrCreate(&permissions[i], models.Permission{Code: permissions[i].Code})
	}

	// Create roles
	adminRole := models.Role{
		Name:        "Administrator",
		Code:        models.RoleAdmin,
		Description: "Full system access",
		IsActive:    true,
	}
	db.FirstOrCreate(&adminRole, models.Role{Code: models.RoleAdmin})
	db.Model(&adminRole).Association("Permissions").Replace(permissions)

	doctorPermissions := []models.Permission{}
	doctorPermissionCodes := []string{
		models.PermissionViewPatients, models.PermissionCreatePatients, models.PermissionUpdatePatients,
		models.PermissionViewEncounters, models.PermissionCreateEncounters, models.PermissionUpdateEncounters,
		models.PermissionViewOrders, models.PermissionCreateOrders,
		models.PermissionViewResults,
	}
	for _, code := range doctorPermissionCodes {
		for _, perm := range permissions {
			if perm.Code == code {
				doctorPermissions = append(doctorPermissions, perm)
			}
		}
	}

	doctorRole := models.Role{
		Name:        "Doctor",
		Code:        models.RoleDoctor,
		Description: "Medical doctor with clinical access",
		IsActive:    true,
	}
	db.FirstOrCreate(&doctorRole, models.Role{Code: models.RoleDoctor})
	db.Model(&doctorRole).Association("Permissions").Replace(doctorPermissions)

	nurseRole := models.Role{
		Name:        "Nurse",
		Code:        models.RoleNurse,
		Description: "Nursing staff",
		IsActive:    true,
	}
	db.FirstOrCreate(&nurseRole, models.Role{Code: models.RoleNurse})

	receptionistRole := models.Role{
		Name:        "Receptionist",
		Code:        models.RoleReceptionist,
		Description: "Front desk staff",
		IsActive:    true,
	}
	db.FirstOrCreate(&receptionistRole, models.Role{Code: models.RoleReceptionist})

	logger.Info("Roles and permissions seeded successfully")
}

func seedUsers(db *database.DB) {
	logger.Info("Seeding users...")

	// Check if admin user exists
	var existingUser models.User
	if err := db.Where("email = ?", "admin@hospital-emr.com").First(&existingUser).Error; err == nil {
		logger.Info("Admin user already exists, skipping...")
		return
	}

	// Hash password
	passwordHash, err := encryption.HashPassword("password123")
	if err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		return
	}

	// Get admin role
	var adminRole models.Role
	if err := db.Where("code = ?", models.RoleAdmin).First(&adminRole).Error; err != nil {
		logger.Errorf("Failed to find admin role: %v", err)
		return
	}

	// Create admin user
	adminUser := models.User{
		Email:        "admin@hospital-emr.com",
		PasswordHash: passwordHash,
		FirstName:    "System",
		LastName:     "Administrator",
		PhoneNumber:  "+1234567890",
		Status:       models.UserStatusActive,
		MFAEnabled:   false,
		Department:   "Administration",
	}

	if err := db.Create(&adminUser).Error; err != nil {
		logger.Errorf("Failed to create admin user: %v", err)
		return
	}

	// Assign admin role
	db.Model(&adminUser).Association("Roles").Append(&adminRole)

	logger.Info("Admin user created successfully")
	logger.Info("Email: admin@hospital-emr.com")
	logger.Info("Password: admin123")
	logger.Warn("Please change the default password immediately!")

	// Create sample doctor
	doctorPassword, _ := encryption.HashPassword("password123")
	var doctorRole models.Role
	db.Where("code = ?", models.RoleDoctor).First(&doctorRole)

	doctorUser := models.User{
		Email:         "doctor@hospital-emr.com",
		PasswordHash:  doctorPassword,
		FirstName:     "Budi",
		LastName:      "Santoso",
		PhoneNumber:   "+6281234567890",
		Status:        models.UserStatusActive,
		MFAEnabled:    false,
		LicenseNumber: "STR-123456",
		Specialty:     "Penyakit Dalam",
		Department:    "Poli Umum",
	}

	if err := db.Create(&doctorUser).Error; err == nil {
		db.Model(&doctorUser).Association("Roles").Append(&doctorRole)
		logger.Info("Sample doctor created: doctor@hospital-emr.com / doctor123")
	}

	// Create sample patients
	seedSamplePatients(db, adminUser.ID)
}

func seedSamplePatients(db *database.DB, createdBy uuid.UUID) {
	logger.Info("Seeding sample patients...")

	patients := []models.Patient{
		{
			MRN:         "MRN000001",
			FirstName:   "Siti",
			LastName:    "Aminah",
			DateOfBirth: time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderFemale,
			BloodType:   "A+",
			Email:       "siti.aminah@email.com",
			PhoneNumber: "+6281234567891",
			Address:     "Jl. Sudirman No. 123",
			City:        "Jakarta Pusat",
			State:       "DKI Jakarta",
			ZipCode:     "10220",
			Country:     "Indonesia",
			Status:      models.PatientStatusActive,
		},
		{
			MRN:         "MRN000002",
			FirstName:   "Ahmad",
			LastName:    "Rizki",
			DateOfBirth: time.Date(1990, 8, 22, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderMale,
			BloodType:   "O+",
			Email:       "ahmad.rizki@email.com",
			PhoneNumber: "+6281234567892",
			Address:     "Jl. Asia Afrika No. 45",
			City:        "Bandung",
			State:       "Jawa Barat",
			ZipCode:     "40111",
			Country:     "Indonesia",
			Status:      models.PatientStatusActive,
		},
		{
			MRN:         "MRN000003",
			FirstName:   "Dewi",
			LastName:    "Sartika",
			DateOfBirth: time.Date(1978, 11, 30, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderFemale,
			BloodType:   "B-",
			Email:       "dewi.sartika@email.com",
			PhoneNumber: "+6281234567893",
			Address:     "Jl. Malioboro No. 10",
			City:        "Yogyakarta",
			State:       "DIY",
			ZipCode:     "55213",
			Country:     "Indonesia",
			Status:      models.PatientStatusActive,
		},
		{
			MRN:         "MRN000004",
			FirstName:   "Bambang",
			LastName:    "Pamungkas",
			DateOfBirth: time.Date(1992, 3, 10, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderMale,
			BloodType:   "AB+",
			Email:       "bambang.pamungkas@email.com",
			PhoneNumber: "+6281234567894",
			Address:     "Jl. Pahlawan No. 88",
			City:        "Surabaya",
			State:       "Jawa Timur",
			ZipCode:     "60174",
			Country:     "Indonesia",
			Status:      models.PatientStatusActive,
		},
		{
			MRN:         "MRN000005",
			FirstName:   "Ratna",
			LastName:    "Sari",
			DateOfBirth: time.Date(1965, 1, 5, 0, 0, 0, 0, time.UTC),
			Gender:      models.GenderFemale,
			BloodType:   "O-",
			Email:       "ratna.sari@email.com",
			PhoneNumber: "+6281234567895",
			Address:     "Jl. Gajah Mada No. 20",
			City:        "Semarang",
			State:       "Jawa Tengah",
			ZipCode:     "50134",
			Country:     "Indonesia",
			Status:      models.PatientStatusActive,
		},
	}

	for i := range patients {
		patients[i].CreatedBy = createdBy
		patients[i].UpdatedBy = createdBy
		db.FirstOrCreate(&patients[i], models.Patient{MRN: patients[i].MRN})
	}

	logger.Info("Sample patients seeded successfully")
}

func seedAppointments(db *database.DB) {
	logger.Info("Seeding appointments...")

	// Get a doctor
	var doctor models.User
	if err := db.Where("email = ?", "doctor@hospital-emr.com").First(&doctor).Error; err != nil {
		logger.Errorf("Failed to find doctor for appointments: %v", err)
		return
	}

	// Get patients
	var patients []models.Patient
	if err := db.Find(&patients).Error; err != nil {
		logger.Errorf("Failed to find patients for appointments: %v", err)
		return
	}

	if len(patients) == 0 {
		logger.Warn("No patients found, skipping appointment seeding")
		return
	}

	// Create appointments
	appointments := []models.Appointment{}
	now := time.Now()

	// 1. Upcoming appointments
	appointments = append(appointments, models.Appointment{
		AppointmentNumber: "APT-001",
		PatientID:         patients[0].ID,
		ProviderID:        doctor.ID,
		AppointmentType:   models.AppointmentTypeConsultation,
		Status:            models.AppointmentStatusScheduled,
		StartTime:         now.Add(24 * time.Hour), // Tomorrow
		EndTime:           now.Add(24*time.Hour + 30*time.Minute),
		Duration:          30,
		Department:        "Poli Umum",
		Location:          "Ruang 101",
		ReasonForVisit:    "Pemeriksaan Rutin",
		Notes:             "Pasien minta jadwal pagi",
	})

	appointments = append(appointments, models.Appointment{
		AppointmentNumber: "APT-002",
		PatientID:         patients[1].ID,
		ProviderID:        doctor.ID,
		AppointmentType:   models.AppointmentTypeFollowUp,
		Status:            models.AppointmentStatusConfirmed,
		StartTime:         now.Add(48 * time.Hour), // Day after tomorrow
		EndTime:           now.Add(48*time.Hour + 15*time.Minute),
		Duration:          15,
		Department:        "Poli Umum",
		Location:          "Ruang 102",
		ReasonForVisit:    "Kontrol Tekanan Darah",
	})

	// 2. Past appointments (Completed)
	appointments = append(appointments, models.Appointment{
		AppointmentNumber: "APT-003",
		PatientID:         patients[0].ID,
		ProviderID:        doctor.ID,
		AppointmentType:   models.AppointmentTypeWellness,
		Status:            models.AppointmentStatusCompleted,
		StartTime:         now.Add(-7 * 24 * time.Hour), // Last week
		EndTime:           now.Add(-7*24*time.Hour + 45*time.Minute),
		Duration:          45,
		Department:        "Poli Umum",
		Location:          "Ruang 101",
		ReasonForVisit:    "Cek Kesehatan Tahunan",
		CheckedInAt:       func() *time.Time { t := now.Add(-7*24*time.Hour - 10*time.Minute); return &t }(),
	})

	// 3. Cancelled appointment
	appointments = append(appointments, models.Appointment{
		AppointmentNumber: "APT-004",
		PatientID:         patients[1].ID,
		ProviderID:        doctor.ID,
		AppointmentType:   models.AppointmentTypeConsultation,
		Status:            models.AppointmentStatusCancelled,
		StartTime:         now.Add(-2 * 24 * time.Hour), // 2 days ago
		EndTime:           now.Add(-2*24*time.Hour + 30*time.Minute),
		Duration:          30,
		Department:        "Poli Umum",
		Location:          "Ruang 101",
		ReasonForVisit:    "Gejala Flu",
		CancelledAt:       func() *time.Time { t := now.Add(-3 * 24 * time.Hour); return &t }(),
		CancellationReason: "Pasien sudah merasa baikan",
	})

	for i := range appointments {
		appointments[i].CreatedBy = doctor.ID
		appointments[i].UpdatedBy = doctor.ID
		if err := db.FirstOrCreate(&appointments[i], models.Appointment{AppointmentNumber: appointments[i].AppointmentNumber}).Error; err != nil {
			logger.Errorf("Failed to seed appointment %s: %v", appointments[i].AppointmentNumber, err)
		}
	}

	logger.Info("Appointments seeded successfully")
}

func seedEncounters(db *database.DB) {
	logger.Info("Seeding encounters...")

	// Get a doctor
	var doctor models.User
	if err := db.Where("email = ?", "doctor@hospital-emr.com").First(&doctor).Error; err != nil {
		logger.Errorf("Failed to find doctor for encounters: %v", err)
		return
	}

	// Get patients
	var patients []models.Patient
	if err := db.Find(&patients).Error; err != nil {
		logger.Errorf("Failed to find patients for encounters: %v", err)
		return
	}

	if len(patients) == 0 {
		logger.Warn("No patients found, skipping encounter seeding")
		return
	}

	now := time.Now()

	// Create Encounters
	encounters := []models.Encounter{
		{
			EncounterNumber: "ENC-001",
			PatientID:       patients[0].ID,
			ProviderID:      doctor.ID,
			EncounterType:   models.EncounterTypeWellness,
			Status:          models.EncounterStatusCompleted,
			Priority:        models.PriorityRoutine,
			Department:      "Poli Umum",
			Location:        "Ruang 101",
			AdmissionDate:   now.Add(-7 * 24 * time.Hour),
			DischargeDate:   func() *time.Time { t := now.Add(-7*24*time.Hour + 45*time.Minute); return &t }(),
			ChiefComplaint:  "Cek Kesehatan Tahunan",
			ReasonForVisit:  "Kunjungan Sehat",
		},
		{
			EncounterNumber: "ENC-002",
			PatientID:       patients[1].ID,
			ProviderID:      doctor.ID,
			EncounterType:   models.EncounterTypeOutpatient,
			Status:          models.EncounterStatusInProgress,
			Priority:        models.PriorityUrgent,
			Department:      "IGD",
			Location:        "Bed 1",
			AdmissionDate:   now.Add(-1 * time.Hour),
			ChiefComplaint:  "Sakit Kepala Hebat",
			ReasonForVisit:  "Migrain",
		},
	}

	for i := range encounters {
		encounters[i].CreatedBy = doctor.ID
		encounters[i].UpdatedBy = doctor.ID
		
		if err := db.FirstOrCreate(&encounters[i], models.Encounter{EncounterNumber: encounters[i].EncounterNumber}).Error; err != nil {
			logger.Errorf("Failed to seed encounter %s: %v", encounters[i].EncounterNumber, err)
			continue
		}

		// Add Vital Signs for the first encounter
		if encounters[i].EncounterNumber == "ENC-001" {
			temp := 37.0
			hr := 72
			rr := 16
			bpSys := 120
			bpDia := 80
			o2 := 98.0
			weight := 70.0
			height := 175.0
			bmi := 22.8

			vitalSign := models.VitalSign{
				EncounterID:            encounters[i].ID,
				PatientID:              encounters[i].PatientID,
				MeasuredAt:             encounters[i].AdmissionDate.Add(5 * time.Minute),
				Temperature:            &temp,
				TemperatureUnit:        "C",
				HeartRate:              &hr,
				RespiratoryRate:        &rr,
				BloodPressureSystolic:  &bpSys,
				BloodPressureDiastolic: &bpDia,
				OxygenSaturation:       &o2,
				Weight:                 &weight,
				Height:                 &height,
				BMI:                    &bmi,
				RecordedBy:             doctor.ID,
			}
			vitalSign.CreatedBy = doctor.ID
			vitalSign.UpdatedBy = doctor.ID
			db.Create(&vitalSign)

			// Add Diagnosis
			diagnosis := models.Diagnosis{
				EncounterID:   encounters[i].ID,
				ICD10Code:     "Z00.00",
				Description:   "Pemeriksaan medis umum dewasa tanpa temuan abnormal",
				DiagnosisType: models.DiagnosisTypePrimary,
				Status:        "Active",
				DiagnosedBy:   doctor.ID,
			}
			diagnosis.CreatedBy = doctor.ID
			diagnosis.UpdatedBy = doctor.ID
			db.Create(&diagnosis)
		}
	}

	logger.Info("Encounters seeded successfully")
}
