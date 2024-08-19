package ab

import (
	"fmt"
	"testing"
)

// Example usage
func TestPredicates(t *testing.T) {
	predicates := NewPredicates()
	predicates.GetPredicateDefaults("predicates.txt")

	// Example of getting a value
	topicValue := predicates.Get("topic")
	fmt.Println("Topic:", topicValue)

	// Example of putting a new value
	predicates.Put("newKey", "newValue")

	// Example of getting the new value
	newValue := predicates.Get("newKey")
	fmt.Println("New Value:", newValue)
}
