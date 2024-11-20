package authrepository

import (
	"context"

	database "github.com/raganrrlaw/server/db/database"
	"github.com/raganrrlaw/server/internal/types"
)

type AuthRepository interface {
	GetById(context.Context, string) (*types.AuthToken, error)
	GetAssociatedTokens(context.Context, string) (*[]types.AuthToken, error)
	AddToken(context.Context, *types.AuthToken) (*types.AuthToken, error)
	RemoveToken(context.Context, []string) error
	RemoveTokensOfUser(context.Context, string) error
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
func (ar *AuthRepo) GetById(ctx context.Context, id string) (*types.AuthToken, error) {
	query := "SELECT id, user_id, role, token FROM auth_token WHERE id = $1"
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
func (ar *AuthRepo) AddToken(ctx context.Context, authToken *types.AuthToken) (*types.AuthToken, error) {
	query := "INSERT INTO auth_token (user_id, role, token) VALUES ($1, $2, $3) RETURNING id"

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
func (ar *AuthRepo) GetAssociatedTokens(ctx context.Context, userId string) (*[]types.AuthToken, error) {
	query := "SELECT id, user_id, role, token FROM auth_token WHERE user_id = $1"

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
func (ar *AuthRepo) RemoveToken(ctx context.Context, tokenIds []string) error {
	query := "DELETE FROM auth_token WHERE id = ANY($1)"

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
func (ar *AuthRepo) RemoveTokensOfUser(ctx context.Context, userId string) error {
	query := "DELETE FROM auth_token WHERE user_id = $1"
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
