package services

import (
	"net/http"

	"github.com/raganrrlaw/server/database"
	"github.com/raganrrlaw/server/repositories"
)

type UserService struct {
  repo *repositories.UserRepo
}

func NewUserService(store *database.DBPool) *UserService {
  return &UserService{
    repo: repositories.NewUserRepo(store),
  }
}

func (s *UserService) RegisterRoutes(r *http.ServeMux) {
  r.HandleFunc("GET /user", s.GetAllUsersHandler)
  r.HandleFunc("GET /user/{id}", s.GetUserHandler)
  r.HandleFunc("POST /user", s.CreateUserHandler)
  r.HandleFunc("PATCH /user/{id}", s.UpdateUserHandler)
  r.HandleFunc("DELETE /user/{id}", s.UpdateUserHandler)
}

func (s *UserService) GetAllUsersHandler (w http.ResponseWriter, r *http.Request) {
}

func (s *UserService) GetUserHandler (w http.ResponseWriter, r *http.Request) {
}

func (s *UserService) CreateUserHandler (w http.ResponseWriter, r *http.Request) {
}

func (s *UserService) UpdateUserHandler (w http.ResponseWriter, r *http.Request) {
}

func (s *UserService) DeleteUserHandler (w http.ResponseWriter, r *http.Request) {
}
