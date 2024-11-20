package storerepository

import (
	"context"
	"fmt"

	database "github.com/raganrrlaw/server/db/database"
	"github.com/raganrrlaw/server/internal/types"
)

type StoreRepository interface {
	Add(context.Context, *types.StoreSignUpPayload) (*types.Store, error)
	GetById(context.Context, string) (*types.Store, error)
	GetBy(context.Context, string, any) (*types.Store, error)
	GetAll(context.Context) (*[]types.Store, error)
	Remove(context.Context, []string) error
	Update(context.Context, string, *types.StoreUpdatePayload) (*types.Store, error)
	UpdateLocation(context.Context, string, *types.GeoPoint) (*types.Store, error)
}

type StoreRepo struct {
	storage *database.Storage
}

func NewStoreRepo(storage *database.Storage) StoreRepository {
	return &StoreRepo{storage: storage}
}

func (repo *StoreRepo) Add(ctx context.Context, s *types.StoreSignUpPayload) (*types.Store, error) {
	query := `
    INSERT INTO store (store_username, store_name, store_address, store_email, store_contact_number, password_digest, store_web_url)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, store_username, store_name, store_address, store_email, store_contact_number, password_digest, store_web_url
  `
	var store types.Store

	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		s.Username,
		s.StoreName,
		s.Address,
		s.Email,
		s.ContactNumber,
		s.Password,
		s.WebURL,
	)

	if err := row.Scan(
		&store.Id,
		&store.StoreUsername,
		&store.StoreName,
		&store.StoreEmail,
		&store.StoreContactNumber,
		&store.Password,
		&store.StoreWebURL,
	); err != nil {
		return nil, err
	} else {
		return &store, nil
	}
}

func (repo *StoreRepo) GetById(ctx context.Context, id string) (*types.Store, error) {
	query := `
    SELECT id, store_username, store_name, store_address, store_email,
		       store_contact_number, store_web_url, password_digest,
		       ST_Y(store_location) AS latitude, ST_X(store_location) AS longitude
		FROM store
		WHERE id = $1
  `
	var store types.Store

	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		id,
	)

	if err := row.Scan(
		&store.Id,
		&store.StoreUsername,
		&store.StoreName,
		&store.StoreAddress,
		&store.StoreEmail,
		&store.StoreContactNumber,
		&store.StoreWebURL,
		&store.Password,
		&store.StoreLocation.Latitude,
		&store.StoreLocation.Longitude,
	); err != nil {
		return nil, err
	} else {
		return &store, nil
	}
}

func (repo *StoreRepo) GetBy(ctx context.Context, key string, value any) (*types.Store, error) {

	query := fmt.Sprintf(
		`SELECT id, store_username, store_name, store_address, store_email,
		       store_contact_number, store_web_url, password_digest,
		       ST_Y(store_location) AS latitude, ST_X(store_location) AS longitude
		FROM store
		WHERE %s = $1`, key,
	)

	var store types.Store

	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		value,
	)

	if err := row.Scan(
		&store.Id,
		&store.StoreUsername,
		&store.StoreName,
		&store.StoreAddress,
		&store.StoreEmail,
		&store.StoreContactNumber,
		&store.StoreWebURL,
		&store.Password,
		&store.StoreLocation.Latitude,
		&store.StoreLocation.Longitude,
	); err != nil {
		return nil, err
	} else {
		return &store, nil
	}
}

func (repo *StoreRepo) GetAll(ctx context.Context) (*[]types.Store, error) {
	query :=
		`SELECT id, store_username, store_name, store_address, store_email,
		       store_contact_number, store_web_url, password_digest,
		       ST_Y(store_location) AS latitude, ST_X(store_location) AS longitude
		FROM store`

	if rows, err := repo.storage.Pool.Query(
		ctx,
		query,
	); err != nil {
		return nil, err
	} else {
		var stores []types.Store
		for rows.Next() {
			var store types.Store
			if err := rows.Scan(
				&store.Id,
				&store.StoreUsername,
				&store.StoreName,
				&store.StoreAddress,
				&store.StoreEmail,
				&store.StoreContactNumber,
				&store.StoreWebURL,
				&store.Password,
				&store.StoreLocation.Latitude,
				&store.StoreLocation.Longitude,
			); err != nil {
				return nil, err
			}
			stores = append(stores, store)
		}
		return &stores, nil
	}
}

func (repo *StoreRepo) Remove(ctx context.Context, ids []string) error {
	query := `DELETE FROM store WHERE id = ANY($1)`
	if _, err := repo.storage.Pool.Exec(
		ctx,
		query,
		ids,
	); err != nil {
		return err
	}
	return nil
}

func (repo *StoreRepo) Update(ctx context.Context, id string, s *types.StoreUpdatePayload) (*types.Store, error) {
	query := `
    UPDATE store SET
      store_name = $1,
      store_email = $2,
      store_contact_number = $3,
      store_username = $4 WHERE id = $5
    RETURNING
      id,
      store_username,
      store_name,
      store_address,
      store_email,
      store_contact_number,
      password_digest,
      store_web_url
  `
	var store types.Store
	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		s.StoreName,
		s.StoreEmail,
		s.StoreContactNumber,
		s.StoreUsername,
		id,
	).Scan(
		&store.Id,
		&store.StoreUsername,
		&store.StoreName,
		&store.StoreAddress,
		&store.StoreContactNumber,
		&store.Password,
	); err != nil {
		return nil, err
	} else {
		return &store, nil
	}
}

func (repo *StoreRepo) UpdateLocation(ctx context.Context, id string, location *types.GeoPoint) (*types.Store, error) {
	query := `
    UPDATE store SET
      store_location = ST_SetSRID(ST_MakePoint($1, $2), 4326)
    WHERE id = $3
		RETURNING
      id,
      store_username,
      store_name,
      store_address,
      store_email,
		  store_contact_number,
      store_web_url,
      password_digest,
		  ST_Y(store_location) AS latitude,
      ST_X(store_location) AS longitude
  `
	var store types.Store
	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		location.Latitude,
		location.Longitude,
		id,
	).Scan(
		&store.Id,
		&store.StoreUsername,
		&store.StoreName,
		&store.StoreAddress,
		&store.StoreEmail,
		&store.StoreContactNumber,
		&store.Password,
		&store.StoreWebURL,
		&store.StoreLocation.Latitude,
		&store.StoreLocation.Longitude,
	); err != nil {
		return nil, err
	} else {
		return &store, nil
	}
}
