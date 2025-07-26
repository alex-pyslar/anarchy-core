package domain

import "time"

// Location represents a player's 3D coordinates in the game world.
type Location struct {
	PlayerID  string    `db:"player_id"`  // Идентификатор игрока
	X         float64   `db:"x"`          // Координата X
	Y         float64   `db:"y"`          // Координата Y
	Z         float64   `db:"z"`          // Координата Z
	UpdatedAt time.Time `db:"updated_at"` // Время последнего обновления
}
