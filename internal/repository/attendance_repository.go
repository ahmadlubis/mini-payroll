package repository

import (
	"payslip-system/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type attendancePeriodRepository struct {
	db *gorm.DB
}

func NewAttendancePeriodRepository(db *gorm.DB) IAttendancePeriodRepository {
	return &attendancePeriodRepository{db: db}
}

func (r *attendancePeriodRepository) GetByID(id uuid.UUID) (*models.AttendancePeriod, error) {
	var period models.AttendancePeriod
	if err := r.db.First(&period, id).Error; err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *attendancePeriodRepository) GetAll() ([]models.AttendancePeriod, error) {
	var periods []models.AttendancePeriod
	if err := r.db.Order("start_date DESC").Find(&periods).Error; err != nil {
		return nil, err
	}
	return periods, nil
}

func (r *attendancePeriodRepository) GetActive() (*models.AttendancePeriod, error) {
	var period models.AttendancePeriod
	now := time.Now()
	if err := r.db.Where("start_date <= ? AND end_date >= ? AND is_processed = false", now, now).First(&period).Error; err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *attendancePeriodRepository) Create(period *models.AttendancePeriod) error {
	return r.db.Create(period).Error
}

func (r *attendancePeriodRepository) Update(period *models.AttendancePeriod) error {
	return r.db.Save(period).Error
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) IAttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) GetByUserAndPeriod(userID, periodID uuid.UUID) ([]models.Attendance, error) {
	var attendances []models.Attendance
	if err := r.db.Where("user_id = ? AND attendance_period_id = ?", userID, periodID).Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *attendanceRepository) GetByUserAndDate(userID uuid.UUID, date time.Time) (*models.Attendance, error) {
	var attendance models.Attendance
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nextDay := dateOnly.Add(24 * time.Hour)

	if err := r.db.Where("user_id = ? AND date >= ? AND date < ?", userID, dateOnly, nextDay).First(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepository) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *attendanceRepository) CountWorkingDaysInPeriod(startDate, endDate time.Time) int {
	count := 0
	for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.Add(24 * time.Hour) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			count++
		}
	}
	return count
}
