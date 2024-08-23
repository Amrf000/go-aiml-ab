package ab

import (
	"aiml/external/go-dom"
	"fmt"
	"math/rand"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type AIMLProcessor struct {
}

var DEBUG bool
var Extension AIMLProcessorExtension

func CategoryProcessor(n dom.Node, categories *[]*Category, topic, aimlFile, language string) {
	var pattern, that, template string
	children := n.GetChildNodes()
	pattern = "*"
	that = "*"
	template = ""

	for j := 0; j < children.GetLength(); j++ {
		m := children.Item(j)
		mName := m.GetNodeName()
		switch mName {
		case "#text":
			//nn := 0
			//nn++
		case "pattern":
			pattern = NodeToString(m)
		case "that":
			that = NodeToString(m)
		case "topic":
			topic = NodeToString(m)
		case "template":
			template = NodeToString(m)
		default:
			panic(fmt.Errorf("categoryProcessor: unexpected %s in %s\n", mName, NodeToString(m)))
		}
	}

	pattern = TrimTag(pattern, "pattern")
	that = TrimTag(that, "that")
	topic = TrimTag(topic, "topic")
	pattern = CleanPattern(pattern)
	that = CleanPattern(that)
	topic = CleanPattern(topic)
	template = TrimTag(template, "template")

	if JpTokenize {
		pattern = TokenizeSentence(pattern)
		that = TokenizeSentence(that)
		topic = TokenizeSentence(topic)
	}

	c := NewCategory(0, pattern, that, topic, template, aimlFile)
	if template == "" {
		panic(fmt.Errorf("Category %s discarded due to blank or missing <template>.\n", c.InputThatTopic()))
	} else {
		*categories = append(*categories, c)
	}
}
func CleanPattern(pattern string) string {
	re := regexp.MustCompile(`(\r\n|\n\r|\r|\n)`)
	pattern = re.ReplaceAllString(pattern, " ")
	pattern = strings.ReplaceAll(pattern, "  ", " ")
	return strings.TrimSpace(pattern)
}

func TrimTag(s, tagName string) string {
	stag := "<" + tagName + ">"
	etag := "</" + tagName + ">"
	if strings.HasPrefix(s, stag) && strings.HasSuffix(s, etag) {
		s = s[len(stag):]
		s = s[:len(s)-len(etag)]
	}
	return strings.TrimSpace(s)
}

func AIMLToCategories(directory, aimlFile string) []*Category {
	var categories []*Category

	root, err := ParseFile(filepath.Join(directory, aimlFile))
	if err != nil {
		fmt.Printf("AIMLToCategories: %v\n", err)
		return nil
	}

	language := DefaultLanguage
	if root.GetAttributes().GetLength() > 0 {
		XMLAttributes := root.GetAttributes()
		for i := 0; i < XMLAttributes.GetLength(); i++ {
			if XMLAttributes.Item(i).GetNodeName() == "language" {
				language = XMLAttributes.Item(i).GetValue()
			}
		}
	}

	nodelist := root.GetChildNodes()
	for i := 0; i < nodelist.GetLength(); i++ {
		n := nodelist.Item(i)
		nodeName := n.GetNodeName()
		if nodeName == "category" {
			//if len(n.Children) > 2 {
			//	nn := 0
			//	nn++
			//}
			CategoryProcessor(n, &categories, "*", aimlFile, language)
		} else if nodeName == "topic" {
			topic, _ := n.(dom.Element).GetAttribute("name")
			children := n.GetChildNodes()
			for j := 0; j < children.GetLength(); j++ {
				m := children.Item(j)
				if m.GetNodeName() == "category" {
					CategoryProcessor(m, &categories, topic, aimlFile, language)
				}
			}
		}
	}

	return categories
}

var sraiCount int
var repeatCount int

func CheckForRepeat(input string, chatSession *Chat) int {
	if input == chatSession.InputHistory.Get(1) {
		return 1
	}
	return 0
}

func Respond(input, that, topic string, chatSession *Chat) string {
	if false {
		return "Repeat!"
	}
	return RespondWithSrCnt(input, that, topic, chatSession, 0)
}

func RespondWithSrCnt(input, that, topic string, chatSession *Chat, srCnt int) string {
	trace(fmt.Sprintf("input: %s, that: %s, topic: %s, chatSession: %#v, srCnt: %d", input, that, topic, chatSession, srCnt))
	var response string
	if input == "" {
		input = NullInput
	}
	sraiCount = srCnt
	response = DefaultBotResponse
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in respondWithSrCnt", r)
		}
	}()
	leaf := chatSession.Bot.Brain.MatchRaw(input, that, topic)
	if leaf == nil {
		return response
	}
	ps := NewParseState(0, chatSession, input, that, topic, leaf)
	template := leaf.Category.GetTemplate()
	response = EvalTemplate(template, ps)
	return response
}
func CapitalizeString(input string) string {
	chars := []rune(strings.ToLower(input))
	found := false
	for i := 0; i < len(chars); i++ {
		if !found && unicode.IsLetter(chars[i]) {
			chars[i] = unicode.ToUpper(chars[i])
			found = true
		} else if unicode.IsSpace(chars[i]) {
			found = false
		}
	}
	return string(chars)
}

