package ab

import (
	"fmt"
	"strconv"
	"strings"
)

type Category struct {
	Pattern           string
	That              string
	Topic             string
	Template          string
	Filename          string
	ActivationCnt     int
	CategoryNumber    int      // Note: This corresponds to categoryCnt in Java
	Matches           *AIMLSet // Assuming AIMLSet is a type that handles set operations
	ValidationMessage string
}

var CategoryCnt = 0

func NewCategory(activationCnt int, pattern, that, topic, template, filename string) *Category {
	if FixExcelCsv {
		pattern = FixCSV(pattern)
		that = FixCSV(that)
		topic = FixCSV(topic)
		template = FixCSV(template)
		filename = FixCSV(filename)
	}

	//if strings.HasPrefix(pattern, "<pattern") {
	//	nn := 0
	//	nn++
	//}

	return &Category{
		Pattern:        strings.ToUpper(strings.TrimSpace(pattern)),
		That:           strings.ToUpper(strings.TrimSpace(that)),
		Topic:          strings.ToUpper(strings.TrimSpace(topic)),
		Template:       strings.ReplaceAll(template, "& ", " and "),
		Filename:       filename,
		ActivationCnt:  activationCnt,
		CategoryNumber: CategoryCnt,
	}
}

func NewCategoryFromPatternThatTopic(patternThatTopic, template, filename string, activationCnt int) *Category {
	thatStart := strings.Index(patternThatTopic, "<THAT>")
	topicStart := strings.Index(patternThatTopic, "<TOPIC>")
	if thatStart == -1 || topicStart == -1 {
		return nil // Error handling if <THAT> or <TOPIC> not found
	}
	that := patternThatTopic[thatStart+len("<THAT>") : topicStart]
	topic := patternThatTopic[topicStart+len("<TOPIC>"):]
	pattern := patternThatTopic[:thatStart]

	return NewCategory(activationCnt, pattern, that, topic, template, filename)
}

func (c *Category) GetMatches(bot *Bot) *AIMLSet {
	if c.Matches != nil {
		return c.Matches
	}
	return NewAIMLSet("No Matches", bot)
}

func (c *Category) InputThatTopic() string {
	return InputThatTopic(c.Pattern, c.That, c.Topic) // Assuming GraphmasterInputThatTopic function exists
}

func (c *Category) AddMatch(input string, bot *Bot) {
	if c.Matches == nil {
		setName := strings.ReplaceAll(c.InputThatTopic(), "*", "STAR")
		setName = strings.ReplaceAll(setName, "_", "UNDERSCORE")
		setName = strings.ReplaceAll(setName, " ", "-")
		setName = strings.ReplaceAll(setName, "<THAT>", "THAT")
		setName = strings.ReplaceAll(setName, "<TOPIC>", "TOPIC")
		c.Matches = NewAIMLSet(setName, bot)
	}
	c.Matches.Add(input)
}

func TemplateToLine(template string) string {
	result := template
	result = strings.ReplaceAll(result, "\r\n", "\\#Newline")
	result = strings.ReplaceAll(result, AimlifSplitChar, AimlifSplitCharName)
	return result
}

func LineToTemplate(line string) string {
	result := strings.ReplaceAll(line, "\\#Newline", "\n")
	result = strings.ReplaceAll(result, AimlifSplitCharName, AimlifSplitChar)
	return result
}

func CategoryToIF(category *Category) string {
	c := AimlifSplitChar
	return fmt.Sprintf("%d%s%s%s%s%s%s", category.ActivationCnt, c, category.Pattern, c, category.That, c,
		category.Topic, c, TemplateToLine(category.Template), c, category.Filename)
}

