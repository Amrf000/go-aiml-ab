package dom

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestNoNS1(t *testing.T) {
	input := `<note>
<to>Tove</to>
<from>Jani</from>
<heading>Reminder</heading>
<body>weekend!</body>
</note>`
	dec := xml.NewDecoder(strings.NewReader(input))
	doc, err := Parse(dec)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root := doc.GetDocumentElement()
	if root.GetNodeName() != "note" {
		t.Errorf("Bad root name")
	}
	el := root.GetFirstChild().GetNextSibling()
	if el.GetNodeName() != "to" && el.GetFirstChild().(Text).GetValue() != "Tove" {
		t.Errorf("Bad to")
	}
	el = el.GetNextSibling().GetNextSibling()
	if el.GetNodeName() != "from" && el.GetFirstChild().(Text).GetValue() != "Jani" {
		t.Errorf("Bad from")
	}
	el = el.GetNextSibling().GetNextSibling()
	if el.GetNodeName() != "heading" && el.GetFirstChild().(Text).GetValue() != "Reminder" {
		t.Errorf("Bad heading")
	}
	el = el.GetNextSibling().GetNextSibling()
	if el.GetNodeName() != "body" && el.GetFirstChild().(Text).GetValue() != "weekend" {
		t.Errorf("Bad body")
	}
	if el.GetNextSibling().GetNextSibling() != nil {
		t.Errorf("Extra data")
	}
}

func TestCDATA(t *testing.T) {
	input := `<note>
<!CDATA[
characters
]]>
<to>Tove</to>
<from>Jani</from>
</note>`
	dec := xml.NewDecoder(strings.NewReader(input))
	doc, err := Parse(dec)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root := doc.GetDocumentElement()
	if root.GetNodeName() != "note" {
		t.Errorf("Bad root name")
	}
	el := root.GetFirstChild().GetNextSibling()
	if el.GetNodeType() != TEXT_NODE {
		t.Errorf("Not cdata: %s", el.GetNodeName())
	}
	if el.(Text).GetValue() != "\ncharacters\n" {
		t.Errorf("Bad text: %T %v", el, el)
	}
	el = el.GetNextSibling().GetNextSibling()
	if el.GetNodeName() != "to" && el.GetFirstChild().(Text).GetValue() != "Tove" {
		t.Errorf("Bad to")
	}
	el = el.GetNextSibling().GetNextSibling()
	if el.GetNodeName() != "from" && el.GetFirstChild().(Text).GetValue() != "Jani" {
		t.Errorf("Bad from")
	}
	if el.GetNextSibling().GetNextSibling() != nil {
		t.Errorf("Extra data")
	}
}

func TestNS1(t *testing.T) {
	input := `

<root xmlns:h="http://www.w3.org/TR/html4/"
xmlns:f="https://test.com/furniture">

<h:table>
  <h:tr>
    <h:td attr="value">Apples</h:td>
    <h:td>Bananas</h:td>
  </h:tr>
</h:table>

<f:table>
  <f:name>African Coffee Table</f:name>
</f:table>

</root>`

	const hSpace = "http://www.w3.org/TR/html4/"
	const fSpace = "https://test.com/furniture"

	dec := xml.NewDecoder(strings.NewReader(input))
	doc, err := Parse(dec)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	root := doc.GetDocumentElement()
	name := root.(Element).GetQName()
	if name.Local != "root" ||
		name.Space != "" ||
		name.Prefix != "" {
		t.Errorf("Bad root name %v", name)
	}

	el := root.GetFirstChild().GetNextSibling().(Element)
	name = el.GetQName()
	if name.Local != "table" || name.Space != hSpace || name.Prefix != "h" {
		t.Errorf("Bad table: %v", name)
	}

	el = el.GetFirstChild().GetNextSibling().(Element)
	name = el.GetQName()
	if name.Local != "tr" || name.Space != hSpace || name.Prefix != "h" {
		t.Errorf("Bad tr: %v", name)
	}
	el = el.GetFirstChild().GetNextSibling().(Element)
	name = el.GetQName()
	if name.Local != "td" || name.Space != hSpace || name.Prefix != "h" {
		t.Errorf("Bad td: %v", name)
	}

	if v, ok := el.GetAttributeNS("", "attr"); v != "value" || !ok {
		t.Errorf("Wrong attr: %v", el.(*BasicElement).attributes)
	}

	if v, ok := el.GetAttribute("attr"); v != "value" || !ok {
		t.Errorf("Wrong attr: %v", el.(*BasicElement).attributes)
	}

	txt := el.GetFirstChild().(Text)
	if txt.GetValue() != "Apples" {
		t.Errorf("Wrong text: %v", txt)
	}

	el = el.GetNextSibling().GetNextSibling().(Element)
	name = el.GetQName()
	if name.Local != "td" || name.Space != hSpace || name.Prefix != "h" {
		t.Errorf("Bad td: %v", name)
	}
	txt = el.GetFirstChild().(Text)
	if txt.GetValue() != "Bananas" {
		t.Errorf("Wrong text: %v", txt)
	}

	el = el.GetParentElement().GetParentElement().GetNextSibling().GetNextSibling().(Element)
	name = el.GetQName()
	if name.Local != "table" || name.Space != fSpace || name.Prefix != "f" {
		t.Errorf("Bad table: %v", name)
	}
	el = el.GetFirstChild().GetNextSibling().(Element)
	name = el.GetQName()
	if name.Local != "name" || name.Space != fSpace || name.Prefix != "f" {
		t.Errorf("Bad name: %v", name)
	}

	if el.GetRootNode() != doc {
		t.Errorf("root")
	}

	if !doc.Contains(el) {
		t.Errorf("Contains")
	}
}

func TestAutoClose(t *testing.T) {
	input := `<note>
<to>Tove<br></to>
<from>Jani</from>
</note>`
	dec := xml.NewDecoder(strings.NewReader(input))
	dec.AutoClose = xml.HTMLAutoClose
	dec.Strict = false
	_, err := Parse(dec)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}
