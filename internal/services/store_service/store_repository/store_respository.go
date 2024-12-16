package storerepository

import (
	"context"
	"fmt"

	database "github.com/raganrrlaw/server/db/database"
	"github.com/raganrrlaw/server/internal/types/search"
	"github.com/raganrrlaw/server/internal/types/store"
)

type StoreRepository interface {
	GetById(context.Context, string) (*store.Store, error)
	GetBy(context.Context, string, any) ([]store.Store, error)
	Get(context.Context) ([]store.Store, error)
	Stats(context.Context) ([]store.StoreStat, error)
	GetTotalNumberOfStores(context.Context) (uint, error)
}

type StoreRepo struct {
	storage *database.Storage
}

func NewStoreRepo(storage *database.Storage) StoreRepository {
	return &StoreRepo{storage: storage}
}

func (repo *StoreRepo) GetTotalNumberOfStores(ctx context.Context) (uint, error) {
	query := `
	SELECT COUNT(id) FROM store;
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

func (repo *StoreRepo) GetById(ctx context.Context, id string) (*store.Store, error) {
	query := `
    SELECT  
		id, 
    	username, 
    	name, 
    	address, 
    	email,
    	contact_number, 
      discounts,
      web_url, 
    	password_digest,
    	COALESCE(ST_Y(GEOMETRY(location_point)), 0) AS latitude, 
    	COALESCE(ST_X(GEOMETRY(location_point)), 0) AS longitude,
      created_at,
      updated_at
	FROM store
	WHERE id = $1
  `
	var store store.Store

	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		id,
	)

	if err := row.Scan(
		&store.Id,
		&store.Username,
		&store.Name,
		&store.Address,
		&store.Email,
		&store.ContactNumber,
		&store.Discounts,
		&store.WebUrl,
		&store.PasswordDigest,
		&store.LocationPoint.Latitude,
		&store.LocationPoint.Longitude,
		&store.CreatedAt,
		&store.UpdatedAt,
	); err != nil {
		return nil, err
	} else {
		return &store, nil
	}
}

func (repo *StoreRepo) GetBy(ctx context.Context, key string, value any) ([]store.Store, error) {
	params := ctx.Value(search.SearchKey).(search.Paginate)
	query := fmt.Sprintf(
		`SELECT	
 		id, 
    	username, 
    	name, 
    	address, 
    	email,
    	contact_number, 
      discounts,
      web_url, 
    	password_digest,
    	COALESCE(ST_Y(GEOMETRY(location_point)), 0) AS latitude, 
    	COALESCE(ST_X(GEOMETRY(location_point)), 0) AS longitude,
      created_at,
      updated_at
	FROM store
		WHERE %s = $1
    OFFSET $2 
    LIMIT $3`,
		key,
	)

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

	defer rows.Close()

	var stores []store.Store

	for rows.Next() {
		var store store.Store
		if err := rows.Scan(
			&store.Id,
			&store.Username,
			&store.Name,
			&store.Address,
			&store.Email,
			&store.ContactNumber,
			&store.Discounts,
			&store.WebUrl,
			&store.PasswordDigest,
			&store.LocationPoint.Latitude,
			&store.LocationPoint.Longitude,
			&store.CreatedAt,
			&store.UpdatedAt,
		); err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	return stores, nil
}

func (repo *StoreRepo) Get(ctx context.Context) ([]store.Store, error) {
	params := ctx.Value(search.SearchKey).(search.Paginate)
	query := `
    SELECT  
			id, 
    	username, 
    	name, 
    	address, 
    	email,
    	contact_number, 
      discounts,
      web_url, 
    	password_digest,
    	COALESCE(ST_Y(GEOMETRY(location_point)), 0) AS latitude, 
    	COALESCE(ST_X(GEOMETRY(location_point)), 0) AS longitude,
      created_at,
      updated_at
	FROM store
  OFFSET $1
  LIMIT $2`
	rows, err := repo.storage.Pool.Query(
		ctx,
		query,
		params.Offset,
		params.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stores []store.Store
	for rows.Next() {
		var store store.Store
		if err := rows.Scan(
			&store.Id,
			&store.Username,
			&store.Name,
			&store.Address,
			&store.Email,
			&store.ContactNumber,
			&store.Discounts,
			&store.WebUrl,
			&store.PasswordDigest,
			&store.LocationPoint.Latitude,
			&store.LocationPoint.Longitude,
			&store.CreatedAt,
			&store.UpdatedAt,
		); err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	return stores, nil
}

func (repo *StoreRepo) Stats(ctx context.Context) ([]store.StoreStat, error) {
	query := `
	SELECT
    	s.name AS name, 
    	s.address AS address, 
      	s.discounts AS discounts,
    	COALESCE(ST_Y(GEOMETRY(s.location_point)), 0) AS latitude, 
    	COALESCE(ST_X(GEOMETRY(s.location_point)), 0) AS longitude,
      	s.created_at AS created_at,
		COUNT(sp.product_id) AS total_product_count,
		COUNT(CASE WHEN sp.stock_quantity = 'LIMITED_QUANTITY' THEN 1 END) AS products_with_limited_quantity,
		COUNT(CASE WHEN sp.stock_quantity = 'AVAILABLE' THEN 1 END) AS fully_available_products,
		COUNT(CASE WHEN sp.stock_quantity = 'NOT_AVAILABLE' THEN 1 END) AS out_of_stock_products
	FROM store s
	JOIN store_product sp ON sp.store_id = s.id
	GROUP BY s.id, s.name, s.address, s.discounts, s.location_point, s.created_at
	`
	rows, err := repo.storage.Pool.Query(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var storeStats []store.StoreStat
	for rows.Next() {
		var storeStat store.StoreStat
		if err := rows.Scan(
			&storeStat.Name,
			&storeStat.Address,
			&storeStat.Discounts,
			&storeStat.Point.Latitude,
			&storeStat.Point.Longitude,
			&storeStat.CreatedAt,
			&storeStat.TotalNumberOfProducts,
			&storeStat.NumberOfLimitedStockProducts,
			&storeStat.NumberOfAvailableProducts,
			&storeStat.NumberOfOutOfStockProducts,
		); err != nil {
			return nil, err
		}
		storeStats = append(storeStats, storeStat)
	}
	return storeStats, nil
}
