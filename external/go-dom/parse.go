package dom

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// Parses an XML document.
//
// If decoder.Strict is false, the parser looks at decoder.AutoClose
// to handle auto-closing HTML tags. Otherwise it is a strict XML
// parser.
func Parse(decoder *xml.Decoder) (ret Document, resultErr error) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				resultErr = e
			} else {
				resultErr = fmt.Errorf("%v", err)
			}
		}
	}()

	ret = NewDocument()

	interner := make(map[string]string)
	intern := func(s string) string {
		existing, ok := interner[s]
		if ok {
			return existing
		}
		interner[s] = s
		return s
	}

	autoClose := func(name xml.Name) bool {
		if decoder.Strict {
			return false
		}
		for _, str := range decoder.AutoClose {
			if strings.EqualFold(str, name.Local) {
				return true
			}
		}
		return false
	}

	elementStack := make([]xml.Name, 0, 16)

	var parent *BasicElement
	autoCloseSeen := false

	closeAutoClose := func() {
		if !autoCloseSeen {
			return
		}
		autoCloseSeen = false
		elementStack = elementStack[:len(elementStack)-1]
		par := parent.GetParentNode()
		if _, ok := par.(*BasicDocument); ok {
			parent = nil
		} else {
			parent = par.(*BasicElement)
		}
	}

	isSpaceOrEmpty := func(s string) bool {
		for _, x := range s {
			if !unicode.IsSpace(x) {
				return false
			}
		}
		return true
	}

	for {
		tok, err := decoder.RawToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch token := tok.(type) {
		case xml.StartElement:

			closeAutoClose()

			elementStack = append(elementStack, token.Name)
			newElement := ret.CreateElement(intern(token.Name.Local)).(*BasicElement) // Create an empty element for now
			newElement.name.Prefix = intern(token.Name.Space)
			if parent == nil {
				ret.AppendChild(newElement) // This is the document element
			} else {
				parent.AppendChild(newElement)
			}

			// First, create all attributes without namespaces
			for _, attr := range token.Attr {
				newAttr := ret.CreateAttribute(intern(attr.Name.Local)).(*BasicAttr)
				newAttr.name.Prefix = intern(attr.Name.Space)
				newAttr.value = attr.Value
				newAttr.parent = newElement
				newElement.attributes.attrs = append(newElement.attributes.attrs, newAttr)
			}

			// Now process all xmlns attributes
			for _, attr := range newElement.attributes.attrs {
				if attr.name.Prefix == xmlnsPrefix {
					attr.name.Space = xmlnsURL
					if newElement.name.Prefix == attr.name.Local {
						newElement.name.Space = intern(attr.value)
						attr.value = intern(attr.value)
					}
				} else if len(attr.name.Prefix) == 0 && attr.name.Local == xmlnsPrefix {
					if len(newElement.name.Prefix) == 0 {
						newElement.name.Space = intern(attr.value)
						attr.value = intern(attr.value)
					}
				} else {
					// If attr has prefix, then we have to find namespace
					if len(attr.name.Prefix) > 0 {
						// Is namespace defined here?
						for _, a := range newElement.attributes.attrs {
							if a.name.Prefix == xmlnsPrefix && a.name.Local == attr.name.Prefix {
								attr.name.Space = a.value
								break
							}
						}
						if len(attr.name.Space) == 0 && parent != nil {
							attr.name.Space = parent.LookupNamespaceURI(attr.name.Prefix)
						}
					}
				}
			}
			newElement.attributes.mapAttrs = make(map[xml.Name]*BasicAttr)
			for _, a := range newElement.attributes.attrs {
				newElement.attributes.mapAttrs[a.name.Name] = a
			}
			// If namespace is not yet resolved, resolve it
			if len(newElement.name.Space) == 0 {
				newElement.name.Space = newElement.LookupNamespaceURI(newElement.name.Prefix)
			}

			parent = newElement
			if autoClose(token.Name) {
				autoCloseSeen = true
			}

		case xml.EndElement:
			if len(elementStack) == 0 {
				return nil, &xml.SyntaxError{
					Msg: "Extra objects before document",
				}
			}
			if autoCloseSeen {
				if elementStack[len(elementStack)-1] == token.Name {
					autoCloseSeen = false
					elementStack = elementStack[:len(elementStack)-1]
					break
				}
				closeAutoClose()
			}

			last := elementStack[len(elementStack)-1]
			if last.Space != token.Name.Space || !strings.EqualFold(last.Local, token.Name.Local) {
				return nil, &xml.SyntaxError{
					Msg: fmt.Sprintf("Mismatched closing tag %s", token.Name.Local),
				}
			}
			elementStack = elementStack[:len(elementStack)-1]
			par := parent.GetParentNode()
			if _, ok := par.(*BasicDocument); ok {
				parent = nil
			} else {
				parent = par.(*BasicElement)
			}

		case xml.CharData:
			if len(elementStack) == 0 {
				// charData must be only spaces
				if !isSpaceOrEmpty(string(token)) {
					return nil, &xml.SyntaxError{
						Msg: "Extra characters before document",
					}
				}
			} else {
				newNode := ret.CreateTextNode(string(token))
				parent.AppendChild(newNode)
			}

		case xml.Comment:
			if len(elementStack) == 0 {
				ret.AppendChild(ret.CreateComment(string(token)))
			} else {
				parent.AppendChild(ret.CreateComment(string(token)))
			}

		case xml.ProcInst:
			closeAutoClose()
			newNode := ret.CreateProcessingInstruction(token.Target, string(token.Inst))
			if parent == nil {
				ret.AppendChild(newNode)
			} else {
				parent.AppendChild(newNode)
			}

		case xml.Directive:
			content := string(token)
			if strings.HasPrefix(content, "CDATA[") && strings.HasSuffix(content, "]]") {
				if len(elementStack) == 0 {
					return nil, &xml.SyntaxError{
						Msg: "CDATA before document",
					}
				}
				newNode := ret.CreateTextNode(string(content[6 : len(content)-2]))
				parent.AppendChild(newNode)
			} else {
				documentType, ok, err := ParseDocumentType([]byte(token))
				if err != nil {
					return nil, err
				}
				if ok {
					parent.AppendChild(documentType)
				}
			}
		}
	}
	return ret, nil
}
