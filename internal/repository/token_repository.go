package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/randhir/aegis-core/internal/models"
)

func CreateRefreshToken(userID uuid.UUID, tokenID uuid.UUID, token string, expiresAt time.Time) (*models.RefreshToken, error) {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, token, expires_at, created_at
	`

	var refreshToken models.RefreshToken
	err := DB.QueryRow(query, tokenID, userID, token, expiresAt).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &refreshToken, nil
}

func GetRefreshTokenByToken(token string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1
	`

	var refreshToken models.RefreshToken
	err := DB.QueryRow(query, token).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return &refreshToken, nil
}

func DeleteRefreshToken(tokenID uuid.UUID) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE id = $1
	`

	result, err := DB.Exec(query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("refresh token not found")
	}

	return nil
}

func GetRefreshTokenByID(tokenID uuid.UUID) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE id = $1
	`

	var refreshToken models.RefreshToken
	err := DB.QueryRow(query, tokenID).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return &refreshToken, nil
}

