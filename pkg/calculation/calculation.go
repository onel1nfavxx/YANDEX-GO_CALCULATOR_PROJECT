package calculation

import (
	"strconv"
	"strings"

	"github.com/onel1nfavxx/YANDEX-GO_CALCULATOR_PROJECT/custom_errors"
)

type Operator struct {
	symbol     string
	precedence int
	assoc      string
}

var operators = map[string]Operator{
	"+": {"*", 1, "L"},
	"-": {"*", 1, "L"},
	"*": {"*", 2, "L"},
	"/": {"*", 2, "L"},
}

func precedence(op string) int {
	if operator, exists := operators[op]; exists {
		return operator.precedence
	}
	return -1
}

func Calc(expression string) (float64, error) {
	output := []string{}
	stack := []string{}
	numStack := []float64{}

	tokens := strings.Split(expression, "")

	for _, token := range tokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			numStack = append(numStack, num)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
		} else {
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(token) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}

	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	for _, token := range output {
		switch token {
		case "+":
			if len(numStack) < 2 {
				return 0, custom_errors.ErrInvalidExpression
			}
			b := numStack[len(numStack)-1]
			numStack = numStack[:len(numStack)-1]
			a := numStack[len(numStack)-1]
			numStack = numStack[:len(numStack)-1]
			numStack = append(numStack, a+b)
		case "-":
			if len(numStack) < 2 {
				return 0, custom_errors.ErrInvalidExpression
			}
			b := numStack[len(numStack)-1]
			numStack = numStack[:len(numStack)-1]
			a := numStack[len(numStack)-1]
			numStack = numStack[:len(numStack)-1]
			numStack = append(numStack, a-b)
		case "*":
			if len(numStack) < 2 {
				return 0, custom_errors.ErrInvalidExpression
			}
			b := numStack[len(numStack)-1]
			numStack = numStack[:len(numStack)-1]
			a := numStack[len(numStack)-1]
			numStack = numStack[:len(numStack)-1]
			numStack = append(numStack, a*b)
		case "/":
			if len(numStack) < 2 {
				return 0, custom_errors.ErrInvalidExpression
			}
			b := numStack[len(numStack)-1]
			numStack = numStack[:len(numStack)-1]
			a := numStack[len(numStack)-1]
			numStack = numStack[:len(numStack)-1]
			if b == 0 {
				return 0, custom_errors.ErrDivisionByZero
			}
			numStack = append(numStack, a/b)
		default:
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, custom_errors.ErrInvalidExpression
			}
			numStack = append(numStack, num)
		}
	}

	if len(numStack) != 1 {
		return 0, custom_errors.ErrInvalidExpression
	}

	return numStack[0], nil
}
