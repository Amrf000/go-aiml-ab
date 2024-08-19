package main

import (
	"aiml/ab"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	ab.SetRootPathFromSystem()
	ab.Extension = ab.NewPCAIMLProcessorExtension()
	MainFunction(os.Args)
}

func MainFunction(args []string) {
	botName := "alice2"
	ab.JpTokenize = false
	traceMode := true
	action := "chat"
	fmt.Println(ab.ProgramNameVersion)

	for _, s := range args {
		splitArg := strings.Split(s, "=")
		if len(splitArg) >= 2 {
			option := splitArg[0]
			value := splitArg[1]
			if option == "bot" {
				botName = value
			}
			if option == "action" {
				action = value
			}
			if option == "trace" {
				if value == "true" {
					traceMode = true
				} else {
					traceMode = false
				}
			}
			if option == "morph" {
				if value == "true" {
					ab.JpTokenize = true
				} else {
					ab.JpTokenize = false
				}
			}
		}
	}

	if traceMode {
		fmt.Println("Working Directory = " + ab.RootPath)
	}
	ab.EnableShortCuts = true
	bot := ab.NewBotWithAction(botName, ab.RootPath, action)
	if ab.MakeVerbsSetsMapsFlag {
		ab.MakeVerbSetsMaps(bot)
	}
	if len(bot.Brain.GetCategories()) < ab.BrainPrintSize {
		bot.Brain.Printgraph()
	}
	if traceMode {
		fmt.Println("Action = '" + action + "'")
	}
	switch action {
	case "chat", "chat-app":
		doWrites := action != "chat-app"
		ab.TestChat(bot, doWrites, traceMode)
	case "ab":
		ab.TestAB(bot, ab.AbSampleFile)
	case "aiml2csv", "csv2aiml":
		Convert(bot, action)
	case "abwq":
		ab := ab.NewAB(bot, ab.AbSampleFile)
		ab.Abwq()
	case "test":
		ab.RunTests(bot, traceMode)
	case "shadow":
		traceMode = false
		bot.ShadowChecker()
	default:
		fmt.Println("Unrecognized action " + action)
	}
}

func Convert(bot *ab.Bot, action string) {
	if action == "aiml2csv" {
		bot.WriteAIMLIFFiles()
	} else if action == "csv2aiml" {
		bot.WriteAIMLFiles()
	}
}

func GetGloss(bot *ab.Bot, filename string) {
	fmt.Println("getGloss")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	defer file.Close()

	GetGlossFromInputStream(bot, file)
}

func GetGlossFromInputStream(bot *ab.Bot, in *os.File) {
	fmt.Println("getGlossFromInputStream")
	scanner := bufio.NewScanner(in)
	def := make(map[string]string)
	var word, gloss string
	for scanner.Scan() {
		strLine := scanner.Text()
		if strings.Contains(strLine, "<entry word") {
			start := strings.Index(strLine, "<entry word=\"") + len("<entry word=\"")
			end := strings.Index(strLine, "#")
			word = strLine[start:end]
			word = strings.ReplaceAll(word, "_", " ")
			fmt.Println(word)
		} else if strings.Contains(strLine, "<gloss>") {
			gloss = strings.ReplaceAll(strLine, "<gloss>", "")
			gloss = strings.ReplaceAll(gloss, "</gloss>", "")
			gloss = strings.TrimSpace(gloss)
			fmt.Println(gloss)
		}
		if word != "" && gloss != "" {
			word = strings.ToLower(strings.TrimSpace(word))
			if len(gloss) > 2 {
				gloss = strings.ToUpper(string(gloss[0])) + gloss[1:]
			}
			definition, exists := def[word]
			if exists {
				definition = definition + "; " + gloss
			} else {
				definition = gloss
			}
			def[word] = definition
			word = ""
			gloss = ""
		}
	}

	d := ab.NewCategory(0, "WNDEF *", "*", "*", "unknown", "wndefs0.aiml")
	bot.Brain.AddCategory(d)
	for word, gloss := range def {
		gloss += "."
		c := ab.NewCategory(0, "WNDEF "+word, "*", "*", gloss, "wndefs0.aiml")
		fmt.Println(c.InputThatTopic() + ":" + c.GetTemplate() + ":" + c.GetFilename())
		node := bot.Brain.FindNode(c)
		if node != nil {
			node.Category.SetTemplate(node.Category.GetTemplate() + "," + gloss)
		}
		bot.Brain.AddCategory(c)
	}
}

func SraixCache(filename string, chatSession ab.Chat) {
	limit := 1000
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() && count < limit {
		strLine := scanner.Text()
		fmt.Println("Human: " + strLine)
		response := chatSession.MultiSentenceRespond(strLine)
		fmt.Println("Robot: " + response)
		count++
	}
}
