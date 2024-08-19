package ab

import (
	"fmt"
	"testing"
)

func TestHistory(t *testing.T) {
	h := NewHistoryWithName("Test")
	h.Add("Item 1")
	h.Add("Item 2")
	h.Add("Item 3")

	fmt.Println("Printing History:")
	h.PrintHistory()
}
