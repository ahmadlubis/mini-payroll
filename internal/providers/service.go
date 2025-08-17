package providers

import (
	"payslip-system/internal/domains"
	"payslip-system/internal/repository"
	"payslip-system/internal/service"
)

type Services struct {
	Auth          domains.IAuthService
	Attendance    domains.IAttendanceService
	Overtime      domains.IOvertimeService
	Reimbursement domains.IReimbursementService
	Payroll       domains.IPayrollService
	Admin         domains.IAdminService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Auth:          service.NewAuthService(repos),
		Attendance:    service.NewAttendanceService(repos),
		Overtime:      service.NewOvertimeService(repos),
		Reimbursement: service.NewReimbursementService(repos),
		Payroll:       service.NewPayrollService(repos),
		Admin:         service.NewAdminService(repos),
	}
}
