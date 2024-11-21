package productrepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/raganrrlaw/server/db/database"
	"github.com/raganrrlaw/server/internal/types"
)

type ProductRepository interface {
	Add(ctx context.Context, storeId string, product *types.ProductPayload) (*types.Product, error)
	GetById(ctx context.Context, storeId string, productId string) (*types.Product, error)
	GetAll(ctx context.Context, storeId string) (*[]types.Product, error)
	Remove(ctx context.Context, storeId string, productIds []string) error
	Update(ctx context.Context, storeId string, productId string, product *types.ProductPayload) (*types.Product, error)
	GetProductDiscounts(ctx context.Context, storeId string, productId string) (*[]types.Discount, error)
	AddProductDiscount(ctx context.Context, storeId string, productId string, payload *types.DiscountPayload) (*types.Discount, error)
	UpdateDiscount(ctx context.Context, storeId string, productId string, discountId string, payload *types.DiscountPayload) (*types.Discount, error)
	RemoveDiscount(ctx context.Context, storeId string, productId string, discountId []string) error
	RemoveAllDiscounts(ctx context.Context, storeId string, productId string) error
	ProductCategories(ctx context.Context) (*[]types.ProductCategory, error)
}

type ProductRepo struct {
	storage *database.Storage
}

func NewProductRepo(storage *database.Storage) ProductRepository {
	return &ProductRepo{storage: storage}
}

func (repo *ProductRepo) Add(ctx context.Context, storeId string, payload *types.ProductPayload) (*types.Product, error) {
	q1 := `
		INSERT INTO product (name, description, brand, category_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, name, description, brand, category_id
	`
	var product types.Product
	if tx, err := repo.storage.Pool.BeginTx(
		ctx,
		pgx.TxOptions{
			IsoLevel:       pgx.Serializable,
			DeferrableMode: pgx.NotDeferrable,
			AccessMode:     pgx.ReadWrite,
		},
	); err != nil {
		return nil, err
	} else {
		row := tx.QueryRow(
			ctx,
			q1,
			payload.Name,
			payload.Description,
			payload.Brand,
			payload.CategoryId,
		)
		if err := row.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Brand,
			&product.Category.Id,
		); err != nil {
			tx.Rollback(ctx)
			return nil, err
		} else {
			q2 := `
				INSERT INTO store_product (store_id, product_id, price, stock_quantity, unit_of_measurement)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, store_id, product_id, price, stock_quantity, unit_of_measurement, currency
			`
			var productInfo types.ProductInfo
			if err := tx.QueryRow(
				ctx,
				q2,
				storeId,
				product.Id,
				payload.Price,
				payload.Quantity,
				payload.UnitOfMeasure,
			).Scan(
				&productInfo.Id,
				&productInfo.StoreId,
				&productInfo.ProductId,
				&productInfo.Price,
				&productInfo.StockQuantity,
				&productInfo.UnitOfMeasure,
				&productInfo.Currency,
			); err != nil {
				tx.Rollback(ctx)
				return nil, err
			} else {
				product.ProductInFo = productInfo
				if count := len(payload.Discount); count > 0 {
					var discounts []types.Discount

					values := []string{}
					args := []interface{}{}
					for i, row := range payload.Discount {
						values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
						args = append(args, product.Id, storeId, row.DiscountType, row.Value, row.StartDate, row.EndDate)
					}
					q3 := fmt.Sprintf(`
					INSERT INTO discount (product_id, store_id, discount_type, value, start_date, end_date) 
					VALUES %s 
					RETURNING id, product_id, store_id, discount_type, value, start_date, end_date
					`, strings.Join(values, ""))
					if rows, err := tx.Query(ctx, q3, args...); err != nil {
						tx.Rollback(ctx)
						return nil, err
					} else {
						for rows.Next() {
							var discount types.Discount
							if err := rows.Scan(
								&discount.Id,
								&discount.ProductId,
								&discount.StoreId,
								&discount.DiscountType,
								&discount.Value,
								&discount.StartDate,
								&discount.EndDate,
							); err != nil {
								tx.Rollback(ctx)
								return nil, err
							} else {
								discounts = append(discounts, discount)
							}
						}
						product.Discount = discounts
					}
				} else {
					product.Discount = []types.Discount{}
				}
				q3 := `SELECT id, category, parent_category_id FROM category WHERE id = $1`
				if err := tx.QueryRow(
					ctx,
					q3,
					product.Category.Id,
				).Scan(
					&product.Category.Id,
					&product.Category.Category,
					&product.Category.ParentCategoryId,
				); err != nil {
					tx.Rollback(ctx)
					return nil, err
				} else {
					if err := tx.Commit(ctx); err != nil {
						return nil, err
					} else {
						fmt.Println("product: ", product)
						return &product, nil
					}
				}
			}
		}
	}
}

