package repositories

import (
	"context"

	"github.com/raganrrlaw/server/database"
	"github.com/raganrrlaw/server/types"
)

/*
Here in this repos after performing retrieval of the records from the database, you can perform
formatting and creating the return types that you want using the type definitions of each entity
*/
type StoreRepository interface {
	GetStore(ctx context.Context, id string) (*types.Store, error)
	GetAllStores(ctx context.Context) (*[]types.Store, error)
	CreateStore(ctx context.Context, data types.Store) (*types.Store, error)
	UpdateStore(ctx context.Context, id string, data types.Store) (*types.Store, error)
	DeleteStore(ctx context.Context, ids []string) error
}

type StoreRepo struct {
	store *database.DBPool
}

func NewStoreRepo(dataStore *database.DBPool) *StoreRepo {
	return &StoreRepo{
		store: dataStore,
	}
}

func (repo *StoreRepo) GetStore(ctx context.Context, id string) (*types.Store, error) {
	return nil, nil
}

func (repo *StoreRepo) GetAllStores(ctx context.Context) (*[]types.Store, error) {
	return nil, nil
}

func (repo *StoreRepo) CreateStore(ctx context.Context, data types.Store) (*types.Store, error) {
	return nil, nil
}

func (repo *StoreRepo) UpdateStore(ctx context.Context, id string, data types.Store) (*types.Store, error) {
	return nil, nil
}

func (repo *StoreRepo) DeleteStore(ctx context.Context, id string) (*types.Store, error) {
	return nil, nil
}
