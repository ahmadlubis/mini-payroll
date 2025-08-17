package service_test

import (
	"testing"
	"time"

	"payslip-system/internal/models"
	"payslip-system/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPayrollService_GeneratePayslip(t *testing.T) {
	db, cleanup := test.SetupTestDB()
	defer cleanup()

	repos, services := test.SetupTestServices(db)

	// Create test user
	testUser := &models.User{
		Username: "testuser",
		Role:     "employee",
		IsActive: true,
	}
	salary := 6000000.0
	testUser.Salary = &salary

	err := repos.User.Create(testUser)
	require.NoError(t, err)

	// Create test attendance period
	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now().AddDate(0, 0, -1)
	period := &models.AttendancePeriod{
		StartDate:   startDate,
		EndDate:     endDate,
		IsProcessed: false,
	}
	err = repos.AttendancePeriod.Create(period)
	require.NoError(t, err)

	// Add attendance records
	for i := 0; i < 20; i++ {
		date := startDate.AddDate(0, 0, i)
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			attendance := &models.Attendance{
				UserID:             testUser.ID,
				AttendancePeriodID: period.ID,
				Date:               date,
				CheckInTime:        time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location()),
			}
			repos.Attendance.Create(attendance)
		}
	}

	// Add overtime record
	overtime := &models.Overtime{
		UserID:             testUser.ID,
		AttendancePeriodID: period.ID,
		Date:               startDate.AddDate(0, 0, 1),
		Hours:              2.0,
	}
	repos.Overtime.Create(overtime)

	// Add reimbursement
	reimbursement := &models.Reimbursement{
		UserID:             testUser.ID,
		AttendancePeriodID: period.ID,
		Amount:             100000.0,
		Description:        "Transportation",
	}
	repos.Reimbursement.Create(reimbursement)

	// Test payslip generation
	payslip, err := services.Payroll.GeneratePayslip(testUser.ID, period.ID)
	require.NoError(t, err)
	require.NotNil(t, payslip)

	// Verify calculations
	assert.Equal(t, salary, payslip.BaseSalary)
	assert.Greater(t, payslip.AttendanceDays, 0)
	assert.Greater(t, payslip.AttendanceAmount, float64(0))
	assert.Equal(t, float64(2), payslip.OvertimeHours)
	assert.Greater(t, payslip.OvertimeAmount, float64(0))
	assert.Equal(t, float64(100000), payslip.ReimbursementAmount)

	// Total should be sum of all components
	expectedTotal := payslip.AttendanceAmount + payslip.OvertimeAmount + payslip.ReimbursementAmount
	assert.Equal(t, expectedTotal, payslip.TotalAmount)
}
