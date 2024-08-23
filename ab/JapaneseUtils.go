package ab

import (
	"aiml/external/go-dom"
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

func recursEval(element dom.Element) string {
	switch element.GetNodeName() {
	case "#text":
		return tokenizeFragment(element.(dom.Text).GetValue())
	case "sentence":
		return evalTagContent(element)
	default:
		return genericXML(element)
	}
}

func genericXML(element dom.Element) string {
	result := evalTagContent(element)
	return unevaluatedXML(result, element)
}

func evalTagContent(element dom.Element) string {
	var result string
	childList := element.GetChildNodes()
	for i := 0; i < childList.GetLength(); i++ {
		child := childList.Item(i).(dom.Element)
		result += recursEval(child)
	}
	return result
}

//func unevaluatedXML(result string, element dom.Element) string {
//	var buffer bytes.Buffer
//	buffer.WriteString(" <")
//	buffer.WriteString(element.GetNodeName())
//	for _, attr := range element.GetAttributes() {
//		buffer.WriteString(fmt.Sprintf(" %s=\"%s\"", attr.Name, attr.Value))
//	}
//	if result == "" {
//		buffer.WriteString("/> ")
//	} else {
//		buffer.WriteString(">")
//		buffer.WriteString(result)
//		buffer.WriteString("</")
//		buffer.WriteString(element.GetNodeName())
//		buffer.WriteString("> ")
//	}
//	return buffer.String()
//}

func unevaluatedXML(result string, node dom.Element) string {
	nodeName := node.GetNodeName()
	attributes := ""
	if node.GetAttributes().GetLength() > 0 {
		XMLAttributes := node.GetAttributes()
		for i := 0; i < XMLAttributes.GetLength(); i++ {
			attributes += " " + XMLAttributes.Item(i).GetNodeName() + "=\"" + XMLAttributes.Item(i).GetValue() + "\""
		}
	}
	if result == "" {
		return " <" + nodeName + attributes + "/> "
	} else {
		return " <" + nodeName + attributes + ">" + result + "</" + nodeName + "> " // add spaces
	}
}
