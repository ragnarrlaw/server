package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raganrrlaw/server/internal/types"
	"github.com/rs/cors"
)

type Middleware func(http.Handler) http.Handler

func AppendMiddlewareStack(fs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(fs) - 1; i >= 0; i-- {
			next = fs[i](next)
		}
		return next
	}
}

func LogRequestDetailsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		method := r.Method
		endpoint := r.URL.Path

		var coloredMethod string
		switch method {
		case http.MethodGet:
			coloredMethod = color.New(color.BgGreen, color.FgBlack).Sprint(method)
		case http.MethodPost:
			coloredMethod = color.New(color.BgBlue, color.FgWhite).Sprint(method)
		case http.MethodPut:
			coloredMethod = color.New(color.BgYellow, color.FgBlack).Sprint(method)
		case http.MethodDelete:
			coloredMethod = color.New(color.BgRed, color.FgWhite).Sprint(method)
		default:
			coloredMethod = method
		}

		log.Printf("%s\t%s\t%s\n", clientIP, coloredMethod, endpoint)
		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		// AllowCredentials: false, // set this up later
	})
	return c.Handler(next)
}

// ValidateAccessTokens checks the validity of the access token
// P.S. - doesn't validate the user credentials and the claims
// in the token, this should be handled by the handler in service layer
func ValidateAccessTokens(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Expecting the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		accessToken := parts[1]
		claims := &types.TokenClaims{}
		token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_TOKEN_SECRET")), nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to the context
		ctx := context.WithValue(r.Context(), types.UserIDKey, claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ValidateRefreshTokens checks the validity of the Refresh token
// P.S. - doesn't validate the user credentials and the claims
// in the token, this should be handled by the handler in service layer
func ValidateRefreshTokens(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "Refresh token is required", http.StatusUnauthorized)
			return
		}

		claims := &types.TokenClaims{}
		token, err := jwt.ParseWithClaims(refreshToken.Value, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_TOKEN_SECRET")), nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), types.UserIDKey, claims.UserId)
		ctx = context.WithValue(ctx, refreshToken, refreshToken.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid content type. Only application/json is allowed", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}
