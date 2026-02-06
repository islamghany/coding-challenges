package main

import (
	"calc/calculator"
	"testing"
)

func TestStack(t *testing.T) {

	s := calculator.NewStack[int](2)
	s.Push(10)
	f := s.Pop()
	if f != 10 {
		t.Error("Expected 10 but got ", f)
	}
	s.Push(20)
	s.Push(30)
	f = s.Pop()
	if f != 30 {
		t.Error("Expected 30 but got ", f)
	}
	f = s.Pop()
	if f != 20 {
		t.Error("Expected 20 but got ", f)
	}
	f = s.Pop()
	if f != 0 {
		t.Error("Expected 0 but got ", f)
	}
	s.Push(40)
	f = s.Peek()
	if f != 40 {
		t.Error("Expected 40 but got ", f)
	}
	s.Clear()
	f = s.Pop()
	if f != 0 {
		t.Error("Expected 0 but got ", f)
	}
}

func TestCalculator(t *testing.T) {
	testCases := []struct {
		input  string
		output float64
	}{
		{"1 + 3 * (50) - 2", 149},
		{"1 + 3 * (50 - 2)", 145},
		{"1 + 3 * 50 - 2", 149},
		{"1 + 3 * 50 - 2 * 2", 147},
		{"1 + 3 * 50 - 2 * 2 + 1", 148},
		{"1 + 3 * 50 - 2 * 2 + 1 * 2", 149},
		{"1+19-19", 1},
		{"1+19-19*2", -18},
		{"1+19-19*2+1", -17},
		{"1+19-19*2+1*2", -16},
		{"1+19-19*2+1*2-1", -17},
		{"1+19-19*2+1*2-1*2**", 0},
		{"0011", 11},
		{"0011+0011", 22},
		{"-10", 0},
		// sin, cos, tan, sqrt
		{"s(90)", 0.893997},
		{"c(0)", 1},
		{"t(45)", 1.619775},
		{"q(16)", 4},
		// power
		{"2^3", 8},
		{"2^3^2", 512},
		{"2^3^2^2", 65536},
	}
	calc := calculator.NewCalculator()
	for _, tc := range testCases {
		r, _ := calc.Evaluate(tc.input)
		// float comparison

		if r-tc.output > 0.0001 {
			t.Errorf("Expected %f but got %f", tc.output, r)
		}

	}

}
