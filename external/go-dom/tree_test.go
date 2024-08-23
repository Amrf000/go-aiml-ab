package dom

import (
	"testing"
)

func TestInsertChildAfter(t *testing.T) {
	doc := &BasicDocument{}
	parent := &basicNode{
		ownerDocument: doc,
	}

	if parent.firstChild() != nil {
		t.Fail()
	}
	if parent.lastChild() != nil {
		t.Fail()
	}

	// Insert the first node
	newChild1 := &basicNode{
		ownerDocument: doc,
	}
	insertChildAfter(parent, newChild1, nil)
	if parent.firstChild() != newChild1 {
		t.Fail()
	}
	if parent.lastChild() != newChild1 {
		t.Fail()
	}
	if newChild1.nextSibling() != nil {
		t.Fail()
	}
	if newChild1.prevSibling() != nil {
		t.Fail()
	}

	// Insert second node
	newChild2 := &basicNode{
		ownerDocument: doc,
	}
	insertChildAfter(parent, newChild2, nil)
	if parent.firstChild() != newChild2 {
		t.Fail()
	}
	if parent.lastChild() != newChild1 {
		t.Fail()
	}
	if newChild1.nextSibling() != nil {
		t.Fail()
	}
	if newChild1.prevSibling() != newChild2 {
		t.Fail()
	}
	if newChild2.nextSibling() != newChild1 {
		t.Fail()
	}
	if newChild2.prevSibling() != nil {
		t.Fail()
	}

	// Insert third node
	newChild3 := &basicNode{
		ownerDocument: doc,
	}
	insertChildAfter(parent, newChild3, nil)
	if parent.firstChild() != newChild3 {
		t.Fail()
	}
	if parent.lastChild() != newChild1 {
		t.Fail()
	}
	if newChild1.nextSibling() != nil {
		t.Fail()
	}
	if newChild1.prevSibling() != newChild2 {
		t.Fail()
	}
	if newChild2.nextSibling() != newChild1 {
		t.Fail()
	}
	if newChild2.prevSibling() != newChild3 {
		t.Fail()
	}
	if newChild3.nextSibling() != newChild2 {
		t.Fail()
	}
	if newChild3.prevSibling() != nil {
		t.Fail()
	}
}

func TestInsertChildBefore(t *testing.T) {
	doc := &BasicDocument{}
	parent := &basicNode{
		ownerDocument: doc,
	}

	if parent.firstChild() != nil {
		t.Fail()
	}
	if parent.lastChild() != nil {
		t.Fail()
	}

	// Insert the first node
	newChild1 := &basicNode{
		ownerDocument: doc,
	}
	insertChildBefore(parent, newChild1, nil)
	if parent.firstChild() != newChild1 {
		t.Fail()
	}
	if parent.lastChild() != newChild1 {
		t.Fail()
	}
	if newChild1.nextSibling() != nil {
		t.Fail()
	}
	if newChild1.prevSibling() != nil {
		t.Fail()
	}

	// Insert second node
	newChild2 := &basicNode{
		ownerDocument: doc,
	}
	insertChildBefore(parent, newChild2, nil)
	if parent.firstChild() != newChild1 {
		t.Fail()
	}
	if parent.lastChild() != newChild2 {
		t.Fail()
	}
	if newChild1.nextSibling() != newChild2 {
		t.Fail()
	}
	if newChild1.prevSibling() != nil {
		t.Fail()
	}
	if newChild2.nextSibling() != nil {
		t.Fail()
	}
	if newChild2.prevSibling() != newChild1 {
		t.Fail()
	}

	// Insert third node
	newChild3 := &basicNode{
		ownerDocument: doc,
	}
	insertChildBefore(parent, newChild3, nil)
	if parent.firstChild() != newChild1 {
		t.Fail()
	}
	if parent.lastChild() != newChild3 {
		t.Fail()
	}
	if newChild1.nextSibling() != newChild2 {
		t.Fail()
	}
	if newChild1.prevSibling() != nil {
		t.Fail()
	}
	if newChild2.nextSibling() != newChild3 {
		t.Fail()
	}
	if newChild2.prevSibling() != newChild1 {
		t.Fail()
	}
	if newChild3.nextSibling() != nil {
		t.Fail()
	}
	if newChild3.prevSibling() != newChild2 {
		t.Fail()
	}
}

