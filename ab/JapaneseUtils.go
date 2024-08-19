package ab

import (
	"bytes"
	"fmt"
	"github.com/subchen/go-xmldom"
	"strings"
)

func tokenizeFragment(fragment string) string {
	//System.out.println("buildFragment "+fragment);
	result := ""
	items := Parse(fragment)
	for _, e := range items {
		result += e.Surface + " "
		//
		// System.out.println("Feature "+e.feature+" Surface="+e.surface);
	}
	return strings.TrimSpace(result)
}

func TokenizeSentence(sentence string) string {
	if !JpTokenize {
		return sentence
	}

	result := tokenizeXML(sentence)
	result = strings.ReplaceAll(result, "$ ", "$")
	result = strings.ReplaceAll(result, "  ", " ")
	result = strings.ReplaceAll(result, "anon ", "anon")
	return strings.TrimSpace(result)
}

func tokenizeXML(xmlExpression string) string {
	//System.out.println("tokenizeXML "+xmlExpression);
	response := TemplateFailed
	xmlExpression = "<sentence>" + xmlExpression + "</sentence>"
	root, _ := ParseString(xmlExpression)
	response = recursEval(root)
	return TrimTag(response, "sentence")
}

func recursEval(element *xmldom.Node) string {
	switch element.Name {
	case "#text":
		return tokenizeFragment(element.Text)
	case "sentence":
		return evalTagContent(element)
	default:
		return genericXML(element)
	}
}

func genericXML(element *xmldom.Node) string {
	result := evalTagContent(element)
	return unevaluatedXML(result, element)
}

func evalTagContent(element *xmldom.Node) string {
	var result string
	for _, child := range element.Children {
		result += recursEval(child)
	}
	return result
}

func unevaluatedXML(result string, element *xmldom.Node) string {
	var buffer bytes.Buffer
	buffer.WriteString(" <")
	buffer.WriteString(element.Name)
	for _, attr := range element.Attributes {
		buffer.WriteString(fmt.Sprintf(" %s=\"%s\"", attr.Name, attr.Value))
	}
	if result == "" {
		buffer.WriteString("/> ")
	} else {
		buffer.WriteString(">")
		buffer.WriteString(result)
		buffer.WriteString("</")
		buffer.WriteString(element.Name)
		buffer.WriteString("> ")
	}
	return buffer.String()
}
