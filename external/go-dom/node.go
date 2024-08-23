package dom

const (
	xmlURL      = "http://www.w3.org/XML/1998/namespace"
	xmlnsURL    = "http://www.w3.org/2000/xmlns"
	xmlnsPrefix = "xmlns"
	xmlPrefix   = "xml"
)

type NodeType uint

const ELEMENT_NODE NodeType = 1
const ATTRIBUTE_NODE NodeType = 2
const TEXT_NODE NodeType = 3
const CDATA_SECTION_NODE NodeType = 4
const PROCESSING_INSTRUCTION_NODE NodeType = 7
const COMMENT_NODE NodeType = 8
const DOCUMENT_NODE NodeType = 9
const DOCUMENT_TYPE_NODE NodeType = 10
const DOCUMENT_FRAGMENT_NODE NodeType = 11

type Node interface {
	// Returns a live NodeList containing all the children of this node
	// (including elements, text and comments). NodeList being live means
	// that if the children of the Node change, the NodeList object is
	// automatically updated.
	GetChildNodes() NodeList

	// Returns a Node representing the first direct child node of the
	// node, or null if the node has no child.
	GetFirstChild() Node

	// Returns a Node representing the last direct child node of the node,
	// or null if the node has no child.
	GetLastChild() Node

	// Returns a Node representing the next node in the tree, or null if
	// there isn't such node.
	GetNextSibling() Node

	// Returns a String containing the name of the Node. The structure of
	// the name will differ with the node type. E.g. An HTMLElement will
	// contain the name of the corresponding tag, like 'audio' for an
	// HTMLAudioElement, a Text node will have the '#text' string, or a
	// Document node will have the '#document' string.
	GetNodeName() string

	// Returns an unsigned short representing the type of the
	// node.
	GetNodeType() NodeType

	// Returns the Document that this node belongs to. If the node is
	// itself a document, returns null.
	GetOwnerDocument() Document

	// Returns a Node that is the parent of this node. If there is no such
	// node, like if this node is the top of the tree or if doesn't
	// participate in a tree, this property returns null.
	GetParentNode() Node

	// Returns an Element that is the parent of this node. If the node has
	// no parent, or if that parent is not an Element, this property
	// returns null.
	GetParentElement() Element

	// Returns a Node representing the previous node in the tree, or null
	// if there isn't such node.
	GetPreviousSibling() Node

	// Adds the specified childNode argument as the last child to the
	// current node. If the argument referenced an existing node on the
	// DOM tree, the node will be detached from its current position and
	// attached at the new position.
	//
	// Returns a Node that is the appended child (aChild), except when aChild
	// is a DocumentFragment, in which case the empty DocumentFragment
	// is returned.
	AppendChild(Node) Node

	// Returns a boolean value indicating whether or not the element has
	// any child nodes.
	HasChildNodes() bool

	// Inserts a Node before the reference node as a child of a
	// specified parent node. Returns the added child (unless newNode is
	// a DocumentFragment, in which case the empty DocumentFragment is
	// returned).
	//
	// Algorithm: To ensure pre-insertion validity of a node into a parent
	// before a child, run these steps:
	//
	// If parent is not a Document, DocumentFragment, or Element node,
	// then throw a "HierarchyRequestError" DOMException.
	//
	// If node is a host-including inclusive ancestor of parent, then
	// throw a "HierarchyRequestError" DOMException.
	//
	// If child is non-null and its parent is not parent, then throw a
	// "NotFoundError" DOMException.
	//
	// If node is not a DocumentFragment, DocumentType, Element, or
	// CharacterData node, then throw a "HierarchyRequestError"
	// DOMException.
	//
	// If either node is a Text node and parent is a document, or node is
	// a doctype and parent is not a document, then throw a
	// "HierarchyRequestError" DOMException.
	//
	// If parent is a document, and any of the statements below, switched
	// on the interface node implements, are true, then throw a
	// "HierarchyRequestError" DOMException.
	//
	// DocumentFragment: If node has more than one element child or has a
	// Text node child. Otherwise, if node has one element child and
	// either parent has an element child, child is a doctype, or child is
	// non-null and a doctype is following child.
	//
	// Element: parent has an element child, child is a doctype, or child
	// is non-null and a doctype is following child.
	//
	// DocumentType: parent has a doctype child, child is non-null and an
	// element is preceding child, or child is null and parent has an
	// element child.
	//
	// To pre-insert a node into a parent before a child, run these steps:

	// Ensure pre-insertion validity of node into parent before child.

	// Let referenceChild be child.

	// If referenceChild is node, then set referenceChild to node’s next sibling.

	// Insert newNode into node before referenceChild.

	// Return node.
	InsertBefore(newNode, referenceNode Node) Node

	// Removes a child node from the current element, which must be a
	// child of the current node.
	RemoveChild(Node)

	// Clone a Node, and optionally, all of its contents.
	//
	// Returns the new Node cloned. The cloned node has no parent and is not
	// part of the document, until it is added to another node that is
	// part of the document, using Node.appendChild() or a similar
	// method.
	CloneNode(deep bool) Node

	// Returns true or false value indicating whether or not a node is a
	// descendant of the calling node.
	Contains(Node) bool

	// Returns the object's root
	GetRootNode() Node

	// Accepts a namespace URI as an argument and returns a boolean value
	// with a value of true if the namespace is the default namespace on
	// the given node or false if not.
	IsDefaultNamespace(uri string) bool

	// Returns a boolean value which indicates whether or not two nodes
	// are of the same type and all their defining data points match.
	IsEqualNode(Node) bool

	// Returns a boolean value indicating whether or not the two nodes are
	// the same (that is, they reference the same object).
	IsSameNode(Node) bool

	// Returns a string  containing the prefix for a given namespace
	// URI, if present, and "" if not. When multiple prefixes are
	// possible, the result is implementation-dependent.
	LookupPrefix(string) string

	// Accepts a prefix and returns the namespace URI associated with it
	// on the given node if found (and "" if not).
	LookupNamespaceURI(string) string

	// Clean up all the text nodes under this element (merge adjacent,
	// remove empty).
	Normalize()

	// // Replaces one child Node of the current one with the second one
	// // given in parameter.
	// ReplaceChild(newChild, oldChild Node)

	treeNode() *tnode
	cloneNode(owner Document, deep bool) Node
}

type NodeList interface {
	GetLength() int
	Item(int) Node
}
