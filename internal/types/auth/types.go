package auth

type ContextKey string

const (
	IDKey           ContextKey = "id"
	RefreshTokenKey ContextKey = "refresh_token"
	AccessTokenKey  ContextKey = "access_token"
	RoleKey         ContextKey = "role"
)
