package ab

import (
	"fmt"
	"github.com/subchen/go-xmldom"
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

func (p *PCAIMLProcessorExtension) NewContact(node *xmldom.Node, ps *ParseState) string {
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

func (p *PCAIMLProcessorExtension) ContactID(node *xmldom.Node, ps *ParseState) string {
	// Implement logic to handle <contactid> tag
	displayName := EvalTagContent(node, ps, nil)
	result := ContactId(displayName)
	return result
}

func (p *PCAIMLProcessorExtension) MultipleIds(node *xmldom.Node, ps *ParseState) string {
	contactName := EvalTagContent(node, ps, nil)
	result := MultipleIds(contactName)
	//System.out.println("multipleIds("+contactName+")="+result);
	return result
}
func (p *PCAIMLProcessorExtension) DisplayName(node *xmldom.Node, ps *ParseState) string {
	id := EvalTagContent(node, ps, nil)
	result := DisplayName(id)
	//System.out.println("displayName("+id+")="+result);
	return result
}

func (p *PCAIMLProcessorExtension) DialNumber(node *xmldom.Node, ps *ParseState) string {
	childList := node.Children
	id := "unknown"
	itype := "unknown"
	for i := 0; i < len(childList); i++ {
		if childList[i].Name == "id" {
			id = EvalTagContent(childList[i], ps, nil)
		}
		if childList[i].Name == "type" {
			itype = EvalTagContent(childList[i], ps, nil)
		}
	}
	result := DialNumber(itype, id)
	//System.out.println("dialNumber("+id+")="+result);
	return result
}

func (p *PCAIMLProcessorExtension) EmailAddress(node *xmldom.Node, ps *ParseState) string {
	childList := node.Children
	id := "unknown"
	itype := "unknown"
	for i := 0; i < len(childList); i++ {
		if childList[i].Name == "id" {
			id = EvalTagContent(childList[i], ps, nil)
		}
		if childList[i].Name == "type" {
			itype = EvalTagContent(childList[i], ps, nil)
		}
	}
	result := EmailAddress(itype, id)
	//System.out.println("emailAddress("+id+")="+result);
	return result
}

func (p *PCAIMLProcessorExtension) ContactBirthday(node *xmldom.Node, ps *ParseState) string {
	id := EvalTagContent(node, ps, nil)
	result := Birthday(id)
	//System.out.println("birthday("+id+")="+result);
	return result
}

func (p *PCAIMLProcessorExtension) RecursEval(node *xmldom.Node, ps *ParseState) string {
	nodeName := node.Name
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
