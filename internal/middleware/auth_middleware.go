package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/randhir/aegis-core/internal/cache"
	"github.com/randhir/aegis-core/internal/logger"
	"github.com/randhir/aegis-core/internal/utils"
	"go.uber.org/zap"
)

type AuthContext struct {
	UserID string
	Email  string
	Role   string
}

const AuthContextKey = "auth_context"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Check if token is blacklisted in Redis
		isBlacklisted, err := cache.IsAccessTokenBlacklisted(tokenString)
		if err != nil || isBlacklisted {
			logger.Warn("Authorization failed: token blacklisted or invalid",
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			logger.Warn("Authorization failed: invalid token",
				zap.String("path", c.Request.URL.Path),
				zap.String("error", err.Error()),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		authContext := AuthContext{
			UserID: claims.UserID,
			Email:  claims.Email,
			Role:   claims.Role,
		}

		c.Set(AuthContextKey, authContext)
		c.Next()
	}
}

func GetAuthContext(c *gin.Context) (*AuthContext, bool) {
	authCtx, exists := c.Get(AuthContextKey)
	if !exists {
		return nil, false
	}

	authContext, ok := authCtx.(AuthContext)
	if !ok {
		return nil, false
	}

	return &authContext, true
}

