package custom_errors

import "errors"

var (
	ErrEmptyString          = errors.New("input string is empty")
	ErrEndWithOperator      = errors.New("expression can't end with operation")
	ErrStartWithOperator    = errors.New("expression can't start with operation")
	ErrDivisionByZero       = errors.New("division by zero is not allowed")
	ErrTwoConsecOperation   = errors.New("expression can't hold two consecutive operation")
	ErrNoClosingParenthesis = errors.New("missing )")
	ErrInvInputs            = errors.New("invalid input")
	ErrInvalidExpression    = errors.New("invalid expression")
)
