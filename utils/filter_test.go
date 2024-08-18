package utils

import (
	"math/rand"
	"testing"
)

func TestFilter(t *testing.T) {
	var numbers []int
	for i := 0; i < 100; i++ {
		numbers = append(numbers, rand.Intn(1000))
	}

	evenNumbers := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})

	for _, n := range evenNumbers {
		if n%2 != 0 {
			t.Errorf("Expected %d to be even", n)
		}
	}
}
