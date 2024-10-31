package storerepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	database "github.com/raganrrlaw/server/internal/db"
	"github.com/raganrrlaw/server/internal/types"
)

type StoreRepository interface {
	Add(context.Context, *types.StoreSignUpPayload) (*types.Store, error)
	GetById(context.Context, string) (*types.Store, error)
	GetBy(context.Context, string, any) (*types.Store, error)
	GetAll(context.Context) (*[]types.Store, error)
	Remove(context.Context, []string) error
	Update(context.Context, string, *types.StoreUpdatePayload) (*types.Store, error)
}

type StoreRepo struct {
	storage *database.Storage
}

func NewStoreRepo(storage *database.Storage) StoreRepository {
	return &StoreRepo{storage: storage}
}

func (repo *StoreRepo) Add(ctx context.Context, s *types.StoreSignUpPayload) (*types.Store, error) {
	query := "INSERT INTO store (store_name, store_address, store_email, store_contact_number, password_digest) VALUES ($1, $2, $3, $4, $5) RETURNING id, store_name, store_address, store_email, store_contact_number, password_digest"
	var store types.Store
	err := repo.storage.GetRow(
		ctx,
		query,
		[]interface{}{
			s.StoreName,
			s.Address,
			s.Email,
			s.ContactNumber,
			s.Password,
		},
		[]interface{}{
			&store.Id,
			&store.StoreName,
			&store.StoreEmail,
			&store.StoreContactNumber,
			&store.Password,
		},
	)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (repo *StoreRepo) GetById(ctx context.Context, id string) (*types.Store, error) {
	query := "SELECT id, store_username, store_name, store_address, store_email, store_contact_number, password_digest FROM store WHERE id = $1"
	var store types.Store
	err := repo.storage.GetRow(
		ctx, query, []interface{}{
			id,
		},
		[]interface{}{
			&store.Id,
			&store.StoreUsername,
			&store.StoreName,
			&store.StoreAddress,
			&store.StoreEmail,
			&store.StoreContactNumber,
			&store.Password,
		})
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (repo *StoreRepo) GetBy(ctx context.Context, key string, value any) (*types.Store, error) {
	query := fmt.Sprintf("SELECT id, store_username, store_name, store_address, store_email, store_contact_number, password_digest FROM store WHERE %s = $1", key)
	var store types.Store
	err := repo.storage.GetRow(
		ctx,
		query,
		[]interface{}{
			value,
		},
		[]interface{}{
			&store.Id,
			&store.StoreUsername,
			&store.StoreName,
			&store.StoreAddress,
			&store.StoreEmail,
			&store.StoreContactNumber,
			&store.Password,
		})
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (repo *StoreRepo) GetAll(ctx context.Context) (*[]types.Store, error) {
	query := "SELECT id, store_username, store_name, store_address, store_email, store_contact_number, password_digest FROM store"
	rows, err := repo.storage.GetRows(ctx, query)
	if err != nil {
		return nil, err
	}

	var stores []types.Store
	for _, row := range rows {
		store := types.Store{
			Id:                 row["id"].(uuid.UUID),
			StoreUsername:      row["store_username"].(string),
			StoreName:          row["store_name"].(string),
			StoreAddress:       row["store_address"].(string),
			StoreEmail:         row["store_email"].(string),
			StoreContactNumber: row["store_contact_number"].(string),
			Password:           row["password_digest"].(string),
		}
		stores = append(stores, store)
	}
	return &stores, nil
}

func (repo *StoreRepo) Remove(ctx context.Context, ids []string) error {
	query := "DELETE FROM store WHERE id = ANY($1)"
	err := repo.storage.Exec(ctx, query, ids)
	if err != nil {
		return err
	}
	return nil
}

/*
*

	stores can only update the store username, store name, store contact number and the store email using
	using the store api, the rest should happen through the authentication apis
*/
func (repo *StoreRepo) Update(ctx context.Context, id string, s *types.StoreUpdatePayload) (*types.Store, error) {
	query := "UPDATE store SET store_name = $1, store_email = $2, store_contact_number = $3, store_username = $4 WHERE id = $5 RETURNING id, store_name, store_address, store_email, store_contact_number, password_digest, store_username"
	var store types.Store
	err := repo.storage.GetRow(
		ctx,
		query,
		[]interface{}{
			s.StoreName,
			s.StoreEmail,
			s.StoreContactNumber,
			s.StoreUsername,
			id,
		},
		[]interface{}{
			&store.Id,
			&store.StoreName,
			&store.StoreAddress,
			&store.StoreContactNumber,
			&store.Password,
			&store.StoreUsername,
		},
	)
	if err != nil {
		return nil, err
	}
	return &store, nil
}
