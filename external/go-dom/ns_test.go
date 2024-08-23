package dom

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestMultiplePrefixForSameNS(t *testing.T) {
	input := `<h:note xmlns:h="http://www.w3.org/TR/html4/"
xmlns:t="https://test.com/t">
<t:to>Tove</t:to>
<t:from>Jani</t:from>
<h:body>weekend!</h:body>
</h:note>`
	hspace := "http://www.w3.org/TR/html4/"
	dec := xml.NewDecoder(strings.NewReader(input))
	doc, err := Parse(dec)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root := doc.GetDocumentElement()
	if root.GetNodeName() != "h:note" {
		t.Errorf("Bad root name")
	}
	qn := root.GetQName()
	if qn.Prefix != "h" || qn.Local != "note" || qn.Space != hspace {
		t.Errorf("Wrong root qname: %v", qn)
	}
}
