package recommenderrepository

import (
	"context"

	"github.com/raganrrlaw/server/db/database"
)

type RecommenderRepository interface {
	Search(context.Context, string) error
}

type RecommenderRepo struct {
	storage *database.Storage
}

func NewRecommenderRepo(storage *database.Storage) RecommenderRepository {
	return &RecommenderRepo{storage: storage}
}

func (repo *RecommenderRepo) Search(ctx context.Context, list string) error {
	return nil
}

func GetStores(ctx context.Context, radius float64) {}
