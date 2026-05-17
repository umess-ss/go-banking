package services

import (
	"context"
	"errors"
	"go-banking/internal/models"
	"go-banking/internal/repository"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, request models.RegisterRequest) (models.UserResponse, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	if request.Name == "" {
		return models.UserResponse{}, errors.New("name is required")
	}
	if request.Email == "" {
		return models.UserResponse{}, errors.New("email is required")
	}
	if request.Password == "" {
		return models.UserResponse{}, errors.New("password is required")
	}

	if len(request.Password) < 6 {
		return models.UserResponse{}, errors.New("password must be at least 6 characters")
	}

	existingUser, _ := s.userRepo.FindByEmail(ctx, request.Email)
	if existingUser != nil {
		return models.UserResponse{}, errors.New("user already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserResponse{}, errors.New("failed to hash password")
	}

	user := models.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return models.UserResponse{}, errors.New("failed to create user")
	}

	return models.UserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}
