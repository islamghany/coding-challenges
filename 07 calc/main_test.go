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
		output int
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
	}
	calc := calculator.NewCalculator()
	for idx, tc := range testCases {
		r, _ := calc.Evaluate(tc.input)
		if r != tc.output {
			t.Errorf("%d Expected %d but got %d", idx, tc.output, r)
		}
	}

}
