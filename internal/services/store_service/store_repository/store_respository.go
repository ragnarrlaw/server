package storerepository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	database "github.com/ragnarrlaw/server/db/database"
	"github.com/ragnarrlaw/server/internal/types/geocode"
	"github.com/ragnarrlaw/server/internal/types/search"
	"github.com/ragnarrlaw/server/internal/types/store"
)

type StoreRepository interface {
	// store
	GetById(context.Context, string) (*store.Store, error)
	GetBy(context.Context, string, any) ([]store.Store, error)
	Get(context.Context) ([]store.Store, error)
	SearchStores(context.Context, *store.FilterCriteria) ([]store.Store, error)
	Stats(context.Context) ([]store.StoreStat, error)
	GetTotalNumberOfStores(context.Context) (uint, error)

	// store products
	GetStoreProducts(context.Context, string) (json.RawMessage, error)
	GetTotalNumberOfProductsInStore(context.Context, string) (uint, error)
	CreateStoreProduct(context.Context, string, *store.StoreProductPayload) (json.RawMessage, error)
	UpdateStoreProduct(context.Context, string, string, *store.StoreProductPayload) (json.RawMessage, error)
	DeleteStoreProducts(context.Context, string, []string) error

	// store discounts
	GetStoreDiscounts(context.Context, uuid.UUID) ([]store.StoreDiscount, error)
	GetTotalNumberOfDiscountsInStore(context.Context, uuid.UUID) (uint, error)
	CreateStoreDiscount(context.Context, uuid.UUID, *store.StoreDiscountPayload) (*store.StoreDiscount, error)
	UpdateStoreDiscount(context.Context, uuid.UUID, *store.StoreDiscountPayload) (*store.StoreDiscount, error)
	DeleteStoreDiscounts(context.Context, uuid.UUID, []uuid.UUID) error
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

func (repo *StoreRepo) SearchStores(ctx context.Context, params *store.FilterCriteria) ([]store.Store, error) {
	var query bytes.Buffer
	var args []interface{}
	argIndex := 1

	query.WriteString(`
    SELECT
        id,
        username,
        name,
        address,
        contact_number,
        email,
        web_url,
        password_digest,
        discounts,
        COALESCE(ST_Y(GEOMETRY(location_point)), 0) AS latitude,
        COALESCE(ST_X(GEOMETRY(location_point)), 0) AS longitude,
        created_at,
        updated_at
    FROM store
    WHERE 1=1
    `)

	if params.Name != "" {
		query.WriteString(fmt.Sprintf(" AND name ILIKE '%%%s%%'", params.Name))
		args = append(args, params.Name)
		argIndex++
	}

	if len(params.Ids) > 0 {
		query.WriteString(fmt.Sprintf(" AND id = ANY($%d)", argIndex))
		args = append(args, params.Ids)
		argIndex++
	}

	if params.Feature.Type == geocode.FeatureType {
		if params.Feature.Geometry != nil {
			if gj, err := json.Marshal(*params.Feature.Geometry); err != nil {
				return nil, err
			} else {
				switch params.Feature.Geometry.Type {
				case geocode.PointType:
					query.WriteString(fmt.Sprintf(" AND ST_DWithin(location_point::geography, ST_GeomFromGeoJSON($%d)::geography, $%d)", argIndex, argIndex+1))
					args = append(args, gj, params.Radius)
					argIndex += 2
				case geocode.LineStringType:
					query.WriteString(fmt.Sprintf(" AND ST_DWithin(location_point::geography, ST_GeomFromGeoJSON($%d)::geography, $%d)", argIndex, argIndex+1))
					args = append(args, gj, params.Radius)
					argIndex += 2
				case geocode.PolygonType:
					query.WriteString(fmt.Sprintf(" AND ST_Contains(ST_GeomFromGeoJSON($%d)::geography, location_point::geography)", argIndex))
					args = append(args, gj)
					argIndex++
				default:
					return nil, fmt.Errorf("unsupported geometry type")
				}
			}
		}
	}

	query.WriteString(fmt.Sprintf(" OFFSET $%d LIMIT $%d", argIndex, argIndex+1))
	args = append(args, params.Offset, params.Limit)

	finalQuery := query.String()
	rows, err := repo.storage.Pool.Query(
		ctx,
		finalQuery,
		args...,
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
			&store.ContactNumber,
			&store.Email,
			&store.WebUrl,
			&store.PasswordDigest,
			&store.Discounts,
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

// store products--------------------------------------------------------------------------------------
func (repo *StoreRepo) GetStoreProducts(ctx context.Context, storeId string) (json.RawMessage, error) {
	params := ctx.Value(search.SearchKey).(search.Paginate)
	query := `
SELECT
    jsonb_build_object(
        'items', jsonb_agg(
            jsonb_build_object(
                'productId', sp.product_id,
                'productName', p.name,
                'quantity', p.product_quantity,
                'currency', sp.currency,
                'stockQuantity', sp.stock_quantity,
                'pricePerUnit', sp.price_per_unit,
                'listedUnitOfMeasure', sp.listed_unit_of_measure,
                'imageUrl', p.image_url,
                'standardUnitOfMeasure', p.unit_of_measure,
                'updatedAt', sp.updated_at
            )
        ),
        'numberOfItems', COUNT(sp.product_id)
    ) AS storeItems
FROM 
    store_product sp
JOIN product p ON sp.product_id = p.id
WHERE sp.store_id = $1
OFFSET $2
LIMIT $3;
  `
	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		storeId,
		params.Offset,
		params.Limit,
	)

	var jsonData json.RawMessage
	if err := row.Scan(&jsonData); err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (repo *StoreRepo) GetTotalNumberOfProductsInStore(ctx context.Context, storeId string) (uint, error) {
	query := `
	SELECT COUNT(product_id) FROM store_product WHERE store_id = $1;
	`
	var total uint
	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		storeId,
	)
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (repo *StoreRepo) CreateStoreProduct(ctx context.Context, storeId string, payload *store.StoreProductPayload) (json.RawMessage, error) {
	query := `
WITH inserted_product AS (
    INSERT INTO store_product (
        store_id,
        product_id,
        price_per_unit,
        listed_unit_of_measure,
        stock_quantity
    )
    VALUES (
        $1, $2, $3, $4, $5
    )
    RETURNING *
)
SELECT
    jsonb_build_object(
        'productId', sp.product_id,
        'productName', p.name,
        'quantity', p.product_quantity,
        'currency', sp.currency,
        'stockQuantity', sp.stock_quantity,
        'pricePerUnit', sp.price_per_unit,
        'listedUnitOfMeasure', sp.listed_unit_of_measure,
        'imageUrl', p.image_url,
        'standardUnitOfMeasure', p.unit_of_measure,
        'updatedAt', sp.updated_at
    ) AS storeItem
FROM
    inserted_product sp
JOIN product p ON sp.product_id = p.id;
`
	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		storeId,
		payload.ProductId,
		payload.PricePerUnit,
		payload.ListUnitOfMeasure,
		payload.StockQuantity,
	)

	var jsonData json.RawMessage
	if err := row.Scan(&jsonData); err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (repo *StoreRepo) UpdateStoreProduct(ctx context.Context, storeId string, productId string, payload *store.StoreProductPayload) (json.RawMessage, error) {
	query := `
WITH updated_product AS (
	UPDATE store_product
	SET
		price_per_unit = $1,
		listed_unit_of_measure = $2,
		stock_quantity = $3
	WHERE
		store_id = $4
		AND product_id = $5
	RETURNING *
)
SELECT
	jsonb_build_object(
		'productId', sp.product_id,
		'productName', p.name,
		'quantity', p.product_quantity,
		'currency', sp.currency,
		'stockQuantity', sp.stock_quantity,
		'pricePerUnit', sp.price_per_unit,
		'listedUnitOfMeasure', sp.listed_unit_of_measure,
		'imageUrl', p.image_url,
		'standardUnitOfMeasure', p.unit_of_measure,
		'updatedAt', sp.updated_at
	) AS storeItem
FROM
	updated_product sp
JOIN product p ON sp.product_id = p.id;
`
	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		payload.PricePerUnit,
		payload.ListUnitOfMeasure,
		payload.StockQuantity,
		storeId,
		productId,
	)

	var jsonData json.RawMessage
	if err := row.Scan(&jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (repo *StoreRepo) DeleteStoreProducts(ctx context.Context, storeId string, productIds []string) error {
	query := `
		DELETE FROM store_product
		WHERE store_id = $1 AND product_id = ANY($2)
	`

	if _, err := repo.storage.Pool.Exec(ctx, query, storeId, productIds); err != nil {
		return err
	}

	return nil
}

// store discount handlers----------------------------------------------------------------
func (repo *StoreRepo) GetStoreDiscounts(ctx context.Context, storeId uuid.UUID) ([]store.StoreDiscount, error) {
	params := ctx.Value(search.SearchKey).(search.Paginate)
	query := `
	SELECT 
		id, 
		store_id, 
		discount, 
		applicable_tags,
		start_date,
		end_date,
		created_at, 
		updated_at
	FROM store_discount
	WHERE store_id = $1
	OFFSET $1
	LIMIT $2;
	`
	rows, err := repo.storage.Pool.Query(
		ctx,
		query,
		storeId,
		params.Offset,
		params.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var storeDiscounts []store.StoreDiscount
	for rows.Next() {
		var storeDiscount store.StoreDiscount
		if err := rows.Scan(
			&storeDiscount.Id,
			&storeDiscount.StoreId,
			&storeDiscount.Discount,
			&storeDiscount.ApplicableTags,
			&storeDiscount.StartDate,
			&storeDiscount.EndDate,
			&storeDiscount.CreatedAt,
			&storeDiscount.UpdatedAt,
		); err != nil {
			return nil, err
		}
		storeDiscounts = append(storeDiscounts, storeDiscount)
	}
	return storeDiscounts, nil
}

func (repo *StoreRepo) GetTotalNumberOfDiscountsInStore(ctx context.Context, storeId uuid.UUID) (uint, error) {
	query := `
	SELECT COUNT(id) FROM store_discount WHERE store_id = $1;
	`
	var total uint
	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		storeId,
	)
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (repo *StoreRepo) CreateStoreDiscount(ctx context.Context, storeId uuid.UUID, payload *store.StoreDiscountPayload) (*store.StoreDiscount, error) {
	query := `
		INSERT INTO store_discount (
			store_id,
			discount,
			applicable_tags,
			start_date,
			end_date
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			store_id,
			discount,
			applicable_tags,
			start_date,
			end_date,
			created_at,
			updated_at;
	`
	var storeDiscount store.StoreDiscount
	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		storeId,
		payload.Discount,
		payload.ApplicableTags,
		payload.StartDate,
		payload.EndDate,
	).Scan(
		&storeDiscount.Id,
		&storeDiscount.StoreId,
		&storeDiscount.Discount,
		&storeDiscount.ApplicableTags,
		&storeDiscount.StartDate,
		&storeDiscount.EndDate,
		&storeDiscount.CreatedAt,
		&storeDiscount.UpdatedAt,
	); err != nil {
		return nil, err
	} else {
		return &storeDiscount, nil
	}
}

func (repo *StoreRepo) UpdateStoreDiscount(ctx context.Context, storeId uuid.UUID, payload *store.StoreDiscountPayload) (*store.StoreDiscount, error) {
	query := `
		UPDATE store_discount
		SET
			discount = $1,
			applicable_tags = $2,
			start_date = $3,
			end_date = $4
		WHERE
			id = $5 AND store_id = $6
		RETURNING
			id,
			store_id,
			discount,
			applicable_tags,
			start_date,
			end_date,
			created_at,
			updated_at;
	`
	var storeDiscount store.StoreDiscount
	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		payload.Discount,
		payload.ApplicableTags,
		payload.StartDate,
		payload.EndDate,
		payload.Id,
		payload.StoreId,
	).Scan(
		&storeDiscount.Id,
		&storeDiscount.StoreId,
		&storeDiscount.Discount,
		&storeDiscount.ApplicableTags,
		&storeDiscount.StartDate,
		&storeDiscount.EndDate,
		&storeDiscount.CreatedAt,
		&storeDiscount.UpdatedAt,
	); err != nil {
		return nil, err
	} else {
		return &storeDiscount, nil
	}
}

func (repo *StoreRepo) DeleteStoreDiscounts(ctx context.Context, storeId uuid.UUID, ids []uuid.UUID) error {
	query := `
		DELETE FROM store_discount
		WHERE store_id = $1 AND id = ANY($2);
	`
	if _, err := repo.storage.Pool.Exec(
		ctx,
		query,
		storeId,
		ids,
	); err != nil {
		return err
	}
	return nil
}
