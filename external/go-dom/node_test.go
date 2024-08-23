package dom

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestNodeList(t *testing.T) {
	input := `<root><el1/><el2/><el3/><el4/></root>`
	dec := xml.NewDecoder(strings.NewReader(input))
	doc, err := Parse(dec)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	elements := make([]Node, 4)
	elements[0] = doc.GetDocumentElement().GetFirstChild()
	elements[1] = elements[0].GetNextSibling()
	elements[2] = elements[1].GetNextSibling()
	elements[3] = elements[2].GetNextSibling()

	list := doc.GetDocumentElement().GetChildNodes()
	if list.GetLength() != 4 {
		t.Errorf("Wrong length: %d", list.GetLength())
	}
	for i := 0; i < len(elements); i++ {
		if list.Item(i) != elements[i] {
			t.Errorf("Wrong item %d", i)
		}
	}
	if list.Item(len(elements)) != nil {
		t.Errorf("Wrong limit")
	}
	doc.GetDocumentElement().RemoveChild(elements[1])
	elements = append(elements[:1], elements[2:]...)
	if list.GetLength() != 3 {
		t.Errorf("Wrong length: %d", list.GetLength())
	}
	for i := 0; i < len(elements); i++ {
		if list.Item(i) != elements[i] {
			t.Errorf("Wrong item %d", i)
		}
	}

}
