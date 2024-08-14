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

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*types.User, error)
	GetAllUsers(ctx context.Context) (*[]types.User, error)
	CreateUser(ctx context.Context, data types.User) (*types.User, error)
	UpdateUser(ctx context.Context, id string, data types.User) (*types.User, error)
	DeleteUser(ctx context.Context, ids []string) error
}

type UserRepo struct {
	store *database.DBPool
}

func NewUserRepo(dataStore *database.DBPool) *UserRepo {
	return &UserRepo{
		store: dataStore,
	}
}

func (repo *StoreRepo) GetUser(ctx context.Context, id string) (*types.Store, error) {
	return nil, nil
}

func (repo *StoreRepo) GetAllUsers(ctx context.Context) (*[]types.Store, error) {
	return nil, nil
}

func (repo *StoreRepo) CreateUser(ctx context.Context, data types.Store) (*types.Store, error) {
	return nil, nil
}

func (repo *StoreRepo) UpdateUser(ctx context.Context, id string, data types.Store) (*types.Store, error) {
	return nil, nil
}

func (repo *StoreRepo) DeleteUser(ctx context.Context, id string) (*types.Store, error) {
	return nil, nil
}
