package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/randhir/aegis-core/internal/logger"
	"github.com/randhir/aegis-core/internal/middleware"
	"github.com/randhir/aegis-core/internal/service"
	"github.com/randhir/aegis-core/internal/utils"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, utils.ErrInvalidRequest)
		return
	}

	// Validate email format
	if !utils.ValidateEmail(req.Email) {
		middleware.ErrorResponse(c, &utils.AppError{Message: "invalid email format", StatusCode: http.StatusBadRequest})
		return
	}

	// Validate password length
	if !utils.ValidatePassword(req.Password) {
		middleware.ErrorResponse(c, &utils.AppError{Message: "password must be at least 8 characters long", StatusCode: http.StatusBadRequest})
		return
	}

	err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		logger.Warn("User registration failed",
			zap.String("email", req.Email),
			zap.String("error", err.Error()),
		)
		middleware.ErrorResponse(c, err)
		return
	}

	logger.Info("User registered successfully",
		zap.String("email", req.Email),
	)

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, utils.ErrInvalidRequest)
		return
	}

	accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		logger.Warn("Login failed",
			zap.String("email", req.Email),
			zap.String("error", err.Error()),
		)
		middleware.ErrorResponse(c, err)
		return
	}

	logger.Info("User logged in successfully",
		zap.String("email", req.Email),
	)

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

