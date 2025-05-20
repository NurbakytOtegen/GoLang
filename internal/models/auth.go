package models

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json: "role" binding: "required"`
}

// AuthResponse represents the authentication response with token
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
