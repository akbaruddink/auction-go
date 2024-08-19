package main

import "testing"

func TestDollarsToCents(t *testing.T) {
	tests := []struct {
		dollar float32
		cents  int
	}{
		{dollar: 0.0, cents: 0},
		{dollar: 1.0, cents: 100},
		{dollar: 2.5, cents: 250},
		{dollar: 10.99, cents: 1099},
		{dollar: 100.0, cents: 10000},
	}

	for _, test := range tests {
		result := dollarsToCents(test.dollar)
		if result != test.cents {
			t.Errorf("Expected dollarsToCents(%f) to return %d, but got %d", test.dollar, test.cents, result)
		}
	}
}
