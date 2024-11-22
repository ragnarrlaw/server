package authservice

import (
	"net/http"

	"github.com/raganrrlaw/server/internal/middleware"
)

func (as *AuthService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /login", as.LoginHandler)
	router.HandleFunc("POST /user/register", as.UserSignUpHandler)
	router.HandleFunc("POST /store/register", as.StoreSignUpHandler)
	router.Handle("POST /logout",
		middleware.ValidateAccessTokens(
			middleware.ValidateRefreshTokens(
				http.HandlerFunc(as.LogoutHandler),
			),
		),
	)
	router.Handle("POST /refresh", middleware.ValidateRefreshTokens(
		http.HandlerFunc(as.RefreshAccessTokenHandler)),
	)
}
