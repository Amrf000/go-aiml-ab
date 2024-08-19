package ab

import (
	"fmt"
	"testing"
)

func TestsNodemapper(t *testing.T) {
	// Example usage
	node := &Nodemapper{
		Category:     nil,
		Height:       10, // Assume MagicNumbers.max_graph_height is defined elsewhere
		StarBindings: nil,
		Map:          nil,
		Key:          "",
		Value:        nil,
		ShortCut:     false,
		Sets:         []string{}, // Initialize an empty string slice
	}

	// Modify node fields as needed
	node.Key = "someKey"
	node.Value = &Nodemapper{} // Example of setting value to another Nodemapper

	// Print example fields
	fmt.Println("Node Key:", node.Key)
	fmt.Println("Node Height:", node.Height)
}
