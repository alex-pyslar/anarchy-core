package service

import (
	"encoding/json"
	"sync"

	"anarchy-core/internal/domain"
	"anarchy-core/internal/util"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client.
type Client struct {
	UserID   string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte // Канал для отправки сообщений клиенту
}

// WebSocketService manages WebSocket connections and broadcasts.
type WebSocketService struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	logger     *util.Logger
	mu         sync.Mutex
}

// NewWebSocketService creates a new WebSocketService.
func NewWebSocketService(logger *util.Logger) *WebSocketService {
	return &WebSocketService{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

// Run starts the WebSocket service's main loop.
func (s *WebSocketService) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()
			s.logger.Info("Client registered: %s (ID: %s)", client.Username, client.UserID)
			// Optionally send current game state to new client
			// s.SendAllPlayerLocations(client)
		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.Send)
				s.logger.Info("Client unregistered: %s (ID: %s)", client.Username, client.UserID)
			}
			s.mu.Unlock()
		case message := <-s.broadcast:
			s.mu.Lock()
			for client := range s.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(s.clients, client)
					s.logger.Error("Failed to send message to client %s, unregistering.", client.Username)
				}
			}
			s.mu.Unlock()
		}
	}
}

// RegisterClient registers a new WebSocket client.
func (s *WebSocketService) RegisterClient(client *Client) {
	s.register <- client
}

// UnregisterClient unregisters a WebSocket client.
func (s *WebSocketService) UnregisterClient(client *Client) {
	s.unregister <- client
}

// BroadcastMessage sends a message to all connected clients.
func (s *WebSocketService) BroadcastMessage(message []byte) {
	s.broadcast <- message
}

// PlayerLocationUpdate represents a message about player location change.
type PlayerLocationUpdate struct {
	Type      string  `json:"type"`
	PlayerID  string  `json:"player_id"`
	Username  string  `json:"username"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Z         float64 `json:"z"`
	Timestamp string  `json:"timestamp"`
}

// NotifyPlayerLocationChange sends a player's location update to all clients.
func (s *WebSocketService) NotifyPlayerLocationChange(playerID, username string, loc *domain.Location) {
	update := PlayerLocationUpdate{
		Type:      "player_location_update",
		PlayerID:  playerID,
		Username:  username,
		X:         loc.X,
		Y:         loc.Y,
		Z:         loc.Z,
		Timestamp: loc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"), // ISO 8601 format
	}

	message, err := json.Marshal(update)
	if err != nil {
		s.logger.Error("Failed to marshal player location update: %v", err)
		return
	}
	s.BroadcastMessage(message)
}

// InitialStateMessage represents the initial state of the game.
type InitialStateMessage struct {
	Type      string                 `json:"type"`
	Locations []PlayerLocationUpdate `json:"locations"`
}

// SendAllPlayerLocations sends the current locations of all players to a specific client.
func (s *WebSocketService) SendAllPlayerLocations(client *Client, locations []domain.Location) {
	updates := make([]PlayerLocationUpdate, len(locations))
	for i, loc := range locations {
		// In a real app, you might need to fetch username here if not readily available
		// For simplicity, assuming username can be derived or passed
		updates[i] = PlayerLocationUpdate{
			Type:      "player_location_update",
			PlayerID:  loc.PlayerID,
			Username:  "unknown", // Placeholder: you'll need to fetch username from user service
			X:         loc.X,
			Y:         loc.Y,
			Z:         loc.Z,
			Timestamp: loc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	initialState := InitialStateMessage{
		Type:      "initial_state",
		Locations: updates,
	}

	message, err := json.Marshal(initialState)
	if err != nil {
		s.logger.Error("Failed to marshal initial state message: %v", err)
		return
	}

	select {
	case client.Send <- message:
	default:
		s.logger.Error("Failed to send initial state to client %s, client channel is full.", client.Username)
		s.UnregisterClient(client) // Unregister if cannot send
	}
}
