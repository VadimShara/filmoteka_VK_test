package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Service struct {
	cfg *Config
}

func New(cfg *Config) *Service {
	return &Service{
		cfg: cfg,
	}
}

func (s *Service) CreateToken(ctx context.Context, user UserData) (Token, error) {
	tokenData := TokenData{
		UserData:  user,
		TokenType: AccessToken,
	}

	token, err := s.generateToken(tokenData)
	if err != nil {
		return Token{}, err
	}

	return Token{AccessToken: token}, nil
}

func (s *Service) ValidateToken(tokenString string) (*UserData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role := fmt.Sprint(claims["role"])
		username := fmt.Sprint(claims["username"])
		user := &UserData{
			Username: username,
			Role:     role,
		}

		return user, nil
	}

	return nil, ErrInvalidToken
}

func (s *Service) generateToken(tokenData TokenData) (string, error) {
	claims := &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(s.cfg.AccessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ID:        uuid.New().String(),
		},
		User:      tokenData.Username,
		Role:      tokenData.Role,
		TokenType: tokenData.TokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}
