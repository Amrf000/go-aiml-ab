package ab

import "aiml/external/go-dom"

type AIMLProcessorExtension interface {
	ExtensionTagSet() map[string]bool
	RecursEval(node dom.Node, ps *ParseState) string
}
