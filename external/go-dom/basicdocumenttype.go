package dom

import (
	"unicode"
	"unicode/utf8"
)

type BasicDocumentType struct {
	basicNode

	name     string
	publicID string
	systemID string

	defn string
}

func (dt *BasicDocumentType) GetNodeType() NodeType { return DOCUMENT_TYPE_NODE }
func (dt *BasicDocumentType) GetName() string       { return dt.name }
func (dt *BasicDocumentType) GetPublicID() string   { return dt.publicID }
func (dt *BasicDocumentType) GetSystemID() string   { return dt.systemID }
func (dt *BasicDocumentType) GetDefinition() string { return dt.defn }

// ParseDocumentType parses a document type starting with <!DOCTYPE ...
// If the input is not a doctype, returns nil,false,nil
func ParseDocumentType(content []byte) (DocumentType, bool, error) {
	nextToken := func(in []byte) (string, []byte) {
		// Skip leading spaces.
		start := 0
		for width := 0; start < len(in); start += width {
			var r rune
			r, width = utf8.DecodeRune(in[start:])
			if !unicode.IsSpace(r) {
				break
			}
		}
		// Scan the token
		for width, i := 0, start; i < len(in); i += width {
			var r rune
			r, width = utf8.DecodeRune(in[i:])
			if r == '[' || r == ']' || r == '(' || r == ')' || r == '<' || r == '>' || r == ',' || r == '|' {
				if i == start {
					return string(in[start : i+width]), in[i+width:]
				}
				return string(in[start:i]), in[i:]
			}
			if unicode.IsSpace(r) {
				return string(in[start:i]), in[i:]
			}
		}
		// We are at the end
		return string(in), in[len(in):]
	}

	tok, rem := nextToken(content)
	if tok != "!DOCTYPE" {
		return nil, false, nil
	}
	ret := &BasicDocumentType{}
	tok, rem = nextToken(rem)
	ret.name = tok

	tok, rem = nextToken(rem)
	if tok == "PUBLIC" {
		tok, rem = nextToken(rem)
		ret.publicID = tok
		tok, rem = nextToken(rem)
		ret.systemID = tok
	} else if tok == "SYSTEM" {
		tok, rem = nextToken(rem)
		ret.systemID = tok
	} else if tok == "[" {
		ret.defn = "[" + string(rem)
	} else {
		return nil, true, ErrDOM{
			Typ: SYNTAX_ERR,
			Msg: "Document type syntax error",
		}
	}
	return ret, true, nil
}
