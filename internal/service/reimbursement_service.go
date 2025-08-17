package service

import (
	"errors"
	"fmt"
	"payslip-system/internal/models"
	"payslip-system/internal/repository"

	"github.com/google/uuid"
)

type ReimbursementService interface {
	SubmitReimbursement(userID uuid.UUID, amount float64, description, ipAddress, requestID string) error
}

type reimbursementService struct {
	repos *repository.Repositories
}

func NewReimbursementService(repos *repository.Repositories) ReimbursementService {
	return &reimbursementService{repos: repos}
}

func (s *reimbursementService) SubmitReimbursement(userID uuid.UUID, amount float64, description, ipAddress, requestID string) error {
	if amount <= 0 {
		return errors.New("reimbursement amount must be greater than 0")
	}

	if description == "" {
		return errors.New("reimbursement description is required")
	}

	// Get active attendance period
	period, err := s.repos.AttendancePeriod.GetActive()
	if err != nil {
		return errors.New("no active attendance period found")
	}

	// Create reimbursement record
	reimbursement := &models.Reimbursement{
		BaseModel: models.BaseModel{
			CreatedBy: &userID,
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:             userID,
		AttendancePeriodID: period.ID,
		Amount:             amount,
		Description:        description,
	}

	if err := s.repos.Reimbursement.Create(reimbursement); err != nil {
		return fmt.Errorf("failed to create reimbursement record: %w", err)
	}

	// Create audit log
	createAuditLog("reimbursements", reimbursement.ID, "INSERT", nil, reimbursement, &userID, ipAddress, requestID, s.repos)

	return nil
}
