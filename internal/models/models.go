package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty" gorm:"type:uuid"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty" gorm:"type:uuid"`
	IPAddress string     `json:"ip_address,omitempty"`
	RequestID string     `json:"request_id,omitempty"`
}

// User represents both employees and admins
type User struct {
	BaseModel
	Username string   `json:"username" gorm:"unique;not null"`
	Password string   `json:"-" gorm:"not null"`
	Role     string   `json:"role" gorm:"not null;default:'employee'"` // 'admin' or 'employee'
	Salary   *float64 `json:"salary,omitempty"`                        // Only for employees
	IsActive bool     `json:"is_active" gorm:"default:true"`
}

// AttendancePeriod represents payroll periods set by admin
type AttendancePeriod struct {
	BaseModel
	StartDate   time.Time  `json:"start_date" gorm:"not null"`
	EndDate     time.Time  `json:"end_date" gorm:"not null"`
	IsProcessed bool       `json:"is_processed" gorm:"default:false"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
}

// Attendance represents employee attendance records
type Attendance struct {
	BaseModel
	UserID             uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	AttendancePeriodID uuid.UUID  `json:"attendance_period_id" gorm:"type:uuid;not null"`
	Date               time.Time  `json:"date" gorm:"not null"`
	CheckInTime        time.Time  `json:"check_in_time" gorm:"not null"`
	CheckOutTime       *time.Time `json:"check_out_time,omitempty"`

	// Relationships
	User             User             `json:"user,omitempty"`
	AttendancePeriod AttendancePeriod `json:"attendance_period,omitempty"`
}

// Overtime represents employee overtime records
type Overtime struct {
	BaseModel
	UserID             uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	AttendancePeriodID uuid.UUID `json:"attendance_period_id" gorm:"type:uuid;not null"`
	Date               time.Time `json:"date" gorm:"not null"`
	Hours              float64   `json:"hours" gorm:"not null"`

	// Relationships
	User             User             `json:"user,omitempty"`
	AttendancePeriod AttendancePeriod `json:"attendance_period,omitempty"`
}

// Reimbursement represents employee reimbursement requests
type Reimbursement struct {
	BaseModel
	UserID             uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	AttendancePeriodID uuid.UUID `json:"attendance_period_id" gorm:"type:uuid;not null"`
	Amount             float64   `json:"amount" gorm:"not null"`
	Description        string    `json:"description" gorm:"not null"`

	// Relationships
	User             User             `json:"user,omitempty"`
	AttendancePeriod AttendancePeriod `json:"attendance_period,omitempty"`
}

// Payroll represents processed payroll for a period
type Payroll struct {
	BaseModel
	AttendancePeriodID uuid.UUID `json:"attendance_period_id" gorm:"type:uuid;not null"`
	TotalAmount        float64   `json:"total_amount" gorm:"not null"`
	ProcessedBy        uuid.UUID `json:"processed_by" gorm:"type:uuid;not null"`

	// Relationships
	AttendancePeriod AttendancePeriod `json:"attendance_period,omitempty"`
	ProcessedByUser  User             `json:"processed_by_user,omitempty" gorm:"foreignKey:ProcessedBy"`
	PayrollItems     []PayrollItem    `json:"payroll_items,omitempty"`
}

// PayrollItem represents individual employee payroll calculation
type PayrollItem struct {
	BaseModel
	PayrollID           uuid.UUID `json:"payroll_id" gorm:"type:uuid;not null"`
	UserID              uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	BaseSalary          float64   `json:"base_salary" gorm:"not null"`
	AttendanceDays      int       `json:"attendance_days" gorm:"not null"`
	WorkingDays         int       `json:"working_days" gorm:"not null"`
	AttendanceAmount    float64   `json:"attendance_amount" gorm:"not null"`
	OvertimeHours       float64   `json:"overtime_hours" gorm:"not null"`
	OvertimeAmount      float64   `json:"overtime_amount" gorm:"not null"`
	ReimbursementAmount float64   `json:"reimbursement_amount" gorm:"not null"`
	TotalAmount         float64   `json:"total_amount" gorm:"not null"`

	// Relationships
	Payroll User `json:"payroll,omitempty"`
	User    User `json:"user,omitempty"`
}

// AuditLog represents audit trail for significant changes
type AuditLog struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TableName string     `json:"table_name" gorm:"not null"`
	RecordID  uuid.UUID  `json:"record_id" gorm:"type:uuid;not null"`
	Action    string     `json:"action" gorm:"not null"` // 'INSERT', 'UPDATE', 'DELETE'
	OldValues string     `json:"old_values,omitempty" gorm:"type:jsonb"`
	NewValues string     `json:"new_values,omitempty" gorm:"type:jsonb"`
	UserID    *uuid.UUID `json:"user_id,omitempty" gorm:"type:uuid"`
	IPAddress string     `json:"ip_address"`
	RequestID string     `json:"request_id"`
	CreatedAt time.Time  `json:"created_at"`
}

// BeforeCreate hook for all models with BaseModel
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
