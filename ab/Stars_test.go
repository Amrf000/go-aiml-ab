package ab

import (
	"fmt"
	"testing"
)

func TestStars(t *testing.T) {
	stars := NewStars()
	stars.Add("Star 1")
	stars.Add("Star 2")
	stars.Add("Star 3")

	fmt.Println("Star at index 1:", stars.Star(1))
	fmt.Println("Star at index 5:", stars.Star(5))
}
