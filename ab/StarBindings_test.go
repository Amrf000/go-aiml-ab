package ab

import (
	"fmt"
	"testing"
)

func TestStarBindings(t *testing.T) {
	starBindings := NewStarBindings()
	starBindings.InputStars.Add("Input Star 1")
	starBindings.InputStars.Add("Input Star 2")

	starBindings.ThatStars.Add("That Star 1")

	starBindings.TopicStars.Add("Topic Star 1")
	starBindings.TopicStars.Add("Topic Star 2")
	starBindings.TopicStars.Add("Topic Star 3")

	fmt.Println("Input Star at index 1:", starBindings.InputStars.Star(1))
	fmt.Println("That Star at index 0:", starBindings.ThatStars.Star(0))
	fmt.Println("Topic Star at index 2:", starBindings.TopicStars.Star(2))
}
