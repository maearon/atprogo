package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	DID          string    `json:"did"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// UserRepository handles user data access
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(ctx context.Context, user *User, password string) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user into database
	query := `
		INSERT INTO users (did, username, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	now := time.Now()
	err = r.db.QueryRow(ctx, query,
		user.DID,
		user.Username,
		user.Email,
		string(hashedPassword),
		now,
		now,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	user.PasswordHash = string(hashedPassword)

	return nil
}

// GetUserByUsername gets a user by username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, did, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	var user User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.DID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByDID gets a user by DID
func (r *UserRepository) GetUserByDID(ctx context.Context, did string) (*User, error) {
	query := `
		SELECT id, did, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE did = $1
	`
	var user User
	err := r.db.QueryRow(ctx, query, did).Scan(
		&user.ID,
		&user.DID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// VerifyPassword verifies a user's password
func (r *UserRepository) VerifyPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
