package validator

import (
	validation "github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type Validator interface {
	Validate(s interface{}) *ErrorsResponse
}

type CustomValidator struct {
	validate *validation.Validate
}

type ErrorResponse struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type ErrorsResponse []ErrorResponse

func NewValidator() Validator {
	validate := validation.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{validate: validate}
}

func (c *CustomValidator) Validate(s interface{}) *ErrorsResponse {
	err := c.validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors ErrorsResponse
	for _, err := range err.(validation.ValidationErrors) {
		var message string
		switch err.Tag() {
		case "required":
			message = "required"
		case "email":
			message = "must be a valid email"
		case "min":
			message = "must be at least " + err.Param() + " characters"
		case "gte":
			message = "must be greater than or equal to " + err.Param()
		case "lte":
			message = "must be less than or equal to " + err.Param()
		case "latitude":
			message = "must be a valid latitude"
		case "longitude":
			message = "must be a valid longitude"
		case "uuid":
			message = "must be a valid uuid"
		}
		errors = append(errors, ErrorResponse{
			Field:  err.Field(),
			Reason: message,
		})
	}

	return &errors
}
