package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/ragnarrlaw/server/internal/types/product"
)

type StoreProductPayload struct {
	ProductId         uuid.UUID             `json:"productId"`
	PricePerUnit      float64               `json:"pricePerUnit"`
	Currency          string                `json:"currency"`
	ListUnitOfMeasure product.UnitOfMeasure `json:"listUnitOfMeasure"`
	StockQuantity     product.StockQuantity `json:"stockQuantity"`
}

type StoreDiscountPayload struct {
	Id             uuid.UUID `json:"id"`
	StoreId        uuid.UUID `json:"storeId"`
	Discount       string    `json:"discount"`
	ApplicableTags []string  `json:"applicableTags"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
}
