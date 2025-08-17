package repository

import (
	"payslip-system/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) IAuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditLogRepository) GetByTableAndRecord(tableName string, recordID uuid.UUID) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	if err := r.db.Where("table_name = ? AND record_id = ?", tableName, recordID).
		Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
