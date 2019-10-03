// +build integration_tests

package main_test

import (
	"testing"

	m "github.com/fbrubbo/go-basics/13_testing"
)

type sumIntegration struct {
	x, y     int
	expected int
}

func TestSumIntegration(t *testing.T) {
	tests := []sumIntegration{
		sumIntegration{1, 1, 2},
		sumIntegration{1, 4, 5},
	}

	for _, test := range tests {
		if result := m.Sum(test.x, test.y); result != test.expected {
			t.Fatalf("Test failed! %d + %d = %d but expected %d", test.x, test.y, result, test.expected)
		}
	}
}
