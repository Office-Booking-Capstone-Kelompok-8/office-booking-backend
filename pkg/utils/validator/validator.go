package validator

import (
	"mime/multipart"
	"reflect"
	"strings"

	validation "github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateStruct(s interface{}) *ErrorsResponse
}

type CustomValidator struct {
	Validate *validation.Validate
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

	validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if file, ok := field.Interface().(multipart.File); ok {
			return file
		}
		return nil
	}, multipart.FileHeader{})

	validate.RegisterValidation("multipartImage", isValidMultipartImage)

	return &CustomValidator{Validate: validate}
}

func (c *CustomValidator) ValidateStruct(s interface{}) *ErrorsResponse {
	err := c.Validate.Struct(s)
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
		case "multipartImage":
			message = "must be a valid image"
		default:
			message = "invalid"
		}
		errors = append(errors, ErrorResponse{
			Field:  err.Field(),
			Reason: message,
		})
	}

	return &errors
}

func isValidMultipartImage(fl validation.FieldLevel) bool {
	fileHeader, ok := fl.Field().Interface().(*multipart.FileHeader)
	if !ok {
		return false
	}
	types := fileHeader.Header.Get("Content-Type")
	isImage := strings.HasPrefix(types, "image/")
	return isImage
}
