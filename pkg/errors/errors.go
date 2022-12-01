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

	// ErrUserHasReservation is returned when the user has reservation
	ErrUserHasReservation = errors.New("user has active reservation")

	//	ErrPictureServiceFailed is returned when the picture service failed
	ErrPictureServiceFailed = errors.New("picture service failed")

	// ErrPicureLimitExceeded is returned when the picture limit is exceeded
	ErrPicureLimitExceeded = errors.New("picture limit exceeded")

	// ErrPictureNotFound is returned when the picture is not found
	ErrPictureNotFound = errors.New("picture not found")

	// ErrFacilityNotFound is returned when the picture is invalid
	ErrFacilityNotFound = errors.New("facility not found")

	// ErrInvalidCategoryID is returned when the category for facility is not found
	ErrInvalidCategoryID = errors.New("facility category id is invalid")

	// ErrInavalidCityID is returned when the city is not found
	ErrInavalidCityID = errors.New("city id is invalid")

	// ErrInvalidDistrictID is returned when the district is not found
	ErrInvalidDistrictID = errors.New("district id is invalid")

	// ErrNotPublishWorthy is returned when the building is not publish worthy (e.g. no picture)
	ErrNotPublishWorthy = errors.New("building is not publish worthy")
)
