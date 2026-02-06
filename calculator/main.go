package main

import (
	"calc/calculator"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Please provide an expression to evaluate")
	}
	exp := args[0]
	fmt.Println("Expression:", exp)

	calc := calculator.NewCalculator()
	r, err := calc.Evaluate(exp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r)
}
