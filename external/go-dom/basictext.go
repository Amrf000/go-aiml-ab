package dom

type basicChardata struct {
	basicNode
	text string
}

func (cd *basicChardata) GetValue() string     { return cd.text }
func (cd *basicChardata) SetValue(text string) { cd.text = text }

func (cd *basicChardata) AppendChild(Node) Node {
	panic(ErrHierarchyRequest("AppendChild", "Invalid node type: character data node"))
}

func (cd *basicChardata) HasChildNodes() bool { return false }

func (cd *basicChardata) InsertBefore(newNode, referenceNode Node) Node {
	panic(ErrHierarchyRequest("InsertBefore", "Invalid node type: character data node"))
}

func (cs *basicChardata) RemoveChild(Node) {
	panic(ErrHierarchyRequest("RemoveChild", "Invalid node type: character data node"))
}

func (cs *basicChardata) Normalize() {}

type BasicText struct {
	basicChardata
}

var _ Text = &BasicText{}

// Returns "#text"
func (cd *BasicText) GetNodeName() string { return "#text" }

// Returns TEXT_NODE
func (cd *BasicText) GetNodeType() NodeType { return TEXT_NODE }

func (cd *BasicText) IsEqualNode(node Node) bool {
	n, ok := node.(*BasicText)
	if !ok {
		return false
	}
	return n.text == cd.text
}

// Returns a boolean value indicating whether or not the two nodes are
// the same (that is, they reference the same object).
func (cd *BasicText) IsSameNode(node Node) bool { return node == cd }

func (cd *BasicText) CloneNode(deep bool) Node {
	return cd.cloneNode(cd.ownerDocument, deep)
}

func (cd *BasicText) cloneNode(owner Document, deep bool) Node {
	return owner.CreateTextNode(cd.text)
}

type BasicComment struct {
	basicChardata
}

var _ Comment = &BasicComment{}

// Returns "#comment"
func (cd *BasicComment) GetNodeName() string { return "#comment" }

// Returns COMMENT_NODE
func (cd *BasicComment) GetNodeType() NodeType { return COMMENT_NODE }

func (cd *BasicComment) IsEqualNode(node Node) bool {
	n, ok := node.(*BasicComment)
	if !ok {
		return false
	}
	return n.text == cd.text
}

// Returns a boolean value indicating whether or not the two nodes are
// the same (that is, they reference the same object).
func (cd *BasicComment) IsSameNode(node Node) bool { return node == cd }

func (cd *BasicComment) CloneNode(deep bool) Node {
	return cd.cloneNode(cd.ownerDocument, deep)
}

func (cd *BasicComment) cloneNode(owner Document, deep bool) Node {
	return owner.CreateComment(cd.text)
}

type BasicProcessingInstruction struct {
	basicChardata
	target string
}

var _ ProcessingInstruction = &BasicProcessingInstruction{}

// Returns target
func (p *BasicProcessingInstruction) GetNodeName() string { return p.target }

// Returns PROCESSING_INSTRUCTION_NODE
func (p *BasicProcessingInstruction) GetNodeType() NodeType { return PROCESSING_INSTRUCTION_NODE }

func (p *BasicProcessingInstruction) IsEqualNode(node Node) bool {
	n, ok := node.(*BasicProcessingInstruction)
	if !ok {
		return false
	}
	return n.target == p.target && n.text == p.text
}

// Returns a boolean value indicating whether or not the two nodes are
// the same (that is, they reference the same object).
func (p *BasicProcessingInstruction) IsSameNode(node Node) bool { return node == p }

func (p *BasicProcessingInstruction) GetTarget() string { return p.target }

func (p *BasicProcessingInstruction) SetTarget(t string) { p.target = t }

func (p *BasicProcessingInstruction) CloneNode(deep bool) Node {
	return p.cloneNode(p.ownerDocument, deep)
}

func (p *BasicProcessingInstruction) cloneNode(owner Document, deep bool) Node {
	return owner.CreateProcessingInstruction(p.target, p.text)
}
