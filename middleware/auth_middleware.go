package middleware

import (
	"minecraft-easyserver/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware creates a JWT authentication middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Use 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token required",
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := authService.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Store claims in context for use in handlers
		c.Set("claims", claims)
		c.Next()
	}
}

// WebSocketAuthMiddleware creates a WebSocket authentication middleware
// that validates JWT token from query parameter
func WebSocketAuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return func(c *gin.Context) {
		// For WebSocket connections, get token from query parameter
		tokenString := c.Query("token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token required in query parameter",
			})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := authService.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Store claims in context for use in handlers
		c.Set("claims", claims)
		c.Next()
	}
}