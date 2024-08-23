package ab

import (
	"aiml/external/go-dom"
	"fmt"
)

type PCAIMLProcessorExtension struct {
	ExtensionTagNames map[string]bool
}

func NewPCAIMLProcessorExtension() *PCAIMLProcessorExtension {
	return &PCAIMLProcessorExtension{
		ExtensionTagNames: map[string]bool{
			"contactid":       true,
			"multipleids":     true,
			"displayname":     true,
			"dialnumber":      true,
			"emailaddress":    true,
			"contactbirthday": true,
			"addinfo":         true,
		},
	}
}

func (p *PCAIMLProcessorExtension) ExtensionTagSet() map[string]bool {
	return p.ExtensionTagNames
}

func (p *PCAIMLProcessorExtension) NewContact(node dom.Node, ps *ParseState) string {
	// Implement logic to handle <addinfo> tag
	displayName := EvalTagContent(node, ps, nil)
	phoneType := "unknown"
	dialNumber := "unknown"
	emailType := "unknown"
	emailAddress := "unknown"
	birthday := "unknown"

	fmt.Printf("Adding new contact %s %s %s %s %s %s\n", displayName, phoneType, dialNumber, emailType, emailAddress, birthday)

	// Assuming Contact struct and relevant methods are defined similarly as in Java
	// contact := NewContact(displayName, phoneType, dialNumber, emailType, emailAddress, birthday)
	return ""
}

func (p *PCAIMLProcessorExtension) ContactID(node dom.Node, ps *ParseState) string {
	// Implement logic to handle <contactid> tag
	displayName := EvalTagContent(node, ps, nil)
	result := ContactId(displayName)
	return result
}

func (p *PCAIMLProcessorExtension) MultipleIds(node dom.Node, ps *ParseState) string {
	contactName := EvalTagContent(node, ps, nil)
	result := MultipleIds(contactName)
	//System.out.println("multipleIds("+contactName+")="+result);
	return result
}
func (p *PCAIMLProcessorExtension) DisplayName(node dom.Node, ps *ParseState) string {
	id := EvalTagContent(node, ps, nil)
	result := DisplayName(id)
	//System.out.println("displayName("+id+")="+result);
	return result
}

func (p *PCAIMLProcessorExtension) DialNumber(node dom.Node, ps *ParseState) string {
	childList := node.GetChildNodes()
	id := "unknown"
	itype := "unknown"
	for i := 0; i < childList.GetLength(); i++ {
		nodeName := childList.Item(i).GetNodeName()
		if nodeName == "id" {
			id = EvalTagContent(childList.Item(i), ps, nil)
		}
		if nodeName == "type" {
			itype = EvalTagContent(childList.Item(i), ps, nil)
		}
	}
	result := DialNumber(itype, id)
	//System.out.println("dialNumber("+id+")="+result);
	return result
}

func (p *PCAIMLProcessorExtension) EmailAddress(node dom.Node, ps *ParseState) string {
	childList := node.GetChildNodes()
	id := "unknown"
	itype := "unknown"
	for i := 0; i < childList.GetLength(); i++ {
		nodeName := childList.Item(i).GetNodeName()
		if nodeName == "id" {
			id = EvalTagContent(childList.Item(i), ps, nil)
		}
		if nodeName == "type" {
			itype = EvalTagContent(childList.Item(i), ps, nil)
		}
	}
	result := EmailAddress(itype, id)
	//System.out.println("emailAddress("+id+")="+result);
	return result
}

func (p *PCAIMLProcessorExtension) ContactBirthday(node dom.Node, ps *ParseState) string {
	id := EvalTagContent(node, ps, nil)
	result := Birthday(id)
	//System.out.println("birthday("+id+")="+result);
	return result
}

func (p *PCAIMLProcessorExtension) RecursEval(node dom.Node, ps *ParseState) string {
	nodeName := node.GetNodeName()
	switch nodeName {
	case "contactid":
		return p.ContactID(node, ps)
	case "multipleids":
		return p.MultipleIds(node, ps)
	case "dialnumber":
		return p.DialNumber(node, ps)
	case "addinfo":
		return p.NewContact(node, ps)
	case "displayname":
		return p.DisplayName(node, ps)
	case "emailaddress":
		return p.EmailAddress(node, ps)
	case "contactbirthday":
		return p.ContactBirthday(node, ps)
	default:
		return GenericXML(node, ps)
	}
}
