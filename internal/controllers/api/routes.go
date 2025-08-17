package api

import (
	"payslip-system/internal/middleware"
	"payslip-system/internal/repository"
	"payslip-system/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, services *service.Services) {
	handlers := NewHandlers(services)

	// Get repository for middleware (needed for auth middleware)
	repos := &repository.Repositories{} // This would be passed from main

	// Public routes
	public := r.Group("/api/v1")
	{
		public.GET("/health", handlers.HealthCheck)
		public.POST("/login", handlers.Login)
	}

	// Protected routes (requires authentication)
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(repos))
	{
		// Employee routes
		employee := protected.Group("/employee")
		employee.Use(middleware.EmployeeMiddleware())
		{
			employee.POST("/attendance", handlers.SubmitAttendance)
			employee.POST("/overtime", handlers.SubmitOvertime)
			employee.POST("/reimbursement", handlers.SubmitReimbursement)
			employee.GET("/payslip/:period_id", handlers.GeneratePayslip)
		}

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/attendance-period", handlers.CreateAttendancePeriod)
			admin.POST("/payroll/:period_id/process", handlers.ProcessPayroll)
			admin.GET("/payroll/:period_id/summary", handlers.GeneratePayrollSummary)
		}
	}
}

// SetupRoutesWithRepos is a helper function that accepts repositories
func SetupRoutesWithRepos(r *gin.Engine, services *service.Services, repos *repository.Repositories) {
	handlers := NewHandlers(services)

	// Public routes
	public := r.Group("/api/v1")
	{
		public.GET("/health", handlers.HealthCheck)
		public.POST("/login", handlers.Login)
	}

	// Protected routes (requires authentication)
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(repos))
	{
		// Employee routes
		employee := protected.Group("/employee")
		employee.Use(middleware.EmployeeMiddleware())
		{
			employee.POST("/attendance", handlers.SubmitAttendance)
			employee.POST("/overtime", handlers.SubmitOvertime)
			employee.POST("/reimbursement", handlers.SubmitReimbursement)
			employee.GET("/payslip/:period_id", handlers.GeneratePayslip)
		}

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/attendance-period", handlers.CreateAttendancePeriod)
			admin.POST("/payroll/:period_id/process", handlers.ProcessPayroll)
			admin.GET("/payroll/:period_id/summary", handlers.GeneratePayrollSummary)
		}
	}
}
