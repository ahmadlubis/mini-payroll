package api

import (
	"net/http"
	"time"

	"payslip-system/internal/middleware"
	"payslip-system/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handlers struct {
	services *service.Services
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{services: services}
}

// Login Request/Response
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		Role     string    `json:"role"`
	} `json:"user"`
}

// Auth handlers
func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _, err := h.services.Auth.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := LoginResponse{
		Token: token,
	}
	response.User.ID = user.ID
	response.User.Username = user.Username
	response.User.Role = user.Role

	c.JSON(http.StatusOK, response)
}

// Attendance requests
type SubmitAttendanceRequest struct {
	Date        string `json:"date" binding:"required"`          // YYYY-MM-DD format
	CheckInTime string `json:"check_in_time" binding:"required"` // HH:MM format
}

func (h *Handlers) SubmitAttendance(c *gin.Context) {
	var req SubmitAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	clientIP := "-"
	// Use this safer approach:
	clientIPRaw, exists := c.Get("client_ip")
	if exists {
		clientIPStr, ok := clientIPRaw.(string)
		if ok {
			clientIP = clientIPStr
		}
	}
	requestID := "-"
	// Use this safer approach:
	requestIDRaw, exists := c.Get("request_id")
	if exists {
		requestIDStr, ok := requestIDRaw.(string)
		if ok {
			requestID = requestIDStr
		}
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	// Parse check-in time
	checkInTime, err := time.Parse("15:04", req.CheckInTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check-in time format, use HH:MM"})
		return
	}

	// Combine date and time
	checkInDateTime := time.Date(date.Year(), date.Month(), date.Day(),
		checkInTime.Hour(), checkInTime.Minute(), 0, 0, date.Location())

	if err := h.services.Attendance.SubmitAttendance(userID, date, checkInDateTime, clientIP, requestID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Attendance submitted successfully"})
}

// Overtime requests
type SubmitOvertimeRequest struct {
	Date  string  `json:"date" binding:"required"` // YYYY-MM-DD format
	Hours float64 `json:"hours" binding:"required,gt=0,lte=3"`
}

func (h *Handlers) SubmitOvertime(c *gin.Context) {
	var req SubmitOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	clientIP := c.MustGet("client_ip").(string)
	requestID := c.MustGet("request_id").(string)

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	if err := h.services.Overtime.SubmitOvertime(userID, date, req.Hours, clientIP, requestID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Overtime submitted successfully"})
}

// Reimbursement requests
type SubmitReimbursementRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"required"`
}

func (h *Handlers) SubmitReimbursement(c *gin.Context) {
	var req SubmitReimbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	clientIP := c.MustGet("client_ip").(string)
	requestID := c.MustGet("request_id").(string)

	if err := h.services.Reimbursement.SubmitReimbursement(userID, req.Amount, req.Description, clientIP, requestID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Reimbursement submitted successfully"})
}

// Payslip generation
func (h *Handlers) GeneratePayslip(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	periodIDStr := c.Param("period_id")

	periodID, err := uuid.Parse(periodIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period ID"})
		return
	}

	payslip, err := h.services.Payroll.GeneratePayslip(userID, periodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payslip)
}

// Admin handlers
type CreateAttendancePeriodRequest struct {
	StartDate string `json:"start_date" binding:"required"` // YYYY-MM-DD format
	EndDate   string `json:"end_date" binding:"required"`   // YYYY-MM-DD format
}

func (h *Handlers) CreateAttendancePeriod(c *gin.Context) {
	var req CreateAttendancePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.MustGet("user_id").(uuid.UUID)
	clientIP := c.MustGet("client_ip").(string)
	requestID := c.MustGet("request_id").(string)

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format, use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format, use YYYY-MM-DD"})
		return
	}

	period, err := h.services.Admin.CreateAttendancePeriod(startDate, endDate, adminID, clientIP, requestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, period)
}

func (h *Handlers) ProcessPayroll(c *gin.Context) {
	periodIDStr := c.Param("period_id")
	periodID, err := uuid.Parse(periodIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period ID"})
		return
	}

	adminID := c.MustGet("user_id").(uuid.UUID)
	clientIP := c.MustGet("client_ip").(string)
	requestID := c.MustGet("request_id").(string)

	if err := h.services.Payroll.ProcessPayroll(periodID, adminID, clientIP, requestID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payroll processed successfully"})
}

func (h *Handlers) GeneratePayrollSummary(c *gin.Context) {
	periodIDStr := c.Param("period_id")
	periodID, err := uuid.Parse(periodIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period ID"})
		return
	}

	summary, err := h.services.Payroll.GeneratePayrollSummary(periodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// Health check
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "payslip-system",
	})
}
