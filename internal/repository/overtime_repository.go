package repository

import (
	"payslip-system/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type overtimeRepository struct {
	db *gorm.DB
}

func NewOvertimeRepository(db *gorm.DB) IOvertimeRepository {
	return &overtimeRepository{db: db}
}

func (r *overtimeRepository) GetByUserAndPeriod(userID, periodID uuid.UUID) ([]models.Overtime, error) {
	var overtimes []models.Overtime
	if err := r.db.Where("user_id = ? AND attendance_period_id = ?", userID, periodID).Find(&overtimes).Error; err != nil {
		return nil, err
	}
	return overtimes, nil
}

func (r *overtimeRepository) GetByUserAndDate(userID uuid.UUID, date time.Time) (*models.Overtime, error) {
	var overtime models.Overtime
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nextDay := dateOnly.Add(24 * time.Hour)

	if err := r.db.Where("user_id = ? AND date >= ? AND date < ?", userID, dateOnly, nextDay).First(&overtime).Error; err != nil {
		return nil, err
	}
	return &overtime, nil
}

func (r *overtimeRepository) Create(overtime *models.Overtime) error {
	return r.db.Create(overtime).Error
}
