package repository

import (
	"payslip-system/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repositories struct {
	DB               *gorm.DB
	User             IUserRepository
	AttendancePeriod IAttendancePeriodRepository
	Attendance       IAttendanceRepository
	Overtime         IOvertimeRepository
	Reimbursement    IReimbursementRepository
	Payroll          IPayrollRepository
	AuditLog         IAuditLogRepository
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

//go:generate mockgen -destination=mocks/mocks.go -source=init.go IUserRepository, IAttendancePeriodRepository, IAttendanceRepository, IOvertimeRepository, IPayrollRepository, IReimbursementRepository, IAuditLogRepository
type IUserRepository interface {
	GetByID(id uuid.UUID) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAllEmployees() ([]models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
}

type IAttendancePeriodRepository interface {
	GetByID(id uuid.UUID) (*models.AttendancePeriod, error)
	GetAll() ([]models.AttendancePeriod, error)
	GetActive() (*models.AttendancePeriod, error)
	Create(period *models.AttendancePeriod) error
	Update(period *models.AttendancePeriod) error
}

type IAttendanceRepository interface {
	GetByUserAndPeriod(userID, periodID uuid.UUID) ([]models.Attendance, error)
	GetByUserAndDate(userID uuid.UUID, date time.Time) (*models.Attendance, error)
	Create(attendance *models.Attendance) error
	CountWorkingDaysInPeriod(startDate, endDate time.Time) int
}

type IOvertimeRepository interface {
	GetByUserAndPeriod(userID, periodID uuid.UUID) ([]models.Overtime, error)
	GetByUserAndDate(userID uuid.UUID, date time.Time) (*models.Overtime, error)
	Create(overtime *models.Overtime) error
}

type IReimbursementRepository interface {
	GetByUserAndPeriod(userID, periodID uuid.UUID) ([]models.Reimbursement, error)
	Create(reimbursement *models.Reimbursement) error
}

type IPayrollRepository interface {
	GetByPeriodID(periodID uuid.UUID) (*models.Payroll, error)
	GetPayrollItemsByPeriodAndUser(periodID, userID uuid.UUID) (*models.PayrollItem, error)
	GetAllPayrollItemsByPeriod(periodID uuid.UUID) ([]models.PayrollItem, error)
	Create(payroll *models.Payroll) error
	CreatePayrollItem(item *models.PayrollItem) error
}

type IAuditLogRepository interface {
	Create(log *models.AuditLog) error
	GetByTableAndRecord(tableName string, recordID uuid.UUID) ([]models.AuditLog, error)
}
