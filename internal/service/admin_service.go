package service

import (
	"errors"
	"fmt"
	"payslip-system/internal/models"
	"payslip-system/internal/repository"
	"time"

	"github.com/google/uuid"
)

type adminService struct {
	repos *repository.Repositories
}

func NewAdminService(repos *repository.Repositories) *adminService {
	return &adminService{repos: repos}
}

func (s *adminService) CreateAttendancePeriod(startDate, endDate time.Time, adminID uuid.UUID, ipAddress, requestID string) (*models.AttendancePeriod, error) {
	if endDate.Before(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	period := &models.AttendancePeriod{
		BaseModel: models.BaseModel{
			CreatedBy: &adminID,
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		StartDate:   startDate,
		EndDate:     endDate,
		IsProcessed: false,
	}

	if err := s.repos.AttendancePeriod.Create(period); err != nil {
		return nil, fmt.Errorf("failed to create attendance period: %w", err)
	}

	// Create audit log
	createAuditLog("attendance_periods", period.ID, "INSERT", nil, period, &adminID, ipAddress, requestID, s.repos)

	return period, nil
}
