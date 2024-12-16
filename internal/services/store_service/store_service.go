package storeservice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	storerepository "github.com/raganrrlaw/server/internal/services/store_service/store_repository"
	"github.com/raganrrlaw/server/internal/types/search"
)

type StoreService struct {
	storeRepository storerepository.StoreRepository
}

func NewStoreService(storeRepository storerepository.StoreRepository) *StoreService {
	return &StoreService{
		storeRepository: storeRepository,
	}
}

func (ss *StoreService) GetAllStoresHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if stores, err := ss.storeRepository.Get(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		if totalRecords, err := ss.storeRepository.GetTotalNumberOfStores(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if b, err := json.Marshal(search.Payload{
				Data: stores,
				Meta: search.Meta{
					Total: totalRecords,
				},
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(b)
			}
		}
	}
}

func (ss *StoreService) GetStoreStatHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	stats, err := ss.storeRepository.Stats(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, err := ss.storeRepository.GetTotalNumberOfStores(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(search.Payload{
		Data: stats,
		Meta: search.Meta{
			Total: total,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (ss *StoreService) GetStoreHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

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
