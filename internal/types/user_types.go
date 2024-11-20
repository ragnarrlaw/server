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

func (u *User) GetId() string {
	return u.Id.String()
}

func (u *User) GetUsername() string {
	return u.Username
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

/*
* Password updates, and email, and contact number updates are handled by the authentication services
 */
type UserUpdatePayload struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	// ProfilePic    string `json:"profile_pic"` // this should be a file
}

func (uup *UserUpdatePayload) String() string {
	return fmt.Sprintf("UserUpdatePayload: { Username: %s, FirstName: %s, LastName: %s }", uup.Username, uup.FirstName, uup.LastName)
}

/*
  - The preferences field is a jsonb field
  - payload example:
    { "userId": "uuid", "preferences": { "theme": "dark", "language": "en", preferredPaymentMethod: {card: [list of cards]} || 'cash' } }
  - Use string type for unstructured data unmarshaling
*/
type UserPreferences struct {
	UserId      uuid.UUID `json:"userId" db:"user_id"`
	Preferences string    `json:"preferences" db:"preferences"`
	CreatedAt   string    `json:"createdAt" db:"created_at"`
	UpdateAt    string    `json:"updatedAt" db:"updated_at"`
}

func (up *UserPreferences) String() string {
	return fmt.Sprintf("UserPreferences: { UserId: %s, Preferences: %v, CreatedAt: %s, UpdatedAt: %s }", up.UserId, up.Preferences, up.CreatedAt, up.UpdateAt)
}

type UserInputProductList struct {
	Id        uuid.UUID `json:"id" db:"id"`
	UserId    uuid.UUID `json:"userId" db:"user_id"`
	Items     string    `json:"items" db:"items"` // contains a json string
	CreatedAt string    `json:"createdAt" db:"created_at"`
	UpdatedAt string    `json:"updatedAt" db:"updated_at"`
}

func (upl *UserInputProductList) String() string {
	return fmt.Sprintf("UserProductList: { UserId: %s, ProductsList: %s, CreatedAt: %s, UpdatedAt: %s }", upl.UserId, upl.Items, upl.CreatedAt, upl.UpdatedAt)
}

type UserInputListPayload struct {
	List string `json:"list"`
}

func (d *UserInputListPayload) String() string {
	return fmt.Sprintf("UserInputListPayload: { List: %v }", d.List)
}
