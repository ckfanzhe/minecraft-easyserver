package handlers

import (
	"minecraft-easyserver/models"
	"minecraft-easyserver/services"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *services.AuthService
	rateLimiter *services.RateLimiterService
}

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler() *AuthHandler {
	rateLimiter := services.NewRateLimiterService()
	rateLimiter.StartCleanupRoutine()
	return &AuthHandler{
		authService: services.NewAuthService(),
		rateLimiter: rateLimiter,
	}
}

// getSecureClientIP returns the real client IP address securely
// It only uses the direct connection IP to prevent IP spoofing
func (h *AuthHandler) getSecureClientIP(c *gin.Context) string {
	// Get the remote address from the connection
	remoteAddr := c.Request.RemoteAddr
	
	// Extract IP from "IP:port" format
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return host
	}
	
	// Fallback: remove port if SplitHostPort fails
	if idx := strings.LastIndex(remoteAddr, ":"); idx != -1 {
		return remoteAddr[:idx]
	}
	
	return remoteAddr
}

// Login handles login requests
func (h *AuthHandler) Login(c *gin.Context) {
	// Use secure IP detection to prevent spoofing
	clientIP := h.getSecureClientIP(c)

	// Check if IP is blocked first - reject all requests during block period
	if h.rateLimiter.IsBlocked(clientIP) {
		blockTime := h.rateLimiter.GetBlockTimeRemaining(clientIP)
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many failed login attempts. Please try again later.",
			"blocked_until": time.Now().Add(blockTime).Format(time.RFC3339),
			"retry_after_seconds": int(blockTime.Seconds()),
		})
		return
	}

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Try authentication
	response, err := h.authService.Login(req.Password)
	if err != nil {
		// Record failed attempt
		h.rateLimiter.RecordFailedAttempt(clientIP)
		remainingAttempts := h.rateLimiter.GetRemainingAttempts(clientIP)
		
		errorResponse := gin.H{
			"error": err.Error(),
		}
		
		// Add rate limiting info if user is getting close to limit
		if remainingAttempts <= 2 {
			errorResponse["remaining_attempts"] = remainingAttempts
			if remainingAttempts == 0 {
				errorResponse["warning"] = "Account will be temporarily blocked after next failed attempt"
			}
		}
		
		c.JSON(http.StatusUnauthorized, errorResponse)
		return
	}

	// Record successful attempt (clears failed attempts and any blocks)
	h.rateLimiter.RecordSuccessfulAttempt(clientIP)
	c.JSON(http.StatusOK, response)
}

// ChangePassword handles password change requests
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	response, err := h.authService.ChangePassword(req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !response.Success {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, response)
}