package authservice

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authRepository "github.com/raganrrlaw/server/internal/services/auth_service/auth_repository"
	storeRepository "github.com/raganrrlaw/server/internal/services/store_service/store_repository"
	userRepository "github.com/raganrrlaw/server/internal/services/user_service/user_repository"
	"github.com/raganrrlaw/server/internal/types"
	"github.com/raganrrlaw/server/internal/utils"
)

type AuthService struct {
	userRepo  userRepository.UserRepository
	storeRepo storeRepository.StoreRepository
	authRepo  authRepository.AuthRepository
}

func NewAuthService(authRepo authRepository.AuthRepository, userRepo userRepository.UserRepository, storeRepo storeRepository.StoreRepository) *AuthService {
	return &AuthService{
		authRepo:  authRepo,
		userRepo:  userRepo,
		storeRepo: storeRepo,
	}
}

func (as *AuthService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	} else {
		switch payload.Role {
		case types.UserEntity:
			{
				if user, err := as.userRepo.GetBy(r.Context(), "username", payload.Username); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					if ok := utils.CheckPasswordHash(user.Password, payload.Password); !ok {
						http.Error(w, "Invalid credentials", http.StatusUnauthorized)
					} else {
						access_token := generateToken(&user, time.Minute*15, types.UserEntity)
						refresh_token := generateToken(&user, time.Hour*24*7, types.UserEntity)
						if _, err := as.authRepo.AddToken(r.Context(), &types.AuthToken{
							UserId: user.Id,
							Role:   types.UserEntity,
							Token:  refresh_token,
						}, types.UserEntity); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						} else {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusCreated)
							http.SetCookie(w, &http.Cookie{
								Name:     "token",
								Value:    refresh_token,
								HttpOnly: false,
								Secure:   false,
								Path:     "/",
								Expires:  time.Now().Add(time.Hour * 24 * 7),
							})
							encoder := json.NewEncoder(w)
							if err := encoder.Encode(types.Token{
								Type:  "Bearer",
								Token: access_token,
							}); err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
							}
						}
					}
				}
			}
		case types.StoreEntity:
			{
				if store, err := as.storeRepo.GetBy(r.Context(), "store_username", payload.Username); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					if ok := utils.CheckPasswordHash(store.Password, payload.Password); !ok {
						http.Error(w, "Invalid credentials", http.StatusUnauthorized)
					} else {
						access_token := generateToken(&store, time.Minute*15, types.StoreEntity)
						refresh_token := generateToken(&store, time.Hour*24*7, types.StoreEntity)
						if _, err := as.authRepo.AddToken(r.Context(), &types.AuthToken{
							UserId: store.Id,
							Role:   types.StoreEntity,
							Token:  refresh_token,
						}, types.StoreEntity); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						} else {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusCreated)
							http.SetCookie(w, &http.Cookie{
								Name:     "refresh_token",
								Value:    refresh_token,
								HttpOnly: true,
								Expires:  time.Now().Add(time.Hour * 24 * 7),
							})
							encoder := json.NewEncoder(w)
							if err := encoder.Encode(types.Token{
								Type:  "Bearer",
								Token: access_token,
							}); err != nil {
								http.Error(w, err.Error(), http.StatusInternalServerError)
							}
						}
					}
				}
			}
		default:
			{
				http.Error(w, "Invalid role", http.StatusBadRequest)
			}
		}
	}
}

// TODO: Add the user role here
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

	id := r.Context().Value(types.IDKey).(string)
	role := r.Context().Value(types.RoleKey).(types.EntityType)

	if err := as.authRepo.RemoveTokensOfUser(r.Context(), id, role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err := as.userRepo.GetById(r.Context(), id)
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
	id := r.Context().Value(types.IDKey).(string)
	if id == "" {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
	} else {
		role := r.Context().Value(types.RoleKey).(types.EntityType)
		switch role {
		case types.UserEntity:
			{
				user, err := as.userRepo.GetById(r.Context(), id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					access_token := generateToken(&user, time.Minute*15, types.UserEntity)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					encoder := json.NewEncoder(w)
					if err := encoder.Encode(types.Token{
						Type:  "Bearer",
						Token: access_token,
					}); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
			}
		case types.StoreEntity:
			{
				if store, err := as.storeRepo.GetById(r.Context(), id); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					access_token := generateToken(&store, time.Minute*15, types.StoreEntity)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					encoder := json.NewEncoder(w)
					if err := encoder.Encode(types.Token{
						Type:  "Bearer",
						Token: access_token,
					}); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
			}
		default:
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}

func generateToken[T types.Identifiable](t *T, duration time.Duration, role types.EntityType) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"_id":  (*t).GetId(),
		"sub":  (*t).GetUsername(),
		"exp":  time.Now().Add(duration).Unix(),
		"role": string(role),
	})
	if tokenString, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN_SECRET"))); err != nil {
		log.Printf(">>>> Error: %s\n", err.Error())
		return ""
	} else {
		return tokenString
	}
}