func (repo *ProductRepo) GetById(ctx context.Context, storeId string, productId string) (*types.Product, error) {
	query := `
		SELECT p.id, p.name, p.description, p.brand,
			pc.id, pc.category, pc.parent_category_id,
			pi.id, pi.store_id, pi.product_id, pi.price, pi.stock_quantity, pi.unit_of_measurement, pi.currency
		FROM product p
		JOIN category pc ON p.category_id = pc.id
		JOIN store_product pi ON p.id = pi.product_id
		WHERE pi.store_id = $1 AND pi.product_id = $2
	`
	row := repo.storage.Pool.QueryRow(
		ctx,
		query,
		storeId,
		productId,
	)
	var product types.Product
	if err := row.Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.Brand,
		&product.Category.Id,
		&product.Category.Category,
		&product.Category.ParentCategoryId,
		&product.ProductInFo.Id,
		&product.ProductInFo.StoreId,
		&product.ProductInFo.ProductId,
		&product.ProductInFo.Price,
		&product.ProductInFo.StockQuantity,
		&product.ProductInFo.UnitOfMeasure,
		&product.ProductInFo.Currency,
	); err != nil {
		return nil, err
	} else {
		query := `
			SELECT id, discount_type, value, start_date, end_date
			FROM discount
			WHERE product_id = $1 AND store_id = $2
		`
		var discounts []types.Discount
		if rows, err := repo.storage.Pool.Query(ctx, query, productId, storeId); err != nil {
			return nil, err
		} else {
			for rows.Next() {
				var discount types.Discount
				if err := rows.Scan(
					&discount.Id,
					&discount.DiscountType,
					&discount.Value,
					&discount.StartDate,
					&discount.EndDate,
				); err != nil {
					return nil, err
				} else {
					discounts = append(discounts, discount)
				}
			}
			product.Discount = discounts
			return &product, nil
		}
	}
}

func (repo *ProductRepo) GetAll(ctx context.Context, storeId string) (*[]types.Product, error) {
	query := `
		SELECT p.id, p.name, p.description, p.brand,
			pc.id, pc.category, pc.parent_category_id,
			pi.id, pi.store_id, pi.product_id, pi.price, pi.stock_quantity, pi.unit_of_measurement, pi.currency
		FROM product p
		JOIN category pc ON p.category_id = pc.id
		JOIN store_product pi ON p.id = pi.product_id
		WHERE pi.store_id = $1
	`
	if rows, err := repo.storage.Pool.Query(ctx, query, storeId); err != nil {
		return nil, err
	} else {
		var products []types.Product
		for rows.Next() {
			var product types.Product
			if err := rows.Scan(
				&product.Id,
				&product.Name,
				&product.Description,
				&product.Brand,
				&product.Category.Id,
				&product.Category.Category,
				&product.Category.ParentCategoryId,
				&product.ProductInFo.Id,
				&product.ProductInFo.StoreId,
				&product.ProductInFo.ProductId,
				&product.ProductInFo.Price,
				&product.ProductInFo.StockQuantity,
				&product.ProductInFo.UnitOfMeasure,
				&product.ProductInFo.Currency,
			); err != nil {
				return nil, err
			} else {
				query := `
					SELECT id, discount_type, value, start_date, end_date
					FROM discount
					WHERE product_id = $1 AND store_id = $2
				`
				var discounts []types.Discount
				if rows, err := repo.storage.Pool.Query(ctx, query, product.Id, storeId); err != nil {
					return nil, err
				} else {
					for rows.Next() {
						var discount types.Discount
						if err := rows.Scan(
							&discount.Id,
							&discount.DiscountType,
							&discount.Value,
							&discount.StartDate,
							&discount.EndDate,
						); err != nil {
							return nil, err
						} else {
							discounts = append(discounts, discount)
						}
					}
					product.Discount = discounts
					products = append(products, product)
				}
			}
		}
		return &products, nil
	}
}

func (repo *ProductRepo) Remove(ctx context.Context, storeId string, productIds []string) error {
	query := `DELETE FROM product WHERE id = ANY($1)`
	if _, err := repo.storage.Pool.Exec(ctx, query, productIds); err != nil {
		return err
	} else {
		return nil
	}
}

func (repo *ProductRepo) Update(ctx context.Context, storeId string, productId string, p *types.ProductPayload) (*types.Product, error) {

	if tx, err := repo.storage.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		DeferrableMode: pgx.NotDeferrable,
		AccessMode:     pgx.ReadWrite,
	}); err != nil {

	} else {
		q1 := `
		UPDATE product SET 
			name = COALESCE(NULLIF($1, name), name) ,
			description = COALESCE(NULLIF($2, description), description),
			brand = COALESCE(NULLIF($3, brand), brand)
		WHERE id = $4
		RETURNING id, name, description, brand, category_id
	`
		var product types.Product
		if err := tx.QueryRow(
			ctx,
			q1,
			p.Name,
			p.Description,
			p.Brand,
			productId,
		).Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Brand,
			&product.Category.Id,
		); err != nil {
			tx.Rollback(ctx)
			return nil, err
		} else {
			q2 := `
			UPDATE store_product SET
				price = COALESCE(NULLIF($1, price), price),
				stock_quantity = COALESCE(NULLIF($2, stock_quantity), stock_quantity),
				unit_of_measurement = COALESCE(NULLIF($3, unit_of_measurement), unit_of_measurement)
			WHERE product_id = $4 AND store_id = $5
			RETURNING id, store_id, product_id, price, stock_quantity, unit_of_measurement, currency
			`
			var productInfo types.ProductInfo
			if err := tx.QueryRow(
				ctx,
				q2,
				p.Price,
				p.Quantity,
				p.UnitOfMeasure,
				productId,
				storeId,
			).Scan(
				&productInfo.Id,
				&productInfo.StoreId,
				&productInfo.ProductId,
				&productInfo.Price,
				&productInfo.StockQuantity,
				&productInfo.UnitOfMeasure,
				&productInfo.Currency,
			); err != nil {
				tx.Rollback(ctx)
				return nil, err
			} else {
				product.ProductInFo = productInfo
				tx.Commit(ctx)
				return &product, nil
			}
		}
	}

	return nil, nil
}

