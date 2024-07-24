package calculator

import (
	"fmt"
	"math"
	"strconv"
)

type Calculator struct {
	exp string
}

func NewCalculator() *Calculator {
	return &Calculator{}
}

func (calc *Calculator) Evaluate(exp string) (int, error) {
	calc.exp = exp
	rpn := calc.generatePostfixNotation()
	return calc.evaluatePostfixNotation(rpn)
}

func (calc *Calculator) generatePostfixNotation() []string {
	output := NewStack[string](10)
	operators := NewStack[string](10)
	isPrevNum := false
	for _, c := range calc.exp {
		isNum := false
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			output.Push(string(c))
			isNum = true
		case '+', '-', '*', '/', '%', '^':
			calc.handleOperator(string(c), output, operators)
			break
		case '(':
			operators.Push(string(c))
			break
		case ')':
			calc.handleClosingBracket(output, operators)
			break
		case ' ':
			break
		default:
			return nil
		}
		if isNum && isPrevNum {
			o1 := output.Pop()
			o2 := output.Pop()
			output.Push(o2 + o1)
		}
		isPrevNum = isNum
	}
	for !operators.IsEmpty() {
		output.Push(operators.Pop())
	}
	return output.ToArray()
}

func (calc *Calculator) handleOperator(op string, output, operators *Stack[string]) {
	for !operators.IsEmpty() {
		op2 := operators.Peek()
		if calc.precedence(op) <= calc.precedence(op2) {
			output.Push(operators.Pop())
		} else {
			break
		}
	}
	operators.Push(op)
}

func (calc *Calculator) precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/", "%":
		return 2
	case "^":
		return 3
	}
	return 0
}

func (calc *Calculator) handleClosingBracket(output, operators *Stack[string]) {
	for !operators.IsEmpty() {
		op := operators.Pop()
		if op == "(" {
			break
		}
		output.Push(op)
	}
}

func (calc *Calculator) evaluatePostfixNotation(rpn []string) (int, error) {
	fmt.Println(rpn)
	stack := NewStack[int](10)
	for _, tok := range rpn {
		switch tok {
		case "+", "-", "*", "/", "%", "^":
			if stack.length < 2 {
				return 0, fmt.Errorf("Invalid expression")
			}
			op2 := stack.Pop()
			op1 := stack.Pop()
			result := calc.evaluate(op1, op2, tok)
			stack.Push(result)
		default:
			val, err := strconv.Atoi(tok)
			if err != nil {
				return 0, err
			}
			stack.Push(val)

		}
	}
	if stack.length != 1 {
		return 0, fmt.Errorf("Invalid expression")
	}
	return stack.Pop(), nil
}

func (calc *Calculator) evaluate(op1, op2 int, op string) int {
	switch op {
	case "+":
		return op1 + op2
	case "-":
		return op1 - op2
	case "*":
		return op1 * op2
	case "/":
		return op1 / op2
	case "%":
		return op1 % op2
	case "^":
		return int(math.Pow(float64(op1), float64(op2)))
	}
	return 0
}