func CategoryToAIML(category *Category) string {
	topicStart, topicEnd := "", ""
	thatStatement := ""
	result := ""
	pattern := category.Pattern
	if strings.Contains(pattern, "<SET>") || strings.Contains(pattern, "<BOT") {
		splitPattern := strings.Split(pattern, " ")
		var rpattern strings.Builder
		for _, w := range splitPattern {
			if strings.HasPrefix(w, "<SET>") || strings.HasPrefix(w, "<BOT") || strings.HasPrefix(w, "NAME=") {
				w = strings.ToLower(w)
			}
			rpattern.WriteString(" " + w)
		}
		pattern = strings.TrimSpace(rpattern.String())
	}
	NL := "\n"
	if !strings.EqualFold(category.Topic, "*") {
		topicStart = fmt.Sprintf("<topic name=\"%s\">%s", category.Topic, NL)
		topicEnd = fmt.Sprintf("</topic>%s", NL)
	}
	if !strings.EqualFold(category.That, "*") {
		thatStatement = fmt.Sprintf("<that>%s</that>", category.That)
	}
	result = fmt.Sprintf("%s<category><pattern>%s</pattern>%s%s<template>%s</template>%s</category>%s",
		topicStart, pattern, thatStatement, NL, category.Template, NL, topicEnd)
	return result
}

func (c *Category) Validate() bool {
	validationMessage := ""
	if !c.ValidPatternForm(c.Pattern) {
		validationMessage += "Badly formatted <pattern> "
		return false
	}
	if !c.ValidPatternForm(c.That) {
		validationMessage += "Badly formatted <that> "
		return false
	}
	if !c.ValidPatternForm(c.Topic) {
		validationMessage += "Badly formatted <topic> "
		return false
	}
	if !ValidTemplate(c.Template) {
		validationMessage += "Badly formatted <template> "
		return false
	}
	if !strings.HasSuffix(c.Filename, ".aiml") {
		validationMessage += "Filename suffix should be .aiml "
		return false
	}
	return true
}

func (c *Category) ValidPatternForm(pattern string) bool {
	if len(pattern) < 1 {
		return false
	}
	words := strings.Split(pattern, " ")
	for _, word := range words {
		_ = word // Process each word as needed
		// Add validation logic if necessary
	}
	return true
}

func (c *Category) IncrementActivationCnt() {
	c.ActivationCnt++
}

func (c *Category) SetActivationCnt(cnt int) {
	c.ActivationCnt = cnt
}

func (c *Category) SetFilename(filename string) {
	c.Filename = filename
}

func (c *Category) SetTemplate(template string) {
	c.Template = template
}

func (c *Category) SetPattern(pattern string) {
	c.Pattern = strings.ToUpper(pattern)
}

func (c *Category) SetThat(that string) {
	c.That = strings.ToUpper(that)
}

func (c *Category) SetTopic(topic string) {
	c.Topic = strings.ToUpper(topic)
}

func (c *Category) GetActivationCnt() int {
	return c.ActivationCnt
}

func (c *Category) GetCategoryNumber() int {
	return c.CategoryNumber
}

func (c *Category) GetPattern() string {
	if c.Pattern == "" {
		return "*"
	}
	return c.Pattern
}

func (c *Category) GetThat() string {
	if c.That == "" {
		return "*"
	}
	return c.That
}

func (c *Category) GetTopic() string {
	if c.Topic == "" {
		return "*"
	}
	return c.Topic
}

func (c *Category) GetTemplate() string {
	return c.Template
}

func (c *Category) GetFilename() string {
	if c.Filename == "" {
		return UnknownAimlFile
	}
	return c.Filename
}

func IFToCategory(IF string) *Category {
	split := strings.Split(IF, AimlifSplitChar)
	//System.out.println("Read: "+split);
	ii, _ := strconv.Atoi(split[0])
	return NewCategory(ii, split[1], split[2], split[3], LineToTemplate(split[4]), split[5])
}

// Comparator functions
type ByActivationCnt []*Category

func (a ByActivationCnt) Len() int           { return len(a) }
func (a ByActivationCnt) Less(i, j int) bool { return a[i].ActivationCnt > a[j].ActivationCnt }
func (a ByActivationCnt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByPattern []*Category

func (a ByPattern) Len() int { return len(a) }
func (a ByPattern) Less(i, j int) bool {
	return strings.ToLower(a[i].InputThatTopic()) < strings.ToLower(a[j].InputThatTopic())
}
func (a ByPattern) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type ByCategoryNumber []*Category

func (a ByCategoryNumber) Len() int           { return len(a) }
func (a ByCategoryNumber) Less(i, j int) bool { return a[i].CategoryNumber < a[j].CategoryNumber }
func (a ByCategoryNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
