package services

import (
	"Cars/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// JWT secret key - in production, this should be stored in environment variables
var jwtSecret = []byte("your_jwt_secret_key_here")

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a JWT token for a user
func GenerateToken(user models.User) (string, error) {
	// Set token expiration time (e.g., 24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims with user ID and standard claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token with claims and sign with secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID
func ValidateToken(tokenString string) (uint, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get user ID from claims
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

// RegisterUser registers a new user
func RegisterUser(registerRequest models.RegisterRequest) (*models.User, error) {
	// Check if user with the same email already exists
	var existingUser models.User
	result := DB.Where("email = ?", registerRequest.Email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("user with this email already exists")
	} else if result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	// Hash the password
	hashedPassword, err := HashPassword(registerRequest.Password)
	if err != nil {
		return nil, err
	}

	role := registerRequest.Role
	if role == "" {
		role = "USER"
	}

	// Create new user
	user := models.User{
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: hashedPassword,
		Role:     role,
	}

	// Save user to database
	if err := DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// LoginUser authenticates a user and returns a token
func LoginUser(loginRequest models.LoginRequest) (*models.User, string, error) {
	// Find user by email
	var user models.User
	result := DB.Where("email = ?", loginRequest.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, "", errors.New("user not found")
		}
		return nil, "", result.Error
	}

	// Check password
	if !CheckPasswordHash(loginRequest.Password, user.Password) {
		return nil, "", errors.New("invalid password")
	}

	// Generate token
	token, err := GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := DB.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
