package recommenderrepository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ragnarrlaw/server/db/database"
	"github.com/ragnarrlaw/server/internal/types/recommender"
	uuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

// NOTE: ST_MakePoint => x = longitude $1, y = latitude $2

type RecommenderRepository interface {
	Search(context.Context, *recommender.Search) ([]recommender.Recommendation, error)
	LineStringSearch(context.Context, *recommender.LineStringSearch) ([]recommender.Recommendation, error)
}

type RecommenderRepo struct {
	storage *database.Storage
}

func NewRecommenderRepo(storage *database.Storage) RecommenderRepository {
	return &RecommenderRepo{storage: storage}
}

func (repo *RecommenderRepo) Search(ctx context.Context, payload *recommender.Search) ([]recommender.Recommendation, error) {
	query := `
SELECT 
    jsonb_build_object(
        'id', s.id,
        'name', s.name,
        'address', s.address,
        'items', jsonb_agg(
            jsonb_build_object(
                'productId', sp.product_id,
                'productName', p.name,
                'quantity', p.product_quantity,
                'stockQuantity', sp.stock_quantity,
                'pricePerUnit', sp.price_per_unit,
                'listedUnitOfMeasure', sp.listed_unit_of_measure,
                'imageUrl', p.image_url,
                'standardUnitOfMeasure', p.unit_of_measure
            )
        ),
		'location', json_build_object(
        	'latitude', ST_Y(s.location_point::geometry),
        	'longitude', ST_X(s.location_point::geometry)
		),
		'discounts', s.discounts,
		'numberOfItems', COUNT(sp.product_id)
    ) AS store_json
FROM 
    store s
JOIN 
    store_product sp ON sp.store_id = s.id
JOIN 
    product p ON sp.product_id = p.id 
WHERE 
    ST_DWithin(
        s.location_point::geography, 
        ST_MakePoint($1, $2)::geography, 
        $3
    )
	AND
	sp.product_id = ANY($4)
  	AND
  	(sp.stock_quantity = 'AVAILABLE' OR sp.stock_quantity = 'LIMITED_QUANTITY')
GROUP BY 
    s.id;
  `
	var productIds []uuid.UUID

	for _, p := range payload.ListOfItems {
		productIds = append(productIds, uuid.UUID(p.ProductId))
	}

	rows, err := repo.storage.Pool.Query(
		ctx,
		query,
		payload.UserLocation.Longitude,
		payload.UserLocation.Latitude,
		2000,
		productIds,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var recommendations []recommender.Recommendation

	for rows.Next() {
		var b []byte
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		var recommendation recommender.Recommendation
		if err := json.Unmarshal(b, &recommendation); err != nil {
			return nil, err
		}
		recommendations = append(recommendations, recommendation)
	}
	return recommendations, nil
}

func (repo *RecommenderRepo) LineStringSearch(
	ctx context.Context,
	payload *recommender.LineStringSearch,
) ([]recommender.Recommendation, error) {
	var productIds []uuid.UUID

	for _, p := range payload.ListOfItems {
		productIds = append(productIds, uuid.UUID(p.ProductId))
	}

	var points []string

	for _, location := range payload.Route {
		points = append(points, fmt.Sprintf("ST_MakePoint(%f, %f)", location.Longitude, location.Latitude))
	}

	lineString := fmt.Sprintf("ST_MakeLine(ARRAY[%s]::geometry[])", strings.Join(points, ", "))

	query := fmt.Sprintf(`
SELECT 
    jsonb_build_object(
        'id', s.id,
        'name', s.name,
        'address', s.address,
        'items', jsonb_agg(
            jsonb_build_object(
                'productId', sp.product_id,
                'productName', p.name,
                'quantity', p.product_quantity,
                'stockQuantity', sp.stock_quantity,
                'pricePerUnit', sp.price_per_unit,
                'listedUnitOfMeasure', sp.listed_unit_of_measure,
                'imageUrl', p.image_url,
                'standardUnitOfMeasure', p.unit_of_measure
            )
        ),
		'location', json_build_object(
        	'latitude', ST_Y(s.location_point::geometry),
        	'longitude', ST_X(s.location_point::geometry)
		),
		'discounts', s.discounts,
		'numberOfItems', COUNT(sp.product_id)
    ) AS store_json
FROM 
    store s
JOIN 
    store_product sp ON sp.store_id = s.id
JOIN 
    product p ON sp.product_id = p.id 
WHERE 
    ST_DWithin(
        %s,
        s.location_point::geography,
        $1
    )
	AND
	sp.product_id = ANY($2)
  	AND
  	(sp.stock_quantity = 'AVAILABLE' OR sp.stock_quantity = 'LIMITED_QUANTITY')
GROUP BY 
    s.id;
  `, lineString)

	rows, err := repo.storage.Pool.Query(
		ctx,
		query,
		500,
		productIds,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var recommendations []recommender.Recommendation

	for rows.Next() {
		var b []byte
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		var recommendation recommender.Recommendation
		if err := json.Unmarshal(b, &recommendation); err != nil {
			return nil, err
		}
		recommendations = append(recommendations, recommendation)
	}
	return recommendations, nil
}

/* Implement the preprocessor */
func preprocessor() {
	panic("function not implemented")
}

/* Implement the single store mode */
func singleStoreMode() {
	panic("function not implemented")
}

/* Implement the multi store mode */
func multiStoreMode() {
	panic("function not implemented")
}
