package repository

import (
	"payslip-system/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repositories struct {
	DB               *gorm.DB
	User             UserRepository
	AttendancePeriod AttendancePeriodRepository
	Attendance       AttendanceRepository
	Overtime         OvertimeRepository
	Reimbursement    ReimbursementRepository
	Payroll          PayrollRepository
	AuditLog         AuditLogRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		DB:               db,
		User:             NewUserRepository(db),
		AttendancePeriod: NewAttendancePeriodRepository(db),
		Attendance:       NewAttendanceRepository(db),
		Overtime:         NewOvertimeRepository(db),
		Reimbursement:    NewReimbursementRepository(db),
		Payroll:          NewPayrollRepository(db),
		AuditLog:         NewAuditLogRepository(db),
	}
}

// UserRepository
type UserRepository interface {
	GetByID(id uuid.UUID) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAllEmployees() ([]models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ? AND is_active = true", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ? AND is_active = true", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllEmployees() ([]models.User, error) {
	var employees []models.User
	if err := r.db.Where("role = ? AND is_active = true", "employee").Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// AttendancePeriodRepository
type AttendancePeriodRepository interface {
	GetByID(id uuid.UUID) (*models.AttendancePeriod, error)
	GetAll() ([]models.AttendancePeriod, error)
	GetActive() (*models.AttendancePeriod, error)
	Create(period *models.AttendancePeriod) error
	Update(period *models.AttendancePeriod) error
}

type attendancePeriodRepository struct {
	db *gorm.DB
}

func NewAttendancePeriodRepository(db *gorm.DB) AttendancePeriodRepository {
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

// AttendanceRepository
type AttendanceRepository interface {
	GetByUserAndPeriod(userID, periodID uuid.UUID) ([]models.Attendance, error)
	GetByUserAndDate(userID uuid.UUID, date time.Time) (*models.Attendance, error)
	Create(attendance *models.Attendance) error
	CountWorkingDaysInPeriod(startDate, endDate time.Time) int
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
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

// OvertimeRepository
type OvertimeRepository interface {
	GetByUserAndPeriod(userID, periodID uuid.UUID) ([]models.Overtime, error)
	GetByUserAndDate(userID uuid.UUID, date time.Time) (*models.Overtime, error)
	Create(overtime *models.Overtime) error
}

type overtimeRepository struct {
	db *gorm.DB
}

func NewOvertimeRepository(db *gorm.DB) OvertimeRepository {
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

// ReimbursementRepository
type ReimbursementRepository interface {
	GetByUserAndPeriod(userID, periodID uuid.UUID) ([]models.Reimbursement, error)
	Create(reimbursement *models.Reimbursement) error
}

type reimbursementRepository struct {
	db *gorm.DB
}

func NewReimbursementRepository(db *gorm.DB) ReimbursementRepository {
	return &reimbursementRepository{db: db}
}

func (r *reimbursementRepository) GetByUserAndPeriod(userID, periodID uuid.UUID) ([]models.Reimbursement, error) {
	var reimbursements []models.Reimbursement
	if err := r.db.Where("user_id = ? AND attendance_period_id = ?", userID, periodID).Find(&reimbursements).Error; err != nil {
		return nil, err
	}
	return reimbursements, nil
}

func (r *reimbursementRepository) Create(reimbursement *models.Reimbursement) error {
	return r.db.Create(reimbursement).Error
}

// PayrollRepository
type PayrollRepository interface {
	GetByPeriodID(periodID uuid.UUID) (*models.Payroll, error)
	GetPayrollItemsByPeriodAndUser(periodID, userID uuid.UUID) (*models.PayrollItem, error)
	GetAllPayrollItemsByPeriod(periodID uuid.UUID) ([]models.PayrollItem, error)
	Create(payroll *models.Payroll) error
	CreatePayrollItem(item *models.PayrollItem) error
}

type payrollRepository struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) PayrollRepository {
	return &payrollRepository{db: db}
}

func (r *payrollRepository) GetByPeriodID(periodID uuid.UUID) (*models.Payroll, error) {
	var payroll models.Payroll
	if err := r.db.Where("attendance_period_id = ?", periodID).Preload("PayrollItems.User").First(&payroll).Error; err != nil {
		return nil, err
	}
	return &payroll, nil
}

func (r *payrollRepository) GetPayrollItemsByPeriodAndUser(periodID, userID uuid.UUID) (*models.PayrollItem, error) {
	var item models.PayrollItem
	if err := r.db.Joins("JOIN payrolls ON payroll_items.payroll_id = payrolls.id").
		Where("payrolls.attendance_period_id = ? AND payroll_items.user_id = ?", periodID, userID).
		Preload("User").First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *payrollRepository) GetAllPayrollItemsByPeriod(periodID uuid.UUID) ([]models.PayrollItem, error) {
	var items []models.PayrollItem
	if err := r.db.Joins("JOIN payrolls ON payroll_items.payroll_id = payrolls.id").
		Where("payrolls.attendance_period_id = ?", periodID).
		Preload("User").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *payrollRepository) Create(payroll *models.Payroll) error {
	return r.db.Create(payroll).Error
}

func (r *payrollRepository) CreatePayrollItem(item *models.PayrollItem) error {
	return r.db.Create(item).Error
}

// AuditLogRepository
type AuditLogRepository interface {
	Create(log *models.AuditLog) error
	GetByTableAndRecord(tableName string, recordID uuid.UUID) ([]models.AuditLog, error)
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditLogRepository) GetByTableAndRecord(tableName string, recordID uuid.UUID) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	if err := r.db.Where("table_name = ? AND record_id = ?", tableName, recordID).
		Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
