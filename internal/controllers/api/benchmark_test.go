package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"payslip-system/internal/controllers/api"
	"payslip-system/internal/models"
	"payslip-system/internal/test"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func setupTestRouter() (*gin.Engine, func()) {
	gin.SetMode(gin.TestMode)

	db, cleanup := test.SetupTestDB()
	repos, services := test.SetupTestServices(db)

	r := gin.New()
	api.SetupRoutesWithRepos(r, services, repos)

	return r, cleanup
}

func TestLogin_Integration(t *testing.T) {
	r, cleanup := setupTestRouter()
	defer cleanup()

	// Create test user directly in database
	db, _ := test.SetupTestDB()
	repos, _ := test.SetupTestServices(db)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &models.User{
		Username: "testuser",
		Password: string(hashedPassword),
		Role:     "employee",
		IsActive: true,
	}
	repos.User.Create(testUser)

	// Test login
	loginReq := map[string]string{
		"username": "testuser",
		"password": "password123",
	}

	jsonBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")
}

func TestSubmitAttendance_Integration(t *testing.T) {
	r, cleanup := setupTestRouter()
	defer cleanup()

	// Setup test data
	db, _ := test.SetupTestDB()
	repos, _ := test.SetupTestServices(db)

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &models.User{
		Username: "testuser",
		Password: string(hashedPassword),
		Role:     "employee",
		IsActive: true,
	}
	salary := 5000000.0
	testUser.Salary = &salary
	repos.User.Create(testUser)

	// Create attendance period
	period := &models.AttendancePeriod{
		StartDate:   time.Now().AddDate(0, 0, -10),
		EndDate:     time.Now().AddDate(0, 0, 10),
		IsProcessed: false,
	}
	repos.AttendancePeriod.Create(period)

	// Get auth token
	token := getAuthToken(t, r, "testuser", "password123")

	// Test submit attendance
	attendanceReq := map[string]string{
		"date":          time.Now().Format("2006-01-02"),
		"check_in_time": "09:00",
	}

	jsonBody, _ := json.Marshal(attendanceReq)
	req, _ := http.NewRequest("POST", "/api/v1/employee/attendance", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("client_ip", "192.168.1.100")
	req.Header.Set("user_id", uuid.New().String())
	req.Header.Set("request_id", uuid.New().String())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Skip if today is weekend
	if time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday {
		assert.Equal(t, http.StatusBadRequest, w.Code)
	} else {
		assert.Equal(t, http.StatusCreated, w.Code)
	}
}

func getAuthToken(t *testing.T, r *gin.Engine, username, password string) string {
	loginReq := map[string]string{
		"username": username,
		"password": password,
	}

	jsonBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	return response["token"].(string)
}

// Benchmark tests
func BenchmarkLogin(b *testing.B) {
	r, cleanup := setupTestRouter()
	defer cleanup()

	// Setup test data
	db, _ := test.SetupTestDB()
	repos, _ := test.SetupTestServices(db)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &models.User{
		Username: "benchuser",
		Password: string(hashedPassword),
		Role:     "employee",
		IsActive: true,
	}
	repos.User.Create(testUser)

	loginReq := map[string]string{
		"username": "benchuser",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(loginReq)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
		}
	})
}

func BenchmarkPayslipGeneration(b *testing.B) {
	r, cleanup := setupTestRouter()
	defer cleanup()

	// Setup comprehensive test data
	db, _ := test.SetupTestDB()
	repos, _ := test.SetupTestServices(db)

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &models.User{
		Username: "benchuser",
		Password: string(hashedPassword),
		Role:     "employee",
		IsActive: true,
	}
	salary := 5000000.0
	testUser.Salary = &salary
	repos.User.Create(testUser)

	// Create attendance period
	period := &models.AttendancePeriod{
		StartDate:   time.Now().AddDate(0, -1, 0),
		EndDate:     time.Now().AddDate(0, 0, -1),
		IsProcessed: false,
	}
	repos.AttendancePeriod.Create(period)

	// Add attendance records
	for i := 0; i < 20; i++ {
		date := period.StartDate.AddDate(0, 0, i)
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			attendance := &models.Attendance{
				UserID:             testUser.ID,
				AttendancePeriodID: period.ID,
				Date:               date,
				CheckInTime:        time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location()),
			}
			repos.Attendance.Create(attendance)
		}
	}

	token := getAuthTokenBenchmark(b, r, "benchuser", "password123")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/api/v1/employee/payslip/"+period.ID.String(), nil)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
		}
	})
}

// Helper function for benchmark
func getAuthTokenBenchmark(b *testing.B, r *gin.Engine, username, password string) string {
	loginReq := map[string]string{
		"username": username,
		"password": password,
	}

	jsonBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	return response["token"].(string)
}
