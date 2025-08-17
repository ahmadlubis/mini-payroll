package service_test

import (
	"payslip-system/internal/models"
	"payslip-system/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	db, cleanup := test.SetupTestDB()
	defer cleanup()

	repos, services := test.SetupTestServices(db)

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

	err := repos.User.Create(testUser)
	require.NoError(t, err)

	tests := []struct {
		name        string
		username    string
		password    string
		expectError bool
	}{
		{
			name:        "valid credentials",
			username:    "testuser",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "invalid username",
			username:    "wronguser",
			password:    "password123",
			expectError: true,
		},
		{
			name:        "invalid password",
			username:    "testuser",
			password:    "wrongpassword",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, _, err := services.Auth.Login(tt.username, tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
			}
		})
	}
}
