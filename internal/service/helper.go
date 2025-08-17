package service

import (
	"encoding/json"
	"time"

	"payslip-system/internal/models"
	"payslip-system/internal/repository"

	"github.com/google/uuid"
)

// Helper methods for other services
func createAuditLog(tableName string, recordID uuid.UUID, action string, oldValues, newValues interface{}, userID *uuid.UUID, ipAddress, requestID string, repos *repository.Repositories) {
	var oldJSON, newJSON string

	if oldValues != nil {
		if data, err := json.Marshal(oldValues); err == nil {
			oldJSON = string(data)
		}
	}

	if newValues != nil {
		if data, err := json.Marshal(newValues); err == nil {
			newJSON = string(data)
		}
	}

	log := &models.AuditLog{
		ID:        uuid.New(),
		TableName: tableName,
		RecordID:  recordID,
		Action:    action,
		OldValues: oldJSON,
		NewValues: newJSON,
		UserID:    userID,
		IPAddress: ipAddress,
		RequestID: requestID,
		CreatedAt: time.Now(),
	}

	repos.AuditLog.Create(log)
}
