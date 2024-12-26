package discountprocessor

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ragnarrlaw/server/internal/types/geocode"
	"github.com/ragnarrlaw/server/internal/types/product"
	"github.com/ragnarrlaw/server/internal/types/recommender"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateRecommendations(t *testing.T) {
	recommendations := []recommender.Recommendation{
		{
			Id:       uuid.New(),
			Name:     "Store 1",
			Address:  "123 Main St",
			Location: geocode.Point{Latitude: 40.7128, Longitude: -74.0060},
			Distance: 1.2,
			Discounts: []string{
				`product_id IN ["9f9285c6-a4d3-407e-9bd6-92ed094d0b02"] THEN product_percentage 10`,
			},
			TotalCost:      0,
			DiscountedCost: 0,
			Items: []recommender.RecommendationItem{
				{
					ProductId:             uuid.MustParse("9f9285c6-a4d3-407e-9bd6-92ed094d0b02"),
					Name:                  "Milk",
					Quantity:              "2",
					PricePerUnit:          50,
					ListedUnitOfMeasure:   product.UnitOfMeasureL,
					StockQuantity:         "10",
					ImageUrl:              "http://example.com/milk.jpg",
					StandardUnitOfMeasure: product.UnitOfMeasureL,
					Cost:                  100,
					PurchaseQuantity:      2,
					DiscountedPrice:       0,
				},
			},
			NumberOfItems: 1,
		},
	}

	err := EvaluateRecommendations(recommendations)
	assert.NoError(t, err)
	assert.Equal(t, 90.0, recommendations[0].Items[0].DiscountedPrice)
	assert.Equal(t, 90.0, recommendations[0].DiscountedCost)
}
