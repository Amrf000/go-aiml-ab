package dom

import (
	"bufio"
	"io"
	"unicode/utf8"
)

func Encode(node Node, writer io.Writer) error {
	out := bufio.NewWriter(writer)
	defer out.Flush()
	return encodeNode(node, out)
}

var (
	escQuot = []byte("&#34;") // shorter than "&quot;"
	escApos = []byte("&#39;") // shorter than "&apos;"
	escAmp  = []byte("&amp;")
	escLT   = []byte("&lt;")
	escGT   = []byte("&gt;")
	escTab  = []byte("&#x9;")
	escNL   = []byte("&#xA;")
	escCR   = []byte("&#xD;")
	escFFFD = []byte("\uFFFD") // Unicode replacement character
)

// Decide whether the given rune is in the XML Character Range, per
// the Char production of https://www.xml.com/axml/testaxml.htm,
// Section 2.2 Characters.
func isInCharacterRange(r rune) (inrange bool) {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xD7FF ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}

// escapeText writes to w the properly escaped XML equivalent
// of the plain text data s. If escapeNewline is true, newline
// characters will be escaped.
func escapeText(w io.Writer, s []byte, escapeNewline bool) error {
	var esc []byte
	last := 0
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRune(s[i:])
		i += width
		switch r {
		case '"':
			esc = escQuot
		case '\'':
			esc = escApos
		case '&':
			esc = escAmp
		case '<':
			esc = escLT
		case '>':
			esc = escGT
		case '\t':
			esc = escTab
		case '\n':
			if !escapeNewline {
				continue
			}
			esc = escNL
		case '\r':
			esc = escCR
		default:
			if !isInCharacterRange(r) || (r == 0xFFFD && width == 1) {
				esc = escFFFD
				break
			}
			continue
		}
		if _, err := w.Write(s[last : i-width]); err != nil {
			return err
		}
		if _, err := w.Write(esc); err != nil {
			return err
		}
		last = i
	}
	_, err := w.Write(s[last:])
	return err
}

func encodeNode(node Node, out *bufio.Writer) error {
	space := func() error {
		_, err := out.WriteRune(' ')
		return err
	}
	writeName := func(n Name) error {
		if len(n.Prefix) > 0 {
			if _, err := out.WriteString(n.Prefix); err != nil {
				return err
			}
			if _, err := out.WriteRune(':'); err != nil {
				return err
			}
		}
		_, err := out.WriteString(n.Local)
		return err
	}
	writeAttrValue := func(value string) error {
		if _, err := out.WriteRune('"'); err != nil {
			return err
		}
		if err := escapeText(out, []byte(value), true); err != nil {
			return err
		}
		_, err := out.WriteRune('"')
		return err
	}
	writeCharData := func(value string) error {
		return escapeText(out, []byte(value), false)
	}
	switch ch := node.(type) {
	case *BasicDocument:
		for c := ch.GetFirstChild(); c != nil; c = c.GetNextSibling() {
			if err := encodeNode(c, out); err != nil {
				return err
			}
		}

	case *BasicElement:
		if _, err := out.WriteRune('<'); err != nil {
			return err
		}
		if err := writeName(ch.GetQName()); err != nil {
			return err
		}
		attrs := ch.GetAttributes()
		for i := 0; i < attrs.GetLength(); i++ {
			if err := space(); err != nil {
				return err
			}
			attr := attrs.Item(i)
			if err := writeName(attr.GetQName()); err != nil {
				return err
			}
			if _, err := out.WriteRune('='); err != nil {
				return err
			}
			writeAttrValue(attr.GetValue())
		}
		if _, err := out.WriteRune('>'); err != nil {
			return err
		}

		for c := ch.GetFirstChild(); c != nil; c = c.GetNextSibling() {
			if err := encodeNode(c, out); err != nil {
				return err
			}
		}
		if _, err := out.WriteString("</"); err != nil {
			return err
		}
		if err := writeName(ch.GetQName()); err != nil {
			return err
		}
		if _, err := out.WriteRune('>'); err != nil {
			return err
		}

	case *BasicComment:
		if _, err := out.WriteString("<!--"); err != nil {
			return err
		}
		if err := writeCharData(ch.GetValue()); err != nil {
			return err
		}
		if _, err := out.WriteString("-->"); err != nil {
			return err
		}

	case *BasicText:
		if err := writeCharData(ch.GetValue()); err != nil {
			return err
		}

	case *BasicProcessingInstruction:
		if _, err := out.WriteString("<?"); err != nil {
			return err
		}
		if _, err := out.WriteString(ch.GetTarget()); err != nil {
			return err
		}
		if err := space(); err != nil {
			return err
		}
		if _, err := out.WriteString(ch.GetValue()); err != nil {
			return err
		}
		if _, err := out.WriteString("?>"); err != nil {
			return err
		}

	}
	return nil
}
