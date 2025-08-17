package service

import (
	"errors"
	"payslip-system/internal/models"
	"payslip-system/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repos *repository.Repositories
}

func NewAuthService(repos *repository.Repositories) *authService {
	return &authService{repos: repos}
}

func (s *authService) Login(username, password string) (*models.User, string, error) {
	user, err := s.repos.User.GetByUsername(username)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate token would be imported from middleware
	return user, "", nil
}

func (s *authService) ValidateToken(tokenString string) (*models.User, error) {
	// Token validation logic would be here
	return nil, nil
}
