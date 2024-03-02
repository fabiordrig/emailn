package domains

import (
	"emailn/constants"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(obj interface{}) error {
	validate := validator.New()

	err := validate.Struct(obj)

	if err == nil {
		return nil
	}

	errors := err.(validator.ValidationErrors)

	validateError := errors[0]

	switch validateError.Tag() {
	case "required":
		return constants.ErrRequiredField
	case "email":
		return constants.ErrInvalidEmail
	case "min":
		return constants.ErrStringMinLength
	case "max":
		return constants.ErrStringMaxLength
	}

	return constants.ErrUnknown
}
