package service

import (
	"errors"
	"fmt"
	"payslip-system/internal/domains"
	"payslip-system/internal/models"
	"payslip-system/internal/repository"
	"time"

	"github.com/google/uuid"
)

type payrollService struct {
	repos *repository.Repositories
}

func NewPayrollService(repos *repository.Repositories) *payrollService {
	return &payrollService{repos: repos}
}

func (s *payrollService) GeneratePayslip(userID, periodID uuid.UUID) (*domains.PayslipResponse, error) {
	// Get user
	user, err := s.repos.User.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.Role != "employee" || user.Salary == nil {
		return nil, errors.New("invalid employee or salary not set")
	}

	// Get period
	period, err := s.repos.AttendancePeriod.GetByID(periodID)
	if err != nil {
		return nil, fmt.Errorf("period not found: %w", err)
	}

	// If payroll is processed, get from payroll item
	if period.IsProcessed {
		item, err := s.repos.Payroll.GetPayrollItemsByPeriodAndUser(periodID, userID)
		if err != nil {
			return nil, fmt.Errorf("payroll item not found: %w", err)
		}

		// Get reimbursements
		reimbursements, _ := s.repos.Reimbursement.GetByUserAndPeriod(userID, periodID)

		return &domains.PayslipResponse{
			Employee:            user,
			Period:              period,
			BaseSalary:          item.BaseSalary,
			AttendanceDays:      item.AttendanceDays,
			WorkingDays:         item.WorkingDays,
			AttendanceAmount:    item.AttendanceAmount,
			OvertimeHours:       item.OvertimeHours,
			OvertimeAmount:      item.OvertimeAmount,
			Reimbursements:      reimbursements,
			ReimbursementAmount: item.ReimbursementAmount,
			TotalAmount:         item.TotalAmount,
		}, nil
	}

	// Calculate live payslip
	return s.calculatePayslip(user, period)
}

func (s *payrollService) calculatePayslip(user *models.User, period *models.AttendancePeriod) (*domains.PayslipResponse, error) {
	// Get attendance records
	attendances, _ := s.repos.Attendance.GetByUserAndPeriod(user.ID, period.ID)
	attendanceDays := len(attendances)

	// Calculate working days in period
	workingDays := s.repos.Attendance.CountWorkingDaysInPeriod(period.StartDate, period.EndDate)

	// Calculate attendance amount (prorated)
	baseSalary := *user.Salary
	dailySalary := baseSalary / 30 // Assuming 30 days per month
	attendanceAmount := dailySalary * float64(attendanceDays)

	// Get overtime records
	overtimes, _ := s.repos.Overtime.GetByUserAndPeriod(user.ID, period.ID)
	var overtimeHours float64
	for _, ot := range overtimes {
		overtimeHours += ot.Hours
	}

	// Calculate overtime amount (2x daily salary per hour)
	hourlyRate := dailySalary / 8 // 8 working hours per day
	overtimeAmount := overtimeHours * hourlyRate * 2

	// Get reimbursements
	reimbursements, _ := s.repos.Reimbursement.GetByUserAndPeriod(user.ID, period.ID)
	var reimbursementAmount float64
	for _, r := range reimbursements {
		reimbursementAmount += r.Amount
	}

	// Calculate total
	totalAmount := attendanceAmount + overtimeAmount + reimbursementAmount

	return &domains.PayslipResponse{
		Employee:            user,
		Period:              period,
		BaseSalary:          baseSalary,
		AttendanceDays:      attendanceDays,
		WorkingDays:         workingDays,
		AttendanceAmount:    attendanceAmount,
		OvertimeHours:       overtimeHours,
		OvertimeAmount:      overtimeAmount,
		Reimbursements:      reimbursements,
		ReimbursementAmount: reimbursementAmount,
		TotalAmount:         totalAmount,
	}, nil
}

