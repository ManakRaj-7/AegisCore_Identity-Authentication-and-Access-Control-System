package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/randhir/aegis-core/internal/config"
)

type AccessTokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID  string `json:"user_id"`
	TokenID string `json:"token_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID, email, role string) (string, error) {
	secret := config.AppConfig.JWT.AccessSecret
	if secret == "" {
		return "", errors.New("JWT_ACCESS_SECRET not configured")
	}

	claims := AccessTokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID, tokenID string) (string, error) {
	secret := config.AppConfig.JWT.RefreshSecret
	if secret == "" {
		return "", errors.New("JWT_REFRESH_SECRET not configured")
	}

	claims := RefreshTokenClaims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	secret := config.AppConfig.JWT.AccessSecret
	if secret == "" {
		return nil, errors.New("JWT_ACCESS_SECRET not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

func ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	secret := config.AppConfig.JWT.RefreshSecret
	if secret == "" {
		return nil, errors.New("JWT_REFRESH_SECRET not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

