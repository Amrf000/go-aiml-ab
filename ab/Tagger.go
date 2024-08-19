package ab

// ViterbiNode represents a node in the Viterbi lattice.
type ViterbiNode struct {
	start      int
	length     int
	posID      int
	cost       int
	morphemeID int
	prev       *ViterbiNode
}

// ViterbiNodeList is a list of ViterbiNode.
type ViterbiNodeList []*ViterbiNode

// BOSNodes represents the beginning of sentence node list.
var BOSNodes = ViterbiNodeList{makeBOSEOS()}

// Morpheme represents a morpheme in the text.
type Morpheme struct {
	Surface    string
	Feature    string
	Start      int
	MorphemeID int
}

// makeBOSEOS creates a BOS/EOS node.
func makeBOSEOS() *ViterbiNode {
	// Implementation of makeBOSEOS.
	return &ViterbiNode{}
}

// PartsOfSpeech is a mock for parts of speech lookup.
func PartsOfSpeech(posID int) string {
	// Mock implementation. Replace with actual logic.
	return ""
}

// WordDic and Unknown are mocks for dictionary lookups.
var WordDic, Unknown = &WordDicType{}, &UnknownType{}

// WordDicType represents the dictionary lookup.
type WordDicType struct{}

// Search is a mock for dictionary search.
func (w *WordDicType) Search(text string, i int, callback *MakeLattice) {
	// Mock implementation. Replace with actual logic.
}

// UnknownType represents the unknown word lookup.
type UnknownType struct{}

// Search is a mock for unknown word search.
func (u *UnknownType) Search(text string, i int, callback *MakeLattice) {
	// Mock implementation. Replace with actual logic.
}

// Matrix is a mock for calculating link costs.
var Matrix = &MatrixType{}

// MatrixType represents the matrix used for link cost calculation.
type MatrixType struct{}

// LinkCost is a mock for calculating the link cost between two nodes.
func (m *MatrixType) LinkCost(posID1, posID2 int) int {
	// Mock implementation. Replace with actual logic.
	return 0
}

// Parse parses the input text into a list of morphemes.
func Parse(text string) []*Morpheme {
	return ParseWithResult(text, make([]*Morpheme, 0, len(text)/2))
}

// ParseWithResult parses the input text into a given list of morphemes.
func ParseWithResult(text string, result []*Morpheme) []*Morpheme {
	for vn := parseImpl(text); vn != nil; vn = vn.prev {
		surface := text[vn.start : vn.start+vn.length]
		feature := PartsOfSpeech(vn.posID)
		result = append(result, &Morpheme{Surface: surface, Feature: feature, Start: vn.start, MorphemeID: vn.morphemeID})
	}
	return result
}

// Wakati splits the input text into a list of words.
func Wakati(text string) []string {
	return WakatiWithResult(text, make([]string, 0, len(text)/1))
}

// WakatiWithResult splits the input text into a given list of words.
func WakatiWithResult(text string, result []string) []string {
	for vn := parseImpl(text); vn != nil; vn = vn.prev {
		result = append(result, text[vn.start:vn.start+vn.length])
	}
	return result
}

// parseImpl implements the Viterbi algorithm to parse the text.
func parseImpl(text string) *ViterbiNode {
	len := len(text)
	nodesAry := make([]ViterbiNodeList, len+1)
	nodesAry[0] = BOSNodes

	fn := &MakeLattice{nodesAry: nodesAry}
	for i := 0; i < len; i++ {
		if nodesAry[i] != nil {
			fn.set(i)
			WordDic.Search(text, i, fn)
			Unknown.Search(text, i, fn)
		}
	}

	cur := setMincostNode(makeBOSEOS(), nodesAry[len]).prev
	var head *ViterbiNode
	for cur.prev != nil {
		tmp := cur.prev
		cur.prev = head
		head = cur
		cur = tmp
	}
	return head
}

// setMincostNode finds the node with the minimum cost and sets it as the previous node.
func setMincostNode(vn *ViterbiNode, prevs ViterbiNodeList) *ViterbiNode {
	f := prevs[0]
	vn.prev = f
	minCost := f.cost + Matrix.LinkCost(f.posID, vn.posID)

	for i := 1; i < len(prevs); i++ {
		p := prevs[i]
		cost := p.cost + Matrix.LinkCost(p.posID, vn.posID)

		if cost < minCost {
			minCost = cost
			vn.prev = p
		}
	}
	vn.cost += minCost

	return vn
}

// MakeLattice is a struct that helps build the Viterbi lattice.
type MakeLattice struct {
	nodesAry []ViterbiNodeList
	i        int
	prevs    ViterbiNodeList
	empty    bool
}

// set sets up the lattice at a given index.
func (fn *MakeLattice) set(i int) {
	fn.i = i
	fn.prevs = fn.nodesAry[i]
	fn.nodesAry[i] = nil
	fn.empty = true
}

// Call is the callback function used during dictionary search.
func (fn *MakeLattice) Call(vn *ViterbiNode, isSpace bool) {
	fn.empty = false

	end := fn.i + vn.length
	if fn.nodesAry[end] == nil {
		fn.nodesAry[end] = ViterbiNodeList{}
	}
	ends := fn.nodesAry[end]

	if isSpace {
		ends = append(ends, fn.prevs...)
	} else {
		ends = append(ends, setMincostNode(vn, fn.prevs))
	}
}

// IsEmpty checks if the lattice is empty.
func (fn *MakeLattice) IsEmpty() bool {
	return fn.empty
}
