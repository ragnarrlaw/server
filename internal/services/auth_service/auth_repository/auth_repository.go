package authrepository

import (
	"context"
	"errors"

	database "github.com/raganrrlaw/server/db/database"
	"github.com/raganrrlaw/server/internal/types"
)

type AuthRepository interface {
	GetById(context.Context, string, types.EntityType) (*types.AuthToken, error)
	GetAssociatedTokens(context.Context, string, types.EntityType) (*[]types.AuthToken, error)
	AddToken(context.Context, *types.AuthToken, types.EntityType) (*types.AuthToken, error)
	RemoveToken(context.Context, []string, types.EntityType) error
	RemoveTokensOfUser(context.Context, string, types.EntityType) error
}

type AuthRepo struct {
	User    *types.User
	storage *database.Storage
}

func NewAuthRepo(storage *database.Storage) AuthRepository {
	return &AuthRepo{
		storage: storage,
	}
}

/** GetById fetches a token by id */
func (ar *AuthRepo) GetById(ctx context.Context, id string, role types.EntityType) (*types.AuthToken, error) {
	var query string

	switch role {
	case types.UserEntity:
		{
			query = "SELECT id, user_id, role, token FROM user_auth_token WHERE id = $1"
		}
	case types.StoreEntity:
		{
			query = "SELECT id, store_id, role, token FROM store_auth_token WHERE id = $1"
		}
	default:
		return nil, errors.New("invalid role")
	}

	var authToken types.AuthToken
	row := ar.storage.Pool.QueryRow(
		ctx,
		query,
		id,
	)
	if err := row.Scan(
		&authToken.Id,
		&authToken.UserId,
		&authToken.Role,
		&authToken.Token,
	); err != nil {
		return nil, err
	} else {
		return &authToken, nil
	}
}

/** AddToken adds a new token to the database */
func (ar *AuthRepo) AddToken(ctx context.Context, authToken *types.AuthToken, role types.EntityType) (*types.AuthToken, error) {
	var query string

	switch role {
	case types.UserEntity:
		{
			query = "INSERT INTO user_auth_token (user_id, role, token) VALUES ($1, $2, $3) RETURNING id, user_id, role, token"
		}
	case types.StoreEntity:
		{
			query = "INSERT INTO store_auth_token (store_id, role, token) VALUES ($1, $2, $3) RETURNING id, store_id, role, token"
		}
	default:
		return nil, errors.New("invalid role")
	}

	var token types.AuthToken
	row := ar.storage.Pool.QueryRow(
		ctx,
		query,
		authToken.UserId,
		authToken.Role,
		authToken.Token,
	)
	if err := row.Scan(
		&authToken.Id,
		&authToken.UserId,
		&authToken.Role,
		&authToken.Token,
	); err != nil {
		return nil, err
	} else {
		return &token, nil
	}
}

/** GetAssociatedTokens fetches all tokens associated with a user */
func (ar *AuthRepo) GetAssociatedTokens(ctx context.Context, userId string, role types.EntityType) (*[]types.AuthToken, error) {
	var query string

	switch role {
	case types.UserEntity:
		{
			query = "SELECT id, user_id, role, token FROM user_auth_token WHERE user_id = $1"
		}
	case types.StoreEntity:
		{
			query = "INSERT INTO store_auth_token (store_id, role, token) VALUES ($1, $2, $3) RETURNING id, store_id, role, token"
		}
	default:
		return nil, errors.New("invalid role")
	}

	if rows, err := ar.storage.Pool.Query(
		ctx,
		query,
		userId,
	); err != nil {
		return nil, err
	} else {
		var authTokens []types.AuthToken
		for rows.Next() {
			var authToken types.AuthToken
			if err := rows.Scan(
				&authToken.Id,
				&authToken.UserId,
				&authToken.Role,
				&authToken.Token,
			); err != nil {
				return nil, err
			}
			authTokens = append(authTokens, authToken)
		}
		return &authTokens, nil
	}
}

/** RemoveToken removes a single token associated with a token id  */
func (ar *AuthRepo) RemoveToken(ctx context.Context, tokenIds []string, role types.EntityType) error {
	var query string

	switch role {
	case types.UserEntity:
		{
			query = "DELETE FROM user_auth_token WHERE id = ANY($1)"
		}
	case types.StoreEntity:
		{
			query = "DELETE FROM store_auth_token WHERE id = ANY($1)"
		}
	default:
		return errors.New("invalid role")
	}

	if _, err := ar.storage.Pool.Exec(
		ctx,
		query,
		tokenIds,
	); err != nil {
		return err
	} else {
		return nil
	}
}

/** RemoveTokensOfUser removes all tokens associated with a user */
func (ar *AuthRepo) RemoveTokensOfUser(ctx context.Context, userId string, role types.EntityType) error {
	var query string

	switch role {
	case types.UserEntity:
		{
			query = "DELETE FROM user_auth_token WHERE user_id = $1"
		}
	case types.StoreEntity:
		{
			query = "DELETE FROM store_auth_token WHERE store_id = $1"
		}
	default:
		return errors.New("invalid role")
	}

	if _, err := ar.storage.Pool.Exec(
		ctx,
		query,
		userId,
	); err != nil {
		return err
	} else {
		return nil
	}
}
