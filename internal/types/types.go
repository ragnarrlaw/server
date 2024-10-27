package types

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
  "github.com/google/uuid"
)

// contextKey is a type used for context keys to avoid conflicts
type ContextKey string

const (
	UserIDKey       ContextKey = "user_id"
	RefreshTokenKey ContextKey = "refresh_token"
	AccessTokenKey  ContextKey = "access_token"
)

type User struct {
	Id            uuid.UUID `json:"id" db:"id"`
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

type Token struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

func (t *Token) String() string {
	return fmt.Sprintf("Token: { Type: %s, AccessToken: %s }", t.Type, t.AccessToken)
}

type TokenClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func (tc *TokenClaims) String() string {
	return fmt.Sprintf("TokenClaims: { UserId: %s }", tc.UserId)
}

type AuthToken struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

func (at *AuthToken) String() string {
	return fmt.Sprintf("AuthToken: { Id: %s, UserId: %s, Token: %s }", at.Id, at.UserId, at.Token)
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

func (uup *UserUpdatePayload) String() string {
	return fmt.Sprintf("UserUpdatePayload: { Username: %s, FirstName: %s, LastName: %s, Email: %s, ContactNumber: %s, Password: %s, RePassword: %s }", uup.Username, uup.FirstName, uup.LastName, uup.Email, uup.ContactNumber, uup.Password, uup.RePassword)
}

type UserLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ulp *UserLoginPayload) String() string {
	return fmt.Sprintf("UserLoginPayload: { Username: %s, Password: %s }", ulp.Username, ulp.Password)
}

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

type StoreUpdatePayload struct{}

type Store struct {
	Id                 string `json:"id" db:"id"`
	StoreName          string `json:"store_name" db:"store_name"`
	StoreAddress       string `json:"store_address" db:"store_address"`
	StoreEmail         string `json:"store_email" db:"store_email"`
	StoreContactNumber string `json:"store_contact_number" db:"store_contact_number"`
	Password           string `json:"password_digest" db:"password_digest"`
	StoreLocation      string `json:"store_location" db:"store_location"`
}
