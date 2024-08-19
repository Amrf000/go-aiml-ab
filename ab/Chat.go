package ab

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Chat struct {
	Bot             *Bot
	DoWrites        bool
	CustomerID      string
	ThatHistory     *History //[*History]
	RequestHistory  *History //[string]
	ResponseHistory *History //[string]
	InputHistory    *History //[string]
	Predicates      Predicates
	TripleStore     *TripleStore
}

var MatchTrace string = ""
var LocationKnown = false
var Longitude string
var Latitude string

func NewChat(bot *Bot) *Chat {
	return NewChatWithOptions(bot, true, "0")
}
func NewChatWithBool(bot *Bot, doWrites bool) *Chat {
	return NewChatWithOptions(bot, doWrites, "0")
}

// customerId = MagicStrings.default_Customer_id;
func NewChatWithOptions(bot *Bot, doWrites bool, customerId string) *Chat {
	chat := Chat{Bot: bot, Predicates: Predicates{}}
	chat.TripleStore = NewTripleStore("anon", &chat)
	chat.ThatHistory = NewHistoryWithName("that")
	chat.RequestHistory = NewHistoryWithName("request")
	chat.ResponseHistory = NewHistoryWithName("response")
	chat.InputHistory = NewHistoryWithName("input")
	chat.DoWrites = doWrites
	chat.CustomerID = customerId
	chat.AddPredicates()
	chat.AddTriples()
	chat.Predicates["topic"] = DefaultTopic
	chat.Predicates["jsenabled"] = JsEnabled
	if TraceMode {
		fmt.Println("Chat Session Created for bot " + bot.Name)
	}
	return &chat
}

func (chat *Chat) AddPredicates() {
	err := chat.Predicates.GetPredicateDefaults(chat.Bot.ConfigPath + "/predicates.txt")
	if err != nil {
		fmt.Println("Error loading predicates:", err)
	}
}

func (chat *Chat) AddTriples() int {
	tripleCnt := 0
	if TraceMode {
		fmt.Println("Loading Triples from " + chat.Bot.ConfigPath + "/triples.txt")
	}
	filePath := chat.Bot.ConfigPath + "/triples.txt"
	if fileExists(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening triples file:", err)
			return tripleCnt
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			triple := strings.Split(line, ":")
			if len(triple) >= 3 {
				subject := strings.TrimSpace(triple[0])
				predicate := strings.TrimSpace(triple[1])
				object := strings.TrimSpace(triple[2])
				chat.TripleStore.AddTriple(subject, predicate, object)
				tripleCnt++
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading triples file:", err)
		}
	}

	if TraceMode {
		fmt.Println("Loaded", tripleCnt, "triples")
	}
	return tripleCnt
}

func (chat *Chat) Chat() {
	logFile := chat.Bot.LogPath + "/log_" + chat.CustomerID + ".txt"
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	request := "SET PREDICATES"
	response := chat.MultiSentenceRespond(request)
	scanner := bufio.NewScanner(os.Stdin)
	for request != "quit" {
		fmt.Print("Human: ")
		if scanner.Scan() {
			request = scanner.Text()
		}
		response = chat.MultiSentenceRespond(request)
		fmt.Println("Robot:", response)

		if _, err := writer.WriteString("Human: " + request + "\n"); err != nil {
			fmt.Println("Error writing to log file:", err)
		}
		if _, err := writer.WriteString("Robot: " + response + "\n"); err != nil {
			fmt.Println("Error writing to log file:", err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}

func (chat *Chat) Respond(input string, contextThatHistory *History) string {
	repetition := true
	for i := 0; i < RepetitionCount; i++ {
		if chat.InputHistory.Get(i) == nil || strings.ToUpper(input) != strings.ToUpper(chat.InputHistory.Get(i).(string)) {
			repetition = false
		}
	}
	if input == NullInput {
		repetition = false
	}

	chat.InputHistory.Add(input)

	if repetition {
		input = RepetitionDetected
	}

	response := Respond(input, "", chat.Predicates["topic"], chat)
	normResponse := chat.Bot.PreProcessor.Normalize(response)

	if JpTokenize {
		normResponse = TokenizeSentence(normResponse)
	}

	sentences := chat.Bot.PreProcessor.SentenceSplit(normResponse)
	for _, sentence := range sentences {
		that := sentence
		if strings.TrimSpace(that) == "" {
			that = DefaultThat
		}
		contextThatHistory.Add(that)
	}

	return strings.TrimSpace(response) + " "
}

func (chat *Chat) MultiSentenceRespond(request string) string {
	var response strings.Builder
	MatchTrace = ""

	normalized := chat.Bot.PreProcessor.Normalize(request)
	normalized = TokenizeSentence(normalized)
	sentences := chat.Bot.PreProcessor.SentenceSplit(normalized)

	contextThatHistory := NewHistoryWithName("contextThat")
	for _, sentence := range sentences {
		reply := chat.Respond(sentence, contextThatHistory)
		response.WriteString(" " + reply)
	}

	chat.RequestHistory.Add(request)
	chat.ResponseHistory.Add(response.String())
	chat.ThatHistory.Add(contextThatHistory)

	responseStr := response.String()
	responseStr = strings.ReplaceAll(responseStr, "\n", "")
	return strings.TrimSpace(responseStr)
}

func SetMatchTrace(newMatchTrace string) {
	MatchTrace = newMatchTrace
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
