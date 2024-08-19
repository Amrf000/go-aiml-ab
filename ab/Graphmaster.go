package ab

import (
	"fmt"
	"golang.org/x/exp/slices"
	"regexp"
	"strconv"
	"strings"
)

type Graphmaster struct {
	Bot         *Bot
	Name        string
	Root        *Nodemapper
	MatchCount  int
	UpgradeCnt  int
	Vocabulary  map[string]bool
	ResultNote  string
	CategoryCnt int
}

var EnableShortCuts = false

func NewGraphmaster(bot *Bot) *Graphmaster {
	return NewGraphmasterWithName(bot, "brain")
}

func NewGraphmasterWithName(bot *Bot, name string) *Graphmaster {
	this := &Graphmaster{}
	this.Root = NewNodemapper()
	this.Bot = bot
	this.Name = name
	this.Vocabulary = make(map[string]bool)
	return this
}

func InputThatTopic(input, that, topic string) string {
	return strings.TrimSpace(input) + " <THAT> " + strings.TrimSpace(that) + " <TOPIC> " + strings.TrimSpace(topic)
}

var botPropRegex = regexp.MustCompile(`(?i)<bot name="(.*?)"/>`)

func (g *Graphmaster) ReplaceBotProperties(pattern string) string {
	if strings.Contains(pattern, "<B") {
		matches := botPropRegex.FindAllStringSubmatch(pattern, -1)
		for _, match := range matches {
			propName := strings.ToLower(match[1])
			property := strings.ToUpper(g.Bot.Properties.Get(propName))
			pattern = strings.Replace(pattern, match[0], property, 1)
		}
	}
	return pattern
}

func (g *Graphmaster) AddCategory(category *Category) {
	inputThatTopic := InputThatTopic(category.GetPattern(), category.GetThat(), category.GetTopic())
	fmt.Println("addCategory: " + inputThatTopic)
	inputThatTopic = g.ReplaceBotProperties(inputThatTopic)
	p := SentenceToPath(inputThatTopic)
	g.AddPath(p, category)
	g.CategoryCnt++
}

func (g *Graphmaster) AddPath(path *Paths, category *Category) {
	g.AddPathRecursive(g.Root, path, category)
}

func (g *Graphmaster) AddPathRecursive(node *Nodemapper, path *Paths, category *Category) {
	if path == nil {
		node.Category = category
		node.Height = 0
	} else if EnableShortCuts && g.ThatStarTopicStar(path) {
		node.Category = category
		node.Height = min(4, node.Height)
		node.ShortCut = true
	} else if ContainsKey(node, path.Word) {
		if strings.HasPrefix(path.Word, "<SET>") {
			g.AddSets(path.Word, g.Bot, node, category.GetFilename())
		}
		nextNode := Get(node, path.Word)
		g.AddPathRecursive(nextNode, path.Next, category)
		offset := 1
		if path.Word == "#" || path.Word == "^" {
			offset = 0
		}
		node.Height = min(offset+nextNode.Height, node.Height)
	} else {
		nextNode := NewNodemapper()
		if strings.HasPrefix(path.Word, "<SET>") {
			g.AddSets(path.Word, g.Bot, node, category.GetFilename())
		}
		//node.MapData[path.Word] = nextNode
		if node.Key != "" {
			Upgrade(node)
			g.UpgradeCnt++
		}
		Put(node, path.Word, nextNode)
		g.AddPathRecursive(nextNode, path.Next, category)
		offset := 1
		if path.Word == "#" || path.Word == "^" {
			offset = 0
		}
		node.Height = min(offset+nextNode.Height, node.Height)
	}
}

func (g *Graphmaster) AddSets(setType string, bot *Bot, node *Nodemapper, filename string) {
	setName := strings.ToLower(TagTrim(setType, "SET"))
	if _, exists := bot.SetMap[setName]; exists {
		if node.Sets == nil {
			node.Sets = []string{}
		}
		if !slices.Contains(node.Sets, setName) {
			node.Sets = append(node.Sets, setName)
		}
	} else {
		panic(fmt.Errorf("No AIML Set found for <set>%s</set> in %s %s\n", setName, bot.Name, filename))
	}
}

