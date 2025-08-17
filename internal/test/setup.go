package test

import (
	"log"

	"payslip-system/internal/config"
	"payslip-system/internal/database"
	"payslip-system/internal/repository"
	"payslip-system/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB() (*gorm.DB, func()) {
	// Load test environment
	cfg := config.Load()

	testDBURL := cfg.DatabaseURL
	if testDBURL == "" {
		testDBURL = "postgres://payslip_user:payslip_password@localhost/payslip_test_db?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(testDBURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		// Clean up test data
		db.Exec("TRUNCATE TABLE audit_logs CASCADE")
		db.Exec("TRUNCATE TABLE payroll_items CASCADE")
		db.Exec("TRUNCATE TABLE payrolls CASCADE")
		db.Exec("TRUNCATE TABLE reimbursements CASCADE")
		db.Exec("TRUNCATE TABLE overtimes CASCADE")
		db.Exec("TRUNCATE TABLE attendances CASCADE")
		db.Exec("TRUNCATE TABLE attendance_periods CASCADE")
		db.Exec("TRUNCATE TABLE users CASCADE")

		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return db, cleanup
}

func SetupTestServices(db *gorm.DB) (*repository.Repositories, *service.Services) {
	repos := repository.NewRepositories(db)
	services := service.NewServices(repos)
	return repos, services
}
