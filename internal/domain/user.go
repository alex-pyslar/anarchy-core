package domain

import "time"

// User represents a user in the system.
type User struct {
	ID           string    `db:"id"`            // Уникальный идентификатор пользователя (UUID)
	Username     string    `db:"username"`      // Имя пользователя
	PasswordHash string    `db:"password_hash"` // Хеш пароля
	CreatedAt    time.Time `db:"created_at"`    // Время создания
}
