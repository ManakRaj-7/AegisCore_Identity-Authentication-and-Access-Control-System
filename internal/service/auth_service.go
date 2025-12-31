package service

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/randhir/aegis-core/internal/repository"
	"github.com/randhir/aegis-core/internal/utils"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(email, password string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	if !utils.ValidateEmail(email) {
		return &utils.AppError{Message: "invalid email format", StatusCode: 400}
	}

	if !utils.ValidatePassword(password) {
		return &utils.AppError{Message: "password must be at least 8 characters long", StatusCode: 400}
	}

	exists, err := repository.UserExistsByEmail(email)
	if err != nil {
		return utils.ErrInternalError
	}
	if exists {
		return utils.ErrConflict
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return utils.ErrInternalError
	}

	_, err = repository.CreateUser(email, passwordHash, "USER")
	if err != nil {
		if err.Error() == "email already exists" {
			return utils.ErrConflict
		}
		return utils.ErrInternalError
	}

	return nil
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)

	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return "", "", utils.ErrInvalidCredentials
	}

	if !utils.ComparePassword(user.PasswordHash, password) {
		return "", "", utils.ErrInvalidCredentials
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return "", "", utils.ErrInternalError
	}

	tokenID := uuid.New()
	refreshToken, err := utils.GenerateRefreshToken(user.ID.String(), tokenID.String())
	if err != nil {
		return "", "", utils.ErrInternalError
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = repository.CreateRefreshToken(user.ID, tokenID, refreshToken, expiresAt)
	if err != nil {
		return "", "", utils.ErrInternalError
	}

	return accessToken, refreshToken, nil
}

