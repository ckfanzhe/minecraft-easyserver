package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"minecraft-easyserver/config"
	"minecraft-easyserver/models"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"
)

// AuthService handles authentication operations
type AuthService struct{}

// NewAuthService creates a new auth service instance
func NewAuthService() *AuthService {
	return &AuthService{}
}

// hashPassword creates SHA256 hash of password
func (s *AuthService) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Login validates password and generates JWT token
func (s *AuthService) Login(password string) (*models.LoginResponse, error) {
	// Hash the input password and compare with stored hash
	hashedPassword := s.hashPassword(password)
	if hashedPassword != config.AppConfig.Auth.Password {
		return nil, errors.New("invalid password")
	}

	// Generate JWT token
	token, err := s.generateJWT()
	if err != nil {
		return nil, err
	}

	// Check if using default password (compare with hash of "admin123")
	defaultPasswordHash := s.hashPassword("admin123")
	requirePasswordChange := hashedPassword == defaultPasswordHash

	return &models.LoginResponse{
		Token:                 token,
		Message:               "Login successful",
		RequirePasswordChange: requirePasswordChange,
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

// ValidatePasswordStrength validates password strength
func (s *AuthService) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度至少需要8位")
	}

	// Check for uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return errors.New("密码必须包含至少一个大写字母")
	}

	// Check for lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return errors.New("密码必须包含至少一个小写字母")
	}

	// Check for digit
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		return errors.New("密码必须包含至少一个数字")
	}

	// Check for special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`).MatchString(password)
	if !hasSpecial {
		return errors.New("密码必须包含至少一个特殊字符")
	}

	return nil
}

// ChangePassword changes the user password
func (s *AuthService) ChangePassword(currentPassword, newPassword string) (*models.ChangePasswordResponse, error) {
	// Validate current password by comparing hashes
	currentPasswordHash := s.hashPassword(currentPassword)
	if currentPasswordHash != config.AppConfig.Auth.Password {
		return &models.ChangePasswordResponse{
			Message: "当前密码不正确",
			Success: false,
		}, nil
	}

	// Validate new password strength
	if err := s.ValidatePasswordStrength(newPassword); err != nil {
		return &models.ChangePasswordResponse{
			Message: err.Error(),
			Success: false,
		}, nil
	}

	// Update password in config (store as hash)
	if err := s.updatePasswordInConfig(newPassword); err != nil {
		return &models.ChangePasswordResponse{
			Message: "更新密码失败: " + err.Error(),
			Success: false,
		}, nil
	}

	return &models.ChangePasswordResponse{
		Message: "密码修改成功",
		Success: true,
	}, nil
}

// updatePasswordInConfig updates password in config file
func (s *AuthService) updatePasswordInConfig(newPassword string) error {
	// Hash the new password before storing
	hashedPassword := s.hashPassword(newPassword)

	// Read current config file
	configPath := "config.yml"
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse config
	var configData map[string]interface{}
	if err := yaml.Unmarshal(data, &configData); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// Update password with hash
	if auth, ok := configData["auth"].(map[string]interface{}); ok {
		auth["password"] = hashedPassword
	} else {
		return fmt.Errorf("auth section not found in config")
	}

	// Write back to file
	updatedData, err := yaml.Marshal(configData)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	// Update in-memory config with hash
	config.AppConfig.Auth.Password = hashedPassword

	return nil
}