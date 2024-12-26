package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID `json:"-"`
	Username       string    `json:"username"`
	FistName       string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Email          string    `json:"email"`
	ContactNumber  string    `json:"contactNumber"`
	PasswordDigest string    `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

type UserPreferences struct {
	UserId      uuid.UUID `json:"-"`
	Preferences struct {
		PaymentMethod       string      `json:"paymentMethod"`
		DigitalPaymentCards []string    `json:"digitalPaymentCards"`
		PreferredStores     []uuid.UUID `json:"preferredStores"`
		PreferredBrands     []uuid.UUID `json:"preferredBrands"`
	} `json:"preferences"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type PreviousUserRecommendations struct {
	Id              uuid.UUID `json:"id"`
	UserId          uuid.UUID `json:"-"`
	Recommendations []struct {
		StoreId  uuid.UUID `json:"storeId"`
		Products []struct {
			ProductId               uuid.UUID `json:"productId"`
			Quantity                int       `json:"quantity"`
			OriginalPurchasePrice   float64   `json:"purchasePrice"`
			DiscountedPurchasePrice float64   `json:"discountedPurchasePrice"`
		} `json:"items"`
		TotalPurchasePrice float64 `json:"totalPurchasePrice"`
	} `json:"recommendations"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
