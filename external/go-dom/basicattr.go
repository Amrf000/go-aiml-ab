package dom

type BasicAttr struct {
	basicNode

	name Name

	value string
}

var _ Attr = &BasicAttr{}

// Returns a boolean value indicating whether or not the two nodes are
// the same (that is, they reference the same object).
func (attr *BasicAttr) IsSameNode(node Node) bool { return node == attr }

// A String representing the local part of the qualified name of the
// attribute.
func (attr *BasicAttr) GetLocalName() string {
	return attr.name.Local
}

// Returns the qualified name of an attribute, that is the name of
// the attribute, with the namespace prefix, if any, in front of
// it. For example, if the local name is lang and the namespace
// prefix is xml, the returned qualified name is xml:lang.
func (attr *BasicAttr) GetName() string {
	return attr.name.QName()
}

func (attr *BasicAttr) GetQName() Name {
	return attr.name
}

// The Element the attribute belongs to.
func (attr *BasicAttr) GetOwnerElement() Element {
	if el, ok := attr.parent.(Element); ok {
		return el
	}
	return nil
}

// The attribute's value, a string that can be set and get using this
// property.
func (attr *BasicAttr) GetValue() string {
	return attr.value
}

func (attr *BasicAttr) SetValue(v string) {
	attr.value = v
}

func (attr *BasicAttr) CloneNode(bool) Node {
	return attr.cloneNode(attr.ownerDocument, false)
}

func (attr *BasicAttr) cloneNode(owner Document, deep bool) Node {
	ret := BasicAttr{
		basicNode: basicNode{
			ownerDocument: owner.(*BasicDocument),
		},
		name:  attr.name,
		value: attr.value,
	}
	return &ret
}

func (attr *BasicAttr) GetNodeName() string {
	return attr.GetName()
}

func (attr *BasicAttr) GetNodeType() NodeType { return ATTRIBUTE_NODE }

// Returns a boolean value which indicates whether or not two nodes
// are of the same type and all their defining data points match.
func (attr *BasicAttr) IsEqualNode(node Node) bool {
	a, ok := node.(*BasicAttr)
	if !ok {
		return false
	}
	return a.name == attr.name && a.value == attr.value
}
