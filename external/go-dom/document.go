package dom

type Document interface {
	Node

	// Returns the Element that is a direct child of the document.
	GetDocumentElement() Element

	// Creates a new Attr object and returns it.
	CreateAttribute(string) Attr

	// Creates a new attribute node in a given namespace and returns it.
	CreateAttributeNS(prefix string, ns string, name string) Attr

	// Creates a new comment node and returns it.
	CreateComment(string) Comment

	// Creates a new element with the given tag name.
	CreateElement(string) Element

	// Creates a new element with the given tag name and namespace URI.
	CreateElementNS(prefix string, ns string, tag string) Element

	// Creates a text node.
	CreateTextNode(string) Text

	//Creates a new ProcessingInstruction object.
	CreateProcessingInstruction(target, data string) ProcessingInstruction

	// Assigns missing namespace prefixes, resolve prefix clashes etc.
	NormalizeNamespaces() error

	//	Adopt node from an external document.
	AdoptNode(Node) Node

	// Return the document type node
	GetDocumentType() DocumentType
}
