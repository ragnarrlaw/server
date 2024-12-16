package product

import (
	"time"

	"github.com/google/uuid"
)

type UnitOfMeasure string

const (
	UnitOfMeasureGrams UnitOfMeasure = "g"
	UnitOfMeasureKg    UnitOfMeasure = "kg"
	UnitOfMeasureMl    UnitOfMeasure = "ml"
	UnitOfMeasureL     UnitOfMeasure = "l"
	UnitOfMeasurePcs   UnitOfMeasure = "pcs"
)

type StockQuantity string

const (
	LimitedQuantity StockQuantity = "LIMITED_QUANTITY"
	Available       StockQuantity = "AVAILABLE"
	NotAvailable    StockQuantity = "NOT_AVAILABLE"
)

type Category struct {
	Id               uuid.UUID `json:"id"`
	Category         string    `json:"category"`
	ParentCategoryId uuid.UUID `json:"parentCategoryId"`
	CreatedAt        string
	UpdatedAt        string
}

type Product struct {
	Id              uuid.UUID     `json:"id"`
	Code            string        `json:"code"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	Brand           string        `json:"brand"`
	BrandTags       string        `json:"brandTags"`
	CategoryId      uuid.UUID     `json:"categoryId"`
	Labels          string        `json:"labels"`
	ImageUrl        string        `json:"imageUrl"`
	ProductQuantity string        `json:"productQuantity"`
	ServingSize     string        `json:"servingSize"`
	UnitOfMeasure   UnitOfMeasure `json:"unitOfMeasure"`
	CreatedAt       time.Time     `json:"-"`
	UpdatedAt       time.Time     `json:"-"`
}

type StoreProduct struct {
	Id                uuid.UUID     `json:"id"`
	StoreId           uuid.UUID     `json:"storeId"`
	ProductId         uuid.UUID     `json:"productId"`
	PricePerUnit      float64       `json:"pricePerUnit"`
	Currency          string        `json:"currency"`
	ListUnitOfMeasure UnitOfMeasure `json:"listUnitOfMeasure"`
	StockQuantity     int           `json:"stockQuantity"`
	CreatedAt         time.Time     `json:"-"`
	UpdatedAt         time.Time     `json:"-"`
}
