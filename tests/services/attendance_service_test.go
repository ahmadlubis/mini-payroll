package service_test

import (
	"testing"
	"time"

	"payslip-system/internal/models"
	"payslip-system/tests"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerformanceSubmitAttendanceService(t *testing.T) {
	db, cleanup := tests.SetupTestDB()
	defer cleanup()

	repos, services := tests.SetupTestServices(db)

	// Create test user
	testUser := &models.User{
		Username: "testuser",
		Role:     "employee",
		IsActive: true,
	}
	salary := 5000000.0
	testUser.Salary = &salary

	err := repos.User.Create(testUser)
	require.NoError(t, err)

	// Create test attendance period
	startDate := time.Now().AddDate(0, 0, -10)
	endDate := time.Now().AddDate(0, 0, 10)
	period := &models.AttendancePeriod{
		StartDate:   startDate,
		EndDate:     endDate,
		IsProcessed: false,
	}
	err = repos.AttendancePeriod.Create(period)
	require.NoError(t, err)

	tests := []struct {
		name        string
		date        time.Time
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid weekday attendance",
			date:        getNextWeekday(time.Now()),
			expectError: false,
		},
		{
			name:        "weekend attendance",
			date:        getNextSaturday(time.Now()),
			expectError: true,
			errorMsg:    "cannot submit attendance on weekends",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkInTime := time.Date(tt.date.Year(), tt.date.Month(), tt.date.Day(), 9, 0, 0, 0, tt.date.Location())

			err := services.Attendance.SubmitAttendance(
				testUser.ID,
				tt.date,
				checkInTime,
				"127.0.0.1",
				uuid.New().String(),
			)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func getNextWeekday(t time.Time) time.Time {
	for t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

func getNextSaturday(t time.Time) time.Time {
	for t.Weekday() != time.Saturday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}
