package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/randhir/aegis-core/internal/logger"
	"github.com/randhir/aegis-core/internal/middleware"
	"github.com/randhir/aegis-core/internal/service"
	"github.com/randhir/aegis-core/internal/utils"
	"go.uber.org/zap"
)

type TokenHandler struct {
	tokenService *service.TokenService
}

func NewTokenHandler(tokenService *service.TokenService) *TokenHandler {
	return &TokenHandler{
		tokenService: tokenService,
	}
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *TokenHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, utils.ErrInvalidRequest)
		return
	}

	accessToken, refreshToken, err := h.tokenService.Refresh(req.RefreshToken)
	if err != nil {
		logger.Warn("Token refresh failed",
			zap.String("error", err.Error()),
		)
		middleware.ErrorResponse(c, err)
		return
	}

	logger.Info("Token refreshed successfully")
	c.JSON(http.StatusOK, RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *TokenHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, utils.ErrInvalidRequest)
		return
	}

	// Extract access token from Authorization header if present
	accessToken := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			accessToken = parts[1]
		}
	}

	err := h.tokenService.Logout(req.RefreshToken, accessToken)
	if err != nil {
		logger.Warn("Logout failed",
			zap.String("error", err.Error()),
		)
		middleware.ErrorResponse(c, err)
		return
	}

	logger.Info("User logged out successfully")
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

