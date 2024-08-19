package ab

import (
	"fmt"
	"testing"
)

func TestMagicStrings(t *testing.T) {
	// Example usage
	SetRootPath("new/root/path")
	fmt.Println(RootPath)

	// Alternatively, using default system property
	SetRootPathFromSystem()
	fmt.Println(RootPath)
}
