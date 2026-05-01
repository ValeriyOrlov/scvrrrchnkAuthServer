package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/model"
	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrWeakPassword      = errors.New("min 8 characters")
	ErrShortUsername     = errors.New("min 3 characters")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (*model.User, error) {
	if email == "" || !strings.Contains(email, "@") {
		return nil, ErrInvalidEmail
	}
	if len(username) < 3 {
		return nil, ErrShortUsername
	}
	if len(password) < 8 {
		return nil, ErrWeakPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	newUser, err := s.userRepo.Create(ctx, email, username, string(hashedPassword))
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("registration error: %w", err)
	}
	return &newUser, nil
}
