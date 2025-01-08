package custom_errors

import "errors"

var (
	ErrDivisionByZero = errors.New("division by zero is not allowed")
	ErrInvInputs      = errors.New("invalid input")
)
