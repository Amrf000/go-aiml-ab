package dom

type tnode struct {
	parent Node
	next   Node
	prev   Node
	child  Node
	// ver is used by nodelists to keep track of child list changes
	ver int
}

func (node *tnode) firstChild() Node {
	return node.child
}

func (node *tnode) lastChild() Node {
	if node.child == nil {
		return nil
	}
	return node.child.treeNode().prev
}

func (node *tnode) prevSibling() Node {
	if node.parent == nil {
		return nil
	}
	if node.parent.treeNode().firstChild().treeNode() == node {
		return nil
	}
	return node.prev
}

func (node *tnode) nextSibling() Node {
	if node.parent == nil {
		return nil
	}
	if node.parent.treeNode().child == node.next {
		return nil
	}
	return node.next
}

// Insert child after given node. If after is nil, insert as first node
func insertChildAfter(parent, newChild, after Node) {
	newChildtn := newChild.treeNode()
	newChildtn.parent = parent
	parenttn := parent.treeNode()
	parenttn.ver++
	if after == nil {
		first := parenttn.firstChild()
		// newChild is the new first node
		parenttn.child = newChild
		if first == nil {
			newChildtn.next = newChild
			newChildtn.prev = newChild
			return
		}
		newChildtn.next = first
		ftn := first.treeNode()
		newChildtn.prev = ftn.prev
		ftn.prev.treeNode().next = newChild
		ftn.prev = newChild
		return
	}
	newChildtn.prev = after
	atn := after.treeNode()
	newChildtn.next = atn.next
	atn.next.treeNode().prev = newChild
	atn.next = newChild
}

// Insert child before given node. If before is nil, insert as last node
func insertChildBefore(parent, newChild, before Node) {
	newChildtn := newChild.treeNode()
	newChildtn.parent = parent
	parenttn := parent.treeNode()
	parenttn.ver++
	if before == nil {
		last := parenttn.lastChild()
		// newChild is the last node
		if last == nil {
			newChildtn.next = newChild
			newChildtn.prev = newChild
			parenttn.child = newChild
			return
		}
		newChildtn.prev = last
		ltn := last.treeNode()
		newChildtn.next = ltn.next
		ltn.next.treeNode().prev = newChild
		ltn.next = newChild
		return
	}
	if before == parenttn.child {
		parenttn.child = newChild
	}
	newChildtn.next = before
	btn := before.treeNode()
	newChildtn.prev = btn.prev
	btn.prev.treeNode().next = newChild
	btn.prev = newChild
}

func detachChild(parent, child Node) {
	childtn := child.treeNode()
	if parent != nil {
		parenttn := parent.treeNode()
		parenttn.ver++
		if parenttn.child == child {
			parenttn.child = childtn.next
			if parenttn.child == child {
				parenttn.child = nil
			}
		}
	}
	if childtn.next != nil {
		childtn.next.treeNode().prev = childtn.prev
	}
	if childtn.prev != nil {
		childtn.prev.treeNode().next = childtn.next
	}
	childtn.next = nil
	childtn.prev = nil
}
