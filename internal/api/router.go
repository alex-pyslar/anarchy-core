package api

import (
	"net/http"

	"anarchy-core/internal/api/handler"
	"anarchy-core/internal/auth"
	"anarchy-core/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CustomValidator implements echo.Validator interface.
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the struct.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// SetupRouter sets up all application routes.
func SetupRouter(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	playerMovementHandler *handler.PlayerMovementHandler,
	jwtManager *auth.JWTManager,
	logger *util.Logger,
) {
	// Set up custom validator for Echo
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // In production, specify concrete domains
		AllowMethods: []string{http.MethodGet, http.MethodPost},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is healthy!")
	})

	// Public routes (no authentication required)
	authGroup := e.Group("/auth")
	authGroup.POST("/register", authHandler.RegisterUser)
	authGroup.POST("/login", authHandler.LoginUser)

	// WebSocket route (authenticated via query param or header)
	// The authentication logic is handled inside the WebSocket handler itself
	e.GET("/ws/game", playerMovementHandler.HandleWebSocketConnection)

	// Example of a protected HTTP route (requires JWT token in Authorization header)
	// This shows how to protect regular HTTP endpoints if you add more later.
	// For this project, player movement is via WebSocket, so this is just an example.
	protectedGroup := e.Group("/api")
	protectedGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid Authorization header")
			}
			tokenString := authHeader[7:]

			claims, err := jwtManager.ValidateToken(tokenString)
			if err != nil {
				logger.Error("JWT validation failed for protected route: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}
			// Store user info in context for later use
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			return next(c)
		}
	})

	// Example protected route (not strictly needed for this project's core logic)
	protectedGroup.GET("/profile", func(c echo.Context) error {
		userID := c.Get("userID").(string)
		username := c.Get("username").(string)
		return c.JSON(http.StatusOK, echo.Map{"message": "Welcome to your profile!", "userID": userID, "username": username})
	})
}
