package authservice

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authRepository "github.com/raganrrlaw/server/internal/services/auth_service/auth_repository"
	userRepository "github.com/raganrrlaw/server/internal/services/user_service/user_repository"
	"github.com/raganrrlaw/server/internal/types"
	"github.com/raganrrlaw/server/internal/utils"
)

type AuthService struct {
	userRepo userRepository.UserRepository
	authRepo authRepository.AuthRepository
}

func NewAuthService(authRepo authRepository.AuthRepository, userRepo userRepository.UserRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (as *AuthService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.UserLoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	} else {
		user, err := as.userRepo.GetBy(r.Context(), "username", payload.Username)
		if err != nil {
			if err.Error() == "no rows found" {
				http.Error(w, "No such user found", http.StatusUnauthorized)
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if ok := utils.CheckPasswordHash(payload.Password, user.Password); !ok {
			http.Error(w, "Invalid password", http.StatusBadRequest)
			return
		}

		// TODO: read the env and set the token duration

		access_token := generateToken(user, time.Minute*15)
		if access_token == "" {
			http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
			return
		}
		refresh_token := generateToken(user, time.Hour*30)
		if refresh_token == "" {
			http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
			return
		}

		if _, err := as.authRepo.AddToken(r.Context(), &types.AuthToken{
			UserId: user.Id.String(),
			Token:  refresh_token,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: read the env and set the token duration

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			HttpOnly: true,
			Value:    refresh_token,
			Expires:  time.Now().Add(time.Minute * 30),
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(struct {
			Type        string `json:"type"`
			AccessToken string `json:"access_token"`
		}{
			Type:        "Bearer",
			AccessToken: access_token,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (as *AuthService) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.UserSignUpPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if digest, err := utils.HashPassword(payload.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		payload.Password = digest
		user, err := as.userRepo.Add(r.Context(), &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(user); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func (as *AuthService) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(types.UserIDKey).(string)

	if err := as.authRepo.RemoveTokensOfUser(r.Context(), userId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err := as.userRepo.GetById(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now().Add(-time.Minute),
	})
	w.WriteHeader(http.StatusNoContent)
}

func (as *AuthService) RefreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(types.UserIDKey).(string)
	if userId == "" {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	user, err := as.userRepo.GetById(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	access_token := generateToken(user, time.Minute*15)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(struct {
		Type        string `json:"type"`
		AccessToken string `json:"access_token"`
	}{
		Type:        "Bearer",
		AccessToken: access_token,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func generateToken(user *types.User, duration time.Duration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": user.Id,
		"sub":     user.Username,
		"exp":     time.Now().Add(duration).Unix(),
	})
	if tokenString, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN_SECRET"))); err != nil {
		log.Printf(">>>> Error: %s\n", err.Error())
		return ""
	} else {
		return tokenString
	}
}