func Explode(input string) string {
	var result strings.Builder
	for _, char := range input {
		result.WriteString(" ")
		result.WriteRune(char)
	}
	exploded := result.String()
	return strings.TrimSpace(strings.ReplaceAll(exploded, "  ", " "))
}

func EvalTagContent(node dom.Node, ps *ParseState, ignoreAttributes map[string]bool) string {
	var result strings.Builder
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Something went wrong with evalTagContent", r)
		}
	}()
	childList := node.GetChildNodes()
	for i := 0; i < childList.GetLength(); i++ {
		child := childList.Item(i)
		if ignoreAttributes == nil {
			ignoreAttributes = map[string]bool{}
		}
		if _, ignored := ignoreAttributes[child.GetNodeName()]; !ignored {
			result.WriteString(RecursEval(child, ps))
		}
	}
	//if childList == nil && node.Text != "" {
	//	return node.Text
	//}
	return result.String()
}

var traceCount int

func GenericXML(node dom.Node, ps *ParseState) string {
	evalResult := EvalTagContent(node, ps, nil)
	result := UnevaluatedXML(evalResult, node, ps)
	return result
}

func UnevaluatedXML(resultIn string, nodea dom.Node, ps *ParseState) string {
	node := nodea.(dom.Element)
	nodeName := node.GetNodeName()
	attributes := ""
	if node.GetAttributes().GetLength() > 0 {
		XMLAttributes := node.GetAttributes()
		for i := 0; i < XMLAttributes.GetLength(); i++ {
			attributes += fmt.Sprintf(" %s=\"%s\"", XMLAttributes.Item(i).GetNodeName(), XMLAttributes.Item(i).GetValue())
		}
	}
	result := fmt.Sprintf("<%s%s/>", nodeName, attributes)
	if resultIn != "" {
		result = fmt.Sprintf("<%s%s>%s</%s>", nodeName, attributes, resultIn, nodeName)
	}
	return result
}
func Srai(node dom.Node, ps *ParseState) string {
	sraiCount++
	if sraiCount > MaxRecursionCount || ps.Depth > MaxRecursionDepth {
		return TooMuchRecursion
	}
	response := DefaultBotResponse
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in srai", r)
		}
	}()
	result := EvalTagContent(node, ps, nil)
	result = strings.TrimSpace(result)
	result = strings.ReplaceAll(result, "\r\n", " ")
	result = strings.ReplaceAll(result, "\n\r", " ")
	result = strings.ReplaceAll(result, "\r", " ")
	result = strings.ReplaceAll(result, "\n", " ")
	result = ps.ChatSession.Bot.PreProcessor.Normalize(result)
	result = TokenizeSentence(result)
	topic := ps.ChatSession.Predicates["topic"]
	if TraceMode {
		fmt.Printf("%d. <srai>%s</srai> from %s topic=%s)\n", traceCount, result, ps.Leaf.Category.InputThatTopic(), topic)
		traceCount++
	}
	leaf := ps.ChatSession.Bot.Brain.MatchRaw(result, ps.That, topic)
	if leaf == nil {
		return response
	}
	response = EvalTemplate(leaf.Category.GetTemplate(), NewParseState(ps.Depth+1, ps.ChatSession, ps.Input, ps.That, topic, leaf))
	return strings.TrimSpace(response)
}
func GetAttributeOrTagValue(node dom.Node, ps *ParseState, attributeName string) string {
	var result string
	m, ok := node.(dom.Element).GetAttribute(attributeName)
	if !ok {
		childList := node.GetChildNodes()
		for i := 0; i < childList.GetLength(); i++ {
			child := childList.Item(i)
			if child.GetNodeName() == attributeName {
				result = EvalTagContent(child, ps, nil)
				break
			}
		}
	} else {
		result = m
	}
	return result
}

