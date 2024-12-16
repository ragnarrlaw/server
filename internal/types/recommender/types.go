package recommender

import (
	"context"
	"errors"

	"github.com/google/uuid"
	pgtype "github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/raganrrlaw/server/internal/types/geocode"
)

type ListItem struct {
	ProductId uuid.UUID `json:"productId"`
	Quantity  float64   `json:"quantity"`
}

type Search struct {
	UserLocation geocode.Point `json:"userLocation"`
	ListOfItems  []ListItem    `json:"listOfItems"`
}

type LineStringSearch struct {
	Route       []geocode.Point `json:"route"`
	ListOfItems []ListItem      `json:"listOfItems"`
}

type Recommendation struct {
	Id             uuid.UUID            `json:"id"`
	Name           string               `json:"name"`
	Address        string               `json:"address"`
	Location       geocode.Point        `json:"location"`
	Distance       float64              `json:"distance"`
	TotalCost      float64              `json:"totalCost"`
	DiscountedCost float64              `json:"discountedCost"`
	Items          []RecommendationItem `json:"items"`
	NumberOfItems  int                  `json:"numberOfItems"`
}

func RegisterRecommendationType(ctx context.Context, conn *pgx.Conn) error {
	datatype, err := conn.LoadType(ctx, "recommendation_t")
	if err != nil {
		return err
	}
	conn.TypeMap().RegisterType(datatype)
	return nil
}

type RecommendationItem struct {
	ProductId             uuid.UUID `json:"productId"`
	Name                  string    `json:"productName"`
	Quantity              string    `json:"quantity"`
	PricePerUnit          float64   `json:"pricePerUnit"`
	ListedUnitOfMeasure   string    `json:"listedUnitOfMeasure"`
	StockQuantity         string    `json:"stockQuantity"`
	ImageUrl              string    `json:"imageUrl"`
	StandardUnitOfMeasure string    `json:"standardUnitOfMeasure"`
}

func RegisterRecommendationItemType(ctx context.Context, conn *pgx.Conn) error {
	datatype, err := conn.LoadType(ctx, "recommendation_item_t")
	if err != nil {
		return err
	}
	conn.TypeMap().RegisterType(datatype)
	return nil
}

func (dst *RecommendationItem) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		return errors.New("null values cannot be decoded")
	}
	if err := (pgtype.CompositeFields{
		&dst.ProductId,
		&dst.PricePerUnit,
		&dst.ListedUnitOfMeasure,
		&dst.StockQuantity,
		&dst.Name,
		&dst.ImageUrl,
		&dst.StandardUnitOfMeasure,
	}).DecodeBinary(ci, src); err != nil {
		return err
	}
	return nil
}