func TestInsertMid(t *testing.T) {
	doc := &BasicDocument{}
	parent := &basicNode{
		ownerDocument: doc,
	}
	newChild1 := &basicNode{
		ownerDocument: doc,
	}
	newChild2 := &basicNode{
		ownerDocument: doc,
	}
	newChild3 := &basicNode{
		ownerDocument: doc,
	}

	if parent.firstChild() != nil {
		t.Fail()
	}
	if parent.lastChild() != nil {
		t.Fail()
	}

	insertChildBefore(parent, newChild1, nil)
	insertChildBefore(parent, newChild3, nil)
	insertChildBefore(parent, newChild2, newChild3)

	if parent.firstChild() != newChild1 {
		t.Fail()
	}
	if parent.lastChild() != newChild3 {
		t.Fail()
	}
	if newChild1.nextSibling() != newChild2 {
		t.Fail()
	}
	if newChild1.prevSibling() != nil {
		t.Fail()
	}
	if newChild2.nextSibling() != newChild3 {
		t.Fail()
	}
	if newChild2.prevSibling() != newChild1 {
		t.Fail()
	}
	if newChild3.nextSibling() != nil {
		t.Fail()
	}
	if newChild3.prevSibling() != newChild2 {
		t.Fail()
	}
}
func TestInserAftertMid(t *testing.T) {
	doc := &BasicDocument{}
	parent := &basicNode{
		ownerDocument: doc,
	}
	newChild1 := &basicNode{
		ownerDocument: doc,
	}
	newChild2 := &basicNode{
		ownerDocument: doc,
	}
	newChild3 := &basicNode{
		ownerDocument: doc,
	}

	if parent.firstChild() != nil {
		t.Fail()
	}
	if parent.lastChild() != nil {
		t.Fail()
	}

	insertChildAfter(parent, newChild3, nil)
	insertChildAfter(parent, newChild1, nil)
	insertChildAfter(parent, newChild2, newChild1)

	if parent.firstChild() != newChild1 {
		t.Fail()
	}
	if parent.lastChild() != newChild3 {
		t.Fail()
	}
	if newChild1.nextSibling() != newChild2 {
		t.Fail()
	}
	if newChild1.prevSibling() != nil {
		t.Fail()
	}
	if newChild2.nextSibling() != newChild3 {
		t.Fail()
	}
	if newChild2.prevSibling() != newChild1 {
		t.Fail()
	}
	if newChild3.nextSibling() != nil {
		t.Fail()
	}
	if newChild3.prevSibling() != newChild2 {
		t.Fail()
	}
}

func TestInsertFirst(t *testing.T) {
	doc := &BasicDocument{}
	parent := &basicNode{
		ownerDocument: doc,
	}
	newChild1 := &basicNode{
		ownerDocument: doc,
	}
	newChild2 := &basicNode{
		ownerDocument: doc,
	}
	newChild3 := &basicNode{
		ownerDocument: doc,
	}

	if parent.firstChild() != nil {
		t.Fail()
	}
	if parent.lastChild() != nil {
		t.Fail()
	}

	insertChildBefore(parent, newChild2, nil)
	insertChildBefore(parent, newChild3, nil)
	insertChildBefore(parent, newChild1, newChild2)

	if parent.firstChild() != newChild1 {
		t.Fail()
	}
	if parent.lastChild() != newChild3 {
		t.Fail()
	}
	if newChild1.nextSibling() != newChild2 {
		t.Fail()
	}
	if newChild1.prevSibling() != nil {
		t.Fail()
	}
	if newChild2.nextSibling() != newChild3 {
		t.Fail()
	}
	if newChild2.prevSibling() != newChild1 {
		t.Fail()
	}
	if newChild3.nextSibling() != nil {
		t.Fail()
	}
	if newChild3.prevSibling() != newChild2 {
		t.Fail()
	}
}

func TestDetach(t *testing.T) {
	doc := &BasicDocument{}
	parent := &basicNode{
		ownerDocument: doc,
	}
	newChild1 := &basicNode{
		ownerDocument: doc,
	}
	newChild2 := &basicNode{
		ownerDocument: doc,
	}
	newChild3 := &basicNode{
		ownerDocument: doc,
	}

	insertChildBefore(parent, newChild2, nil)
	insertChildBefore(parent, newChild3, nil)
	insertChildBefore(parent, newChild1, newChild2)

	detachChild(parent, newChild1)
	if parent.firstChild() != newChild2 {
		t.Fail()
	}
	if parent.lastChild() != newChild3 {
		t.Fail()
	}
	if newChild2.nextSibling() != newChild3 {
		t.Fail()
	}
	if newChild2.prevSibling() != nil {
		t.Fail()
	}
	if newChild3.nextSibling() != nil {
		t.Fail()
	}
	if newChild3.prevSibling() != newChild2 {
		t.Fail()
	}
	detachChild(parent, newChild2)
	if parent.firstChild() != newChild3 {
		t.Fail()
	}
	if parent.lastChild() != newChild3 {
		t.Fail()
	}
	if newChild3.nextSibling() != nil {
		t.Fail()
	}
	if newChild3.prevSibling() != nil {
		t.Fail()
	}
	detachChild(parent, newChild3)
	if parent.firstChild() != nil {
		t.Fail()
	}
	if parent.lastChild() != nil {
		t.Fail()
	}
}

func TestDetach2(t *testing.T) {
	doc := &BasicDocument{}
	parent := &basicNode{
		ownerDocument: doc,
	}
	newChild1 := &basicNode{
		ownerDocument: doc,
	}
	newChild2 := &basicNode{
		ownerDocument: doc,
	}
	newChild3 := &basicNode{
		ownerDocument: doc,
	}

	insertChildBefore(parent, newChild2, nil)
	insertChildBefore(parent, newChild3, nil)
	insertChildBefore(parent, newChild1, newChild2)

	detachChild(parent, newChild2)
	if parent.firstChild() != newChild1 {
		t.Fail()
	}
	if parent.lastChild() != newChild3 {
		t.Fail()
	}
	if newChild1.nextSibling() != newChild3 {
		t.Fail()
	}
	if newChild1.prevSibling() != nil {
		t.Fail()
	}
	if newChild3.nextSibling() != nil {
		t.Fail()
	}
	if newChild3.prevSibling() != newChild1 {
		t.Fail()
	}
}
