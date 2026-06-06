package handlers

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/kodokbakar/go-ecommerce-api/internal/response"
)

type validationErrorDetail struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func bindAndValidateJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		response.BadRequest(c, "invalid request body", buildValidationErrorDetails(dst, err))
		return false
	}

	return true
}

func buildValidationErrorDetails(dst any, err error) any {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err.Error()
	}

	details := make([]validationErrorDetail, 0, len(validationErrors))

	for _, fieldError := range validationErrors {
		fieldName := jsonFieldName(dst, fieldError.StructField())

		details = append(details, validationErrorDetail{
			Field:   fieldName,
			Rule:    fieldError.Tag(),
			Message: validationMessage(fieldName, fieldError),
		})
	}

	return details
}

func jsonFieldName(dst any, structFieldName string) string {
	t := reflect.TypeOf(dst)
	if t == nil {
		return strings.ToLower(structFieldName)
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return strings.ToLower(structFieldName)
	}

	field, ok := t.FieldByName(structFieldName)
	if !ok {
		return strings.ToLower(structFieldName)
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return strings.ToLower(structFieldName)
	}

	jsonName := strings.Split(jsonTag, ",")[0]
	if jsonName == "" || jsonName == "-" {
		return strings.ToLower(structFieldName)
	}

	return jsonName
}

func validationMessage(field string, fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, fieldError.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, fieldError.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fieldError.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fieldError.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fieldError.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
