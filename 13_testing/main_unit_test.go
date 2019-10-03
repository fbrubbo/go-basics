// +build unit_tests

package main_test

import (
	"testing"

	log "github.com/sirupsen/logrus"

	m "github.com/fbrubbo/go-basics/13_testing"
)

type sum struct {
	x, y     int
	expected int
}

func TestSum(t *testing.T) {
	tests := []sum{
		sum{1, 1, 2},
		sum{1, 4, 5},
	}

	for i, test := range tests {
		log.Infof("test info %d ", i)
		log.Warn("test warn", i)
		log.Error("test error", i)
		if result := m.Sum(test.x, test.y); result != test.expected {
			t.Fatalf("Test failed! %d + %d = %d but expected %d", test.x, test.y, result, test.expected)
		}
	}
}
