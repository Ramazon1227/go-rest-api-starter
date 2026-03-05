package models

import "time"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	User      *User     `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name,omitempty" example:"John Doe"`
	Email    string `json:"email,omitempty" example:"john@example.com"`
	Phone    string `json:"phone,omitempty" example:"+1234567890"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" example:"oldpassword123"`
	NewPassword     string `json:"new_password" binding:"required" example:"newpassword123"`
}
