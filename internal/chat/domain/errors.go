package domain

import "errors"

// Common errors
var (
	ErrRequired              = errors.New("required value")
	ErrNotFound              = errors.New("not found")
	ErrNegative              = errors.New("negative value")
	ErrLowThenOneValue       = errors.New("low then one value")
	ErrInvalidFormatEmail    = errors.New("must be email type")
	ErrInvalidFormatPassword = errors.New("invalid password format")
	ErrValidation            = errors.New("error validation")
	ErrNoUserInContext       = errors.New("no user in context")
	ErrInvalidUser           = errors.New("user data is invalid")
	ErrAccessDenied          = errors.New("access denied")
	ErrEditDenied            = errors.New("edit denied")
	ErrDbCreationFailed      = errors.New("data base creation failed")
	ErrFailure               = errors.New("failure")
)
