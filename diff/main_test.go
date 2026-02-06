package main

import (
	"testing"
)

func TestLCSBasics(t *testing.T) {
	testCases := []struct {
		s1, s2, expected string
	}{
		{"ABC", "AC", "AC"},
		{"ABCDEF", "ABCDEF", "ABCDEF"},
		{"ABC", "DEF", ""},
		{"AABCXY", "XYZ", "XY"},
		{"", "", ""},
		{"ABCD", "AC", "AC"},
		{"abdcq", "cbdq", "bdq"},
	}

	for _, tc := range testCases {
		actual := LCS(tc.s1, tc.s2)
		if actual != tc.expected {
			t.Errorf("LCS(%s, %s): expected %s, actual %s", tc.s1, tc.s2, tc.expected, actual)
		}
	}

}

func TestLCSMultipleLine(t *testing.T) {
	testCases := []struct {
		arr1, arr2, expected []string
	}{
		{
			[]string{
				"This is a test which contains:",
				"this is the lcs",
			},
			[]string{
				"this is the lcs",
				"we're testing",
			},
			[]string{
				"this is the lcs",
			},
		},
		{
			[]string{
				"Coding Challenges helps you become a better software engineer through that build real applications.",
				"I share a weekly coding challenge aimed at helping software engineers level up their skills through deliberate practice.",
				"I’ve used or am using these coding challenges as exercise to learn a new programming language or technology.",
				"Each challenge will have you writing a full application or tool. Most of which will be based on real world tools and utilities."},
			[]string{
				"Helping you become a better software engineer through coding challenges that build real applications.",
				"I share a weekly coding challenge aimed at helping software engineers level up their skills through deliberate practice.",
				"These are challenges that I’ve used or am using as exercises to learn a new programming language or technology.",
				"Each challenge will have you writing a full application or tool. Most of which will be based on real world tools and utilities.",
			},
			[]string{
				"I share a weekly coding challenge aimed at helping software engineers level up their skills through deliberate practice.",
				"Each challenge will have you writing a full application or tool. Most of which will be based on real world tools and utilities.",
			},
		},
	}

	for _, tc := range testCases {
		actual := LCSArrays(tc.arr1, tc.arr2)
		if len(actual) != len(tc.expected) {
			t.Errorf("LCSMultipleLine(%v, %v): expected %v, actual %v", tc.arr1, tc.arr2, tc.expected, actual)
		}
		for i := 0; i < len(actual); i++ {
			if actual[i] != tc.expected[i] {
				t.Errorf("LCSMultipleLine(%v, %v): expected %v, actual %v", tc.arr1, tc.arr2, tc.expected, actual)
			}
		}
	}
}
