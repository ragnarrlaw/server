package userservice

import (
	"context"
	"encoding/json"
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
func (us *UserService) RemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	/**
	  TODO: WHETHER THE USER EXISTS IN A MIDDLEWARE AND THE PERMISSION LEVEL OF THE USER
	        USERS CAN ONLY DELETE THEIR OWN PROFILES, ADMINS (IF THERE ARE ANY) THEY CAN DELETE
	        USERS AS WELL
	*/
	/** TODO: PAYLOAD CAN BE USED IF NEEDED */

	if userUserWithId := r.PathValue("id"); userUserWithId != "" {
		if err := us.userRepository.Remove(ctx, []string{userUserWithId}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	} else {
		http.Error(w, "user id is required", http.StatusBadRequest)
	}
}

/*
*

	Update the user fields such as username, first name, last name
*/
func (us *UserService) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := &types.UserUpdatePayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/** TODO: PAYLOAD VALIDATION */
	/**
	  TODO: CHECK IF USER EXISTS IN A MIDDLEWARE
	        CHECK WHETH THE ID IN THE CONTEXT MATCHES THE ID GIVEN IN THE URL PARAMETER
	*/

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if userId := r.PathValue("id"); userId != "" {
		user, err := us.userRepository.Update(ctx, userId, payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "user id is required", http.StatusBadRequest)
	}
}

func (us *UserService) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	users, err := us.userRepository.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if b, err := json.Marshal(users); err != nil {
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
