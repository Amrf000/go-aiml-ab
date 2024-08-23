package dom

import (
	"testing"
)

func TestDTDParse(t *testing.T) {
	ret, ok, err := ParseDocumentType([]byte(`!DOCTYPE note
[
<!ELEMENT note (to,from,heading,body)>
<!ELEMENT to (#PCDATA)>
]`))
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf("Not DTD")
	}
	dt := ret.(*BasicDocumentType)
	if dt.GetName() != "note" {
		t.Errorf("Wrong name")
	}
	if dt.defn != `[
<!ELEMENT note (to,from,heading,body)>
<!ELEMENT to (#PCDATA)>
]` {
		t.Errorf("Wrong defn: %s", dt.defn)
	}
}