func Sraix(node dom.Node, ps *ParseState) string {
	attributeNames := map[string]bool{
		"botid": true,
		"host":  true,
	}
	host := GetAttributeOrTagValue(node, ps, "host")
	botid := GetAttributeOrTagValue(node, ps, "botid")
	hint := GetAttributeOrTagValue(node, ps, "hint")
	limit := GetAttributeOrTagValue(node, ps, "limit")
	defaultResponse := GetAttributeOrTagValue(node, ps, "default")
	evalResult := EvalTagContent(node, ps, attributeNames)
	result := SraixSraix(ps.ChatSession, evalResult, defaultResponse, hint, host, botid, "", limit)
	return result
}
func MapNode(node dom.Node, ps *ParseState) string {
	result := DefaultMap
	attributeNames := StringSet("name")
	mapName := GetAttributeOrTagValue(node, ps, "name")
	contents := EvalTagContent(node, ps, attributeNames)
	contents = strings.TrimSpace(contents)
	if mapName == "" {
		result = "<map>" + contents + "</map>"
	} else {
		if mapInstance, ok := ps.ChatSession.Bot.MapMap[mapName]; ok {
			result = mapInstance.Get(strings.ToUpper(contents))
		}
		if result == "" {
			result = DefaultMap
		}
		result = strings.TrimSpace(result)
	}
	return result
}

func SetNode(node dom.Node, ps *ParseState) string {
	attributeNames := StringSet("name", "var")
	predicateName := GetAttributeOrTagValue(node, ps, "name")
	varName := GetAttributeOrTagValue(node, ps, "var")
	result := strings.TrimSpace(EvalTagContent(node, ps, attributeNames))
	result = strings.ReplaceAll(result, "\r\n", " ")
	result = strings.ReplaceAll(result, "\n\r", " ")
	result = strings.ReplaceAll(result, "\r", " ")
	result = strings.ReplaceAll(result, "\n", " ")
	value := strings.TrimSpace(result)

	if predicateName != "" {
		ps.ChatSession.Predicates[predicateName] = result
		trace(fmt.Sprintf("Set predicate %s to %s in %s", predicateName, result, ps.Leaf.Category.InputThatTopic()))
	}

	if varName != "" {
		ps.Vars[varName] = result
		trace(fmt.Sprintf("Set var %s to %s in %s", varName, value, ps.Leaf.Category.InputThatTopic()))
	}

	if ps.ChatSession.Bot.PronounSet[predicateName] {
		result = predicateName
	}

	return result
}
func GetA(node dom.Node, ps *ParseState) string {
	result := DefaultGet
	predicateName := GetAttributeOrTagValue(node, ps, "name")
	varName := GetAttributeOrTagValue(node, ps, "var")
	tupleName := GetAttributeOrTagValue(node, ps, "tuple")

	if predicateName != "" {
		if val, ok := ps.ChatSession.Predicates[predicateName]; ok {
			result = strings.TrimSpace(val)
		}
	} else if varName != "" && tupleName != "" {
		result = TupleGet(tupleName, varName)
	} else if varName != "" {
		if val, ok := ps.Vars[varName]; ok {
			result = strings.TrimSpace(val)
		}
	}
	return result
}

func TupleGet(tupleName, varName string) string {
	result := DefaultGet
	if tuple, ok := TupleMap[tupleName]; ok {
		result = tuple.GetValue(varName)
	}
	return result
}

func Abot(node dom.Node, ps *ParseState) string {
	result := DefaultProperty
	propertyName := GetAttributeOrTagValue(node, ps, "name")
	if propertyName != "" {
		val := ps.ChatSession.Bot.Properties.Get(propertyName)
		result = strings.TrimSpace(val)
	}
	return result
}

func Date(node dom.Node, ps *ParseState) string {
	jformat := GetAttributeOrTagValue(node, ps, "jformat")
	locale := GetAttributeOrTagValue(node, ps, "locale")
	timezone := GetAttributeOrTagValue(node, ps, "timezone")
	return DateCustom(jformat, locale, timezone)
}

func Interval(node dom.Node, ps *ParseState) string {
	style := GetAttributeOrTagValue(node, ps, "style")
	jformat := GetAttributeOrTagValue(node, ps, "jformat")
	from := GetAttributeOrTagValue(node, ps, "from")
	to := GetAttributeOrTagValue(node, ps, "to")

	if style == "" {
		style = "years"
	}
	if jformat == "" {
		jformat = "MMMMMMMMM dd, yyyy"
	}
	if from == "" {
		from = "January 1, 1970"
	}
	if to == "" {
		to = DateCustom(jformat, "", "")
	}

	var result string
	switch style {
	case "years":
		result = strconv.Itoa(GetYearsBetween(from, to, jformat))
	case "months":
		result = strconv.Itoa(GetMonthsBetween(from, to, jformat))
	case "days":
		result = strconv.Itoa(GetDaysBetween(from, to, jformat))
	case "hours":
		result = strconv.Itoa(GetHoursBetween(from, to, jformat))
	default:
		result = "unknown"
	}
	return result
}

