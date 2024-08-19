package authservice

import "net/http"

func (as *AuthService) RegisterRoutes(router *http.ServeMux) {
  router.HandleFunc("POST /login", as.LoginHandler)
  router.HandleFunc("POST /register", as.SignupHandler)
  router.HandleFunc("POST /logout", as.LogoutHandler)
}
