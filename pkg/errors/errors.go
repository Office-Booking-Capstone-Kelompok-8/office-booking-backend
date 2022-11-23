package errors

import "errors"

var (
	// ErrInvalidRequestBody is returned when request body is invalid because it is  syntactically incorrect or doesn't conforms to the schema (e.g. missing required fields)
	ErrInvalidRequestBody = errors.New("invalid request body")

	// ErrDuplicateEmail is returned when the email is already registered
	ErrDuplicateEmail = errors.New("email already used")

	// ErrUserNotFound is returned when the user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidCredentials is returned when the email or password is incorrect
	ErrInvalidCredentials = errors.New("invalid email or password")
)