func GetIndexValue(node dom.Node, ps *ParseState) int {
	index := 0
	value := GetAttributeOrTagValue(node, ps, "index")
	if value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			index = intValue - 1
		}
	}
	return index
}

func InputStar(node dom.Node, ps *ParseState) string {
	result := ""
	index := GetIndexValue(node, ps)
	if index < 0 || index >= len(ps.StarBindings.InputStars.Items) {
		return result
	}
	if ps.StarBindings.InputStars.Items[index] != "" {
		result = strings.TrimSpace(ps.StarBindings.InputStars.Items[index])
	}
	return result
}
func ThatStar(node dom.Node, ps *ParseState) string {
	index := GetIndexValue(node, ps)
	if star := ps.StarBindings.ThatStars.Star(index); star != "" {
		return strings.TrimSpace(star)
	}
	return ""
}

func TopicStar(node dom.Node, ps *ParseState) string {
	index := GetIndexValue(node, ps)
	if star := ps.StarBindings.TopicStars.Star(index); star != "" {
		return strings.TrimSpace(star)
	}
	return ""
}
func Id(node dom.Node, ps *ParseState) string {
	return ps.ChatSession.CustomerID
}

func SizeA(node dom.Node, ps *ParseState) string {
	return strconv.Itoa(len(ps.ChatSession.Bot.Brain.GetCategories()))
}

func Vocabulary(node dom.Node, ps *ParseState) string {
	return strconv.Itoa(len(ps.ChatSession.Bot.Brain.GetVocabulary()))
}

func Program(node dom.Node, ps *ParseState) string {
	return ProgramNameVersion
}
func That(node dom.Node, ps *ParseState) string {
	index := 0
	jndex := 0
	value := GetAttributeOrTagValue(node, ps, "index")
	if value != "" {
		parts := strings.Split(value, ",")
		if len(parts) == 2 {
			var err error
			if index, err = strconv.Atoi(parts[0]); err == nil {
				index--
			}
			if jndex, err = strconv.Atoi(parts[1]); err == nil {
				jndex--
			}
			if err != nil {
				fmt.Println("Error parsing index:", err)
			}
		}
		fmt.Printf("That index=%d,%d\n", index, jndex)
	}

	that := unknownHistoryItem
	if index >= 0 && index < len(ps.ChatSession.ThatHistory.MHistory) {
		hist := ps.ChatSession.ThatHistory.Get(index)
		if hist != nil {
			if hh, ok := hist.(History); ok {
				item := hh.Get(jndex)
				if kk, ok := item.(string); ok && kk != "" {
					that = strings.TrimSpace(kk)
				}
			}
		}
	}
	return that
}
func Input(node dom.Node, ps *ParseState) string {
	index := GetIndexValue(node, ps)
	return ps.ChatSession.InputHistory.GetString(index)
}

func Request(node dom.Node, ps *ParseState) string {
	index := GetIndexValue(node, ps)
	return strings.TrimSpace(ps.ChatSession.RequestHistory.GetString(index))
}

func Response(node dom.Node, ps *ParseState) string {
	index := GetIndexValue(node, ps)
	return strings.TrimSpace(ps.ChatSession.ResponseHistory.GetString(index))
}

func System(node dom.Node, ps *ParseState) string {
	attributeNames := StringSet("timeout")
	evaluatedContents := EvalTagContent(node, ps, attributeNames)
	result := Utils_System(evaluatedContents, SystemFailed)
	return result
}
func Think(node dom.Node, ps *ParseState) string {
	EvalTagContent(node, ps, nil)
	return ""
}

func ExplodeNode(node dom.Node, ps *ParseState) string {
	result := EvalTagContent(node, ps, nil)
	return Explode(result)
}

func Normalize(node dom.Node, ps *ParseState) string {
	result := EvalTagContent(node, ps, nil)
	return ps.ChatSession.Bot.PreProcessor.Normalize(result)
}

func Denormalize(node dom.Node, ps *ParseState) string {
	result := EvalTagContent(node, ps, nil)
	return ps.ChatSession.Bot.PreProcessor.Denormalize(result)
}
func Uppercase(node dom.Node, ps *ParseState) string {
	result := EvalTagContent(node, ps, nil)
	return strings.ToUpper(result)
}

