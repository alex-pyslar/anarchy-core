package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"anarchy-core/internal/auth"
	"anarchy-core/internal/service"
	"anarchy-core/internal/util"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// PlayerMovementHandler handles WebSocket connections and player movement.
type PlayerMovementHandler struct {
	playerService    *service.PlayerService
	websocketService *service.WebSocketService
	jwtManager       *auth.JWTManager
	logger           *util.Logger
	upgrader         websocket.Upgrader
}

// NewPlayerMovementHandler creates a new PlayerMovementHandler.
func NewPlayerMovementHandler(
	playerService *service.PlayerService,
	websocketService *service.WebSocketService,
	jwtManager *auth.JWTManager,
	logger *util.Logger,
) *PlayerMovementHandler {
	return &PlayerMovementHandler{
		playerService:    playerService,
		websocketService: websocketService,
		jwtManager:       jwtManager,
		logger:           logger,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development. In production, restrict this.
				return true
			},
		},
	}
}

// PlayerMovementMessage represents a message from client about player movement.
type PlayerMovementMessage struct {
	Type string  `json:"type"` // "move"
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
}

// HandleWebSocketConnection handles the WebSocket upgrade and message loop.
func (h *PlayerMovementHandler) HandleWebSocketConnection(c echo.Context) error {
	// Extract token from query parameter or header for WebSocket authentication
	tokenString := c.QueryParam("token")
	if tokenString == "" {
		// Fallback to Authorization header if query param is empty
		authHeader := c.Request().Header.Get("Authorization")
		if len(authHeader) > 7 && authHeader[:7] != "Bearer " {
			tokenString = authHeader[7:]
		}
	}

	if tokenString == "" {
		h.logger.Error("WebSocket: No token provided")
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication token required")
	}

	claims, err := h.jwtManager.ValidateToken(tokenString)
	if err != nil {
		h.logger.Error("WebSocket: Invalid token: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
	}

	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upgrade to WebSocket")
	}

	client := &service.Client{
		UserID:   claims.UserID,
		Username: claims.Username,
		Conn:     conn,
		Send:     make(chan []byte, 256), // Буферизованный канал для отправки
	}

	h.websocketService.RegisterClient(client)
	h.logger.Info("WebSocket client connected: %s (ID: %s)", client.Username, client.UserID)

	// Send initial state to the newly connected client
	allLocations, err := h.playerService.GetAllPlayerLocations()
	if err != nil {
		h.logger.Error("Failed to get all player locations for initial state: %v", err)
	} else {
		h.websocketService.SendAllPlayerLocations(client, allLocations)
	}

	// Goroutine for reading messages from the client
	go h.readPump(client)
	// Goroutine for writing messages to the client
	go h.writePump(client)

	return nil // Connection is handled by goroutines
}

// readPump pumps messages from the websocket connection to the broadcast channel.
func (h *PlayerMovementHandler) readPump(client *service.Client) {
	defer func() {
		h.websocketService.UnregisterClient(client)
		client.Conn.Close()
		h.logger.Info("WebSocket client disconnected (readPump): %s", client.Username)
	}()

	client.Conn.SetReadLimit(512)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error("WebSocket read error for client %s: %v", client.Username, err)
			}
			break
		}

		var moveMsg PlayerMovementMessage
		if err := json.Unmarshal(message, &moveMsg); err != nil {
			h.logger.Error("Failed to unmarshal player movement message from client %s: %v", client.Username, err)
			continue
		}

		if moveMsg.Type == "move" {
			loc, err := h.playerService.UpdatePlayerLocation(client.UserID, moveMsg.X, moveMsg.Y, moveMsg.Z)
			if err != nil {
				h.logger.Error("Failed to update player location for client %s: %v", client.Username, err)
				// Optionally send an error message back to the client
				// client.Send <- []byte(`{"type": "error", "message": "Failed to update location"}`)
				continue
			}
			// Notify all other players about the movement
			h.websocketService.NotifyPlayerLocationChange(client.UserID, client.Username, loc)
		} else {
			h.logger.Info("Received unknown message type '%s' from client %s", moveMsg.Type, client.Username)
		}
	}
}

// writePump pumps messages from the WebSocketService's send channel to the websocket connection.
func (h *PlayerMovementHandler) writePump(client *service.Client) {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
		h.logger.Info("WebSocket client disconnected (writePump): %s", client.Username)
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				// The WebSocketService closed the channel.
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				h.logger.Error("Failed to write message to client %s: %v", client.Username, err)
				return
			}
		case <-ticker.C:
			// Send a ping message to keep the connection alive
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				h.logger.Error("Failed to send ping to client %s: %v", client.Username, err)
				return
			}
		}
	}
}
