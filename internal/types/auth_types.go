package types

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// contextKey is a type used for context keys to avoid conflicts
type ContextKey string

const (
	UserIDKey       ContextKey = "user_id"
	RefreshTokenKey ContextKey = "refresh_token"
	AccessTokenKey  ContextKey = "access_token"
)

type Token struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

func (t *Token) String() string {
	return fmt.Sprintf("Token: { Type: %s, AccessToken: %s }", t.Type, t.AccessToken)
}

type TokenClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func (tc *TokenClaims) String() string {
	return fmt.Sprintf("TokenClaims: { UserId: %s }", tc.UserId)
}

type AuthToken struct {
	Id     uuid.UUID `json:"id"`
	UserId string    `json:"user_id"`
	Token  string    `json:"token"`
}

func (at *AuthToken) String() string {
	return fmt.Sprintf("AuthToken: { Id: %s, UserId: %s, Token: %s }", at.Id, at.UserId, at.Token)
}