func Lowercase(node dom.Node, ps *ParseState) string {
	result := EvalTagContent(node, ps, nil)
	return strings.ToLower(result)
}

func Formal(node dom.Node, ps *ParseState) string {
	result := EvalTagContent(node, ps, nil)
	return CapitalizeString(result)
}

// Function implementations
func Sentence(node dom.Node, ps *ParseState) string {
	result := EvalTagContent(node, ps, nil)
	if len(result) > 1 {
		return strings.ToUpper(result[:1]) + result[1:]
	}
	return ""
}

func Person(node dom.Node, ps *ParseState) string {
	result := ""
	if node.GetChildNodes().GetLength() > 0 {
		result = EvalTagContent(node, ps, nil)
	} else {
		result = " " + ps.StarBindings.InputStars.Star(0) + " "
	}
	result = ps.ChatSession.Bot.PreProcessor.Person(result)
	return strings.TrimSpace(result)
}

func Person2(node dom.Node, ps *ParseState) string {
	result := ""
	if node.GetChildNodes().GetLength() > 0 {
		result = EvalTagContent(node, ps, nil)
	} else {
		result = " " + ps.StarBindings.InputStars.Star(0) + " "
	}
	result = ps.ChatSession.Bot.PreProcessor.Person2(result)
	return strings.TrimSpace(result)
}
func Gender(node dom.Node, ps *ParseState) string {
	result := EvalTagContent(node, ps, nil)
	result = " " + result + " "
	result = ps.ChatSession.Bot.PreProcessor.Gender(result)
	return strings.TrimSpace(result)
}

func Random(node dom.Node, ps *ParseState) string {
	childList := node.GetChildNodes()
	var liList []dom.Node
	// setName := GetAttributeOrTagValue(node, ps, "set")

	for i := 0; i < childList.GetLength(); i++ {
		if childList.Item(i).GetNodeName() == "li" {
			liList = append(liList, childList.Item(i))
		}
	}

	if len(liList) == 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(liList))
	if QaTestMode {
		index = 0
	}
	return EvalTagContent(liList[index], ps, nil)
}

func UnevaluatedAIML(node dom.Node, ps *ParseState) string {
	result := LearnEvalTagContent(node, ps)
	return UnevaluatedXML(result, node, ps)
}
func RecursLearn(node dom.Node, ps *ParseState) string {
	nodeName := node.GetNodeName()
	if nodeName == "#text" {
		return node.(dom.Text).GetValue()
	} else if nodeName == "eval" {
		return EvalTagContent(node, ps, nil)
	} else {
		return UnevaluatedAIML(node, ps)
	}
}

