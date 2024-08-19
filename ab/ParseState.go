package ab

type ParseState struct {
	Leaf         *Nodemapper
	Input        string
	That         string
	Topic        string
	ChatSession  *Chat
	Depth        int
	Vars         Predicates
	StarBindings *StarBindings
}

func NewParseState(depth int, chatSession *Chat, input, that, topic string, leaf *Nodemapper) *ParseState {
	return &ParseState{
		Leaf:         leaf,
		Input:        input,
		That:         that,
		Topic:        topic,
		ChatSession:  chatSession,
		Depth:        depth,
		Vars:         NewPredicates(),
		StarBindings: leaf.StarBindings,
	}
}
