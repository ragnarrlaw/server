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
	UserId    uuid.UUID `json:"-"`
	Items     string    `json:"items"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type PreviousUserRecommendations struct {
	Id              uuid.UUID `json:"id"`
	UserId          uuid.UUID `json:"-"`
	Recommendations string    `json:"recommendations"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}
