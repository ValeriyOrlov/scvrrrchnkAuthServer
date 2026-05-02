package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/model"
	"github.com/ValeriyOrlov/scvrrrchnkAuthServer/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrWeakPassword       = errors.New("min 8 characters")
	ErrShortUsername      = errors.New("min 3 characters")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
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

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", fmt.Errorf("find user by email: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", fmt.Errorf("comparing hash and password: %w", err)
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.accessTTL).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	jti := uuid.New().String()
	refreshExpAt := time.Now().Add(s.refreshTTL)
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     refreshExpAt,
		"iat":     time.Now().Unix(),
		"jti":     jti,
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	refreshModel := model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: refreshExpAt,
	}

	err = s.tokenRepo.Create(ctx, &refreshModel)
	if err != nil {
		return "", "", fmt.Errorf("adding refresh token to bd: %w", err)

	}
	return accessToken, refreshToken, nil
}
