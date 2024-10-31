package storeservice

import (
	"net/http"

	storerepository "github.com/raganrrlaw/server/internal/services/store_service/store_repository"
)

type StoreService struct {
	storeRepository storerepository.StoreRepository
}

func NewStoreService(storeRepository storerepository.StoreRepository) *StoreService {
	return &StoreService{storeRepository: storeRepository}
}

func (ss *StoreService) GetAllStoresHandler(w http.ResponseWriter, r *http.Request) {
}

func (ss *StoreService) GetStoreHandler(w http.ResponseWriter, r *http.Request) {}

func (ss *StoreService) AddStoreHandler(w http.ResponseWriter, r *http.Request) {}

func (ss *StoreService) UpdateStoreHandler(w http.ResponseWriter, r *http.Request) {}

func (ss *StoreService) RemoveStoreHandler(w http.ResponseWriter, r *http.Request) {}
