package errs

import (
	"fmt"
)

type ErrorType string

const (
	ErrorTypeRequest       ErrorType = "request"
	ErrorTypeAuthorization ErrorType = "authorization"
	ErrorTypeNotFound      ErrorType = "not_found"
)

type Error struct {
	// The classification of errs encountered.
	Type ErrorType

	// StatusCode is the HTTP status code from the response.
	StatusCode int

	// ErrorMessage is a list of all the errs codes.
	ErrorMessage string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.ErrorMessage, e.StatusCode)
}