func (g *Graphmaster) ThatStarTopicStar(path *Paths) bool {
	tail := PathToSentence(path)
	return strings.TrimSpace(tail) == "<THAT> * <TOPIC> *"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func (g *Graphmaster) ExistsCategory(c *Category) bool {
	return (g.FindNode(c) != nil)
}
func (g *Graphmaster) FindNode(c *Category) *Nodemapper {
	return g.FindNodeWithTopic(c.GetPattern(), c.GetThat(), c.GetTopic())
}

var verbose = false

func (g *Graphmaster) FindNodeWithTopic(input, that, topic string) *Nodemapper {
	result := g.FindNodeWithPath(g.Root, SentenceToPath(InputThatTopic(input, that, topic)))
	if verbose {
		fmt.Printf("findNode %s %v\n", InputThatTopic(input, that, topic), result)
	}
	return result
}

func (g *Graphmaster) FindNodeWithPath(node *Nodemapper, path *Paths) *Nodemapper {
	if path == nil && node != nil {
		if verbose {
			fmt.Println("findNode: path is null, returning node " + node.Category.InputThatTopic())
		}
		return node
	} else if strings.TrimSpace(PathToSentence(path)) == "<THAT> * <TOPIC> *" && node.ShortCut && path.Word == "<THAT>" {
		if verbose {
			fmt.Println("findNode: shortcut, returning " + node.Category.InputThatTopic())
		}
		return node
	} else if ContainsKey(node, path.Word) {
		if verbose {
			fmt.Println("findNode: node contains " + path.Word)
		}
		nextNode := Get(node, strings.ToUpper(path.Word))
		return g.FindNodeWithPath(nextNode, path.Next)
	} else {
		if verbose {
			fmt.Println("findNode: returning null")
		}
		return nil
	}
}

func (g *Graphmaster) MatchRaw(input, that, topic string) *Nodemapper {
	var n *Nodemapper = nil

	inputThatTopic := InputThatTopic(input, that, topic)
	fmt.Println("Matching: " + inputThatTopic)
	p := SentenceToPath(inputThatTopic)
	p.Print()
	n = g.MatchWithTopic(p, inputThatTopic)
	if TraceMode {
		if n != nil {
			//MagicBooleans.trace("in graphmaster.match(), matched "+n.category.inputThatTopic()+" "+n.category.getFilename());
			if TraceMode {
				fmt.Println("Matched: " + n.Category.InputThatTopic() + " " + n.Category.GetFilename())
			}
		} else {
			//MagicBooleans.trace("in graphmaster.match(), no match.");
			if TraceMode {
				fmt.Println("No match.")
			}
		}

	}

	if TraceMode && len(MatchTrace) < MaxTraceLength {
		if n != nil {
			SetMatchTrace(MatchTrace + n.Category.InputThatTopic() + "\n")
		}
	}
	//MagicBooleans.trace("in graphmaster.match(), returning: " + n);
	return n
}

func (g *Graphmaster) MatchWithTopic(path *Paths, inputThatTopic string) *Nodemapper {
	inputStars := make([]string, MaxStars)
	thatStars := make([]string, MaxStars)
	topicStars := make([]string, MaxStars)
	starState := "inputStar"
	MatchTrace = ""
	n := g.Match(path, g.Root, inputThatTopic, starState, 0, inputStars, thatStars, topicStars, MatchTrace)
	if n != nil {
		sb := NewStarBindings()
		for i := 0; inputStars[i] != "" && i < MaxStars; i++ {
			sb.InputStars.Add(inputStars[i])
		}
		for i := 0; thatStars[i] != "" && i < MaxStars; i++ {
			sb.ThatStars.Add(thatStars[i])
		}
		for i := 0; topicStars[i] != "" && i < MaxStars; i++ {
			sb.TopicStars.Add(topicStars[i])
		}
		n.StarBindings = sb
	}
	//if (!n.category.getPattern().contains("*")) System.out.println("adding match "+inputThatTopic);
	if n != nil {
		n.Category.AddMatch(inputThatTopic, g.Bot)
	}
	return n
}

func (g *Graphmaster) Match(path *Paths, node *Nodemapper, inputThatTopic, starState string, starIndex int,
	inputStars, thatStars, topicStars []string, matchTrace string) *Nodemapper {
	var matchedNode *Nodemapper
	if path != nil {
		fmt.Println("Match: Height=" + strconv.Itoa(node.Height) + " Length=" + strconv.Itoa(path.Length) + " Path=" + PathToSentence(path))
	}
	g.MatchCount++
	if matchedNode = g.NullMatch(path, node, matchTrace); matchedNode != nil {
		return matchedNode
	} else if path.Length < node.Height {
		return nil
	} else if matchedNode = g.DollarMatch(path, node, inputThatTopic, starState, starIndex, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
		return matchedNode
	} else if matchedNode = g.SharpMatch(path, node, inputThatTopic, starState, starIndex, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
		return matchedNode
	} else if matchedNode = g.UnderMatch(path, node, inputThatTopic, starState, starIndex, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
		return matchedNode
	} else if matchedNode = g.WordMatch(path, node, inputThatTopic, starState, starIndex, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
		return matchedNode
	} else if matchedNode = g.SetMatch(path, node, inputThatTopic, starState, starIndex, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
		return matchedNode
	} else if matchedNode = g.ShortCutMatch(path, node, inputThatTopic, starState, starIndex, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
		return matchedNode
	} else if matchedNode = g.CaretMatch(path, node, inputThatTopic, starState, starIndex, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
		return matchedNode
	} else if matchedNode = g.StarMatch(path, node, inputThatTopic, starState, starIndex, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
		return matchedNode
	} else {
		return nil
	}
}

func (g *Graphmaster) Fail(mode, trace string) {
	fmt.Println("Match failed (" + mode + ") " + trace)
}

func (g *Graphmaster) NullMatch(path *Paths, node *Nodemapper, matchTrace string) *Nodemapper {
	if path == nil && node != nil && IsLeaf(node) && node.Category != nil {
		return node
	} else {
		g.Fail("null", matchTrace)
		return nil
	}
}

func (g *Graphmaster) ShortCutMatch(path *Paths, node *Nodemapper, inputThatTopic, starState string,
	starIndex int, inputStars []string, thatStars, topicStars []string, matchTrace string) *Nodemapper {
	if node != nil && node.ShortCut && path.Word == "<THAT>" && node.Category != nil {
		tail := strings.TrimSpace(PathToSentence(path))
		//System.out.println("Shortcut tail = "+tail);
		that := strings.TrimSpace(tail[strings.Index(tail, "<THAT>")+len("<THAT>") : strings.Index(tail, "<TOPIC>")])
		topic := strings.TrimSpace(tail[strings.Index(tail, "<TOPIC>")+len("<TOPIC>") : len(tail)])
		//System.out.println("Shortcut that = "+that+" topic = "+topic)
		//System.out.println("Shortcut matched: "+node.category.inputThatTopic());
		thatStars[0] = that
		topicStars[0] = topic
		return node
	} else {
		g.Fail("shortCut", matchTrace)
		return nil
	}
}

func (g *Graphmaster) WordMatch(path *Paths, node *Nodemapper, inputThatTopic, starState string, starIndex int,
	inputStars, thatStars, topicStars []string, matchTrace string) *Nodemapper {
	var matchedNode *Nodemapper

	uword := strings.ToUpper(path.Word)
	if uword == "<THAT>" {
		starIndex = 0
		starState = "thatStar"
	} else if uword == "<TOPIC>" {
		starIndex = 0
		starState = "topicStar"
	}
	//System.out.println("path.next= "+path.next+" node.get="+node.get(uword));
	matchTrace += "[" + uword + "," + uword + "]"
	if path != nil && ContainsKey(node, uword) {
		if matchedNode = g.Match(path.Next, Get(node, uword), inputThatTopic, starState, starIndex, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
			return matchedNode
		}
	}

	g.Fail("word", matchTrace)
	return nil
}
func (g *Graphmaster) DollarMatch(path *Paths, node *Nodemapper, inputThatTopic, starState string, starIndex int,
	inputStars, thatStars, topicStars []string, matchTrace string) *Nodemapper {
	uword := "$" + strings.ToUpper(path.Word)
	var matchedNode *Nodemapper
	if path != nil && ContainsKey(node, uword) {
		if matchedNode = g.Match(path.Next, Get(node, uword), inputThatTopic, starState, starIndex, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
			return matchedNode
		}
	}
	g.Fail("dollar", matchTrace)
	return nil
}
func (g *Graphmaster) StarMatch(path *Paths, node *Nodemapper, input, starState string, starIndex int, inputStars, thatStars, topicStars []string, matchTrace string) *Nodemapper {
	return g.WildMatch(path, node, input, starState, starIndex, inputStars, thatStars, topicStars, "*", matchTrace)
}
func (g *Graphmaster) UnderMatch(path *Paths, node *Nodemapper, input, starState string,
	starIndex int, inputStars, thatStars, topicStars []string, matchTrace string) *Nodemapper {
	return g.WildMatch(path, node, input, starState, starIndex, inputStars, thatStars, topicStars, "_", matchTrace)
}
func (g *Graphmaster) CaretMatch(path *Paths, node *Nodemapper, input, starState string,
	starIndex int, inputStars, thatStars, topicStars []string, matchTrace string) *Nodemapper {
	var matchedNode *Nodemapper
	matchedNode = g.ZeroMatch(path, node, input, starState, starIndex, inputStars, thatStars, topicStars, "^", matchTrace)
	if matchedNode != nil {
		return matchedNode
	} else {
		return g.WildMatch(path, node, input, starState, starIndex, inputStars, thatStars, topicStars, "^", matchTrace)
	}
}

func (g *Graphmaster) SharpMatch(path *Paths, node *Nodemapper, input, starState string, starIndex int,
	inputStars, thatStars, topicStars []string, matchTrace string) *Nodemapper {
	//System.out.println("Entering sharpMatch with path.word="+path.word); NodemapperOperator.printKeys(node);
	var matchedNode *Nodemapper
	matchedNode = g.ZeroMatch(path, node, input, starState, starIndex, inputStars, thatStars, topicStars, "#", matchTrace)
	if matchedNode != nil {
		return matchedNode
	} else {
		return g.WildMatch(path, node, input, starState, starIndex, inputStars, thatStars, topicStars, "#", matchTrace)
	}
}

func (g *Graphmaster) ZeroMatch(path *Paths, node *Nodemapper, input, starState string, starIndex int,
	inputStars, thatStars, topicStars []string, wildcard, matchTrace string) *Nodemapper {
	// System.out.println("Entering zeroMatch on "+path.word+" "+NodemapperOperator.get(node, wildcard));
	matchTrace += "[" + wildcard + ",]"
	if path != nil && ContainsKey(node, wildcard) {
		//System.out.println("Zero match calling setStars Prop "+MagicStrings.null_star+" = "+bot.properties.get(MagicStrings.null_star));
		g.SetStars(g.Bot.Properties.Get(NullStar), starIndex, starState, inputStars, thatStars, topicStars)
		nextNode := Get(node, wildcard)
		return g.Match(path, nextNode, input, starState, starIndex+1, inputStars, thatStars, topicStars, matchTrace)
	} else {
		g.Fail("zero "+wildcard, matchTrace)
		return nil
	}

}

func (g *Graphmaster) WildMatch(path *Paths, node *Nodemapper, input, starState string, starIndex int,
	inputStars, thatStars, topicStars []string, wildcard, matchTrace string) *Nodemapper {
	var matchedNode *Nodemapper
	if path.Word == "<THAT>" || path.Word == "<TOPIC>" {
		g.Fail("wild1 "+wildcard, matchTrace)
		return nil
	}

	if path != nil && ContainsKey(node, wildcard) {
		matchTrace += "[" + wildcard + "," + path.Word + "]"
		var currentWord string
		var starWords string
		var pathStart *Paths
		currentWord = path.Word
		starWords = currentWord + " "
		pathStart = path.Next
		nextNode := Get(node, wildcard)
		if IsLeaf(nextNode) && !nextNode.ShortCut {
			matchedNode = nextNode
			starWords = PathToSentence(path)
			//System.out.println(starIndex+". starwords="+starWords);
			g.SetStars(starWords, starIndex, starState, inputStars, thatStars, topicStars)
			return matchedNode
		} else {
			for path = pathStart; path != nil && currentWord != "<THAT>" && currentWord != "<TOPIC>"; path = path.Next {
				matchTrace += "[" + wildcard + "," + path.Word + "]"
				if matchedNode = g.Match(path, nextNode, input, starState, starIndex+1, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
					g.SetStars(starWords, starIndex, starState, inputStars, thatStars, topicStars)
					return matchedNode
				} else {
					currentWord = path.Word
					starWords += currentWord + " "
				}
			}
			g.Fail("wild2 "+wildcard, matchTrace)
			return nil
		}
	}

	g.Fail("wild3 "+wildcard, matchTrace)
	return nil
}

func (g *Graphmaster) SetMatch(path *Paths, node *Nodemapper, input, starState string, starIndex int, inputStars, thatStars, topicStars []string, matchTrace string) *Nodemapper {
	if DEBUG {
		fmt.Printf("Graphmaster.setMatch(path: %v , node: %v, input: "+
			" %s , starState: %s, starIndex: %d, inputStars, thatStars, topicStars, matchTrace: %s, ) \n",
			path, node, input, starState, starIndex, matchTrace)
	}
	if node.Sets == nil || path.Word == "<THAT>" || path.Word == "<TOPIC>" {
		return nil
	}
	if DEBUG {
		fmt.Printf("in Graphmaster.setMatch, setMatch sets = %v\n", node.Sets)
	}
	for _, setName := range node.Sets {
		if DEBUG {
			fmt.Printf("in Graphmaster.setMatch, setMatch trying type %s\n", setName)
		}
		nextNode := Get(node, "<SET>"+strings.ToUpper(setName)+"</SET>")
		aimlSet := g.Bot.SetMap[setName]
		//System.out.println(aimlSet.setName + "="+ aimlSet);
		var matchedNode *Nodemapper
		var bestMatchedNode *Nodemapper = nil
		currentWord := path.Word
		starWords := currentWord + " "
		length := 1
		matchTrace += "[<set>" + setName + "</set>," + path.Word + "]"
		if DEBUG {
			fmt.Println("in Graphmaster.setMatch, setMatch starWords =\"" + starWords + "\"")
		}
		for qath := path.Next; qath != nil && currentWord != "<THAT>" && currentWord != "<TOPIC>" && length <= aimlSet.MaxLength; qath = qath.Next {
			if DEBUG {
				fmt.Println("in Graphmaster.setMatch, qath.word = " + qath.Word)
			}
			phrase := strings.ToUpper(g.Bot.PreProcessor.Normalize(strings.TrimSpace(starWords)))
			if DEBUG {
				fmt.Println("in Graphmaster.setMatch, setMatch trying \"" + phrase + "\" in " + setName)
			}
			if aimlSet.Contains(phrase) {
				if matchedNode = g.Match(qath, nextNode, input, starState, starIndex+1, inputStars, thatStars, topicStars, matchTrace); matchedNode != nil {
					g.SetStars(starWords, starIndex, starState, inputStars, thatStars, topicStars)
					if DEBUG {
						fmt.Println("in Graphmaster.setMatch, setMatch found " + phrase + " in " + setName)
					}
					bestMatchedNode = matchedNode
				}
			}
			//    else if (qath.word.equals("<THAT>") || qath.word.equals("<TOPIC>")) return null;

			length = length + 1
			currentWord = qath.Word
			starWords += currentWord + " "

		}
		if bestMatchedNode != nil {
			return bestMatchedNode
		}
	}
	g.Fail("set", matchTrace)
	return nil
}

func (g *Graphmaster) SetStars(starWords string, starIndex int, starState string, inputStars, thatStars, topicStars []string) {
	if starIndex < MaxStars {
		//System.out.println("starWords="+starWords);
		starWords = strings.TrimSpace(starWords)
		if starState == "inputStar" {
			inputStars[starIndex] = starWords
		} else if starState == "thatStar" {
			thatStars[starIndex] = starWords
		} else if starState == "topicStar" {
			topicStars[starIndex] = starWords
		}
	}
}
func (g *Graphmaster) Printgraph() {
	g.PrintgraphPartial(g.Root, "")
}
func (g *Graphmaster) PrintgraphPartial(node *Nodemapper, partial string) {
	if node == nil {
		fmt.Println("Null graph")
	} else {
		template := ""
		if IsLeaf(node) || node.ShortCut {
			ss := node.Category.GetTemplate()
			template = TemplateToLine(ss)
			template = template[0:min(16, len(template))]
			if node.ShortCut {
				fmt.Println(partial + "(" + strconv.Itoa(Size(node)) + "[" + strconv.Itoa(node.Height) + "])--<THAT>-->X(1)--*-->X(1)--<TOPIC>-->X(1)--*-->" + template + "...")
			} else {
				fmt.Println(partial + "(" + strconv.Itoa(Size(node)) + "[" + strconv.Itoa(node.Height) + "]) " + template + "...")
			}
		}
		for _, key := range KeySet(node) {
			//System.out.println(key);
			g.PrintgraphPartial(Get(node, key), partial+"("+strconv.Itoa(Size(node))+"["+strconv.Itoa(node.Height)+"])--"+key+"-->")
		}
	}
}
func (g *Graphmaster) GetCategories() []*Category {
	categories := []*Category{}
	GetCategoriesTo(g.Root, categories)
	//for (Category c : categories) System.out.println("getCategories: "+c.inputThatTopic()+" "+c.getTemplate());
	return categories
}
func GetCategoriesTo(node *Nodemapper, categories []*Category) {
	if node == nil {
		return
	} else {
		//String template = "";
		if IsLeaf(node) || node.ShortCut {
			if node.Category != nil {
				categories = append(categories, node.Category)
			} // node.category == null when the category is deleted.
		}
		for _, key := range KeySet(node) {
			//System.out.println(key);
			GetCategoriesTo(Get(node, key), categories)
		}
	}
}

var leafCnt int
var nodeCnt int
var nodeSize int64
var singletonCnt int
var shortCutCnt int
var naryCnt int

func (g *Graphmaster) NodeStats() {
	leafCnt = 0
	nodeCnt = 0
	nodeSize = 0
	singletonCnt = 0
	shortCutCnt = 0
	naryCnt = 0
	g.NodeStatsGraph(g.Root)
	v := strconv.FormatFloat(float64(float32(nodeSize)/float32(nodeCnt)), 'E', -1, 32)
	resultNote := g.Bot.Name + " (" + g.Name + "): " + strconv.Itoa(len(g.GetCategories())) + " categories " + strconv.Itoa(nodeCnt) + " nodes " + strconv.Itoa(singletonCnt) +
		" singletons " + strconv.Itoa(leafCnt) + " leaves " + strconv.Itoa(shortCutCnt) + " shortcuts " + strconv.Itoa(naryCnt) +
		" n-ary " + strconv.Itoa(int(nodeSize)) + " branches " + v + " average branching "
	if TraceMode {
		fmt.Println(resultNote)
	}
}
func (g *Graphmaster) NodeStatsGraph(node *Nodemapper) {
	if node != nil {
		//System.out.println("Counting "+node.key+ " size="+NodemapperOperator.size(node));
		nodeCnt++
		nodeSize += int64(Size(node))
		if Size(node) == 1 {
			singletonCnt += 1
		}
		if IsLeaf(node) && !node.ShortCut {
			leafCnt++
		}
		if Size(node) > 1 {
			naryCnt += 1
		}
		if node.ShortCut {
			shortCutCnt += 1
		}
		for _, key := range KeySet(node) {
			g.NodeStatsGraph(Get(node, key))
		}
	}
}

func (g *Graphmaster) GetVocabulary() []string {
	vocabulary := map[string]bool{}
	g.GetBrainVocabulary(g.Root)
	for key, _ := range g.Bot.SetMap {
		for _, val := range g.Bot.SetMap[key].Set {
			vocabulary[val] = true
		}
	}
	ret := []string{}
	for key, _ := range vocabulary {
		ret = append(ret, key)
	}
	return ret
}
func (g *Graphmaster) GetBrainVocabulary(node *Nodemapper) {
	if node != nil {
		//System.out.println("Counting "+node.key+ " size="+NodemapperOperator.size(node));
		for _, key := range KeySet(node) {
			g.Vocabulary[key] = true
			g.GetBrainVocabulary(Get(node, key))
		}
	}
}
