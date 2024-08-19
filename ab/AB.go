package ab

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strings"
)

type AB struct {
	ShuffleMode         bool
	SortMode            bool
	FilterAtomicMode    bool
	FilterWildMode      bool
	OfferAliceResponses bool
	Logfile             string
	RunCompletedCnt     int
	Bot                 *Bot
	Alice               *Bot
	Passed              *AIMLSet
	TestSet             *AIMLSet
	InputGraph          *Graphmaster
	PatternGraph        *Graphmaster
	DeletedGraph        *Graphmaster
	SuggestedCategories []Category
}

var Limit = 500000

func NewAB(bot *Bot, sampleFile string) *AB {
	ab := &AB{
		ShuffleMode:         true,
		SortMode:            false,
		FilterAtomicMode:    false,
		FilterWildMode:      false,
		OfferAliceResponses: true,
		Logfile:             "/data/abSampleFile",
		Bot:                 bot,
		InputGraph:          NewGraphmasterWithName(bot, "input"),
		DeletedGraph:        NewGraphmasterWithName(bot, "deleted"),
		PatternGraph:        NewGraphmasterWithName(bot, "pattern"),
		SuggestedCategories: make([]Category, 0),
		Passed:              NewAIMLSet("passed", bot),
		TestSet:             NewAIMLSet("1000", bot),
	}
	ab.ReadDeletedIFCategories()
	return ab
}

func (ab *AB) Productivity(runCompletedCnt int, timer *Timer) {
	time := timer.ElapsedTimeMins()
	fmt.Printf("Completed %d in %.2f min. Productivity %.2f cat/min\n", runCompletedCnt, time, float32(runCompletedCnt)/float32(time))
}

func (ab *AB) ReadDeletedIFCategories() {
	ab.Bot.ReadCertainIFCategories(ab.DeletedGraph, "/data/deletedAimlFile")
	if TraceMode {
		fmt.Printf("--- DELETED CATEGORIES -- read %d deleted categories\n", len(ab.DeletedGraph.GetCategories()))
	}
}

func (ab *AB) WriteDeletedIFCategories() {
	fmt.Println("--- DELETED CATEGORIES -- write")
	ab.Bot.WriteCertainIFCategories(ab.DeletedGraph, "/data/deletedAimlFile")
	fmt.Printf("--- DELETED CATEGORIES -- write %d deleted categories\n", len(ab.DeletedGraph.GetCategories()))
}

func (ab *AB) SaveCategory(pattern, template, filename string) {
	that := "*"
	topic := "*"
	c := NewCategory(0, pattern, that, topic, template, filename)
	if ok := c.Validate(); ok {
		ab.Bot.Brain.AddCategory(c)
		ab.Bot.WriteAIMLIFFiles()
		ab.RunCompletedCnt++
	} else {
		fmt.Printf("Invalid Category %s\n", c.ValidationMessage)
	}
}

func (ab *AB) DeleteCategory(c *Category) {
	c.Filename = "/data/deletedAimlFile"
	c.Template = "/data/deletedTemplate"
	ab.DeletedGraph.AddCategory(c)
	fmt.Println("--- bot.writeDeletedIFCategories()")
	ab.WriteDeletedIFCategories()
}

func (ab *AB) SkipCategory(c *Category) {
	// Implement skipCategory logic here
}

func (ab *AB) Abwq() {
	timer := NewTimer()
	timer.Start()
	ab.ClassifyInputs(ab.Logfile)
	fmt.Printf("%.2f classifying inputs\n", timer.ElapsedTimeSecs())
	ab.Bot.WriteQuit()
}

func (ab *AB) GraphInputs(filename string) {
	count := 0
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if count >= Limit {
			break
		}
		c := Category{
			Pattern:  line,
			Filename: "/data/unknownAimlFile",
		}
		node := ab.InputGraph.FindNode(&c)
		if node == nil {
			ab.InputGraph.AddCategory(&c)
			c.IncrementActivationCnt()
		} else {
			node.Category.IncrementActivationCnt()
		}
		count++
	}
}

