package validator

import (
	"mime/multipart"
	"net/textproto"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_trans "github.com/go-playground/validator/v10/translations/en"

	validation "github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateJSON(s interface{}) *ErrorsResponse
	ValidateQuery(s interface{}) *ErrorsResponse
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

func (e *ErrorsResponse) AddError(field string, reason string) {
	*e = append(*e, ErrorResponse{
		Field:  field,
		Reason: reason,
	})
}

func NewValidator() Validator {
	validate := validation.New()

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

	customValidator := &CustomValidator{
		Validate: validate,
		trans:    trans,
	}

	customValidator.addTranslation("required_with", "{0} is required when {1} is present")
	customValidator.addTranslation("required_if", "{0} is required when {1}")

	return customValidator
}

func (c *CustomValidator) ValidateStruct(s interface{}) *ErrorsResponse {
	err := c.Validate.Struct(s)
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

func (c *CustomValidator) ValidateJSON(s interface{}) *ErrorsResponse {
	c.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return c.ValidateStruct(s)
}

func (c *CustomValidator) ValidateQuery(s interface{}) *ErrorsResponse {
	c.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return c.ValidateStruct(s)
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

func (c *CustomValidator) addTranslation(tag string, message string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, message, true)
	}

	transFn := func(ut ut.Translator, fe validation.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = c.Validate.RegisterTranslation(tag, c.trans, registerFn, transFn)
}
