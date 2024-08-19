package userservice

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	userrepository "github.com/raganrrlaw/server/internal/services/user_service/user_repository"
	"github.com/raganrrlaw/server/internal/types"
	"github.com/raganrrlaw/server/internal/utils"
)

type UserService struct {
	userRepository userrepository.UserRepository
}

func NewUserService(repo userrepository.UserRepository) *UserService {
	return &UserService{
		userRepository: repo,
	}
}

/** Used mainly for testing purposes */
func (us *UserService) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := &types.UserSignUpPayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	/** TODO: PAYLOAD VALIDATION */
	/** TODO: GENERIC RESPONSE ERROR STRUCTURE ORGANIZATION */

	if digest, err := utils.HashPassword(payload.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		payload.Password = digest
		user, err := us.userRepository.Add(ctx, payload)
		log.Printf(">>>> Created User: %s\n", user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if b, err := json.Marshal(user); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(b)
			}
		}
	}
}

/** Used mainly for testing purposes */
func (us *UserService) RemoveUserHandler(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) GetUserHandler(w http.ResponseWriter, r *http.Request) {}
