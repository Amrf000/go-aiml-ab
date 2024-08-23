package dom

import (
	"encoding/xml"
)

type basicNamedNodeMap struct {
	attrs    []*BasicAttr
	mapAttrs map[xml.Name]*BasicAttr
}

func (m *basicNamedNodeMap) GetLength() int {
	return len(m.attrs)
}

// Returns a Attr identified by a namespace and related local name.
func (m *basicNamedNodeMap) GetNamedItemNS(uri string, name string) Attr {
	if m.mapAttrs == nil {
		return nil
	}
	item, ok := m.mapAttrs[xml.Name{Local: name, Space: uri}]
	if !ok {
		return nil
	}
	return item
}

// Returns the Attr at the given index, or null if the index is higher or equal to the number of nodes
func (m *basicNamedNodeMap) Item(index int) Attr {
	if index < 0 || index >= len(m.attrs) {
		return nil
	}
	return m.attrs[index]
}

// RemoveNamedItemNS removes the Attr identified by the given name
func (m *basicNamedNodeMap) RemoveNamedItemNS(uri string, name string) {
	if m.mapAttrs == nil {
		return
	}
	attr, exists := m.mapAttrs[xml.Name{Local: name, Space: uri}]
	if !exists {
		return
	}
	m.removeAttr(attr)
}

func (m *basicNamedNodeMap) removeAttr(attr Attr) {
	qname := attr.(*BasicAttr).name.Name
	delete(m.mapAttrs, qname)
	w := 0
	for k := range m.attrs {
		if m.attrs[k] != attr {
			m.attrs[w] = m.attrs[k]
			w++
		}
	}
	m.attrs = m.attrs[:w]
}

// Replaces, or adds, the Attr identified in the map by the given namespace and related local name.
func (m *basicNamedNodeMap) setNamedItemNS(owner Node, attr Attr) {
	if attr.GetOwnerElement() != nil && attr.GetOwnerElement() != owner {
		panic(ErrDOM{
			Typ: INUSE_ATTRIBUTE_ERR,
			Msg: "Attribute already in use",
			Op:  "SetNamedItem",
		})
	}
	if m.mapAttrs == nil {
		m.mapAttrs = make(map[xml.Name]*BasicAttr)
	}
	ba := attr.(*BasicAttr)
	qname := ba.name.Name
	existing := m.mapAttrs[qname]
	if existing != nil {
		if existing == attr {
			return
		}
		m.mapAttrs[qname] = ba
		for k := range m.attrs {
			if m.attrs[k] == existing {
				m.attrs[k] = ba
				ba.parent = owner
				return
			}
		}
	}
	m.mapAttrs[qname] = ba
	m.attrs = append(m.attrs, ba)
	ba.parent = owner
}

type BasicNamedNodeMap struct {
	owner *BasicElement
}

var _ NamedNodeMap = &BasicNamedNodeMap{}

func (b *BasicNamedNodeMap) GetLength() int { return b.owner.attributes.GetLength() }

// Returns the Attr at the given index, or null if the index is higher or equal to the number of nodes
func (b *BasicNamedNodeMap) Item(i int) Attr {
	return b.owner.attributes.Item(i)
}

// Returns a Attr identified by a namespace and related local name.
func (b *BasicNamedNodeMap) GetNamedItemNS(uri string, name string) Attr {
	return b.owner.attributes.GetNamedItemNS(uri, name)
}

// Replaces, or adds, the Attr identified in the map by the given namespace and related local name.
func (b *BasicNamedNodeMap) SetNamedItemNS(a Attr) {
	b.owner.attributes.setNamedItemNS(b.owner, a)
}

func (b *BasicNamedNodeMap) RemoveNamedItemNS(uri string, name string) {
	b.owner.attributes.RemoveNamedItemNS(uri, name)
}
