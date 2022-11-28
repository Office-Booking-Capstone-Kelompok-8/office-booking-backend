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

	// ErrPasswordNotMatch is returned when the old password field are different from the saved one
	ErrPasswordNotMatch = errors.New("password not match")

	// ErrInvalidToken is returned when the token is invalid or expired
	ErrInvalidToken = errors.New("invalid or expired JWT")

	// ErrInvalidOTP is returned when the OTP is invalid
	ErrInvalidOTP = errors.New("invalid OTP")

	//	ErrInvalidOTPToken is returned when the OTP token is invalid
	ErrInvalidOTPToken = errors.New("invalid OTP token")

	// ErrNoPermission is returned when the user doesn't have permission to perform the action (e.g. access admin page as a user)
	ErrNoPermission = errors.New("you don't have permission to perform the action")

	//	ErrInvalidQueryParams is returned when the query params is invalid
	ErrInvalidQueryParams = errors.New("invalid query params")

	//	ErrStartDateAfterEndDate is returned when the start date is after the end date
	ErrStartDateAfterEndDate = errors.New("start date must be before end date")

	//	ErrBuildingNotFound is returned when the building is not found
	ErrBuildingNotFound = errors.New("building not found")
)