func LearnEvalTagContent(node dom.Node, ps *ParseState) string {
	var result strings.Builder
	childList := node.GetChildNodes()
	for i := 0; i < childList.GetLength(); i++ {
		child := childList.Item(i)
		result.WriteString(RecursLearn(child, ps))
	}
	return result.String()
}
func Learn(node dom.Node, ps *ParseState) string {
	childList := node.GetChildNodes()
	var pattern, that, template string
	that = "*"

	for i := 0; i < childList.GetLength(); i++ {
		nodeName := childList.Item(i).GetNodeName()
		if nodeName == "category" {
			grandChildList := childList.Item(i).GetChildNodes()
			for j := 0; j < grandChildList.GetLength(); j++ {
				subNodeName := grandChildList.Item(j).GetNodeName()
				switch subNodeName {
				case "pattern":
					pattern = RecursLearn(grandChildList.Item(j).(dom.Element), ps)
				case "that":
					that = RecursLearn(grandChildList.Item(j).(dom.Element), ps)
				case "template":
					template = RecursLearn(grandChildList.Item(j).(dom.Element), ps)
				}
			}

			pattern = strings.TrimPrefix(pattern, "<pattern>")
			pattern = strings.TrimSuffix(pattern, "</pattern>")
			if TraceMode {
				println("Learn Pattern = " + pattern)
			}

			if len(template) >= len("<template></template>") {
				template = strings.TrimPrefix(template, "<template>")
				template = strings.TrimSuffix(template, "</template>")
			}

			if len(that) >= len("<that></that>") {
				that = strings.TrimPrefix(that, "<that>")
				that = strings.TrimSuffix(that, "</that>")
			}

			pattern = strings.ToUpper(pattern)
			pattern = strings.ReplaceAll(pattern, "\n", " ")
			pattern = strings.Join(strings.Fields(pattern), " ")

			that = strings.ToUpper(that)
			that = strings.ReplaceAll(that, "\n", " ")
			that = strings.Join(strings.Fields(that), " ")

			if TraceMode {
				println("Learn Pattern = " + pattern)
				println("Learn That = " + that)
				println("Learn Template = " + template)
			}

			var c *Category = nil

			if node.GetNodeName() == "learn" {
				c = NewCategory(0, pattern, that, "*", template, NullAimlFile)
				ps.ChatSession.Bot.LearnGraph.AddCategory(c)
			} else {
				c = NewCategory(0, pattern, that, "*", template, LearnfAimlFile)
				ps.ChatSession.Bot.LearnfGraph.AddCategory(c)
			}

			ps.ChatSession.Bot.Brain.AddCategory(c)
		}
	}
	return ""
}
func LoopCondition(node dom.Node, ps *ParseState) string {
	loop := true
	result := ""
	loopCnt := 0

	for loop && loopCnt < MaxLoops {
		loopResult := Condition(node, ps)

		if strings.TrimSpace(loopResult) == TooMuchRecursion {
			return TooMuchRecursion
		}

		if strings.Contains(loopResult, "<loop/>") {
			loopResult = strings.ReplaceAll(loopResult, "<loop/>", "")
			loop = true
		} else {
			loop = false
		}

		result += loopResult
		loopCnt++
	}

	if loopCnt >= MaxLoops {
		result = TooMuchLooping
	}

	return result
}
func Condition(node dom.Node, ps *ParseState) string {
	var result string
	childList := node.GetChildNodes() // Assuming method to get child nodes
	var liList []dom.Node
	var predicate, varName, value string
	attributeNames := make(map[string]bool) // Use a map for attribute names
	attributeNames["name"] = true
	attributeNames["var"] = true
	attributeNames["value"] = true

	predicate = GetAttributeOrTagValue(node, ps, "name")
	varName = GetAttributeOrTagValue(node, ps, "var")

	for i := 0; i < childList.GetLength(); i++ {
		if childList.Item(i).GetNodeName() == "li" {
			liList = append(liList, childList.Item(i))
		}
	}

	if len(liList) == 0 {
		value = GetAttributeOrTagValue(node, ps, "value")
		if value != "" {
			if predicate != "" && ps.ChatSession.Predicates[predicate] == value {
				return EvalTagContent(node, ps, attributeNames)
			}
			if varName != "" && ps.Vars[varName] == value {
				return EvalTagContent(node, ps, attributeNames)
			}
		}
	} else {
		for i := 0; i < len(liList) && result == ""; i++ {
			n := liList[i].(dom.Element)
			liPredicate := predicate
			liVarName := varName
			value = GetAttributeOrTagValue(n, ps, "value")

			if liPredicate == "" {
				liPredicate = GetAttributeOrTagValue(n, ps, "name")
			}
			if liVarName == "" {
				liVarName = GetAttributeOrTagValue(n, ps, "var")
			}

			if value != "" {
				if liPredicate != "" && (ps.ChatSession.Predicates[liPredicate] == value || value == "*") {
					return EvalTagContent(n, ps, attributeNames)
				}
				if liVarName != "" && (ps.Vars[liVarName] == value || value == "*") {
					return EvalTagContent(n, ps, attributeNames)
				}
			} else {
				return EvalTagContent(n, ps, attributeNames)
			}
		}
	}

	return result
}

func EvalTagForLoop(node dom.Node) bool {
	childList := node.GetChildNodes() // Placeholder for actual implementation
	for i := 0; i < childList.GetLength(); i++ {
		if childList.Item(i).GetNodeName() == "loop" {
			return true
		}
	}
	return false
}

func DeleteTriple(node dom.Node, ps *ParseState) string {
	subject := GetAttributeOrTagValue(node, ps, "subj")
	predicate := GetAttributeOrTagValue(node, ps, "pred")
	object := GetAttributeOrTagValue(node, ps, "obj")
	return ps.ChatSession.TripleStore.DeleteTriple(subject, predicate, object)
}