func (repo *ProductRepo) GetProductDiscounts(ctx context.Context, storeId string, productId string) (*[]types.Discount, error) {
	query := `
		SELECT
			id, store_id, product_id, discount_type, value, start_date, end_date
		FROM discount
		WHERE store_id = $1 AND product_id = $2
	`
	var discounts []types.Discount
	if rows, err := repo.storage.Pool.Query(
		ctx,
		query,
		storeId,
		productId,
	); err != nil {
		return nil, err
	} else {
		for rows.Next() {
			var discount types.Discount
			if err := rows.Scan(
				&discount.Id,
				&discount.StoreId,
				&discount.ProductId,
				&discount.DiscountType,
				&discount.Value,
				&discount.StartDate,
				&discount.EndDate,
			); err != nil {
				return nil, err
			} else {
				discounts = append(discounts, discount)
			}
		}
		return &discounts, nil
	}
}

func (repo *ProductRepo) AddProductDiscount(ctx context.Context, storeId string, productId string, payload *types.DiscountPayload) (*types.Discount, error) {
	query := `
		INSERT INTO discount (product_id, store_id, discount_type, value, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, product_id, store_id, discount_type, value, start_date, end_date
	`
	var discount types.Discount
	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		productId,
		storeId,
		payload.DiscountType,
		payload.Value,
		payload.StartDate,
		payload.EndDate,
	).Scan(
		&discount.Id,
		&discount.ProductId,
		&discount.StoreId,
		&discount.DiscountType,
		&discount.Value,
		&discount.StartDate,
		&discount.EndDate,
	); err != nil {
		return nil, err
	} else {
		return &discount, nil
	}
}

func (repo *ProductRepo) UpdateDiscount(ctx context.Context, storeId string, productId string, discountId string, p *types.DiscountPayload) (*types.Discount, error) {
	query := `
		UPDATE discount SET
			discount_type = COALESCE(NULLIF($1, discount_type), discount_type),
			value = COALESCE(NULLIF($2, value), value),
			start_date = COALESCE(NULLIF($3, start_date), start_date),
			end_date = COALESCE(NULLIF($4, end_date), end_date)
		WHERE product_id = $5 AND store_id = $6 AND id = $7
		RETURNING id, product_id, store_id, discount_type, value, start_date, end_date
	`
	var discount types.Discount
	if err := repo.storage.Pool.QueryRow(
		ctx,
		query,
		p.DiscountType,
		p.Value,
		p.StartDate,
		p.EndDate,
		productId,
		storeId,
		discountId,
	).Scan(
		&discount.Id,
		&discount.ProductId,
		&discount.StoreId,
		&discount.DiscountType,
		&discount.Value,
		&discount.StartDate,
		&discount.EndDate,
	); err != nil {
		return nil, err
	} else {
		return &discount, nil
	}
}

func (repo *ProductRepo) RemoveDiscount(ctx context.Context, storeId string, productId string, discountId []string) error {
	query := `DELETE FROM discount WHERE product_id = $1 AND store_id = $2 AND id = ANY($3)`
	if _, err := repo.storage.Pool.Exec(ctx, query, productId, storeId, discountId); err != nil {
		return err
	} else {
		return nil
	}
}

func (repo *ProductRepo) RemoveAllDiscounts(ctx context.Context, storeId string, productId string) error {
	query := `DELETE FROM discount WHERE product_id = $1 AND store_id = $2`
	if _, err := repo.storage.Pool.Exec(ctx, query, productId, storeId); err != nil {
		return err
	} else {
		return nil
	}
}

func (repo *ProductRepo) ProductCategories(ctx context.Context) (*[]types.ProductCategory, error) {
	query := `SELECT id, category, parent_category_id FROM category`
	if rows, err := repo.storage.Pool.Query(ctx, query); err != nil {
		return nil, err
	} else {
		var categories []types.ProductCategory
		for rows.Next() {
			var category types.ProductCategory
			if err := rows.Scan(
				&category.Id,
				&category.Category,
				&category.ParentCategoryId,
			); err != nil {
				return nil, err
			} else {
				categories = append(categories, category)
			}
		}
		return &categories, nil
	}
}
