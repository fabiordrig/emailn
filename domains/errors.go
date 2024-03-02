package domains

import "errors"

var (
	ErrAtLeastOneEmailIsRequired = errors.New("at least one email is required")
	ErrNameAndContentAreRequired = errors.New("name and content are required")
	ErrInvalidEmail              = errors.New("invalid email")
)
