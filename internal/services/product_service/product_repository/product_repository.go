package productrepository

import (
	"context"

	"github.com/ragnarrlaw/server/db/database"
	"github.com/ragnarrlaw/server/internal/types/product"
	"github.com/ragnarrlaw/server/internal/types/search"
)

type ProductRepository interface {
	// product related
	Get(ctx context.Context) (*[]product.Product, error)
	GetById(ctx context.Context, productId string) (*product.Product, error)
	SearchProductByName(ctx context.Context, value any) (*[]product.Product, error)
	GetTotalNumberOfProducts(context.Context) (uint, error)
}

type ProductRepo struct {
	storage *database.Storage
}

func NewProductRepo(storage *database.Storage) ProductRepository {
	return &ProductRepo{storage: storage}
}

func (repo *ProductRepo) GetTotalNumberOfProducts(ctx context.Context) (uint, error) {
	query := `
	SELECT COUNT(id) FROM product;
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

func (repo *ProductRepo) Get(ctx context.Context) (*[]product.Product, error) {
	params := ctx.Value(search.SearchKey).(search.Paginate)
	query := `
		SELECT
    id, 
    code, 
    name, 
    description, 
    brand,
    category_id,
    labels,
    image_url,
    product_quantity,
    serving_size,
    unit_of_measure,
    created_at,
    updated_at
		FROM product OFFSET $1 LIMIT $2`
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

	var products []product.Product
	for rows.Next() {
		var p product.Product
		if err := rows.Scan(
			&p.Id,
			&p.Code,
			&p.Name,
			&p.Description,
			&p.Brand,
			&p.CategoryId,
			&p.Labels,
			&p.ImageUrl,
			&p.ProductQuantity,
			&p.ServingSize,
			&p.UnitOfMeasure,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		} else {
			products = append(products, p)
		}
	}
	return &products, nil
}

func (repo *ProductRepo) GetById(ctx context.Context, productId string) (*product.Product, error) {
	query := `
		SELECT
    id, 
    code, 
    name, 
    description, 
    brand,
    category_id,
    labels,
    image_url,
    product_quantity,
    serving_size,
    unit_of_measure,
    created_at,
    updated_at
		FROM product`

	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		productId,
	)
	var p product.Product
	if err := row.Scan(
		&p.Id,
		&p.Code,
		&p.Name,
		&p.Description,
		&p.Brand,
		&p.CategoryId,
		&p.Labels,
		&p.ImageUrl,
		&p.ProductQuantity,
		&p.ServingSize,
		&p.UnitOfMeasure,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return nil, err
	} else {
		return &p, nil
	}
}

func (repo *ProductRepo) SearchProductByName(ctx context.Context, term any) (*[]product.Product, error) {
	params := ctx.Value(search.SearchKey).(search.Paginate)
	query := `
		SELECT
    id, 
    code, 
    name, 
    description, 
    brand,
    category_id,
    labels,
    image_url,
    product_quantity,
    serving_size,
    unit_of_measure,
    created_at,
    updated_at
    FROM product 
		WHERE name ILIKE '%%' || $1 || '%%'
		LIMIT $2;
	`
	rows, err := repo.storage.Pool.Query(
		ctx,
		query,
		term,
		params.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []product.Product
	for rows.Next() {
		var p product.Product
		if err := rows.Scan(
			&p.Id,
			&p.Code,
			&p.Name,
			&p.Description,
			&p.Brand,
			&p.CategoryId,
			&p.Labels,
			&p.ImageUrl,
			&p.ProductQuantity,
			&p.ServingSize,
			&p.UnitOfMeasure,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		} else {
			products = append(products, p)
		}
	}
	return &products, nil
}
