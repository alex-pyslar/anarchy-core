package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"anarchy-core/internal/domain"
	"anarchy-core/internal/util"

	"github.com/jmoiron/sqlx"
)

// PlayerMovementRepositoryPostgres implements domain.PlayerMovementRepository for PostgreSQL.
type PlayerMovementRepositoryPostgres struct {
	db *sqlx.DB
}

// NewPlayerMovementRepositoryPostgres creates a new PlayerMovementRepositoryPostgres.
func NewPlayerMovementRepositoryPostgres(db *sqlx.DB) *PlayerMovementRepositoryPostgres {
	return &PlayerMovementRepositoryPostgres{db: db}
}

// SavePlayerLocation inserts or updates a player's location.
func (r *PlayerMovementRepositoryPostgres) SavePlayerLocation(location *domain.Location) error {
	query := `
		INSERT INTO player_locations (player_id, x, y, z)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (player_id) DO UPDATE
		SET x = EXCLUDED.x, y = EXCLUDED.y, z = EXCLUDED.z, updated_at = NOW()
		RETURNING updated_at`
	err := r.db.QueryRow(query, location.PlayerID, location.X, location.Y, location.Z).Scan(&location.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save player location: %w", err)
	}
	return nil
}

// GetPlayerLocation retrieves a player's location by player ID.
func (r *PlayerMovementRepositoryPostgres) GetPlayerLocation(playerID string) (*domain.Location, error) {
	var location domain.Location
	query := `SELECT player_id, x, y, z, updated_at FROM player_locations WHERE player_id = $1`
	err := r.db.Get(&location, query, playerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.ErrPlayerLocationNotFound
		}
		return nil, fmt.Errorf("failed to get player location: %w", err)
	}
	return &location, nil
}

// GetAllPlayerLocations retrieves all player locations.
func (r *PlayerMovementRepositoryPostgres) GetAllPlayerLocations() ([]domain.Location, error) {
	var locations []domain.Location
	query := `SELECT player_id, x, y, z, updated_at FROM player_locations`
	err := r.db.Select(&locations, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all player locations: %w", err)
	}
	return locations, nil
}
