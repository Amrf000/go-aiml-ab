package ab

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// public class TestAB {
var sample_file = "sample.random.txt"

func TestChat(bot *Bot, doWrites, traceMode bool) {
	chatSession := NewChatWithBool(bot, doWrites)
	bot.Brain.NodeStats()
	TraceMode = traceMode
	textLine := ""
	for {
		textLine = ReadInputTextLineWithPrompt("Human")
		if textLine == "" || len(textLine) < 1 {
			textLine = NullInput
		}
		if textLine == "q" {
			os.Exit(0)
		} else if textLine == "wq" {
			bot.WriteQuit()
			os.Exit(0)
		} else if textLine == "sc" {
			sraixCache("c:/ab/data/sraixdata6.txt", chatSession)
		} else if textLine == "iqtest" {
			TestMultiSentenceRespond()
		} else if textLine == "ab" {
			TestAB(bot, sample_file)
		} else {
			request := textLine
			if TraceMode {
				his := chatSession.ThatHistory.Get(0)
				if his != nil {
					if hi, ok := his.(*History); ok {
						val := hi.Get(0)
						if val != nil {
							if rva, ok := val.(string); ok {
								fmt.Println("STATE=" + request + ":THAT=" + rva + ":TOPIC=" + chatSession.Predicates.Get("topic"))
							}
						}
					}
				}

			}
			response := chatSession.MultiSentenceRespond(request)
			for strings.Contains(response, "&lt;") {
				response = strings.Replace(response, "&lt;", "<", -1)
			}
			for strings.Contains(response, "&gt;") {
				response = strings.Replace(response, "&gt;", ">", -1)
			}
			WriteOutputTextLine("Robot", response)
		}
	}
}
func testBotChat() {
	bot := NewBotWithName("alice")
	fmt.Println(strconv.Itoa(bot.Brain.UpgradeCnt) + " brain upgrades")
	chatSession := NewChat(bot)
	request := "Hello.  How are you?  What is your name?  Tell me about yourself."
	response := chatSession.MultiSentenceRespond(request)
	fmt.Println("Human: " + request)
	fmt.Println("Robot: " + response)
}
func RunTests(bot *Bot, traceMode bool) {
	QaTestMode = true
	chatSession := NewChatWithBool(bot, false)
	bot.Brain.NodeStats()
	TraceMode = traceMode
	testInputRaw, err := os.Open(RootPath + "/data/lognormal-500.txt")
	if err != nil {
		panic(err)
	}
	defer testInputRaw.Close()
	testInput := bufio.NewReader(testInputRaw)
	testOutputRaw, err := os.Open(RootPath + "/data/lognormal-500-out.txt")
	if err != nil {
		panic(err)
	}
	defer testOutputRaw.Close()
	testOutput := bufio.NewWriter(testOutputRaw)
	textLine, err := testInput.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			log.Fatalf("read file line error: %v", err)
			return
		}
	}
	i := 1
	fmt.Print(0)
	for textLine != "" {
		if textLine == "" || len(textLine) < 1 {
			textLine = NullInput
		}
		if textLine == "q" {
			os.Exit(0)
		} else if textLine == "wq" {
			bot.WriteQuit()
			os.Exit(0)
		} else if textLine == "ab" {
			TestAB(bot, sample_file)
		} else if textLine == NullInput {
			testOutput.WriteString("\n")
		} else if strings.HasPrefix(textLine, "#") {
			testOutput.WriteString(textLine + "\n")
		} else {
			request := textLine
			if TraceMode {
				fmt.Println("STATE=" + request + ":THAT=" + chatSession.ThatHistory.Get(0).(*History).Get(0).(string) + ":TOPIC=" + chatSession.Predicates.Get("topic"))
			}
			response := chatSession.MultiSentenceRespond(request)
			for strings.Contains(response, "&lt;") {
				response = strings.Replace(response, "&lt;", "<", -1)
			}
			for strings.Contains(response, "&gt;") {
				response = strings.Replace(response, "&gt;", ">", -1)
			}
			testOutput.WriteString("Robot: " + response + "\n")
		}
		textLine, err = testInput.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			return
		}
		fmt.Print(".")
		if i%10 == 0 {
			fmt.Print(" ")
		}
		if i%100 == 0 {
			fmt.Println("")
			fmt.Printf("%d ", i)
		}
		i++
	}
	fmt.Println("")
}

func TestAB(bot *Bot, sampleFile string) {
	TraceMode = true
	ab := NewAB(bot, sampleFile)
	ab.Ab()
	fmt.Println("Begin Pattern Suggestor Terminal Interaction")
	ab.TerminalInteraction()
}
func testShortCuts() {
}

func sraixCache(filename string, chatSession *Chat) {
	limit := 650000
	CacheSraix = true
	// try {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		strLine := scanner.Text()
		fmt.Println("Human: " + strLine)
		response := chatSession.MultiSentenceRespond(strLine)
		fmt.Println("Robot: " + response)
		count++
		if count >= limit {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
