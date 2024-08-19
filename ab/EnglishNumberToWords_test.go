package ab

import (
	"fmt"
	"testing"
)

func TestEnglishNumberToWords(t *testing.T) {
	fmt.Println("***", convert(0))
	fmt.Println("***", convert(1))
	fmt.Println("***", convert(16))
	fmt.Println("***", convert(100))
	fmt.Println("***", convert(118))
	fmt.Println("***", convert(200))
	fmt.Println("***", convert(219))
	fmt.Println("***", convert(800))
	fmt.Println("***", convert(801))
	fmt.Println("***", convert(1316))
	fmt.Println("***", convert(1000000))
	fmt.Println("***", convert(2000000))
	fmt.Println("***", convert(3000200))
	fmt.Println("***", convert(700000))
	fmt.Println("***", convert(9000000))
	fmt.Println("***", convert(9001000))
	fmt.Println("***", convert(123456789))
	fmt.Println("***", convert(2147483647))
	fmt.Println("***", convert(3000000010))
}
