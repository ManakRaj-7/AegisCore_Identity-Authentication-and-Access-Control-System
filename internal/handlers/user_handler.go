package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/randhir/aegis-core/internal/logger"
	"github.com/randhir/aegis-core/internal/middleware"
	"github.com/randhir/aegis-core/internal/repository"
	"github.com/randhir/aegis-core/internal/utils"
	"go.uber.org/zap"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

type ProfileResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserListResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	authContext, exists := middleware.GetAuthContext(c)
	if !exists {
		middleware.ErrorResponse(c, utils.ErrUnauthorized)
		return
	}

	logger.Info("Profile accessed",
		zap.String("user_id", authContext.UserID),
		zap.String("email", authContext.Email),
	)

	c.JSON(http.StatusOK, ProfileResponse{
		ID:    authContext.UserID,
		Email: authContext.Email,
		Role:  authContext.Role,
	})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	authContext, exists := middleware.GetAuthContext(c)
	if !exists {
		middleware.ErrorResponse(c, utils.ErrUnauthorized)
		return
	}

	users, err := repository.ListUsers()
	if err != nil {
		logger.Error("Failed to fetch users",
			zap.String("admin_id", authContext.UserID),
			zap.Error(err),
		)
		middleware.ErrorResponse(c, utils.ErrInternalError)
		return
	}

	logger.Info("Users list accessed",
		zap.String("admin_id", authContext.UserID),
		zap.Int("user_count", len(users)),
	)

	response := make([]UserListResponse, len(users))
	for i, user := range users {
		response[i] = UserListResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, response)
}

