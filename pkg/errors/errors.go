package errors

import "errors"

var (
	// ErrInvalidRequestBody is returned when request body is invalid because it is  syntactically incorrect or doesn't conforms to the schema (e.g. missing required fields)
	ErrInvalidRequestBody = errors.New("invalid request body")

	// ErrDuplicateEmail is returned when the email is already registered
	ErrDuplicateEmail = errors.New("email already used")
)
