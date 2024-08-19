package ab

import (
	"bytes"
	"encoding/xml"
	"github.com/subchen/go-xmldom"
	"io/ioutil"
	"strings"
)

// Node represents a simple node structure for XML processing.
//type Node struct {
//	XMLName xml.Name
//	Attrs   []xml.Attr `xml:",any,attr"`
//	Content []byte     `xml:",innerxml"`
//	Nodes   []Node     `xml:",any"`
//}

// ParseFile parses an XML file and returns the root node.
func ParseFile(fileName string) (*xmldom.Node, error) {
	xmlData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return ParseXML(xmlData)
}

// ParseString parses an XML string and returns the root node.
func ParseString(xmlString string) (*xmldom.Node, error) {
	return ParseXML([]byte(xmlString))
}

// ParseXML parses XML data and returns the root node.
func ParseXML(xmlData []byte) (*xmldom.Node, error) {
	doc := xmldom.Must(xmldom.ParseXML(string(xmlData)))
	root := doc.Root
	return root, nil
	//var node xmldom.Node
	//decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	//err := decoder.Decode(&node)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &node, nil
}

// NodeToString converts an XML node back to its string representation.
func NodeToString(node *xmldom.Node) string {
	//var buf bytes.Buffer
	//encoder := xml.NewEncoder(&buf)
	//err := encoder.Encode(node)
	//if err != nil {
	//	panic(err)
	//}
	//return buf.String()
	//return node.XMLPretty()

	if node == nil {
		return ""
	}

	// Convert the node to a string without the XML declaration and indentation
	//result := node.XML()
	//return result
	buf := new(bytes.Buffer)
	printXML(buf, node, 0, "")
	return buf.String()
}

func printXML(buf *bytes.Buffer, n *xmldom.Node, level int, indent string) {
	pretty := len(indent) > 0

	if pretty {
		buf.WriteString(strings.Repeat(indent, level))
	}
	buf.WriteByte('<')
	buf.WriteString(n.Name)

	if len(n.Attributes) > 0 {
		for _, attr := range n.Attributes {
			buf.WriteByte(' ')
			buf.WriteString(attr.Name)
			buf.WriteByte('=')
			buf.WriteByte('"')
			xml.Escape(buf, []byte(attr.Value))
			buf.WriteByte('"')
		}
	}

	if len(n.Children) == 0 && len(n.Text) == 0 {
		buf.WriteString(" />")
		if pretty {
			buf.WriteByte('\n')
		}
		return
	}

	buf.WriteByte('>')

	if len(n.Text) > 0 {
		xml.EscapeText(buf, []byte(n.Text))
	}

	if len(n.Children) > 0 {
		if pretty {
			buf.WriteByte('\n')
		}
		for _, c := range n.Children {
			printXML(buf, c, level+1, indent)
		}
	}

	if len(n.Children) > 0 && len(indent) > 0 {
		buf.WriteString(strings.Repeat(indent, level))
	}
	buf.WriteString("</")
	buf.WriteString(n.Name)
	buf.WriteByte('>')

	if pretty {
		buf.WriteByte('\n')
	}
}
