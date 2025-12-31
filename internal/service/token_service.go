package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/randhir/aegis-core/internal/cache"
	"github.com/randhir/aegis-core/internal/repository"
	"github.com/randhir/aegis-core/internal/utils"
)

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

// Refresh generates new access and refresh tokens, invalidating the old refresh token
func (s *TokenService) Refresh(refreshTokenString string) (string, string, error) {
	// Validate refresh token signature and expiry
	claims, err := utils.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", utils.ErrInvalidToken
	}

	// Check if refresh token exists in DB
	dbToken, err := repository.GetRefreshTokenByToken(refreshTokenString)
	if err != nil {
		return "", "", utils.ErrInvalidToken
	}

	// Verify token hasn't expired
	if time.Now().After(dbToken.ExpiresAt) {
		return "", "", utils.ErrInvalidToken
	}

	// Verify token ID matches
	tokenID, err := uuid.Parse(claims.TokenID)
	if err != nil || tokenID != dbToken.ID {
		return "", "", utils.ErrInvalidToken
	}

	// Parse user ID from claims
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", "", utils.ErrInvalidToken
	}

	// Get user details
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return "", "", utils.ErrInternalError
	}

	// Delete old refresh token (rotation)
	err = repository.DeleteRefreshToken(dbToken.ID)
	if err != nil {
		return "", "", utils.ErrInternalError
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return "", "", utils.ErrInternalError
	}

	// Generate new refresh token with new token ID
	newTokenID := uuid.New()
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID.String(), newTokenID.String())
	if err != nil {
		return "", "", utils.ErrInternalError
	}

	// Store new refresh token
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = repository.CreateRefreshToken(user.ID, newTokenID, newRefreshToken, expiresAt)
	if err != nil {
		return "", "", utils.ErrInternalError
	}

	return accessToken, newRefreshToken, nil
}

// Logout invalidates the refresh token and blacklists the access token
func (s *TokenService) Logout(refreshTokenString, accessTokenString string) error {
	// Validate refresh token to get claims
	claims, err := utils.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return utils.ErrInvalidToken
	}

	// Get refresh token from DB
	dbToken, err := repository.GetRefreshTokenByToken(refreshTokenString)
	if err != nil {
		return utils.ErrInvalidToken
	}

	// Verify token ID matches
	tokenID, err := uuid.Parse(claims.TokenID)
	if err != nil || tokenID != dbToken.ID {
		return utils.ErrInvalidToken
	}

	// Delete refresh token from DB
	err = repository.DeleteRefreshToken(dbToken.ID)
	if err != nil {
		return utils.ErrInternalError
	}

	// Blacklist access token in Redis
	if accessTokenString != "" {
		// Parse access token to get expiry
		accessClaims, err := utils.ValidateAccessToken(accessTokenString)
		if err == nil && accessClaims.ExpiresAt != nil {
			expiryTime := accessClaims.ExpiresAt.Time
			err = cache.BlacklistAccessToken(accessTokenString, expiryTime)
			if err != nil {
				// Log error but don't fail logout
				// In production, you might want to handle this differently
			}
		}
	}

	return nil
}

