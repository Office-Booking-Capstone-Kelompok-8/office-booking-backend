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

	//	ErrInvalidDateRange is returned when the date range is invalid
	ErrInvalidDateRange = errors.New("invalid date range")

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

	// ErrPaymentNotFound is returned when the payment is not found
	ErrPaymentNotFound = errors.New("payment not found")

	// ErrFacilityNotFound is returned when the picture is invalid
	ErrFacilityNotFound = errors.New("facility not found")

	// ErrInvalidCategoryID is returned when the category for facility is not found
	ErrInvalidCategoryID = errors.New("facility category id is invalid")

	// ErrInavalidCityID is returned when the city is not found
	ErrInavalidCityID = errors.New("city id is invalid")

	// ErrInvalidDistrictID is returned when the district is not found
	ErrInvalidDistrictID = errors.New("district id is invalid")

	// ErrNotPublishWorthy is returned when the building is not publish worthy (e.g. no picture)
	ErrNotPublishWorthy = errors.New("building is not valid, please check the required fields")

	// ErrBuildingHasReservation is returned when the building has an active reservation
	ErrBuildingHasReservation = errors.New("building has active reservation")

	// ErrInvalidFacilityID is returned when the facility is invalid
	ErrInvalidFacilityID = errors.New("facility id is invalid")

	// ErrBuildingNotAvailable is returned when the building is not available for the given time
	ErrBuildingNotAvailable = errors.New("building is not available")

	// ErrDistrictNotInCity is returned when the district is not in the city
	ErrDistrictNotInCity = errors.New("district is not in the city")

	// ErrReservationNotFound is returned when the reservation is not found
	ErrReservationNotFound = errors.New("reservation not found")

	// ErrReservationActive is returned when the reservation is still active
	ErrReservationActive = errors.New("reservation is still active")

	// ErrStartDateBeforeToday is returned when the start date is before today
	ErrStartDateBeforeToday = errors.New("start date must be after today")

	// ErrInvalidStatus is returned when the status is invalid
	ErrInvalidStatus = errors.New("status id is invalid")

	// ErrInvalidPaymentID is returned when the payment id is invalid
	ErrInvalidPaymentID = errors.New("payment id is invalid")

	// ErrInvalidRole is returned when the role is invalid
	ErrInvalidBuildingID = errors.New("building id is invalid")

	// ErrInvalidUserID is returned when the user id is invalid
	ErrInvalidUserID = errors.New("user id is invalid")

	// ErrInvalidBankID is returned when the bank id is invalid
	ErrInvalidBankID = errors.New("bank id is invalid")

	// ErrEmailNotVerified is returned when the email is not verified when accessing email sensitive page
	ErrEmailNotVerified = errors.New("email not verified")

	// ErrEmailAlreadyVerified is returned when the user email is already verified
	ErrEmailAlreadyVerified = errors.New("email already verified")

	// ErrInvalidPaymentMethodID is returned when the payment method id is invalid
	ErrInvalidPaymentMethodID = errors.New("payment method id is invalid")

	// ErrPaymentMethodNotFound is returned when the payment method is not found
	ErrPaymentMethodNotFound = errors.New("payment method not found")

	// ErrReservationAlreadyPaid is returned when the reservation is already paid
	ErrReservationAlreadyPaid = errors.New("reservation already paid")

	// ErrReservationNotAwaitingPayment is returned when the reservation is not awaiting payment
	ErrReservationNotAwaitingPayment = errors.New("reservation is not awaiting payment")

	// ErrInvalidReviewID is returned when the review is not found
	ErrReviewNotFound = errors.New("review not found")

	// ErrReviewAlreadyExist is returned when the review is already exist
	ErrReviewAlreadyExist = errors.New("review already exist")

	// ErrReservationNotCompleted is returned when the reservation is not completed yet but the user is trying to review it
	ErrReservationNotCompleted = errors.New("reservation is not completed")

	// ErrReviewNotEditable is returned when the review is not editable (e.g. the review is already passed max edit time)
	ErrReviewNotEditable = errors.New("review is not editable")

	// ErrPaymentAlreadyExpired is returned when the payment is already expired
	ErrPaymentAlreadyExpired = errors.New("reservation payment has been expired")
)
