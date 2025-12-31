package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/randhir/aegis-core/internal/models"
)

func CreateUser(email, passwordHash, role string) (*models.User, error) {
	userID := uuid.New()
	query := `
		INSERT INTO users (id, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, password_hash, role, created_at
	`

	var user models.User
	err := DB.QueryRow(query, userID, email, passwordHash, role).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if isUniqueConstraintError(err) {
			return nil, errors.New("email already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func GetUserByID(userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func UserExistsByEmail(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "duplicate key") || contains(errStr, "unique constraint") || contains(errStr, "violates unique constraint")
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func ListUsers() ([]models.User, error) {
	query := `
		SELECT id, email, role, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

