package config

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// FormatValidationError returns properly formatted validation error, which contains description of all validation
// errors for fields (or for field tags value, if fieldTagKey is specified)
func FormatValidationError(originalStructType reflect.Type, originalError error, fieldTagKey string) error {
	if validationErrors, ok := originalError.(validator.ValidationErrors); ok {
		var errorMessage string

		for _, validationError := range validationErrors {
			field, _ := originalStructType.FieldByName(validationError.Field())
			fieldName := field.Tag.Get(fieldTagKey)
			if fieldName == "" {
				fieldName = validationError.Field()
			}

			errorMessage = errorMessage + fmt.Sprintf("validation failed for field '%s' (validation type: %s); ", fieldName, validationError.Tag())
		}

		return fmt.Errorf(errorMessage)
	}

	return fmt.Errorf("validation failed. Original error: %s", originalError)
}
