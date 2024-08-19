package types

type User struct {
	ID            int    `json:"id" db:"id"`
	Username      string `json:"username" db:"username"`
	FirstName     string `json:"first_name" db:"first_name"`
	LastName      string `json:"last_name" db:"last_name"`
	Email         string `json:"email" db:"email"`
	ContactNumber string `json:"contact_number" db:"contact_number"`
	Password      string `json:"password" db:"password"`
}

type UserSignUpPayload struct {
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	ContactNumber string `json:"contact_number"`
	Password      string `json:"password"`
}

type UserLoginPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type StoreRegisterPayload struct {
	StoreName     string `json:"store_name"`
	Address       string `json:"address"`
	Email         string `json:"store_email"`
	ContactNumber string `json:"contact_number"` // this should be a list of numbers
	Password      string `json:"password"`
}

type StoreLoginPayload struct {
	StoreName  string `json:"store_name"`
	StoreEmail string `json:"store_email"`
	Password   string `json:"password"`
}