func AddTriple(node dom.Node, ps *ParseState) string {
	subject := GetAttributeOrTagValue(node, ps, "subj")
	predicate := GetAttributeOrTagValue(node, ps, "pred")
	object := GetAttributeOrTagValue(node, ps, "obj")
	return ps.ChatSession.TripleStore.AddTriple(subject, predicate, object)
}
func Uniq(node dom.Node, ps *ParseState) string {
	vars := make(map[string]bool)
	visibleVars := make(map[string]bool)
	subj := "?subject"
	pred := "?predicate"
	obj := "?object"
	childList := node.GetChildNodes() // Placeholder for actual implementation

	for j := 0; j < childList.GetLength(); j++ {
		childNode := childList.Item(j)
		contents := EvalTagContent(childNode, ps, nil)
		nodeName := childNode.GetNodeName() // Placeholder for actual implementation
		if nodeName == "subj" {
			subj = contents
		} else if nodeName == "pred" {
			pred = contents
		} else if nodeName == "obj" {
			obj = contents
		}
		if strings.HasPrefix(contents, "?") {
			visibleVars[contents] = true
			vars[contents] = true
		}
	}

	partial := Tuple{Name: ""}
	clause := Clause{Subj: subj, Pred: pred, Obj: obj}
	tuples := ps.ChatSession.TripleStore.SelectFromSingleClause(&partial, &clause, true)
	var tupleList string
	for _, tuple := range tuples {
		tupleList = tuple.Name + " " + tupleList
	}
	tupleList = strings.TrimSpace(tupleList)
	if len(tupleList) == 0 {
		tupleList = "NIL"
	}
	var varName string
	for x := range visibleVars {
		varName = x
	}
	firstTuple := FirstWord(tupleList)
	result := TupleGet(firstTuple, varName)
	return result
}
func Select(node dom.Node, ps *ParseState) string {
	clauses := []*Clause{}
	vars := make(map[string]bool)
	visibleVars := make(map[string]bool)
	childList := node.GetChildNodes() // Placeholder for actual implementation

	for i := 0; i < childList.GetLength(); i++ {
		childNode := childList.Item(i)
		nodeName := childNode.GetNodeName() // Placeholder for actual implementation
		if nodeName == "vars" {
			contents := EvalTagContent(childNode, ps, nil)
			splitVars := strings.Fields(contents)
			for _, varName := range splitVars {
				varName = strings.TrimSpace(varName)
				if len(varName) > 0 {
					visibleVars[varName] = true
				}
			}
		} else if nodeName == "q" || nodeName == "notq" {
			affirm := nodeName == "q"
			grandChildList := childNode.GetChildNodes() // Placeholder for actual implementation
			subj, pred, obj := "", "", ""
			for j := 0; j < grandChildList.GetLength(); j++ {
				grandChildNode := grandChildList.Item(j)
				contents := EvalTagContent(grandChildNode, ps, nil)
				grandChildNodeName := grandChildNode.GetNodeName() // Placeholder for actual implementation
				if grandChildNodeName == "subj" {
					subj = contents
				} else if grandChildNodeName == "pred" {
					pred = contents
				} else if grandChildNodeName == "obj" {
					obj = contents
				}
				if strings.HasPrefix(contents, "?") {
					vars[contents] = true
				}
			}
			clause := Clause{Subj: subj, Pred: pred, Obj: obj, Affirm: affirm}
			clauses = append(clauses, &clause)
		}
	}

	tuples := ps.ChatSession.TripleStore.Select(vars, visibleVars, clauses)
	var result strings.Builder
	for _, tuple := range tuples {
		result.WriteString(tuple.Name + " ")
	}
	finalResult := strings.TrimSpace(result.String())
	if len(finalResult) == 0 {
		finalResult = "NIL"
	}
	return finalResult
}
func Subject(node dom.Node, ps *ParseState) string {
	id := EvalTagContent(node, ps, nil)
	ts := ps.ChatSession.TripleStore
	subject := "unknown"
	if triple, exists := ts.IdTriple[id]; exists {
		subject = triple.Subject
	}
	return subject
}

func Predicate(node dom.Node, ps *ParseState) string {
	id := EvalTagContent(node, ps, nil)
	ts := ps.ChatSession.TripleStore
	if triple, exists := ts.IdTriple[id]; exists {
		return triple.Predicate
	}
	return "unknown"
}
func Object(node dom.Node, ps *ParseState) string {
	id := EvalTagContent(node, ps, nil)
	ts := ps.ChatSession.TripleStore
	if triple, exists := ts.IdTriple[id]; exists {
		return triple.Object
	}
	return "unknown"
}

// javascript evaluates a JavaScript script and returns the result
func Javascript(node dom.Node, ps *ParseState) string {
	result := BadJavascript
	script := EvalTagContent(node, ps, nil)

	res, err := EvalScript("JavaScript", script)
	if err != nil {
		panic(err) // Print the error
	} else {
		result = res
	}

	trace("in AIMLProcessor.javascript, returning result: " + result)
	return result
}
func FirstWord(sentence string) string {
	content := strings.TrimSpace(sentence)
	if strings.Contains(content, " ") {
		head := content[:strings.Index(content, " ")]
		return head
	} else if len(content) > 0 {
		return content
	}
	return DefaultListItem
}

