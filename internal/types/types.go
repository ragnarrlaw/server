package types

/**
  Most commonly used types goes here
*/

// used for authentication service - for generics
type Identifiable interface {
	GetId() string
	GetUsername() string
}

// pagination types
type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