var LeafPatternCnt = 0
var StarPatternCnt = 0

func (ab *AB) FindPatterns() {
	ab.FindPatternsTopic(ab.InputGraph.Root, "")
	fmt.Printf("%d Leaf Patterns %d Star Patterns\n", LeafPatternCnt, StarPatternCnt)
}

func (ab *AB) FindPatternsTopic(node *Nodemapper, partialPatternThatTopic string) {
	if IsLeaf(node) {
		if node.Category.GetActivationCnt() > NodeActivationCnt {
			var categoryPatternThatTopic string
			if node.ShortCut {
				categoryPatternThatTopic = partialPatternThatTopic + " <THAT> * <TOPIC> *"
			} else {
				categoryPatternThatTopic = partialPatternThatTopic
			}
			c := Category{
				Pattern:  categoryPatternThatTopic,
				Template: BlankTemplate,
				Filename: UnknownAimlFile,
			}
			if !ab.Bot.Brain.ExistsCategory(&c) && !ab.DeletedGraph.ExistsCategory(&c) {
				ab.PatternGraph.AddCategory(&c)
				ab.SuggestedCategories = append(ab.SuggestedCategories, c)
			}
		}
	}
	if len(node.Map) > NodeSize {
		StarPatternCnt++
		c := Category{
			Pattern:  partialPatternThatTopic + " * <THAT> * <TOPIC> *",
			Template: BlankTemplate,
			Filename: UnknownAimlFile,
		}
		if !ab.Bot.Brain.ExistsCategory(&c) && !ab.DeletedGraph.ExistsCategory(&c) {
			ab.PatternGraph.AddCategory(&c)
			ab.SuggestedCategories = append(ab.SuggestedCategories, c)
		}
	}
	for key := range node.Map {
		value := node.Map[key]
		ab.FindPatternsTopic(value, partialPatternThatTopic+" "+key)
	}
}

func (ab *AB) ClassifyInputs(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	lines := strings.Split(string(data), "\n")
	count := 0
	for _, line := range lines {
		if count >= Limit {
			break
		}
		if line != "" {
			if strings.HasPrefix(line, "Human: ") {
				line = line[len("Human: "):]
			}
			sentences := ab.Bot.PreProcessor.SentenceSplit(line)
			for _, sentence := range sentences {
				if len(sentence) > 0 {
					match := ab.PatternGraph.MatchRaw(sentence, "unknown", "unknown")
					if match == nil {
						fmt.Println(sentence + " null match")
					} else {
						match.Category.IncrementActivationCnt()
					}
					count++
					if count%10000 == 0 {
						fmt.Println(count)
					}
				}
			}
		}
	}
	fmt.Printf("Finished classifying %d inputs\n", count)
}

func (ab *AB) Ab() {
	logFile := ab.Logfile
	TraceMode = false
	EnableExternalSets = false
	if ab.OfferAliceResponses {
		ab.Alice = &Bot{}
	}
	timer := NewTimer()
	ab.Bot.Brain.NodeStats()
	if len(ab.Bot.Brain.GetCategories()) < BrainPrintSize {
		ab.Bot.Brain.Printgraph()
	}
	timer.Start()
	fmt.Println("Graphing inputs")
	ab.GraphInputs(logFile)
	fmt.Printf("%.2f seconds Graphing inputs\n", timer.ElapsedTimeSecs())
	ab.InputGraph.NodeStats()
	if len(ab.InputGraph.GetCategories()) < BrainPrintSize {
		ab.InputGraph.Printgraph()
	}
	timer.Start()
	fmt.Println("Finding Patterns")
	ab.FindPatterns()
	fmt.Printf("%d suggested categories\n", len(ab.SuggestedCategories))
	fmt.Printf("%.2f seconds finding patterns\n", timer.ElapsedTimeSecs())
	timer.Start()
	ab.PatternGraph.NodeStats()
	if len(ab.PatternGraph.GetCategories()) < BrainPrintSize {
		ab.PatternGraph.Printgraph()
	}
	fmt.Println("Classifying Inputs from " + logFile)
	ab.ClassifyInputs(logFile)
	fmt.Printf("%.2f classifying inputs\n", timer.ElapsedTimeSecs())
}

