package types

import uuid "github.com/google/uuid"

type StoreSignUpPayload struct {
	StoreName     string `json:"store_name"`
	Address       string `json:"address"`
	Email         string `json:"store_email"`
	ContactNumber string `json:"contact_number"` // this should be a list of numbers
	Password      string `json:"password"`
}

type StoreLoginPayload struct {
	StoreUsername string `json:"store_username"`
	Password      string `json:"password"`
}

type StoreUpdatePayload struct {
	StoreUsername      string `json:"store_username" db:"store_username"`
	StoreName          string `json:"store_name" db:"store_name"`
	StoreEmail         string `json:"store_email" db:"store_email"`
	StoreContactNumber string `json:"store_contact_number" db:"store_contact_number"`
}

type Store struct {
	Id                 uuid.UUID `json:"id" db:"id"`
	StoreUsername      string    `json:"store_username" db:"store_username"`
	StoreName          string    `json:"store_name" db:"store_name"`
	StoreAddress       string    `json:"store_address" db:"store_address"`
	StoreEmail         string    `json:"store_email" db:"store_email"`
	StoreContactNumber string    `json:"store_contact_number" db:"store_contact_number"`
	Password           string    `json:"password_digest" db:"password_digest"`
	// StoreLocation      string    `json:"store_location" db:"store_location"`
}
