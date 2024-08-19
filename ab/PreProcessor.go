package ab

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type PreProcessor struct {
	NormalCount      int
	DenormalCount    int
	PersonCount      int
	Person2Count     int
	GenderCount      int
	NormalSubs       []string
	NormalPatterns   []*regexp.Regexp
	DenormalSubs     []string
	DenormalPatterns []*regexp.Regexp
	PersonSubs       []string
	PersonPatterns   []*regexp.Regexp
	Person2Subs      []string
	Person2Patterns  []*regexp.Regexp
	GenderSubs       []string
	GenderPatterns   []*regexp.Regexp
}

func NewPreProcessor(bot *Bot) *PreProcessor {
	pp := &PreProcessor{}
	pp.NormalCount = pp.ReadSubstitutions(bot.ConfigPath+"/normal.txt", &pp.NormalPatterns, &pp.NormalSubs)
	pp.DenormalCount = pp.ReadSubstitutions(bot.ConfigPath+"/denormal.txt", &pp.DenormalPatterns, &pp.DenormalSubs)
	pp.PersonCount = pp.ReadSubstitutions(bot.ConfigPath+"/person.txt", &pp.PersonPatterns, &pp.PersonSubs)
	pp.Person2Count = pp.ReadSubstitutions(bot.ConfigPath+"/person2.txt", &pp.Person2Patterns, &pp.Person2Subs)
	pp.GenderCount = pp.ReadSubstitutions(bot.ConfigPath+"/gender.txt", &pp.GenderPatterns, &pp.GenderSubs)
	if TraceMode {
		fmt.Println("Preprocessor:", pp.NormalCount, "norms", pp.PersonCount, "persons", pp.Person2Count, "person2")
	}
	return pp
}

func (pp *PreProcessor) Normalize(request string) string {
	result := pp.Substitute(request, pp.NormalPatterns, pp.NormalSubs, pp.NormalCount)
	result = strings.ReplaceAll(result, "\r\n", " ")
	result = strings.ReplaceAll(result, "\n\r", " ")
	result = strings.ReplaceAll(result, "\r", " ")
	result = strings.ReplaceAll(result, "\n", " ")
	return result
}

func (pp *PreProcessor) Denormalize(request string) string {
	return pp.Substitute(request, pp.DenormalPatterns, pp.DenormalSubs, pp.DenormalCount)
}

func (pp *PreProcessor) Person(input string) string {
	return pp.Substitute(input, pp.PersonPatterns, pp.PersonSubs, pp.PersonCount)
}

func (pp *PreProcessor) Person2(input string) string {
	return pp.Substitute(input, pp.Person2Patterns, pp.Person2Subs, pp.Person2Count)
}

func (pp *PreProcessor) Gender(input string) string {
	return pp.Substitute(input, pp.GenderPatterns, pp.GenderSubs, pp.GenderCount)
}

func (pp *PreProcessor) Substitute(request string, patterns []*regexp.Regexp, subs []string, count int) string {
	result := " " + request + " "
	for i := 0; i < count; i++ {
		replacement := subs[i]
		p := patterns[i]
		result = p.ReplaceAllString(result, replacement)
	}
	result = strings.ReplaceAll(result, "  ", " ")
	return strings.TrimSpace(result)
}

func (pp *PreProcessor) ReadSubstitutionsFromInputStream(in *os.File, patterns *[]*regexp.Regexp, subs *[]string) int {
	scanner := bufio.NewScanner(in)
	subCount := 0
	for scanner.Scan() {
		strLine := scanner.Text()
		strLine = strings.TrimSpace(strLine)
		if !strings.HasPrefix(strLine, TextCommentMark) {
			pattern := regexp.MustCompile(`"(.*?)","(.*?)"`)
			matches := pattern.FindStringSubmatch(strLine)
			if len(matches) == 3 && subCount < MaxSubstitutions {
				*subs = append(*subs, matches[2])
				quotedPattern := regexp.QuoteMeta(matches[1])
				pattern := regexp.MustCompile(quotedPattern)
				*patterns = append(*patterns, pattern)
				subCount++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
	return subCount
}

func (pp *PreProcessor) ReadSubstitutions(filename string, patterns *[]*regexp.Regexp, subs *[]string) int {
	subCount := 0
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return subCount
	}
	defer file.Close()

	subCount = pp.ReadSubstitutionsFromInputStream(file, patterns, subs)
	return subCount
}

func (pp *PreProcessor) SentenceSplit(line string) []string {
	line = strings.ReplaceAll(line, "。", ".")
	line = strings.ReplaceAll(line, "？", "?")
	line = strings.ReplaceAll(line, "！", "!")
	sentences := strings.Split(line, "[\\.\\!\\?]")
	for i := range sentences {
		sentences[i] = strings.TrimSpace(sentences[i])
	}
	return sentences
}

func (pp *PreProcessor) NormalizeFile(infile, outfile string) {
	fin, err := os.Open(infile)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer fin.Close()

	fout, err := os.Create(outfile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer fout.Close()

	scanner := bufio.NewScanner(fin)
	writer := bufio.NewWriter(fout)
	for scanner.Scan() {
		strLine := scanner.Text()
		strLine = strings.TrimSpace(strLine)
		if len(strLine) > 0 {
			norm := pp.Normalize(strLine)
			norm = strings.ToUpper(norm)
			sentences := pp.SentenceSplit(norm)
			for _, s := range sentences {
				fmt.Println(norm + "-->" + s)
				if len(s) > 0 {
					_, err := writer.WriteString(s + "\n")
					if err != nil {
						fmt.Println("Error writing to output file:", err)
						return
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input file:", err)
	}
	writer.Flush()
}
