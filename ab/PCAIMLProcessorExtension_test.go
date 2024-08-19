package ab

import (
	"fmt"
	"testing"
)

// Example usage
func TestPCAIMLProcessorExtension(t *testing.T) {
	extension := NewPCAIMLProcessorExtension()
	fmt.Println("Extension tag names:", extension.ExtensionTagSet())
}
