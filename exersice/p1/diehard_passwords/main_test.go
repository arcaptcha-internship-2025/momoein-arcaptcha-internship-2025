package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiehardPasswords(t *testing.T) {
	tests := []struct {
		n int
		s []int
	}{
		{n: 0, s: []int{}},
		{n: 1, s: []int{2, 3, 5, 7}},
		{n: 3, s: []int{233, 239, 293, 311, 313, 317, 373, 379, 593, 599, 719, 733, 739, 797}},
	}

	for _, test := range tests {
		s := DiehardPasswords(test.n)
		assert.Equal(t, fmt.Sprint(test.s), fmt.Sprint(s))
	}
}
