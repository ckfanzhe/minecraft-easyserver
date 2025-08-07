package services

import (
	"sync"
	"time"
)

// LoginAttempt represents a login attempt record
type LoginAttempt struct {
	Count     int
	LastAttempt time.Time
	BlockedUntil time.Time
}

// RateLimiterService handles rate limiting for login attempts
type RateLimiterService struct {
	attempts map[string]*LoginAttempt
	mutex    sync.RWMutex
	// Configuration
	maxAttempts   int           // Maximum failed attempts before blocking
	blockDuration time.Duration // How long to block after max attempts
	windowDuration time.Duration // Time window for counting attempts
}

// NewRateLimiterService creates a new rate limiter service
func NewRateLimiterService() *RateLimiterService {
	return &RateLimiterService{
		attempts:      make(map[string]*LoginAttempt),
		maxAttempts:   5,               // 5 failed attempts
		blockDuration: 5 * time.Minute, // Block for 5 minutes
		windowDuration: 5 * time.Minute, // 5 minute window
	}
}

// IsBlocked checks if an IP is currently blocked
func (r *RateLimiterService) IsBlocked(ip string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	attempt, exists := r.attempts[ip]
	if !exists {
		return false
	}

	// Check if block period has expired
	if time.Now().After(attempt.BlockedUntil) {
		return false
	}

	return attempt.Count >= r.maxAttempts && !attempt.BlockedUntil.IsZero()
}

// RecordFailedAttempt records a failed login attempt
func (r *RateLimiterService) RecordFailedAttempt(ip string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	attempt, exists := r.attempts[ip]

	if !exists {
		// First attempt
		r.attempts[ip] = &LoginAttempt{
			Count:       1,
			LastAttempt: now,
			BlockedUntil: time.Time{},
		}
		return
	}

	// Check if we're outside the time window, reset counter
	if now.Sub(attempt.LastAttempt) > r.windowDuration {
		attempt.Count = 1
		attempt.LastAttempt = now
		attempt.BlockedUntil = time.Time{}
		return
	}

	// Increment counter
	attempt.Count++
	attempt.LastAttempt = now

	// Block if max attempts reached
	if attempt.Count >= r.maxAttempts {
		attempt.BlockedUntil = now.Add(r.blockDuration)
	}
}

// RecordSuccessfulAttempt clears the failed attempts for an IP
func (r *RateLimiterService) RecordSuccessfulAttempt(ip string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Clear all records for this IP, including any block
	delete(r.attempts, ip)
}

// GetRemainingAttempts returns the number of remaining attempts before blocking
func (r *RateLimiterService) GetRemainingAttempts(ip string) int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	attempt, exists := r.attempts[ip]
	if !exists {
		return r.maxAttempts
	}

	// Check if we're outside the time window or block has expired
	now := time.Now()
	if now.Sub(attempt.LastAttempt) > r.windowDuration || now.After(attempt.BlockedUntil) {
		return r.maxAttempts
	}

	remaining := r.maxAttempts - attempt.Count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetBlockTimeRemaining returns the remaining block time for an IP
func (r *RateLimiterService) GetBlockTimeRemaining(ip string) time.Duration {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	attempt, exists := r.attempts[ip]
	if !exists {
		return 0
	}

	if time.Now().After(attempt.BlockedUntil) {
		return 0
	}

	return attempt.BlockedUntil.Sub(time.Now())
}

// CleanupExpiredEntries removes expired entries to prevent memory leaks
func (r *RateLimiterService) CleanupExpiredEntries() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	for ip, attempt := range r.attempts {
		// Remove entries that are outside the window and not blocked
		if now.Sub(attempt.LastAttempt) > r.windowDuration && now.After(attempt.BlockedUntil) {
			delete(r.attempts, ip)
		}
	}
}

// StartCleanupRoutine starts a background routine to clean up expired entries
func (r *RateLimiterService) StartCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute) // Cleanup every 10 minutes
		defer ticker.Stop()

		for range ticker.C {
			r.CleanupExpiredEntries()
		}
	}()
}