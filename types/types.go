package types

/*
  Add the types definitions here.
*/

type User struct {
	Id        string `json:"id" db:"id"`
	Username  string `json:"username" db:"username"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"passwd" db:"password"`
}

type AuthToken struct {
	Id           string `json:"id" db:"id"`
	UserId       string `json:"user_id" db:"user_id"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}

type GeoLocation struct {
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
	// TODO: Add other fields if there are any
}

type Store struct {
	Id       string      `json:"id" db:"id"`
	Name     string      `json:"name" db:"name"`
	Location GeoLocation `json:"-" db:"location"` // store location as a JSONB or a composite type
}

type Product struct {
	Id       string  `json:"id" db:"id"`
	StoreId  string  `json:"store_id" db:"store_id"`
	Name     string  `json:"name" db:"name"`
	Quantity uint32  `json:"quantity" db:"quantity"`
	Unit     string  `json:"unit_of_measure" db:"unit"`
	Price    float64 `json:"price" db:"price"`
}
