package dom

type BasicElement struct {
	basicNode

	attributes basicNamedNodeMap
	name       Name
}

var _ Element = &BasicElement{}

// Returns a boolean value indicating whether or not the two nodes are
// the same (that is, they reference the same object).
func (el *BasicElement) IsSameNode(node Node) bool { return node == el }

func (el *BasicElement) IsEqualNode(node Node) bool {
	el2, ok := node.(*BasicElement)
	if !ok {
		return false
	}
	if el.name != el2.name {
		return false
	}
	if len(el.attributes.attrs) != len(el2.attributes.attrs) {
		return false
	}
	for _, attr := range el.attributes.attrs {
		attr2, ok := el2.attributes.mapAttrs[attr.name.Name]
		if !ok {
			return false
		}
		if !attr.IsEqualNode(attr2) {
			return false
		}
	}
	return isEqualNode(el, el2)
}

// Returns HTML uppercased qualified name
func (el *BasicElement) GetNodeName() string { return el.getQName() }

// Returns ELEMENT_NODE
func (el *BasicElement) GetNodeType() NodeType { return ELEMENT_NODE }

func (el *BasicElement) getQName() string {
	return el.name.QName()
}

func (el *BasicElement) GetQName() Name {
	return el.name
}

// Returns a String with the name of the tag for the given element.
func (el *BasicElement) GetTagName() string {
	return el.name.QName()
}

// Returns a string representing the namespace prefix of the element,
// or "" if no prefix is specified.
func (el *BasicElement) GetPrefix() string {
	return el.name.Prefix
}

// A string representing the local part of the qualified name of the
// element.
func (el *BasicElement) GetLocalName() string {
	return el.name.Local
}

// Accepts a namespace URI as an argument and returns a boolean value
// with a value of true if the namespace is the default namespace on
// the given node or false if not.
func (el *BasicElement) IsDefaultNamespace(uri string) bool {
	if len(uri) == 0 {
		return false
	}
	return el.LookupNamespaceURI("") == uri
}

// // The namespace URI of the element, or "" if it is no namespace.}
// func (el *BasicElement) GetNamespaceURI() string {
// 	return el.space
// }

func (el *BasicElement) GetFirstElementChild() Element {
	return nextElementSibling(el.GetFirstChild())
}

func (el *BasicElement) GetLastElementChild() Element {
	return prevElementSibling(el.GetLastChild())
}

func (el *BasicElement) GetNextElementSibling() Element {
	return nextElementSibling(el.GetNextSibling())
}

func (el *BasicElement) GetPreviousElementSibling() Element {
	return prevElementSibling(el.GetPreviousSibling())
}

// Returns a string  containing the prefix for a given namespace
// URI, if present, and "" if not. When multiple prefixes are
// possible, the result is implementation-dependent.
func (el *BasicElement) LookupPrefix(uri string) string {
	if el.name.Space == uri && len(el.name.Prefix) > 0 {
		return el.name.Prefix
	}
	for _, attr := range el.attributes.attrs {
		if attr.name.Prefix == xmlnsPrefix && attr.value == uri {
			return attr.name.Local
		}
	}
	if el.parent == nil {
		return ""
	}
	if _, ok := el.parent.(*BasicDocument); ok {
		return ""
	}
	return el.parent.LookupPrefix(uri)
}

// Accepts a prefix and returns the namespace URI associated with it
// on the given node if found (and "" if not). Supplying "" for
// the prefix will return the default namespace.
func (el *BasicElement) LookupNamespaceURI(prefix string) string {
	switch prefix {
	case xmlPrefix:
		return xmlURL
	case xmlnsPrefix:
		return xmlnsURL
	}
	if len(el.name.Space) > 0 && el.name.Prefix == prefix {
		return el.name.Space
	}
	for _, attr := range el.attributes.attrs {
		if attr.name.Space == xmlnsURL {
			if attr.name.Prefix == xmlnsPrefix && attr.name.Local == prefix {
				return attr.value
			}
			if len(prefix) == 0 && len(attr.name.Prefix) == 0 && attr.name.Local == xmlnsPrefix {
				return attr.value
			}
		}
	}
	if el.parent == nil {
		return ""
	}
	if _, ok := el.parent.(*BasicDocument); ok {
		return ""
	}
	return el.parent.LookupNamespaceURI(prefix)
}

func nextElementSibling(start Node) Element {
	if start == nil {
		return nil
	}
	for trc := start; trc != nil; trc = trc.GetNextSibling() {
		if el, ok := trc.(Element); ok {
			return el
		}
	}
	return nil
}

func prevElementSibling(start Node) Element {
	if start == nil {
		return nil
	}
	for trc := start; trc != nil; trc = trc.GetPreviousSibling() {
		if el, ok := trc.(Element); ok {
			return el
		}
	}
	return nil
}

// Removes the element from the children list of its parent.
func (el *BasicElement) Remove() {
	detachChild(el.GetParentNode(), el)
}

func (el *BasicElement) GetAttributes() NamedNodeMap {
	return &BasicNamedNodeMap{
		owner: el,
	}
}

// // 	GetID() string

// Retrieves the value of the named attribute from the current node
// and returns it as a string.
func (el *BasicElement) GetAttribute(name string) (string, bool) {
	attr := el.GetAttributeNode(name)
	if attr == nil {
		return "", false
	}
	return attr.GetValue(), true
}

// Returns an array of attribute names from the current element.
func (el *BasicElement) GetAttributeNames() []string {
	ret := make([]string, el.attributes.GetLength())
	for i := 0; i < len(ret); i++ {
		ret[i] = el.attributes.attrs[i].name.QName()
	}
	return ret
}

