package domains

import (
	"payslip-system/internal/models"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=mocks/mocks.go -source=service.go IAdminService, IAttendanceService, IAuthService, IOvertimeService, IPayrollService, IReimbursementService
type IAdminService interface {
	CreateAttendancePeriod(startDate, endDate time.Time, adminID uuid.UUID, ipAddress, requestID string) (*models.AttendancePeriod, error)
}

type IAttendanceService interface {
	SubmitAttendance(userID uuid.UUID, date time.Time, checkInTime time.Time, ipAddress, requestID string) error
}

type IAuthService interface {
	Login(username, password string) (*models.User, string, error)
	ValidateToken(tokenString string) (*models.User, error)
}

type IOvertimeService interface {
	SubmitOvertime(userID uuid.UUID, date time.Time, hours float64, ipAddress, requestID string) error
}

type IPayrollService interface {
	GeneratePayslip(userID, periodID uuid.UUID) (*PayslipResponse, error)
	GeneratePayrollSummary(periodID uuid.UUID) (*PayrollSummaryResponse, error)
	ProcessPayroll(periodID, adminID uuid.UUID, ipAddress, requestID string) error
}

type IReimbursementService interface {
	SubmitReimbursement(userID uuid.UUID, amount float64, description, ipAddress, requestID string) error
}
