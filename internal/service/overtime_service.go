package service

import (
	"errors"
	"fmt"
	"payslip-system/internal/models"
	"payslip-system/internal/repository"
	"time"

	"github.com/google/uuid"
)

// OvertimeService
type OvertimeService interface {
	SubmitOvertime(userID uuid.UUID, date time.Time, hours float64, ipAddress, requestID string) error
}

type overtimeService struct {
	repos *repository.Repositories
}

func NewOvertimeService(repos *repository.Repositories) OvertimeService {
	return &overtimeService{repos: repos}
}

func (s *overtimeService) SubmitOvertime(userID uuid.UUID, date time.Time, hours float64, ipAddress, requestID string) error {
	// Validate hours (max 3 hours per day)
	if hours <= 0 || hours > 3 {
		return errors.New("overtime hours must be between 0 and 3")
	}

	// Check if overtime already exists for this date
	if _, err := s.repos.Overtime.GetByUserAndDate(userID, date); err == nil {
		return errors.New("overtime already submitted for this date")
	}

	// Get active attendance period
	period, err := s.repos.AttendancePeriod.GetActive()
	if err != nil {
		return errors.New("no active attendance period found")
	}

	// Create overtime record
	overtime := &models.Overtime{
		BaseModel: models.BaseModel{
			CreatedBy: &userID,
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:             userID,
		AttendancePeriodID: period.ID,
		Date:               date,
		Hours:              hours,
	}

	if err := s.repos.Overtime.Create(overtime); err != nil {
		return fmt.Errorf("failed to create overtime record: %w", err)
	}

	// Create audit log
	createAuditLog("overtimes", overtime.ID, "INSERT", nil, overtime, &userID, ipAddress, requestID, s.repos)

	return nil
}
