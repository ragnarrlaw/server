package types

import (
	"fmt"

	uuid "github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID       `json:"id" db:"id"`
	Name        string          `json:"name" db:"name"`
	Description string          `json:"description" db:"description"`
	Brand       string          `json:"brand" db:"brand"`
	Category    ProductCategory `json:"category" db:"category"`
	ProductInFo ProductInfo     `json:"productInfo" db:"product_info"`
	Discount    []Discount      `json:"discount" db:"discount"`
}

func (p *Product) String() string {
	return fmt.Sprintf("Product: { Id: %s, Name: %s, Description: %s, Brand: %s, Category: %s }", p.Id, p.Name, p.Description, p.Brand, p.Category)
}

type ProductInfo struct {
	Id            uuid.UUID `json:"id" db:"id"`
	StoreId       uuid.UUID `json:"storeId" db:"store_id"`
	ProductId     uuid.UUID `json:"productId" db:"product_id"`
	Price         float64   `json:"price" db:"price"`
	StockQuantity int       `json:"stockQuantity" db:"stock_quantity"`
	UnitOfMeasure string    `json:"unitOfMeasure" db:"unit_of_measure"`
	Currency      string    `json:"currency" db:"currency"`
}

func (p *ProductInfo) String() string {
	return fmt.Sprintf("ProductInfo: { Id: %s, StoreId: %s, ProductId: %s, Price: %f, StockQuantity: %d, UnitOfMeasure: %s, Currency: %s }", p.Id, p.StoreId, p.ProductId, p.Price, p.StockQuantity, p.UnitOfMeasure, p.Currency)
}

type DiscountPayload struct {
	DiscountType string  `json:"discountType" db:"discount_type"`
	Value        float64 `json:"value" db:"value"`
	StartDate    string  `json:"startDate" db:"start_date"`
	EndDate      string  `json:"endDate" db:"end_date"`
}

func (dp *DiscountPayload) String() string {
	return fmt.Sprintf("DiscountPayload: { DiscountType: %s, Value: %f, StartDate: %s, EndDate: %s }", dp.DiscountType, dp.Value, dp.StartDate, dp.EndDate)
}

type Discount struct {
	Id           uuid.UUID `json:"id" db:"id"`
	ProductId    uuid.UUID `json:"productId" db:"product_id"`
	StoreId      uuid.UUID `json:"storeId" db:"store_id"`
	DiscountType string    `json:"discountType" db:"discount_type"`
	Value        float64   `json:"value" db:"value"`
	StartDate    string    `json:"startDate" db:"start_date"`
	EndDate      string    `json:"endDate" db:"end_date"`
}

func (d *Discount) String() string {
	return fmt.Sprintf("Discount: { Id: %s, ProductId: %s, StoreId: %s, DiscountType: %s, Value: %f, StartDate: %s, EndDate: %s }", d.Id, d.ProductId, d.StoreId, d.DiscountType, d.Value, d.StartDate, d.EndDate)
}

type ProductCategory struct {
	Id               uuid.UUID `json:"id" db:"id"`
	Category         string    `json:"category" db:"category"`
	ParentCategoryId string    `json:"parentCategoryId" db:"parent_category_id"`
}

func (pc *ProductCategory) String() string {
	return fmt.Sprintf("ProductCategory: { Id: %s, Category: %s, ParentCategoryId: %s }", pc.Id, pc.Category, pc.ParentCategoryId)
}

// for creating and updating the product information
type ProductPayload struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Brand         string  `json:"brand"`
	CategoryId    string  `json:"categoryId"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`
	UnitOfMeasure string  `json:"unitOfMeasure"`
	Discount      []struct {
		DiscountType string  `json:"discountType"`
		Value        float64 `json:"value"`
		StartDate    string  `json:"startDate"`
		EndDate      string  `json:"endDate"`
	} `json:"discount"`
}

func (pcp *ProductPayload) String() string {
	return fmt.Sprintf("ProductCreatePayload: { Name: %s, Description: %s, Brand: %s, CategoryId: %s, Quantity: %d, Price: %f, UnitOfMeasure: %s, Discount: %v }", pcp.Name, pcp.Description, pcp.Brand, pcp.CategoryId, pcp.Quantity, pcp.Price, pcp.UnitOfMeasure, pcp.Discount)
}
