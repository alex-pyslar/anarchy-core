package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"anarchy-core/internal/domain"
	"anarchy-core/internal/util"

	"github.com/jmoiron/sqlx"
)

// UserRepositoryPostgres implements domain.UserRepository for PostgreSQL.
type UserRepositoryPostgres struct {
	db *sqlx.DB
}

// NewUserRepositoryPostgres creates a new UserRepositoryPostgres.
func NewUserRepositoryPostgres(db *sqlx.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

// CreateUser inserts a new user into the database.
func (r *UserRepositoryPostgres) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db.QueryRow(query, user.Username, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		// Check for unique constraint violation
		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return util.ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByUsername retrieves a user by their username.
func (r *UserRepositoryPostgres) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`
	err := r.db.Get(&user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *UserRepositoryPostgres) GetUserByID(id string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, username, password_hash, created_at FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}
