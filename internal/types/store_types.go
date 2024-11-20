package types

import (
	"fmt"

	uuid "github.com/google/uuid"
)

type StoreSignUpPayload struct {
	StoreName     string `json:"storeName"`
	Username      string `json:"username"`
	Address       string `json:"address"`
	Email         string `json:"storeEmail"`
	ContactNumber string `json:"contactNumber"`
	WebURL        string `json:"storeWebUrl"`
	Password      string `json:"password"`
}

type StoreLoginPayload struct {
	Username string     `json:"username"`
	Password string     `json:"password"`
	Role     EntityType `json:"role"`
}

type StoreUpdatePayload struct {
	StoreUsername      string `json:"storeUsername" db:"store_username"`
	StoreName          string `json:"storeName" db:"store_name"`
	StoreEmail         string `json:"storeEmail" db:"store_email"`
	StoreWebURL        string `json:"storeWebUrl" db:"store_web_url"`
	StoreContactNumber string `json:"storeContactNumber" db:"store_contact_number"`
}

type Store struct {
	Id                 uuid.UUID `json:"id" db:"id"`
	StoreUsername      string    `json:"storeUsername" db:"store_username"`
	StoreName          string    `json:"storeName" db:"store_name"`
	StoreAddress       string    `json:"storeAddress" db:"store_address"`
	StoreEmail         string    `json:"storeEmail" db:"store_email"`
	StoreContactNumber string    `json:"storeContactNumber" db:"store_contact_number"`
	StoreWebURL        string    `json:"storeWebUrl" db:"store_web_url"`
	Password           string    `db:"password_digest"`
	StoreLocation      GeoPoint  `json:"storeLocation" db:"store_location"`
}

func (s *Store) String() string {
	return fmt.Sprintf("Store: { Id: %s, StoreUsername: %s, StoreName: %s, StoreAddress: %s, StoreEmail: %s, StoreContactNumber: %s, StoreWebUrl: %s, StoreLocation: {latitude: %v, longitude: %v} }", s.Id, s.StoreUsername, s.StoreName, s.StoreAddress, s.StoreEmail, s.StoreContactNumber, s.StoreWebURL, s.StoreLocation.Latitude, s.StoreLocation.Longitude)
}

func (s *Store) GetId() string {
	return s.Id.String()
}

func (s *Store) GetUsername() string {
	return s.StoreUsername
}
