package service

import (
	"anarchy-core/internal/auth"
	"anarchy-core/internal/domain"
	"anarchy-core/internal/util"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AuthService handles user authentication and registration.
type AuthService struct {
	userRepo   domain.UserRepository
	jwtManager *auth.JWTManager
	logger     *util.Logger
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo domain.UserRepository, jwtManager *auth.JWTManager, logger *util.Logger) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// RegisterUser registers a new user.
func (s *AuthService) RegisterUser(username, password string) (string, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password: %v", err)
		return "", util.ErrInternalServer
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		if errors.Is(err, util.ErrUserAlreadyExists) {
			return "", err
		}
		s.logger.Error("Failed to create user in repository: %v", err)
		return "", util.ErrInternalServer
	}

	// Generate JWT token for the newly registered user
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		s.logger.Error("Failed to generate token for new user: %v", err)
		return "", util.ErrInternalServer
	}

	s.logger.Info("User registered successfully: %s", username)
	return token, nil
}

// LoginUser authenticates a user and returns a JWT token.
func (s *AuthService) LoginUser(username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, util.ErrUserNotFound) {
			return "", util.ErrInvalidCredentials
		}
		s.logger.Error("Failed to get user by username for login: %v", err)
		return "", util.ErrInternalServer
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", util.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		s.logger.Error("Failed to generate token for login: %v", err)
		return "", util.ErrInternalServer
	}

	s.logger.Info("User logged in successfully: %s", username)
	return token, nil
}
