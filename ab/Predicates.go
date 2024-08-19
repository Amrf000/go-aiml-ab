package ab

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Predicates map[string]string

func NewPredicates() Predicates {
	this := make(map[string]string)
	return this
}

func (p Predicates) Put(key, value string) {
	if key == "topic" && JpTokenize {
		value = TokenizeSentence(value)
	}
	if key == "topic" && len(value) == 0 {
		value = DefaultGet
	}
	if value == TooMuchRecursion {
		value = DefaultListItem
	}
	p[key] = value
}

func (p Predicates) Get(key string) string {
	value, ok := p[key]
	if !ok {
		value = DefaultGet
	}
	return value
}

func (p *Predicates) GetPredicateDefaultsFromInputStream(in *os.File) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		strLine := scanner.Text()
		if strings.Contains(strLine, ":") {
			property := strLine[:strings.Index(strLine, ":")]
			value := strLine[strings.Index(strLine, ":")+1:]
			p.Put(property, value)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}

func (p *Predicates) GetPredicateDefaults(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	p.GetPredicateDefaultsFromInputStream(file)
	return nil
}
