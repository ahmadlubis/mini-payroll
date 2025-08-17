package repository

import (
	"payslip-system/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type payrollRepository struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) IPayrollRepository {
	return &payrollRepository{db: db}
}

func (r *payrollRepository) GetByPeriodID(periodID uuid.UUID) (*models.Payroll, error) {
	var payroll models.Payroll
	if err := r.db.Where("attendance_period_id = ?", periodID).Preload("PayrollItems.User").First(&payroll).Error; err != nil {
		return nil, err
	}
	return &payroll, nil
}

func (r *payrollRepository) GetPayrollItemsByPeriodAndUser(periodID, userID uuid.UUID) (*models.PayrollItem, error) {
	var item models.PayrollItem
	if err := r.db.Joins("JOIN payrolls ON payroll_items.payroll_id = payrolls.id").
		Where("payrolls.attendance_period_id = ? AND payroll_items.user_id = ?", periodID, userID).
		Preload("User").First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *payrollRepository) GetAllPayrollItemsByPeriod(periodID uuid.UUID) ([]models.PayrollItem, error) {
	var items []models.PayrollItem
	if err := r.db.Joins("JOIN payrolls ON payroll_items.payroll_id = payrolls.id").
		Where("payrolls.attendance_period_id = ?", periodID).
		Preload("User").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *payrollRepository) Create(payroll *models.Payroll) error {
	return r.db.Create(payroll).Error
}

func (r *payrollRepository) CreatePayrollItem(item *models.PayrollItem) error {
	return r.db.Create(item).Error
}
