package types

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// contextKey is a type used for context keys to avoid conflicts
type ContextKey string

const (
	IDKey           ContextKey = "_id"
	RefreshTokenKey ContextKey = "refresh_token"
	AccessTokenKey  ContextKey = "access_token"
	RoleKey         ContextKey = "_role"
)

type EntityType string

const (
	UserEntity  EntityType = "user"
	StoreEntity EntityType = "store"
)

type LoginPayload struct {
	Username string     `json:"username"`
	Password string     `json:"password"`
	Role     EntityType `json:"role"`
}

type Token struct {
	Type  string `json:"type"`
	Token string `json:"access_token"`
}

func (t *Token) String() string {
	return fmt.Sprintf("Token: { Type: %s, AccessToken: %s }", t.Type, t.Token)
}

type TokenClaims struct {
	Id   string     `json:"user_id"`
	Role EntityType `json:"role"`
	jwt.RegisteredClaims
}

func (tc *TokenClaims) String() string {
	return fmt.Sprintf("TokenClaims: { UserId: %s, Role: %s }", tc.Id, tc.Role)
}

type AuthToken struct {
	Id     uuid.UUID  `json:"id"`
	UserId uuid.UUID  `json:"user_id"`
	Role   EntityType `json:"role"`
	Token  string     `json:"token"`
}

func (at *AuthToken) String() string {
	return fmt.Sprintf("AuthToken: { Id: %s, UserId: %s, Role: %s, Token: %s }", at.Id, at.UserId, at.Role, at.Token)
}
