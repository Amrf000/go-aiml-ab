package ab

import (
	"github.com/subchen/go-xmldom"
)

type AIMLProcessorExtension interface {
	ExtensionTagSet() map[string]bool
	RecursEval(node *xmldom.Node, ps *ParseState) string
}
