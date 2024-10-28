package types

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	Id            uuid.UUID `json:"id" db:"id"`
	Username      string    `json:"username" db:"username"`
	FirstName     string    `json:"firstName" db:"first_name"`
	LastName      string    `json:"lastName" db:"last_name"`
	Email         string    `json:"email" db:"email"`
	ContactNumber string    `json:"contactNumber" db:"contact_number"`
	Password      string    `json:"password" db:"password_digest"`
}

func (u *User) String() string {
	return fmt.Sprintf("User: { Id: %s, Username: %s, FirstName: %s, LastName: %s, Email: %s, ContactNumber: %s }", u.Id, u.Username, u.FirstName, u.LastName, u.Email, u.ContactNumber)
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

func (usp *UserSignUpPayload) String() string {
	return fmt.Sprintf("UserSignUpPayload: { Username: %s, FirstName: %s, LastName: %s, Email: %s, ContactNumber: %s, Password: %s, RePassword: %s }", usp.Username, usp.FirstName, usp.LastName, usp.Email, usp.ContactNumber, usp.Password, usp.RePassword)
}

/** 
  Password updates, and email, and contact number updates are handled by the authentication services
*/
type UserUpdatePayload struct {
	Username      string `json:"username"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	// ProfilePic    string `json:"profile_pic"` // this should be a file
}

func (uup *UserUpdatePayload) String() string {
	return fmt.Sprintf("UserUpdatePayload: { Username: %s, FirstName: %s, LastName: %s }", uup.Username, uup.FirstName, uup.LastName)
}

type UserLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ulp *UserLoginPayload) String() string {
	return fmt.Sprintf("UserLoginPayload: { Username: %s, Password: %s }", ulp.Username, ulp.Password)
}