// restWords returns the rest of the words from a sentence after the first word
func RestWords(sentence string) string {
	content := strings.TrimSpace(sentence)
	if strings.Contains(content, " ") {
		tail := content[strings.Index(content, " ")+1:]
		return tail
	}
	return DefaultListItem
}

// first returns the first word from the evaluated content of a node
func First(node dom.Node, ps *ParseState) string {
	content := EvalTagContent(node, ps, nil)
	return FirstWord(content)
}
func Rest(node dom.Node, ps *ParseState) string {
	content := EvalTagContent(node, ps, nil)
	content = ps.ChatSession.Bot.PreProcessor.Normalize(content)
	return RestWords(content)
}

// resetlearnf deletes Learnf categories and returns a confirmation message
func Resetlearnf(node dom.Node, ps *ParseState) string {
	ps.ChatSession.Bot.DeleteLearnfCategories()
	return "Deleted Learnf Categories"
}

// resetlearn deletes Learn categories and returns a confirmation message
func Resetlearn(node dom.Node, ps *ParseState) string {
	ps.ChatSession.Bot.DeleteLearnCategories()
	return "Deleted Learn Categories"
}
func RecursEval(node dom.Node, ps *ParseState) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in recursEval:", r)
		}
	}()

	nodeName := node.GetNodeName()
	switch nodeName {
	case "#text":
		return node.(dom.Text).GetValue()
	case "#comment":
		return ""
	case "template":
		return EvalTagContent(node, ps, nil)
	case "random":
		return Random(node, ps)
	case "condition":
		return LoopCondition(node, ps)
	case "srai":
		return Srai(node, ps)
	case "sr":
		return RespondWithSrCnt(ps.StarBindings.InputStars.Star(0), ps.That, ps.Topic, ps.ChatSession, sraiCount)
	case "sraix":
		return Sraix(node, ps)
	case "set":
		return SetNode(node, ps)
	case "get":
		return GetA(node, ps)
	case "map":
		return MapNode(node, ps)
	case "bot":
		return Abot(node, ps)
	case "id":
		return Id(node, ps)
	case "size":
		return SizeA(node, ps)
	case "vocabulary":
		return Vocabulary(node, ps)
	case "program":
		return Program(node, ps)
	case "date":
		return Date(node, ps)
	case "interval":
		return Interval(node, ps)
	case "think":
		return Think(node, ps)
	case "system":
		return System(node, ps)
	case "explode":
		return ExplodeNode(node, ps)
	case "normalize":
		return Normalize(node, ps)
	case "denormalize":
		return Denormalize(node, ps)
	case "uppercase":
		return Uppercase(node, ps)
	case "lowercase":
		return Lowercase(node, ps)
	case "formal":
		return Formal(node, ps)
	case "sentence":
		return Sentence(node, ps)
	case "person":
		return Person(node, ps)
	case "person2":
		return Person2(node, ps)
	case "gender":
		return Gender(node, ps)
	case "star":
		return InputStar(node, ps)
	case "thatstar":
		return ThatStar(node, ps)
	case "topicstar":
		return TopicStar(node, ps)
	case "that":
		return That(node, ps)
	case "input":
		return Input(node, ps)
	case "request":
		return Request(node, ps)
	case "response":
		return Response(node, ps)
	case "learn", "learnf":
		return Learn(node, ps)
	case "addtriple":
		return AddTriple(node, ps)
	case "deletetriple":
		return DeleteTriple(node, ps)
	case "javascript":
		return Javascript(node, ps)
	case "select":
		return Select(node, ps)
	case "uniq":
		return Uniq(node, ps)
	case "first":
		return First(node, ps)
	case "rest":
		return Rest(node, ps)
	case "resetlearnf":
		return Resetlearnf(node, ps)
	case "resetlearn":
		return Resetlearn(node, ps)
	default:
		if Extension != nil {
			if _, ok := Extension.ExtensionTagSet()[nodeName]; ok {
				return Extension.RecursEval(node, ps)
			}
		}
		return GenericXML(node, ps)
	}
}

func EvalTemplate(template string, ps *ParseState) string {
	response := "template_failed"
	template = "<template>" + template + "</template>"
	root, err := ParseString(template)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return response
	}
	response = RecursEval(root, ps)
	return response
}

// validTemplate checks if the template is valid
func ValidTemplate(template string) bool {
	fmt.Println("AIMLProcessor.validTemplate(template: " + template + ")")
	template = "<template>" + template + "</template>"
	_, err := ParseString(template)
	if err != nil {
		fmt.Println("Invalid Template", template)
		return false
	}
	return true
}
