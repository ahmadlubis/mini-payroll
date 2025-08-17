package domains

import "payslip-system/internal/models"

type PayslipResponse struct {
	Employee            *models.User             `json:"employee"`
	Period              *models.AttendancePeriod `json:"period"`
	BaseSalary          float64                  `json:"base_salary"`
	AttendanceDays      int                      `json:"attendance_days"`
	WorkingDays         int                      `json:"working_days"`
	AttendanceAmount    float64                  `json:"attendance_amount"`
	OvertimeHours       float64                  `json:"overtime_hours"`
	OvertimeAmount      float64                  `json:"overtime_amount"`
	Reimbursements      []models.Reimbursement   `json:"reimbursements"`
	ReimbursementAmount float64                  `json:"reimbursement_amount"`
	TotalAmount         float64                  `json:"total_amount"`
}

type PayrollSummaryResponse struct {
	Period      *models.AttendancePeriod `json:"period"`
	Employees   []EmployeeSummary        `json:"employees"`
	TotalAmount float64                  `json:"total_amount"`
}

type EmployeeSummary struct {
	Employee    *models.User `json:"employee"`
	TotalAmount float64      `json:"total_amount"`
}
