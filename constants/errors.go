package constants

import "errors"

var (
	ErrAtLeastOneEmailIsRequired = errors.New("at least one email is required")
	ErrRequiredField             = errors.New("missing required field")
	ErrInvalidEmail              = errors.New("invalid email")
	ErrStringMinLength           = errors.New("the minimum length is not met")
	ErrStringMaxLength           = errors.New("the maximum length is exceeded")
	ErrUnknown                   = errors.New("unknown error")
	ErrInternalServer            = errors.New("internal server error")
	ErrUnprocessableEntity       = errors.New("unprocessable entity")
	ErrNotFound                  = errors.New("the resource was not found")
)
