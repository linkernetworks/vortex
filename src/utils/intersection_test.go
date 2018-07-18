package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersection(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	b := []string{"b", "d", "e", "f"}
	ans := []string{"b", "d"}

	c := Intersection(a, b)
	assert.Equal(t, 2, len(c))
	assert.Equal(t, c, ans)
}

func TestIntersections(t *testing.T) {
	a := [][]string{
		{"one", "two", "three", "four"},
		{"one", "two", "five", "six"},
		{"two", "five", "six"},
		{"ten", "seven", "two", "six"},
	}
	ans := []string{"two"}

	c := Intersections(a)
	assert.Equal(t, 1, len(c))
	assert.Equal(t, c, ans)
}
