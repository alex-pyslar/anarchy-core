package domain

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	CreateUser(user *User) error                      // Создать нового пользователя
	GetUserByUsername(username string) (*User, error) // Получить пользователя по имени
	GetUserByID(id string) (*User, error)             // Получить пользователя по ID
}

// PlayerMovementRepository defines the interface for player location data operations.
type PlayerMovementRepository interface {
	SavePlayerLocation(location *Location) error          // Сохранить/обновить местоположение игрока
	GetPlayerLocation(playerID string) (*Location, error) // Получить местоположение игрока по ID
	GetAllPlayerLocations() ([]Location, error)           // Получить все местоположения игроков
}
