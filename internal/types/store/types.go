package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/ragnarrlaw/server/internal/types/geocode"
)

type Store struct {
	Id             uuid.UUID     `json:"id"`
	Username       string        `json:"username"`
	Name           string        `json:"name"`
	Address        string        `json:"address"`
	ContactNumber  string        `json:"contactNumber"`
	Email          string        `json:"email"`
	WebUrl         string        `json:"webUrl"`
	PasswordDigest string        `json:"-"`
	Discounts      []string      `json:"discounts"`
	LocationPoint  geocode.Point `json:"locationPoint"`
	CreatedAt      time.Time     `json:"-"`
	UpdatedAt      time.Time     `json:"-"`
}

type StoreStat struct {
	Name                         string        `json:"name"`
	Address                      string        `json:"address"`
	Discounts                    []string      `json:"discounts"`
	Point                        geocode.Point `json:"location"`
	CreatedAt                    time.Time     `json:"createdAt"`
	TotalNumberOfProducts        int           `json:"totalNumberOfProducts"`
	NumberOfLimitedStockProducts int           `json:"numberOfLimitedStockProducts"`
	NumberOfAvailableProducts    int           `json:"numberOfAvailableProducts"`
	NumberOfOutOfStockProducts   int           `json:"numberOfOutOfStockProducts"`
}

type StoreDiscount struct {
	Id             uuid.UUID `json:"id"`
	StoreId        uuid.UUID `json:"storeId"`
	Discount       string    `json:"discount"`
	ApplicableTags []string  `json:"applicableTags"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type FilterCriteria struct {
	Name    string                 `json:"name"`
	Ids     []string               `json:"ids"`
	Feature geocode.GeoJsonFeature `json:"feature"`
	Radius  float64                `json:"radius"`
	Limit   int                    `json:"limit"`
	Offset  int                    `json:"offset"`
}
