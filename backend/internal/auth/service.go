package auth

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *UserRepository
}

func NewAuthService(userRepo *UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, request RegisterRequest) (UserResponse, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	if request.Name == "" {
		return UserResponse{}, errors.New("name is required")
	}
	if request.Email == "" {
		return UserResponse{}, errors.New("email is required")
	}
	if request.Password == "" {
		return UserResponse{}, errors.New("password is required")
	}

	if len(request.Password) < 6 {
		return UserResponse{}, errors.New("password must be at least 6 characters")
	}

	existingUser, _ := s.userRepo.FindByEmail(ctx, request.Email)
	if existingUser != nil {
		return UserResponse{}, errors.New("user already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, errors.New("failed to hash password")
	}

	user := User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return UserResponse{}, errors.New("failed to create user")
	}

	return UserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, request LoginRequest) (LoginResponse, error) {
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	if request.Email == "" {
		return LoginResponse{}, errors.New("email is required")
	}
	if request.Password == "" {
		return LoginResponse{}, errors.New("password is required")
	}

	user, err := s.userRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		return LoginResponse{}, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return LoginResponse{}, errors.New("invalid email or password")
	}

	token, err := s.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return LoginResponse{}, errors.New("failed to generate access token")
	}

	return LoginResponse{
		AccessToken: token,
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID int64) (UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return UserResponse{}, errors.New("user not found")
	}

	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *AuthService) GenerateJWT(userID int64, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT secret is not set")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
