package userservice

import (
	"net/http"

	"github.com/raganrrlaw/server/internal/services/user_service/user_repository"
)

type UserService struct {
  userRepository *userrepository.UserRepository
}

func NewUserService(repo *userrepository.UserRepository) *UserService {
  return &UserService{
    userRepository: repo,
  }
}

func (us *UserService) AddUserHandler(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) RemoveUserHandler(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) GetUserHandler(w http.ResponseWriter, r *http.Request) {}
