package types

import "fmt"

type User struct {
	Id            string `json:"id" db:"id"`
	Username      string `json:"username" db:"username"`
	FirstName     string `json:"firstName" db:"first_name"`
	LastName      string `json:"lastName" db:"last_name"`
	Email         string `json:"email" db:"email"`
	ContactNumber string `json:"contactNumber" db:"contact_number"`
	Password      string `json:"password" db:"password_digest"`
}

func (u *User) String() string {
	return fmt.Sprintf("User: { Id: %s, Username: %s, FirstName: %s, LastName: %s, Email: %s, ContactNumber: %s }", u.Id, u.Username, u.FirstName, u.LastName, u.Email, u.ContactNumber)
}

type AuthToken struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type UserSignUpPayload struct {
	Username      string `json:"username"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
	ContactNumber string `json:"contactNumber"`
	Password      string `json:"password"`
	RePassword    string `json:"rePassword"`
}

type UserUpdatePayload struct {
	Username      string `json:"username"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
	ContactNumber string `json:"contactNumber"`
	Password      string `json:"password"`
	RePassword    string `json:"rePassword"`
	// ProfilePic    string `json:"profile_pic"` // this should be a file
}

type UserLoginPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type StoreRegisterPayload struct {
	StoreName     string `json:"storeName"`
	Address       string `json:"address"`
	Email         string `json:"storeEmail"`
	ContactNumber string `json:"contactNumber"` // this should be a list of numbers
	Password      string `json:"password"`
}

type StoreLoginPayload struct {
	StoreName  string `json:"storeName"`
	StoreEmail string `json:"storeEmail"`
	Password   string `json:"password"`
}
