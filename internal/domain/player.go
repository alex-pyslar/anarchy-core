package domain

// Player represents a game player. Currently, it's directly linked to a User.
type Player struct {
	ID       string // Player ID (same as User ID)
	Username string // Player's username
}
