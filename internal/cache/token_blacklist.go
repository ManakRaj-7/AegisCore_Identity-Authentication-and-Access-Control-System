package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const blacklistPrefix = "blacklist:access_token:"

// BlacklistAccessToken adds an access token to the Redis blacklist with TTL equal to token expiry
func BlacklistAccessToken(tokenString string, expiryTime time.Time) error {
	if Client == nil {
		return errors.New("redis client not initialized")
	}

	ctx := context.Background()
	key := blacklistPrefix + tokenString

	// Calculate TTL from now until expiry
	ttl := time.Until(expiryTime)
	if ttl <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}

	err := Client.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

// IsAccessTokenBlacklisted checks if an access token is in the Redis blacklist
func IsAccessTokenBlacklisted(tokenString string) (bool, error) {
	if Client == nil {
		return false, errors.New("redis client not initialized")
	}

	ctx := context.Background()
	key := blacklistPrefix + tokenString

	exists, err := Client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}

	return exists > 0, nil
}

