package types

import (
	"errors"
)

/**
client side errors
*/

var (
	ErrNotFound               = errors.New("not found")
	ErrBadRequest             = errors.New("bad request")  // invalid data received from the client, such as invalid parameter formats or
	ErrUnauthorized           = errors.New("unauthorized") // lack of identification data (no token or expired token)
	ErrForbidden              = errors.New("forbidden")    // permission issues and restrictions imposed on clients
	ErrConflict               = errors.New("conflict")     // scenarios where conflict have occurred (trying to create something that already exits, or concurrent updates)
	ErrUnprocessableEntity    = errors.New("unprocessable entity")
	ErrTooManyRequests        = errors.New("too many requests")
	ErrUnsupportedContentType = errors.New("unsupported media type")
)

/**
server side errors
*/

var (
	ErrServiceUnavailable  = errors.New("service unavailable")
	ErrTimeout             = errors.New("gateway timeout")
	ErrInternalServerError = errors.New("internal server error")
)
