package dom

import (
	"encoding/xml"
)

type Name struct {
	xml.Name
	// Namespace prefix
	Prefix string
}

func (name *Name) QName() string {
	if len(name.Prefix) == 0 {
		return name.Local
	}
	return name.Prefix + ":" + name.Local
}
