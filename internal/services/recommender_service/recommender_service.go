package recommenderservice

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	recommenderrepository "github.com/ragnarrlaw/server/internal/services/recommender_service/recommender_repository"
	"github.com/ragnarrlaw/server/internal/types/recommender"
	"github.com/ragnarrlaw/server/internal/types/search"
)

type RecommenderService struct {
	recommenderRepository recommenderrepository.RecommenderRepository
}

func NewRecommenderService(recommenderRepo recommenderrepository.RecommenderRepository) *RecommenderService {
	return &RecommenderService{
		recommenderRepository: recommenderRepo,
	}
}

func (rs *RecommenderService) SearchHandler(w http.ResponseWriter, r *http.Request) {
	payload := &recommender.Search{}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		if recommendations, err := rs.recommenderRepository.Search(ctx, payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			if b, err := json.Marshal(search.Payload{
				Data: recommendations,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(b)
			}
		}
	}
}

func (rs *RecommenderService) LineStringSearchHandler(w http.ResponseWriter, r *http.Request) {
	payload := &recommender.LineStringSearch{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		if recommendations, err := rs.recommenderRepository.LineStringSearch(ctx, payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			if b, err := json.Marshal(search.Payload{
				Data: recommendations,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write(b)
			}
		}
	}
}
