package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

const (
	AccessToken TokenType = "access_token"
)

type (
	CustomClaims struct {
		jwt.RegisteredClaims
		User      string    `json:"user"`
		Role      string    `json:"role"`
		TokenType TokenType `json:"token_type"`
	}

	TokenType string

	TokenData struct {
		UserData
		TokenType TokenType
	}

	UserData struct {
		Username string
		Role     string
	}

	Token struct {
		AccessToken string `json:"access_token"`
	}
)
