package service

import (
	"errors"
	"fmt"
	"payslip-system/internal/models"
	"payslip-system/internal/repository"
	"time"

	"github.com/google/uuid"
)

type AttendanceService interface {
	SubmitAttendance(userID uuid.UUID, date time.Time, checkInTime time.Time, ipAddress, requestID string) error
}

type attendanceService struct {
	repos *repository.Repositories
}

func NewAttendanceService(repos *repository.Repositories) AttendanceService {
	return &attendanceService{repos: repos}
}

func (s *attendanceService) SubmitAttendance(userID uuid.UUID, date time.Time, checkInTime time.Time, ipAddress, requestID string) error {
	// Check if it's weekend
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return errors.New("cannot submit attendance on weekends")
	}

	// Check if attendance already exists for this date
	if _, err := s.repos.Attendance.GetByUserAndDate(userID, date); err == nil {
		return errors.New("attendance already submitted for this date")
	}

	// Get active attendance period
	period, err := s.repos.AttendancePeriod.GetActive()
	if err != nil {
		return errors.New("no active attendance period found")
	}

	// Check if date is within period
	if date.Before(period.StartDate) || date.After(period.EndDate) {
		return errors.New("date is not within active attendance period")
	}

	// Create attendance record
	attendance := &models.Attendance{
		BaseModel: models.BaseModel{
			CreatedBy: &userID,
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:             userID,
		AttendancePeriodID: period.ID,
		Date:               date,
		CheckInTime:        checkInTime,
	}

	if err := s.repos.Attendance.Create(attendance); err != nil {
		return fmt.Errorf("failed to create attendance record: %w", err)
	}

	// Create audit log
	createAuditLog("attendances", attendance.ID, "INSERT", nil, attendance, &userID, ipAddress, requestID, s.repos)

	return nil
}
