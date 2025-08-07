package services

import (
	"errors"
	"minecraft-easyserver/config"
	"minecraft-easyserver/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles authentication operations
type AuthService struct{}

// NewAuthService creates a new auth service instance
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Login validates password and generates JWT token
func (s *AuthService) Login(password string) (*models.LoginResponse, error) {
	// Validate password
	if password != config.AppConfig.Auth.Password {
		return nil, errors.New("invalid password")
	}

	// Generate JWT token
	token, err := s.generateJWT()
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:   token,
		Message: "Login successful",
	}, nil
}

// generateJWT creates a new JWT token
func (s *AuthService) generateJWT() (string, error) {
	now := time.Now()
	claims := &models.JWTClaims{
		Authorized: true,
		Exp:        now.Add(24 * time.Hour).Unix(), // Token expires in 24 hours
		Iat:        now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": claims.Authorized,
		"exp":        claims.Exp,
		"iat":        claims.Iat,
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.Auth.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token
func (s *AuthService) ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.AppConfig.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if token is expired
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("token expired")
			}
		}

		return &models.JWTClaims{
			Authorized: claims["authorized"].(bool),
			Exp:        int64(claims["exp"].(float64)),
			Iat:        int64(claims["iat"].(float64)),
		}, nil
	}

	return nil, errors.New("invalid token")
}