// Retrieves the node representation of the named attribute from the
// current node and returns it as an Attr.
func (el *BasicElement) GetAttributeNode(name string) Attr {
	return el.attributes.GetNamedItemNS("", name)
}

// Retrieves the node representation of the attribute with the
// specified name and namespace, from the current node and returns
// it as an Attr.
func (el *BasicElement) GetAttributeNodeNS(uri, name string) Attr {
	return el.attributes.GetNamedItemNS(uri, name)
}

// Retrieves the value of the attribute with the specified
// namespace and name from the current node and returns it as a
// string.
func (el *BasicElement) GetAttributeNS(uri string, name string) (string, bool) {
	attr := el.GetAttributeNodeNS(uri, name)
	if attr == nil {
		return "", false
	}
	return attr.GetValue(), true
}

// Returns a boolean value indicating if the element has the
// specified attribute or not.
func (el *BasicElement) HasAttribute(name string) bool {
	return el.GetAttributeNode(name) != nil
}

// Returns a boolean value indicating if the element has the
// specified attribute, in the specified namespace, or not.
func (el *BasicElement) HasAttributeNS(uri string, name string) bool {
	return el.GetAttributeNodeNS(uri, name) != nil
}

// Removes the named attribute from the current node.
func (el *BasicElement) RemoveAttribute(name string) {
	el.attributes.RemoveNamedItemNS("", name)
}

// Removes the node representation of the named attribute from the
// current node.
func (el *BasicElement) RemoveAttributeNode(attr Attr) {
	if attr.GetParentElement() == el {
		el.attributes.removeAttr(attr)
	}
}

// Removes the attribute with the specified name and namespace, from
// the current node.
func (el *BasicElement) RemoveAttributeNS(uri string, name string) {
	el.attributes.RemoveNamedItemNS(uri, name)
}

// Sets the value of a named attribute of the current node.
func (el *BasicElement) SetAttribute(name string, value string) {
	existing := el.attributes.GetNamedItemNS("", name)
	if existing != nil {
		existing.SetValue(value)
		return
	}
	attr := el.ownerDocument.CreateAttribute(name)
	attr.SetValue(value)
	el.attributes.setNamedItemNS(el, attr)
}

// Sets the value of the attribute with the specified name and
// namespace, from the current node.
func (el *BasicElement) SetAttributeNS(prefix, uri, name string, value string) {
	existing := el.attributes.GetNamedItemNS(uri, name)
	if existing != nil {
		existing.SetValue(value)
		return
	}
	attr := el.ownerDocument.CreateAttributeNS(prefix, uri, name)
	attr.SetValue(value)
	el.attributes.setNamedItemNS(el, attr)
}

// Sets the node representation of the named attribute from the
// current node.
func (el *BasicElement) SetAttributeNode(attr Attr) {
	el.attributes.setNamedItemNS(el, attr)
}

// Sets the node representation of the attribute with the specified
// name and namespace, from the current node.
func (el *BasicElement) SetAttributeNodeNS(attr Attr) {
	el.attributes.setNamedItemNS(el, attr)
}

func (el *BasicElement) InsertBefore(newNode, referenceNode Node) Node {
	if err := validatePreInsertion(newNode, el, referenceNode, "InsertBefore"); err != nil {
		panic(err)
	}
	return insertBefore(el, newNode, referenceNode)
}

// Append newNode as a child of node
func (el *BasicElement) AppendChild(newNode Node) Node {
	if err := validatePreInsertion(newNode, el, nil, "AppendChild"); err != nil {
		panic(err)
	}
	return insertBefore(el, newNode, nil)
}

// Remove child from node
func (el *BasicElement) RemoveChild(child Node) {
	if child.GetParentNode() != el {
		panic(ErrDOM{
			Typ: NOT_FOUND_ERR,
			Msg: "Wrong parent",
			Op:  "RemoveChild",
		})
	}
	detachChild(el, child)
}

func (el *BasicElement) Normalize() {
	// Combine all text nodes
	for childNode := el.GetFirstChild(); childNode != nil; {
		childNode.Normalize()
		text, ok := childNode.(*BasicText)
		if !ok {
			childNode = childNode.GetNextSibling()
			continue
		}
		// This is a text node. Combine with the next text node
		nextNode := text.GetNextSibling()
		for {
			if nextNode == nil {
				childNode = nil
				break
			}
			nextText, ok := nextNode.(*BasicText)
			if !ok {
				childNode = nextNode
				break
			}
			text.text += nextText.text
			nextNode = nextText.GetNextSibling()
			detachChild(el, nextText)
		}
	}
}

func (el *BasicElement) CloneNode(deep bool) Node {
	return el.cloneNode(el.ownerDocument, deep)
}

func (el *BasicElement) cloneNode(owner Document, deep bool) Node {
	newElement := owner.CreateElement("").(*BasicElement)
	newElement.name = el.name
	for _, attr := range el.attributes.attrs {
		newAttr := attr.cloneNode(owner, deep).(*BasicAttr)
		newElement.attributes.setNamedItemNS(newElement, newAttr)
	}
	if deep {
		for child := el.GetFirstChild(); child != nil; child = child.GetNextSibling() {
			newEl := child.cloneNode(owner, deep)
			newElement.AppendChild(newEl)
		}
	}
	return newElement
}

func (el *BasicElement) getNamespaceInfo() (defaultNS string, definedPrefixes map[string]string) {
	for _, attr := range el.attributes.attrs {
		if len(attr.name.Space) == 0 && attr.name.Local == xmlnsPrefix {
			defaultNS = attr.value
		} else if attr.name.Space == xmlnsURL {
			if definedPrefixes == nil {
				definedPrefixes = make(map[string]string)
			}
			definedPrefixes[attr.name.Local] = attr.value
		}
	}
	return
}
