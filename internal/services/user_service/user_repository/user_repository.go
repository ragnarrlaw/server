package userrepository

import (
	"context"

	database "github.com/raganrrlaw/server/internal/db"
	"github.com/raganrrlaw/server/internal/types"
)

type UserRepository interface {
	Add(context.Context, *types.UserSignUpPayload) (*types.User, error)
	GetById(context.Context, string) (*types.User, error)
	GetAll(context.Context) (*[]types.User, error)
	Remove(context.Context, []string) error
	Update(context.Context, string, *types.UserUpdatePayload) (*types.User, error)
}

type UserRepo struct {
	storage *database.Storage
}

func NewUserRepo(storage *database.Storage) UserRepository {
	return &UserRepo{storage: storage}
}

func (repo *UserRepo) Add(ctx context.Context, u *types.UserSignUpPayload) (*types.User, error) {
	query := "INSERT INTO users (username, first_name, last_name, email, contact_number, password_digest) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, username, first_name, last_name, email, contact_number, password_digest"
	var user types.User
	err := repo.storage.GetRow(
		ctx,
		query,
		[]interface{}{
			u.Username,
			u.FirstName,
			u.LastName,
			u.Email,
			u.ContactNumber,
			u.Password,
		},
		[]interface{}{
			&user.Id,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.ContactNumber,
			&user.Password,
		},
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepo) GetById(ctx context.Context, id string) (*types.User, error) {
	query := "SELECT id, username, first_name, last_name, email, contact_number, password_digest FROM users WHERE id = $1"
	var user types.User
	err := repo.storage.GetRow(
		ctx,
		query,
		[]interface{}{
			id,
		},
		[]interface{}{
			&user.Id,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.ContactNumber,
			&user.Password,
		},
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepo) GetAll(ctx context.Context) (*[]types.User, error) {
	query := "SELECT id, username, first_name, last_name, email, contact_number, password_digest FROM users"
	rows, err := repo.storage.GetRows(ctx, query)
	if err != nil {
		return nil, err
	}

	var users []types.User
	for _, row := range rows {
		user := types.User{
			Id:            row["id"].(string),
			Username:      row["username"].(string),
			FirstName:     row["first_name"].(string),
			LastName:      row["last_name"].(string),
			Email:         row["email"].(string),
			ContactNumber: row["contact_number"].(string),
			Password:      row["password_digest"].(string),
		}
		users = append(users, user)
	}
	return &users, nil
}

func (repo *UserRepo) Remove(ctx context.Context, ids []string) error {
	query := "DELETE FROM users WHERE id = ANY($1)"
	err := repo.storage.Exec(ctx, query, ids)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) Update(ctx context.Context, id string, u *types.UserUpdatePayload) (*types.User, error) {
	query := "UPDATE users SET username = $1, first_name = $2, last_name = $3, email = $4, contact_number = $5, password_digest = $6 WHERE id = $7 RETURNING id, username, first_name, last_name, email, contact_number, password_digest"
	var user types.User
	err := repo.storage.GetRow(
		ctx,
		query,
		[]interface{}{
			u.Username,
			u.FirstName,
			u.LastName,
			u.Email,
			u.ContactNumber,
			u.Password,
			id,
		},
		[]interface{}{
			&user.Id,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.ContactNumber,
			&user.Password,
		},
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
