package ab

import (
	"fmt"
	"strings"
)

type Path []string

type Paths struct {
	Path
	Word   string
	Next   *Paths
	Length int
}

func NewPath() *Paths {
	this := Paths{}
	this.Word = ""
	this.Next = nil
	this.Length = 0
	return &this
}

func SentenceToPath(sentence string) *Paths {
	sentence = strings.TrimSpace(sentence)
	sentence = strings.ReplaceAll(sentence, "><", "> <")
	words := strings.Split(sentence, " ")
	return ArrayToPath(words)
}

func PathToSentence(path *Paths) string {
	var result strings.Builder
	for p := path; p != nil; p = p.Next {
		result.WriteString(" ")
		result.WriteString(p.Word)
	}
	return strings.TrimSpace(result.String())
}

func ArrayToPath(array []string) *Paths {
	//if slices.Contains(array, "<SET>MONTH</SET> <SET>ORDINAL</SET> <SET>NUMBER</SET>") {
	//	nn := 0
	//	nn++
	//}
	var head *Paths = nil
	var tail *Paths = nil
	for i := len(array) - 1; i >= 0; i-- {
		head = NewPath()
		head.Word = array[i]
		head.Next = tail
		if tail == nil {
			head.Length = 1
		} else {
			head.Length = tail.Length + 1
		}
		tail = head
	}
	return head
}

func ArrayToPathRecursive(array []string, index int) *Paths {
	if index >= len(array) {
		return nil
	}
	newPath := NewPath()
	newPath.Word = array[index]
	newPath.Next = ArrayToPathRecursive(array, index+1)

	if newPath.Next == nil {
		newPath.Length = 1
	} else {
		newPath.Length = newPath.Next.Length + 1
	}
	return newPath
}

func (p *Paths) Print() {
	var result strings.Builder
	for ; p != nil; p = p.Next {
		result.WriteString(p.Word)
		result.WriteString(",")
	}
	if result.Len() > 0 {
		//result = result[0:(result.Len() - 1)] // Remove the trailing comma
	}
	fmt.Println(result.String())
}
