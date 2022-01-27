package main

import (
	"testing"
)

func TestIsEven(t *testing.T) {
	test_table := []struct {
		number   int
		expected bool
	}{
		{
			number:   0,
			expected: true,
		},
		{
			number:   -10000,
			expected: true,
		},
		{
			number:   10000,
			expected: true,
		},
		{
			number:   -11,
			expected: false,
		},
		{
			number:   9,
			expected: false,
		},
		{
			number:   1011,
			expected: false,
		},
		{
			number:   22,
			expected: true,
		},
		{
			number:   -10001,
			expected: false,
		},
	}

	for i, test_case := range test_table {
		result := isEven(test_case.number)

		t.Logf("Test %v. isEven(%v), result: %v", i+1, test_case.number, result)

		if result != test_case.expected {
			t.Errorf("Incorrect result. Expected %v, but got %v.", test_case.expected, result)
		}
	}
}
