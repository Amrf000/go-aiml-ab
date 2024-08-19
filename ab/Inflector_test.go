package ab

import (
	"fmt"
	"testing"
)

func TestPluralize(t *testing.T) {
	inflector := Inflector{}
	pairs := [][]string{
		{"dog", "dogs"},
		{"person", "people"},
		{"cats", "cats"},
	}

	for _, pair := range pairs {
		singular := pair[0]
		expected := pair[1]
		actual := inflector.Pluralize(singular)
		if actual != expected {
			t.Errorf("Pluralize %s: expected %s, got %s", pair[0], expected, actual)
		}
	}
}

func TestSingularize(t *testing.T) {
	// Add your test cases for singularize here
}

func TestInflecor(t *testing.T) {
	// Run tests
	resultPluralize := TestPluralize
	resultSingularize := TestSingularize

	fmt.Println(resultPluralize, resultSingularize)
}