func (ab *AB) NonZeroActivationCount(suggestedCategories []Category) []Category {
	result := make([]Category, 0)
	for _, c := range suggestedCategories {
		if c.ActivationCnt > 0 {
			result = append(result, c)
		}
	}
	return result
}

func (ab *AB) TerminalInteraction() {
	firstInteraction := true
	var alicetemplate string
	timer := NewTimer()
	ab.SortMode = !ab.ShuffleMode
	sort.Slice(ab.SuggestedCategories, func(i, j int) bool {
		return ab.SuggestedCategories[i].ActivationCnt > ab.SuggestedCategories[j].ActivationCnt
	})
	topSuggestCategories := make([]Category, 0)
	for i := 0; i < 10000 && i < len(ab.SuggestedCategories); i++ {
		topSuggestCategories = append(topSuggestCategories, ab.SuggestedCategories[i])
	}
	ab.SuggestedCategories = topSuggestCategories
	if ab.ShuffleMode {
		rand.Shuffle(len(ab.SuggestedCategories), func(i, j int) {
			ab.SuggestedCategories[i], ab.SuggestedCategories[j] = ab.SuggestedCategories[j], ab.SuggestedCategories[i]
		})
	}
	timer = NewTimer()
	timer.Start()
	ab.RunCompletedCnt = 0
	filteredAtomicCategories := make([]Category, 0)
	filteredWildCategories := make([]Category, 0)
	for _, c := range ab.SuggestedCategories {
		if !strings.Contains(c.Pattern, "*") {
			filteredAtomicCategories = append(filteredAtomicCategories, c)
		} else {
			filteredWildCategories = append(filteredWildCategories, c)
		}
	}
	var browserCategories []Category
	if ab.FilterAtomicMode {
		browserCategories = filteredAtomicCategories
	} else if ab.FilterWildMode {
		browserCategories = filteredWildCategories
	} else {
		browserCategories = ab.SuggestedCategories
	}
	browserCategories = ab.NonZeroActivationCount(browserCategories)
	for _, c := range browserCategories {
		samplesRaw := c.GetMatches(ab.Bot).Set
		samples := make([]string, len(samplesRaw))
		i := 0
		for _, k := range samplesRaw {
			samples[i] = k
			i++
		}
		rand.Shuffle(len(samples), func(i, j int) {
			samples[i], samples[j] = samples[j], samples[i]
		})
		sampleSize := min(DisplayedInputSampleSize, len(samplesRaw))
		for i := 0; i < sampleSize; i++ {
			fmt.Println("" + samples[i])
		}
		fmt.Printf("[%d] %s\n", c.ActivationCnt, c.InputThatTopic())
		var node *Nodemapper
		if ab.OfferAliceResponses {
			node = ab.Alice.Brain.FindNode(&c)
			if node != nil {
				alicetemplate = node.Category.Template
				displayAliceTemplate := alicetemplate
				displayAliceTemplate = strings.Replace(displayAliceTemplate, "\n", " ", -1)
				if len(displayAliceTemplate) > 200 {
					displayAliceTemplate = displayAliceTemplate[:200]
				}
				fmt.Println("ALICE: " + displayAliceTemplate)
			} else {
				alicetemplate = ""
			}
		}
		textLine := ReadInputTextLine()
		if firstInteraction {
			timer.Start()
			firstInteraction = false
		}
		ab.Productivity(ab.RunCompletedCnt, timer)
		ab.TerminalInteractionStep("", textLine, c, alicetemplate)
	}
	fmt.Println("No more samples")
	ab.Bot.WriteAIMLFiles()
	ab.Bot.WriteAIMLIFFiles()
}

