package database

import (
	"fmt"
	"log"
	"math/rand"
	"payslip-system/internal/models"
	"payslip-system/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.AttendancePeriod{},
		&models.Attendance{},
		&models.Overtime{},
		&models.Reimbursement{},
		&models.Payroll{},
		&models.PayrollItem{},
		&models.AuditLog{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func SeedDatabase(repos *repository.Repositories) error {
	log.Println("Seeding database...")

	// Check if users already exist
	var userCount int64
	if err := repos.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount > 0 {
		log.Println("Database already seeded")
		return nil
	}

	// Create admin user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &models.User{
		Username: "admin",
		Password: string(hashedPassword),
		Role:     "admin",
		IsActive: true,
	}

	if err := repos.DB.Create(admin).Error; err != nil {
		return err
	}

	// Create 100 fake employees
	employees := make([]*models.User, 100)
	for i := 0; i < 100; i++ {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("employee%d", i+1)), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		salary := float64(rand.Intn(5000000-3000000) + 3000000) // Random salary between 3M - 8M

		employees[i] = &models.User{
			Username: fmt.Sprintf("employee%d", i+1),
			Password: string(hashedPassword),
			Role:     "employee",
			Salary:   &salary,
			IsActive: true,
		}
	}

	if err := repos.DB.CreateInBatches(employees, 50).Error; err != nil {
		return err
	}

	// Create a sample attendance period
	startDate := time.Now().AddDate(0, -1, 0) // Last month
	endDate := time.Now().AddDate(0, 0, -1)   // Yesterday

	period := &models.AttendancePeriod{
		StartDate:   startDate,
		EndDate:     endDate,
		IsProcessed: false,
	}

	if err := repos.DB.Create(period).Error; err != nil {
		return err
	}

	log.Println("Database seeded successfully")
	log.Printf("Admin credentials: username=admin, password=admin123")
	log.Printf("Employee credentials: username=employee1-100, password=employee1-100")

	return nil
}
