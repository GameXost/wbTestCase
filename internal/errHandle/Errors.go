package errHandle

import "errors"

var (
	ErrNotFound   = errors.New("order not found")
	ErrValidation = errors.New("validation error")
)
