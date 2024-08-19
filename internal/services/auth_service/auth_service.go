package authservice

import (
	"net/http"

	authrepository "github.com/raganrrlaw/server/internal/services/auth_service/auth_repository"
)

type AuthService struct {
  repo *authrepository.AuthRepository
}

func NewAuthService(repo *authrepository.AuthRepository) *AuthService {
  return &AuthService{
    repo: repo,
  }
}

func (as *AuthService) LoginHandler(w http.ResponseWriter, r *http.Request) {
}

func (as *AuthService) SignupHandler(w http.ResponseWriter, r *http.Request) {
}

func (as *AuthService) LogoutHandler(w http.ResponseWriter, r * http.Request) {
}