func (ab *AB) TerminalInteractionStep(request, textLine string, c Category, alicetemplate string) {
	var template string
	if strings.Contains(textLine, "<pattern>") && strings.Contains(textLine, "</pattern>") {
		index := strings.Index(textLine, "<pattern>") + len("<pattern>")
		jndex := strings.Index(textLine, "</pattern>")
		kndex := jndex + len("</pattern>")
		if index < jndex {
			pattern := textLine[index:jndex]
			c.Pattern = pattern
			textLine = textLine[kndex:]
			fmt.Printf("Got pattern = %s template = %s\n", pattern, textLine)
		}
	}
	var botThinks string
	pronouns := []string{"he", "she", "it", "we", "they"}
	for _, p := range pronouns {
		if strings.Contains(textLine, "<"+p+">") {
			textLine = strings.Replace(textLine, "<"+p+">", "", -1)
			botThinks = "<think><set name=\"" + p + "\"><set name=\"topic\"><star/></set></set></think>"
		}
	}
	if textLine == "q" {
		os.Exit(0)
	} else if textLine == "wq" {
		ab.Bot.WriteQuit()
		os.Exit(0)
	} else if textLine == "skip" || textLine == "" {
		ab.SkipCategory(&c)
	} else if textLine == "s" || textLine == "pass" {
		ab.Passed.Add(request)
		difference := NewAIMLSet("difference", ab.Bot)
		for _, itm := range ab.TestSet.Set {
			difference.Add(itm)
		}
		for _, itm := range ab.Passed.Set {
			difference.Remove(itm)
		}
		difference.WriteAIMLSet()
		ab.Passed.WriteAIMLSet()
	} else if textLine == "a" {
		template = alicetemplate
		var filename string
		if strings.Contains(template, "<sr") {
			filename = ReductionsUpdateAimlFile
		} else {
			filename = PersonalityAimlFile
		}
		ab.SaveCategory(c.Pattern, template, filename)
	} else if textLine == "d" {
		ab.DeleteCategory(&c)
	} else if textLine == "x" {
		template = "<sraix services=\"pannous\">" + strings.Replace(c.Pattern, "*", "<star/>", -1) + "</sraix>"
		template += botThinks
		ab.SaveCategory(c.Pattern, template, SraixAimlFile)
	} else if textLine == "p" {
		template = "<srai>" + InappropriateFilter + "</srai>"
		template += botThinks
		ab.SaveCategory(c.Pattern, template, InappropriateAimlFile)
	} else if textLine == "f" {
		template = "<srai>" + ProfanityFilter + "</srai>"
		template += botThinks
		ab.SaveCategory(c.Pattern, template, ProfanityAimlFile)
	} else if textLine == "i" {
		template = "<srai>" + InsultFilter + "</srai>"
		template += botThinks
		ab.SaveCategory(c.Pattern, template, InsultAimlFile)
	} else if strings.Contains(textLine, "<srai>") || strings.Contains(textLine, "<sr/>") {
		template = textLine
		template += botThinks
		ab.SaveCategory(c.Pattern, template, ReductionsUpdateAimlFile)
	} else if strings.Contains(textLine, "<oob>") {
		template = textLine
		template += botThinks
		ab.SaveCategory(c.Pattern, template, OobAimlFile)
	} else if strings.Contains(textLine, "<set name") || botThinks != "" {
		template = textLine
		template += botThinks
		ab.SaveCategory(c.Pattern, template, PredicatesAimlFile)
	} else if strings.Contains(textLine, "<get name") && !strings.Contains(textLine, "<get name=\"name") {
		template = textLine
		template += botThinks
		ab.SaveCategory(c.Pattern, template, PredicatesAimlFile)
	} else {
		template = textLine
		template += botThinks
		ab.SaveCategory(c.Pattern, template, PersonalityAimlFile)
	}
}

func main() {
	bot := Bot{}
	ab := NewAB(&bot, "sampleFile")
	ab.Ab()
}
