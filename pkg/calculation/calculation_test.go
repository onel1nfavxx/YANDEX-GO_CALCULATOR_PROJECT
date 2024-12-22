package calculation

import (
	"testing"

	"github.com/onel1nfavxx/YANDEX-GO_CALCULATOR_PROJECT/custom_errors"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		expectErr  error
	}{
		{"3 + 5", 8, nil},
		{"10 - 2", 8, nil},
		{"4 * 2", 8, nil},
		{"16 / 2", 8, nil},
		{"3 + 5 * 2", 13, nil},
		{"(1 + 2) * 3", 9, nil},
		{"(1 + 2) / 3", 1, nil},
		{"2 * (3 + 5)", 16, nil},
		{"2 * (3 + 5) - 4", 12, nil},
		{"2 * (3 + 5) / 4", 4, nil},
		{"2 * (3 + 5) / 0", 0, custom_errors.ErrDivisionByZero},
		{"3 +", 0, custom_errors.ErrInvalidExpression},
		{"* 3 + 5", 0, custom_errors.ErrInvalidExpression},
		{"3 + a", 0, custom_errors.ErrInvalidExpression},
		{"", 0, custom_errors.ErrInvalidExpression},
	}

	for _, test := range tests {
		result, err := Calc(test.expression)
		if result != test.expected || (err != nil && err.Error() != test.expectErr.Error()) {
			t.Errorf("Calc(%q) = %v, %v; want %v, %v", test.expression, result, err, test.expected, test.expectErr)
		}
	}
}
