package ab

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type AIMLSet struct {
	Name       string
	MaxLength  int
	Host       string
	Botid      string
	IsExternal bool
	Bot        *Bot
	iSet       map[string]bool
	Set        []string
	InCache    map[string]struct{}
	OutCache   map[string]struct{}
}

func NewAIMLSet(name string, bot *Bot) *AIMLSet {
	set := &AIMLSet{
		Name:      strings.ToLower(name),
		MaxLength: 1,
		Bot:       bot,
		iSet:      map[string]bool{},
		Set:       []string{},
		InCache:   make(map[string]struct{}),
		OutCache:  make(map[string]struct{}),
	}

	if set.Name == NaturalNumberSetName {
		set.MaxLength = 1
	}

	return set
}

// add adds an item to the Set
func (a *AIMLSet) Add(item string) {
	if _, ok := a.iSet[item]; ok {
		return
	}
	a.iSet[item] = true
	a.Set = append(a.Set, item)
}
func removeByValue[T comparable](l []T, item T) []T {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func (a *AIMLSet) Remove(item string) {
	if _, ok := a.iSet[item]; !ok {
		return
	}
	delete(a.iSet, item)
	a.Set = removeByValue(a.Set, item)
}

func (s *AIMLSet) Contains(item string) bool {
	if s.IsExternal && EnableExternalSets {
		if _, ok := s.InCache[item]; ok {
			return true
		}
		if _, ok := s.OutCache[item]; ok {
			return false
		}
		split := strings.Split(item, " ")
		if len(split) > s.MaxLength {
			return false
		}
		query := fmt.Sprintf("%s%s %s", SetMemberString, strings.ToUpper(s.Name), item)
		response := SraixSraix(nil, query, "false", "", s.Host, s.Botid, "", "0")
		if response == "true" {
			s.InCache[item] = struct{}{}
			return true
		} else {
			s.OutCache[item] = struct{}{}
			return false
		}
	} else if s.Name == NaturalNumberSetName {
		numberPattern := regexp.MustCompile("[0-9]+")
		return numberPattern.MatchString(item)
	} else {
		_, ok := s.iSet[item]
		return ok
	}
}

func (s *AIMLSet) WriteAIMLSet() {
	fmt.Println("Writing AIML Set", s.Name)
	filename := fmt.Sprintf("%s/%s.txt", s.Bot.SetsPath, s.Name)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, item := range s.Set {
		fmt.Fprintln(writer, strings.TrimSpace(item))
	}
	writer.Flush()
}

func (s *AIMLSet) ReadAIMLSetFromInputStream(file *os.File, bot *Bot) int {
	scanner := bufio.NewScanner(file)
	cnt := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			cnt++
			if strings.HasPrefix(line, "external") {
				splitLine := strings.Split(line, ":")
				if len(splitLine) >= 4 {
					s.Host = splitLine[1]
					s.Botid = splitLine[2]
					s.MaxLength = atoi(splitLine[3])
					s.IsExternal = true
					fmt.Println("Created external set at", s.Host, s.Botid)
				}
			} else {
				line = strings.ToUpper(strings.TrimSpace(line))
				splitLine := strings.Split(line, " ")
				length := len(splitLine)
				if length > s.MaxLength {
					s.MaxLength = length
				}
				k := strings.TrimSpace(line)
				if _, ok := s.iSet[k]; !ok {
					s.iSet[k] = true
					s.Set = append(s.Set, k)
				}
			}
		}
	}
	return cnt
}
func (s *AIMLSet) ReadAIMLSet(bot *Bot) int {
	fmt.Println("Reading AIML Set", s.Name)
	filename := fmt.Sprintf("%s/%s.txt", s.Bot.SetsPath, s.Name)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(filename, "not found")
		return 0
	}
	defer file.Close()
	cnt := s.ReadAIMLSetFromInputStream(file, bot)
	return cnt
}

func atoi(s string) int {
	val := 0
	for _, c := range s {
		val = val*10 + int(c-'0')
	}
	return val
}