func (s *payrollService) GeneratePayrollSummary(periodID uuid.UUID) (*domains.PayrollSummaryResponse, error) {
	// Get period
	period, err := s.repos.AttendancePeriod.GetByID(periodID)
	if err != nil {
		return nil, fmt.Errorf("period not found: %w", err)
	}

	// Get all employees
	employees, err := s.repos.User.GetAllEmployees()
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}

	var employeeSummaries []domains.EmployeeSummary
	var totalAmount float64

	for _, employee := range employees {
		if period.IsProcessed {
			// Get from payroll items
			item, err := s.repos.Payroll.GetPayrollItemsByPeriodAndUser(periodID, employee.ID)
			if err != nil {
				continue // Skip if no payroll item found
			}
			employeeSummaries = append(employeeSummaries, domains.EmployeeSummary{
				Employee:    &employee,
				TotalAmount: item.TotalAmount,
			})
			totalAmount += item.TotalAmount
		} else {
			// Calculate live
			payslip, err := s.calculatePayslip(&employee, period)
			if err != nil {
				continue
			}
			employeeSummaries = append(employeeSummaries, domains.EmployeeSummary{
				Employee:    &employee,
				TotalAmount: payslip.TotalAmount,
			})
			totalAmount += payslip.TotalAmount
		}
	}

	return &domains.PayrollSummaryResponse{
		Period:      period,
		Employees:   employeeSummaries,
		TotalAmount: totalAmount,
	}, nil
}

func (s *payrollService) ProcessPayroll(periodID, adminID uuid.UUID, ipAddress, requestID string) error {
	// Get period
	period, err := s.repos.AttendancePeriod.GetByID(periodID)
	if err != nil {
		return fmt.Errorf("period not found: %w", err)
	}

	if period.IsProcessed {
		return errors.New("payroll already processed for this period")
	}

	// Get all employees
	employees, err := s.repos.User.GetAllEmployees()
	if err != nil {
		return fmt.Errorf("failed to get employees: %w", err)
	}

	// Start transaction
	tx := s.repos.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create payroll record
	now := time.Now()
	payroll := &models.Payroll{
		BaseModel: models.BaseModel{
			CreatedBy: &adminID,
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		AttendancePeriodID: periodID,
		ProcessedBy:        adminID,
	}

	if err := tx.Create(payroll).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create payroll: %w", err)
	}

	var totalAmount float64

	// Process each employee
	for _, employee := range employees {
		if employee.Salary == nil {
			continue
		}

		payslip, err := s.calculatePayslip(&employee, period)
		if err != nil {
			continue
		}

		// Create payroll item
		item := &models.PayrollItem{
			BaseModel: models.BaseModel{
				CreatedBy: &adminID,
				IPAddress: ipAddress,
				RequestID: requestID,
			},
			PayrollID:           payroll.ID,
			UserID:              employee.ID,
			BaseSalary:          payslip.BaseSalary,
			AttendanceDays:      payslip.AttendanceDays,
			WorkingDays:         payslip.WorkingDays,
			AttendanceAmount:    payslip.AttendanceAmount,
			OvertimeHours:       payslip.OvertimeHours,
			OvertimeAmount:      payslip.OvertimeAmount,
			ReimbursementAmount: payslip.ReimbursementAmount,
			TotalAmount:         payslip.TotalAmount,
		}

		if err := tx.Create(item).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create payroll item: %w", err)
		}

		totalAmount += payslip.TotalAmount
	}

	// Update payroll total
	payroll.TotalAmount = totalAmount
	if err := tx.Save(payroll).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update payroll total: %w", err)
	}

	// Mark period as processed
	period.IsProcessed = true
	period.ProcessedAt = &now
	period.UpdatedBy = &adminID
	period.IPAddress = ipAddress
	period.RequestID = requestID

	if err := tx.Save(period).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update period: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Create audit logs
	createAuditLog("payrolls", payroll.ID, "INSERT", nil, payroll, &adminID, ipAddress, requestID, s.repos)
	createAuditLog("attendance_periods", period.ID, "UPDATE", nil, period, &adminID, ipAddress, requestID, s.repos)

	return nil
}
