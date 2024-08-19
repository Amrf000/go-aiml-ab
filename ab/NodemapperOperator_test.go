package ab

import (
	"fmt"
	"testing"
)

func TestNodemapperOperator(t *testing.T) {
	// Example usage
	node := &Nodemapper{} // Assuming Nodemapper is initialized somehow
	Put(node, "key1", &Nodemapper{})
	fmt.Println(Size(node))
}
