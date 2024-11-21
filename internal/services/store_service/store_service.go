package storeservice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	storerepository "github.com/raganrrlaw/server/internal/services/store_service/store_repository"
	"github.com/raganrrlaw/server/internal/types"
	"github.com/raganrrlaw/server/internal/utils"
)

type StoreService struct {
	storeRepository storerepository.StoreRepository
}

func NewStoreService(storeRepository storerepository.StoreRepository) *StoreService {
	return &StoreService{
		storeRepository: storeRepository,
	}
}

func (ss *StoreService) AddStoreHandler(w http.ResponseWriter, r *http.Request) {
	payload := &types.StoreSignUpPayload{}
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
		store, err := ss.storeRepository.Add(ctx, payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if b, err := json.Marshal(store); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(b)
			}
		}
	}
}

func (ss *StoreService) UpdateStoreHandler(w http.ResponseWriter, r *http.Request) {
	payload := &types.StoreUpdatePayload{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		/** TODO: PAYLOAD VALIDATION */
		/**
		  TODO: CHECK IF USER EXISTS IN A MIDDLEWARE
		  CHECK WHETH THE ID IN THE CONTEXT MATCHES THE ID GIVEN IN THE URL PARAMETER
		*/
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if storeId := r.PathValue("id"); storeId != "" {
			if store, err := ss.storeRepository.Update(ctx, storeId, payload); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				if b, err := json.Marshal(store); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					w.Write(b)
				}
			}
		} else {
			http.Error(w, "store id is required", http.StatusBadRequest)
		}
	}
}

func (ss *StoreService) UpdateStoreLocationHandler(w http.ResponseWriter, r *http.Request) {
	payload := &types.GeoPoint{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		if storeId := r.PathValue("id"); storeId != "" {
			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
			defer cancel()

			if store, err := ss.storeRepository.UpdateLocation(ctx, storeId, payload); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				if b, err := json.Marshal(store); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					w.Write(b)
				}
			}
		} else {
			http.Error(w, "store id is required", http.StatusBadRequest)
		}
	}
}

func (ss *StoreService) RemoveStoreHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	/**
	  TODO: WHETHER THE STORE EXISTS IN A MIDDLEWARE AND THE PERMISSION LEVEL OF THE

	*/
	/** TODO: PAYLOAD CAN BE USED IF NEEDED IN THE BODY OF THE DELETE REQUEST */
	if storeId := r.PathValue("id"); storeId != "" {
		if err := ss.storeRepository.Remove(ctx, []string{storeId}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	} else {
		http.Error(w, "store id is required", http.StatusBadRequest)
	}
}

func (ss *StoreService) GetAllStoresHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if stores, err := ss.storeRepository.GetAll(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		if b, err := json.Marshal(stores); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		}
	}
}

func (ss *StoreService) GetStoreHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	/**
	  TODO: ADD VALIDATIONS TO CHECK WHETHER THE USER EXISTS
	        IN HERE CHECK OF THE NO ROWS FOUND ERROR TO SEND
	        NO RESOURCE FOUND ERROR
	*/
	if storeId := r.PathValue("id"); storeId != "" {
		store, err := ss.storeRepository.GetById(ctx, storeId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if b, err := json.Marshal(store); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(b)
			}
		}
	}
}
