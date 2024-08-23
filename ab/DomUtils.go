package ab

import (
	"aiml/external/go-dom"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

// Node represents a simple node structure for XML processing.
//type Node struct {
//	XMLName xml.Name
//	Attrs   []xml.Attr `xml:",any,attr"`
//	Content []byte     `xml:",innerxml"`
//	Nodes   []Node     `xml:",any"`
//}

// ParseFile parses an XML file and returns the root node.
func ParseFile(fileName string) (dom.Element, error) {
	xmlData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return ParseXML(xmlData, fileName)
}

// ParseString parses an XML string and returns the root node.
func ParseString(xmlString string) (dom.Element, error) {
	return ParseXML([]byte(xmlString), "")
}

// ParseXML parses XML data and returns the root node.
func ParseXML(xmlData []byte, filename string) (dom.Element, error) {
	dec := xml.NewDecoder(bytes.NewReader(xmlData))
	doc, err := dom.Parse(dec)
	if err != nil {
		fmt.Println(filename)
		panic(err)
	}
	root := doc.GetDocumentElement()
	return root, nil
	//var node dom.Node
	//decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	//err := decoder.Decode(&node)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &node, nil
}

// NodeToString converts an XML node back to its string representation.
func NodeToString(node dom.Node) string {
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
	buf := bytes.Buffer{}
	if err := dom.Encode(node, &buf); err != nil {
		panic(err)
		return ""
	}

	return buf.String()
	// Convert the node to a string without the XML declaration and indentation
	//result := node.XML()
	//return result
	//buf := new(bytes.Buffer)
	//printXML(buf, node, 0, "")
	//return buf.String()
}

//func printXML(buf *bytes.Buffer, n dom.Node, level int, indent string) {
//	pretty := len(indent) > 0
//
//	if pretty {
//		buf.WriteString(strings.Repeat(indent, level))
//	}
//	buf.WriteByte('<')
//	buf.WriteString(n.GetNodeName())
//
//	if len(n.Attributes) > 0 {
//		for _, attr := range n.Attributes {
//			buf.WriteByte(' ')
//			buf.WriteString(attr.Name)
//			buf.WriteByte('=')
//			buf.WriteByte('"')
//			xml.Escape(buf, []byte(attr.Value))
//			buf.WriteByte('"')
//		}
//	}
//
//	if len(n.Children) == 0 && len(n.Text) == 0 {
//		buf.WriteString(" />")
//		if pretty {
//			buf.WriteByte('\n')
//		}
//		return
//	}
//
//	buf.WriteByte('>')
//
//	if len(n.Text) > 0 {
//		xml.EscapeText(buf, []byte(n.Text))
//	}
//
//	if len(n.Children) > 0 {
//		if pretty {
//			buf.WriteByte('\n')
//		}
//		for _, c := range n.Children {
//			printXML(buf, c, level+1, indent)
//		}
//	}
//
//	if len(n.Children) > 0 && len(indent) > 0 {
//		buf.WriteString(strings.Repeat(indent, level))
//	}
//	buf.WriteString("</")
//	buf.WriteString(n.Name)
//	buf.WriteByte('>')
//
//	if pretty {
//		buf.WriteByte('\n')
//	}
//}
