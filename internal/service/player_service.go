package service

import (
	"anarchy-core/internal/domain"
	"anarchy-core/internal/util"
	"errors"
)

// PlayerService handles player-related business logic, especially movement.
type PlayerService struct {
	playerMovementRepo domain.PlayerMovementRepository
	logger             *util.Logger
}

// NewPlayerService creates a new PlayerService.
func NewPlayerService(playerMovementRepo domain.PlayerMovementRepository, logger *util.Logger) *PlayerService {
	return &PlayerService{
		playerMovementRepo: playerMovementRepo,
		logger:             logger,
	}
}

// UpdatePlayerLocation updates a player's location in the database.
func (s *PlayerService) UpdatePlayerLocation(playerID string, x, y, z float64) (*domain.Location, error) {
	location := &domain.Location{
		PlayerID: playerID,
		X:        x,
		Y:        y,
		Z:        z,
	}

	err := s.playerMovementRepo.SavePlayerLocation(location)
	if err != nil {
		s.logger.Error("Failed to save player location for player %s: %v", playerID, err)
		return nil, util.ErrInternalServer
	}

	s.logger.Info("Player %s moved to (%.2f, %.2f, %.2f)", playerID, x, y, z)
	return location, nil
}

// GetPlayerLocation retrieves a player's current location.
func (s *PlayerService) GetPlayerLocation(playerID string) (*domain.Location, error) {
	location, err := s.playerMovementRepo.GetPlayerLocation(playerID)
	if err != nil {
		if errors.Is(err, util.ErrPlayerLocationNotFound) {
			return nil, err
		}
		s.logger.Error("Failed to get player location for player %s: %v", playerID, err)
		return nil, util.ErrInternalServer
	}
	return location, nil
}

// GetAllPlayerLocations retrieves all players' current locations.
func (s *PlayerService) GetAllPlayerLocations() ([]domain.Location, error) {
	locations, err := s.playerMovementRepo.GetAllPlayerLocations()
	if err != nil {
		s.logger.Error("Failed to get all player locations: %v", err)
		return nil, util.ErrInternalServer
	}
	return locations, nil
}
