package ab

import (
	"fmt"
	"testing"
)

func TestPath(t *testing.T) {
	sentence := "This is a sample sentence."
	path := SentenceToPath(sentence)
	fmt.Println("Path to sentence:", PathToSentence(path))
	fmt.Println("Printing path:")
	path.Print()
}
