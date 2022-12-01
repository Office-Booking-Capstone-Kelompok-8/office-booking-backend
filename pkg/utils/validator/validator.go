package validator

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	"mime/multipart"
	"net/textproto"
	"reflect"
	"strings"

	validation "github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateStruct(s interface{}) *ErrorsResponse
	ValidateVar(field interface{}, tag string) *ErrorsResponse
}

type CustomValidator struct {
	Validate *validation.Validate
	trans    ut.Translator
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
		// register for multipart.FileHeader, so we can validate the file type from the header instead of the file itself
		// this is because the file is not available in the request body, only the header,
		//  furthermore, we don't want to return multipart.FileHeader again to the validator because it can trigger forever loop
		if value, ok := field.Interface().(multipart.FileHeader); ok {
			return value.Header
		}
		return nil
	}, multipart.FileHeader{})

	validate.RegisterValidation("multipartImage", isValidMultipartImage)

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_trans.RegisterDefaultTranslations(validate, trans)

	return &CustomValidator{Validate: validate, trans: trans}
}

func (c *CustomValidator) ValidateStruct(s interface{}) *ErrorsResponse {
	err := c.Validate.Struct(s)
	if err == nil {
		return nil
	}
	fmt.Println(s)
	var errors ErrorsResponse
	for _, err := range err.(validation.ValidationErrors) {
		errors = append(errors, ErrorResponse{
			Field:  err.Field(),
			Reason: err.Translate(c.trans),
		})
	}

	return &errors
}

func (c *CustomValidator) ValidateVar(field interface{}, tag string) *ErrorsResponse {
	err := c.Validate.Var(field, tag)
	if err == nil {
		return nil
	}
	var errors ErrorsResponse
	for _, err := range err.(validation.ValidationErrors) {
		errors = append(errors, ErrorResponse{
			Field:  err.Field(),
			Reason: err.Translate(c.trans),
		})
	}

	return &errors
}

func isValidMultipartImage(fl validation.FieldLevel) bool {
	header, ok := fl.Field().Interface().(textproto.MIMEHeader)
	if !ok {
		return false
	}
	types := header.Get("Content-Type")
	isImage := strings.HasPrefix(types, "image/")
	return isImage
}
