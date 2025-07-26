package handler

import (
	"errors"
	"net/http"

	"anarchy-core/internal/service"
	"anarchy-core/internal/util"

	"github.com/labstack/echo/v4"
)

// AuthHandler handles HTTP requests for user authentication.
type AuthHandler struct {
	authService *service.AuthService
	logger      *util.Logger
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService, logger *util.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// RegisterRequest represents the request body for user registration.
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=50"`
}

// RegisterUser handles user registration.
func (h *AuthHandler) RegisterUser(c echo.Context) error {
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Error("RegisterUser: Failed to bind request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}
	if err := c.Validate(req); err != nil {
		h.logger.Error("RegisterUser: Validation failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, err := h.authService.RegisterUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, util.ErrUserAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, "User with this username already exists")
		}
		h.logger.Error("RegisterUser: Failed to register user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to register user")
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "User registered successfully", "token": token})
}

// LoginRequest represents the request body for user login.
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginUser handles user login.
func (h *AuthHandler) LoginUser(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		h.logger.Error("LoginUser: Failed to bind request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}
	if err := c.Validate(req); err != nil {
		h.logger.Error("LoginUser: Validation failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, err := h.authService.LoginUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, util.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
		}
		h.logger.Error("LoginUser: Failed to login user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to login user")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Login successful", "token": token})
}
