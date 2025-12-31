package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/randhir/aegis-core/internal/logger"
	"github.com/randhir/aegis-core/internal/utils"
	"go.uber.org/zap"
)

// ErrorHandler is a middleware that handles errors consistently
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			appErr := utils.ToAppError(err)

			// Log error (but never log sensitive data)
			logger.Error("Request failed",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Int("status", appErr.StatusCode),
				zap.String("error", appErr.Message),
			)

			// Return standardized error response
			c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			c.Abort()
		}
	}
}

// ErrorResponse sends a standardized error response
func ErrorResponse(c *gin.Context, err error) {
	appErr := utils.ToAppError(err)
	c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
}

