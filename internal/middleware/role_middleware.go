package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/randhir/aegis-core/internal/logger"
	"go.uber.org/zap"
)

func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authContext, exists := GetAuthContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if authContext.Role != requiredRole {
			logger.Warn("Authorization failed: insufficient permissions",
				zap.String("user_id", authContext.UserID),
				zap.String("required_role", requiredRole),
				zap.String("user_role", authContext.Role),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

