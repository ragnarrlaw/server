package userservice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	userrepository "github.com/raganrrlaw/server/internal/services/user_service/user_repository"
	"github.com/raganrrlaw/server/internal/types/search"
)

type UserService struct {
	userRepository userrepository.UserRepository
}

func NewUserService(repo userrepository.UserRepository) *UserService {
	return &UserService{
		userRepository: repo,
	}
}

func (us *UserService) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	users, err := us.userRepository.Get(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalUsers, err := us.userRepository.GetTotalNumberOfUsers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if b, err := json.Marshal(search.Payload{
		Data: users,
		Meta: search.Meta{
			Total: totalUsers,
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func (us *UserService) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	/**
	  TOOD: ADD VALIDATIONS TO CHECK WHETHER THE USER EXISTS
	        IN HERE CHECK OF THE NO ROWS FOUND ERROR TO SEND
	        NO RESOURCE FOUND ERROR
	*/

	if userId := r.PathValue("id"); userId != "" {
		user, err := us.userRepository.GetById(ctx, userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if b, err := json.Marshal(user); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(b)
			}
		}
	} else {
		http.Error(w, "user is required", http.StatusBadRequest)
	}
}
