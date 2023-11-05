package domain

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type ContextKey string

const (
	ContextKeyAuth  ContextKey = "authentication"
	ContextKeyToken ContextKey = "token"
)

type AuthConfig struct {
	ModelPath          string
	PolicyPath         string
	AccessTokenExpire  time.Duration
	RefreshTokenExpire time.Duration
	Secret             string
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Pwd   string `json:"p[assword]" validate:"required"`
}

type RefreshTokenRequest struct {
	AccessToken  string `json:"-" validate:"required,jwt"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	AccessToken  string `json:"access_token" validate:"required,jwt"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type AuthenticationResponse struct {
	User *User
}

type AuthorizationResponse struct {
	GrantAccess bool `json:"grant_access"`
}

type UserClaims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}
