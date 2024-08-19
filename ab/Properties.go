package ab

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Properties map[string]string

func (p Properties) Get(key string) string {
	if value, ok := p[key]; ok {
		return value
	}
	return DefaultProperty
}

func (p Properties) GetPropertiesFromInputStream(in *os.File) int {
	cnt := 0
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		strLine := scanner.Text()
		if strings.Contains(strLine, ":") {
			property := strLine[:strings.Index(strLine, ":")]
			value := strLine[strings.Index(strLine, ":")+1:]
			p[property] = value
			cnt++
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
	return cnt
}

func (p Properties) GetProperties(filename string) int {
	cnt := 0
	fmt.Println("Get Properties:", filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return cnt
	}
	defer file.Close()

	cnt = p.GetPropertiesFromInputStream(file)
	return cnt
}
