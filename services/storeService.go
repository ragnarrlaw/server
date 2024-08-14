package services

import (
	"net/http"

	"github.com/raganrrlaw/server/database"
	"github.com/raganrrlaw/server/repositories"
)

type StoreService struct {
	repo *repositories.StoreRepo
}

func NewStoreService(store *database.DBPool) *StoreService {
	return &StoreService{
		repo: repositories.NewStoreRepo(store),
	}
}

func (s *StoreService) RegisterRoutes(r *http.ServeMux) {
	storeRouter := http.NewServeMux()
	storeRouter.HandleFunc("GET /store", s.GetAllStoreHandler)
	storeRouter.HandleFunc("GET /store/{id}", s.GetStoreHandler)
	storeRouter.HandleFunc("POST /store", s.CreateStoreHandler)
	storeRouter.HandleFunc("PATCH /store/{id}", s.UpdateStoreHandler)
	storeRouter.HandleFunc("DELETE /store/{id}", s.UpdateStoreHandler)

	// you can wrap this in auth middleware, or more
	r.Handle("/", storeRouter)
}

func (s *StoreService) GetAllStoreHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *StoreService) GetStoreHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *StoreService) CreateStoreHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *StoreService) UpdateStoreHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *StoreService) DeleteStoreHandler(w http.ResponseWriter, r *http.Request) {
}
