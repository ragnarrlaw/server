package search

type ContextKey string

const (
	SearchKey ContextKey = "_search"
)

type Paginate struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Meta struct {
	Total uint `json:"total"`
}

type Payload struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}
