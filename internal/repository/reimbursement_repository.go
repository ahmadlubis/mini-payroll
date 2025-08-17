package repository

import (
	"payslip-system/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reimbursementRepository struct {
	db *gorm.DB
}

func NewReimbursementRepository(db *gorm.DB) IReimbursementRepository {
	return &reimbursementRepository{db: db}
}

func (r *reimbursementRepository) GetByUserAndPeriod(userID, periodID uuid.UUID) ([]models.Reimbursement, error) {
	var reimbursements []models.Reimbursement
	if err := r.db.Where("user_id = ? AND attendance_period_id = ?", userID, periodID).Find(&reimbursements).Error; err != nil {
		return nil, err
	}
	return reimbursements, nil
}

func (r *reimbursementRepository) Create(reimbursement *models.Reimbursement) error {
	return r.db.Create(reimbursement).Error
}
