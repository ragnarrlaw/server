package userrepository

import (
	"context"

	database "github.com/raganrrlaw/server/internal/db"
	"github.com/raganrrlaw/server/internal/types"
)

type UserRepository interface {
	Add(context.Context, *types.User) (*types.User, error)
	GetById(context.Context, string) (*types.User, error)
	GetAll(context.Context) (*[]types.User, error)
	Remove(context.Context, []string) error
	Update(context.Context, *types.User) (*types.User, error)
}

type UserRepo struct {
	storage *database.Storage
}

func NewUserRepo(storage *database.Storage) UserRepository {
	return &UserRepo{storage: storage}
}

func (repo *UserRepo) Add(ctx context.Context, u *types.User) (*types.User, error) {
	return nil, nil
}

func (repo *UserRepo) GetById(ctx context.Context, id string) (*types.User, error) {
	query := "SELECT id, name, email FROM users WHERE id = $1"
	var user types.User
	err := repo.storage.GetRow(
		ctx,
		query, []string{id}, &user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepo) GetAll(ctx context.Context) (*[]types.User, error) {
	return nil, nil
}

func (repo *UserRepo) Remove(ctx context.Context, ids []string) error {
	return nil
}

func (repo *UserRepo) Update(ctx context.Context, u *types.User) (*types.User, error) {
	return nil, nil
}
