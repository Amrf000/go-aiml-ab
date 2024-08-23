package dom

type Attr interface {
	Node

	// A String representing the local part of the qualified name of the
	// attribute.
	GetLocalName() string

	// Returns the qualified name of an attribute, that is the name of
	// the attribute, with the namespace prefix, if any, in front of
	// it. For example, if the local name is lang and the namespace
	// prefix is xml, the returned qualified name is xml:lang.
	GetName() string

	GetQName() Name

	// The Element the attribute belongs to.
	GetOwnerElement() Element

	// // A String representing the URI of the namespace of the attribute, or
	// // null if there is no namespace.
	// GetNamespaceURI() string

	// // A String representing the namespace prefix of the attribute, or
	// // null if a namespace without prefix or no namespace are specified.
	// GetPrefix() string

	// The attribute's value, a string that can be set and get using this
	// property.
	GetValue() string

	SetValue(string)
}
