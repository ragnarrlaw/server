package authrepository

import (
	"context"
	"fmt"
	"log"

	uuid "github.com/google/uuid"
	database "github.com/raganrrlaw/server/internal/db"
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
	query := "SELECT id, user_id, token FROM auth_token WHERE id = $1"
	var authToken types.AuthToken
	if err := ar.storage.GetRow(
		ctx,
		query,
		[]interface{}{id},
		[]interface{}{
			&authToken.Id,
			&authToken.UserId,
			&authToken.Token,
		},
	); err != nil {
		log.Printf("Error fetching auth token by id: %v", err)
		return nil, fmt.Errorf("failed to fetch auth token by id: %w", err)
	}
	return &authToken, nil
}

/** AddToken adds a new token to the database */
func (ar *AuthRepo) AddToken(ctx context.Context, authToken *types.AuthToken) (*types.AuthToken, error) {
	query := "INSERT INTO auth_token (user_id, token) VALUES ($1, $2) RETURNING id"
	if err := ar.storage.GetRow(
		ctx,
		query,
		[]interface{}{
			authToken.UserId,
			authToken.Token,
		},
		[]interface{}{
			&authToken.Id,
		},
	); err != nil {
		log.Printf("Error adding auth token: %v", err)
		return nil, fmt.Errorf("failed to add auth token: %w", err)
	}
	return authToken, nil
}

/** GetAssociatedTokens fetches all tokens associated with a user */
func (ar *AuthRepo) GetAssociatedTokens(ctx context.Context, userId string) (*[]types.AuthToken, error) {
	query := "SELECT id, user_id, token FROM auth_token WHERE user_id = $1"
	rows, err := ar.storage.GetRows(ctx, query, userId)
	if err != nil {
		log.Printf("Error fetching associated tokens: %v", err)
		return nil, fmt.Errorf("failed to fetch associated tokens: %w", err)
	}
	var authTokens []types.AuthToken
	for _, row := range rows {
		authToken := types.AuthToken{
			Id:     row["id"].(uuid.UUID),
			UserId: row["user_id"].(string),
			Token:  row["token"].(string),
		}
		authTokens = append(authTokens, authToken)
	}
	return &authTokens, nil
}

/** RemoveToken removes a single token associated with a token id  */
func (ar *AuthRepo) RemoveToken(ctx context.Context, tokenIds []string) error {
	query := "DELETE FROM auth_token WHERE id = ANY($1)"
	if err := ar.storage.Exec(ctx, query, tokenIds); err != nil {
		log.Printf("Error removing token: %v", err)
		return fmt.Errorf("failed to remove token: %w", err)
	}
	return nil
}

/** RemoveTokensOfUser removes all tokens associated with a user */
func (ar *AuthRepo) RemoveTokensOfUser(ctx context.Context, userId string) error {
	query := "DELETE FROM auth_token WHERE user_id = $1"
	if err := ar.storage.Exec(ctx, query, userId); err != nil {
		log.Printf("Error removing tokens of user: %v", err)
		return fmt.Errorf("failed to remove tokens of user: %w", err)
	}
	return nil
}
