package dto

import "github.com/nantapop-kj/go-auth-service/models"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=72"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type TokenResponse struct {
	AccessToken      string       `json:"access_token"`
	TokenType        string       `json:"token_type"`
	ExpiresIn        int64        `json:"expires_in"`
	RefreshToken     string       `json:"refresh_token"`
	RefreshExpiresIn int64        `json:"refresh_expires_in"`
	User             UserResponse `json:"user"`
}

func NewUserResponse(u *models.Users) UserResponse {
	return UserResponse{ID: u.ID.String(), Email: u.Email}
}
