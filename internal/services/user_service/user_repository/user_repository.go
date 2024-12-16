package userrepository

import (
	"context"
	"fmt"

	database "github.com/raganrrlaw/server/db/database"
	"github.com/raganrrlaw/server/internal/types/search"
	"github.com/raganrrlaw/server/internal/types/user"
)

type UserRepository interface {
	// user
	GetById(context.Context, string) (*user.User, error)
	GetBy(context.Context, string, any) ([]user.User, error)
	Get(context.Context) ([]user.User, error)
	GetTotalNumberOfUsers(context.Context) (uint, error)
}

type UserRepo struct {
	storage *database.Storage
}

func NewUserRepo(storage *database.Storage) UserRepository {
	return &UserRepo{storage: storage}
}

func (repo *UserRepo) GetTotalNumberOfUsers(ctx context.Context) (uint, error) {
	query := `
	SELECT COUNT(id) FROM users;
	`
	var total uint
	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
	)
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (repo *UserRepo) GetById(ctx context.Context, id string) (*user.User, error) {
	query := `
		SELECT id, username, first_name, last_name, email, contact_number, password_digest, created_at, updated_at
		FROM users WHERE id = $1
	`
	var user user.User

	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.Id,
		&user.Username,
		&user.FistName,
		&user.LastName,
		&user.Email,
		&user.ContactNumber,
		&user.PasswordDigest,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func (repo *UserRepo) GetBy(ctx context.Context, field string, value any) ([]user.User, error) {
	params := ctx.Value(search.SearchKey).(search.Paginate)
	query := fmt.Sprintf(`
		SELECT id, username, first_name, last_name, email, contact_number, password_digest, created_at, updated_at
		FROM users WHERE %s = $1 OFFSET $2 LIMIT $3
	`,
		field,
	)
	var users []user.User

	rows, err := repo.storage.Pool.Query(
		ctx,
		query,
		value,
		params.Offset,
		params.Limit,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user user.User
		if err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.FistName,
			&user.LastName,
			&user.Email,
			&user.ContactNumber,
			&user.PasswordDigest,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (repo *UserRepo) Get(ctx context.Context) ([]user.User, error) {
	params := ctx.Value(search.SearchKey).(search.Paginate)
	query := `
		SELECT id, username, first_name, last_name, email, contact_number, password_digest, created_at, updated_at
		FROM users OFFSET $1 LIMIT $2
	`
	if rows, err := repo.storage.Pool.Query(
		ctx,
		query,
		params.Offset,
		params.Limit,
	); err != nil {
		return nil, err
	} else {
		var users []user.User
		for rows.Next() {
			var user user.User
			if err := rows.Scan(
				&user.Id,
				&user.Username,
				&user.FistName,
				&user.LastName,
				&user.Email,
				&user.ContactNumber,
				&user.PasswordDigest,
				&user.CreatedAt,
				&user.UpdatedAt,
			); err != nil {
				return nil, err
			}
			users = append(users, user)
		}
		return users, nil
	}
}
