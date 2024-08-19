package ab

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AIMLMap struct {
	MapName    string
	Host       string
	Botid      string
	IsExternal bool
	Inflector  *Inflector
	Bot        *Bot
	obj        map[string]string
}

func NewAIMLMap(name string, bot *Bot) *AIMLMap {
	return &AIMLMap{
		MapName:   name,
		Bot:       bot,
		Inflector: GetInstance(),
		obj:       make(map[string]string),
	}
}

func (am *AIMLMap) Get(key string) string {
	var value string

	switch am.MapName {
	case "successor":
		number, err := strconv.Atoi(key)
		if err == nil {
			return strconv.Itoa(number + 1)
		}
		return DefaultMap
	case "predecessor":
		number, err := strconv.Atoi(key)
		if err == nil {
			return strconv.Itoa(number - 1)
		}
		return DefaultMap
	case "singular":
		return strings.ToLower(am.Inflector.Singularize(key))
	case "plural":
		return strings.ToLower(am.Inflector.Pluralize(key))
	default:
		if am.IsExternal && EnableExternalSets {
			query := strings.ToUpper(fmt.Sprintf("%s %s", am.MapName, key))
			response := SraixSraix(nil, query, DefaultMap, "", am.Host, am.Botid, "", "0")
			fmt.Printf("External %s(%s)=%s\n", am.MapName, key, response)
			value = response
		} else {
			value = am.Get(key)
		}
		if value == "" {
			value = DefaultMap
		}
		return value
	}
}

func (am *AIMLMap) Put(key, value string) {
	am.obj[key] = value
}

func (am *AIMLMap) WriteAIMLMap() {
	fmt.Printf("Writing AIML Map %s\n", am.MapName)
	filePath := fmt.Sprintf("%s/%s.txt", am.Bot.MapsPath, am.MapName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for p, v := range am.obj {
		p = strings.TrimSpace(p)
		_, err := writer.WriteString(fmt.Sprintf("%s:%s\n", p, strings.TrimSpace(v)))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			writer.Flush()
			return
		}
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
	}
}

func (am *AIMLMap) ReadAIMLMapFromInputStream(file *os.File, bot *Bot) int {
	cnt := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			splitLine := strings.Split(line, ":")
			if len(splitLine) >= 2 {
				cnt++
				if strings.HasPrefix(line, RemoteMapKey) && len(splitLine) >= 3 {
					am.Host = splitLine[1]
					am.Botid = splitLine[2]
					am.IsExternal = true
					fmt.Printf("Created external map at %s %s\n", am.Host, am.Botid)
				} else {
					key := strings.ToUpper(splitLine[0])
					value := splitLine[1]
					am.Put(key, value)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return cnt
}

func (am *AIMLMap) ReadAIMLMap(bot *Bot) int {
	cnt := 0
	fmt.Printf("Reading AIML Map %s/%s.txt\n", bot.MapsPath, am.MapName)
	filePath := fmt.Sprintf("%s/%s.txt", bot.MapsPath, am.MapName)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("%s not found\n", filePath)
		return cnt
	}
	defer file.Close()
	cnt = am.ReadAIMLMapFromInputStream(file, bot)
	return cnt
}
