package domain

import "errors"

var (
	ErrInvalidJobFields error = errors.New("invalid job fields")
	ErrJobNotFound		error = errors.New("job not found")
)
