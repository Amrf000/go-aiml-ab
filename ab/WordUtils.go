package ab

import (
	"fmt"
	"github.com/yanyiwu/gojieba"
	"regexp"
	"strings"
)

var jieba *gojieba.Jieba

func InitJieba() {
	jieba = gojieba.NewJieba()
}
func DeInitJieba() {
	jieba.Free()
}
func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}
func SplitWords(s string) []string {
	re, err := regexp.Compile("([\u4e00-\u9fa5]+)")
	if err != nil {
		fmt.Println(err)
		return strings.Split(s, " ")
	}
	if !re.MatchString(s) {
		return strings.Split(s, " ")
	}
	if strings.Contains(s, "陈述您的") {
		fmt.Printf("%s", s)
	}
	//return strings.Split(s, " ")
	oldWords := strings.Split(s, " ")
	newWords := []string{}
	for _, oldWord := range oldWords {
		if strings.Contains(strings.ToUpper(oldWord), "<SET>") {
			newWords = append(newWords, oldWord)
			continue
		}
		if !re.MatchString(oldWord) {
			newWords = append(newWords, oldWord)
			continue
		}
		words := jieba.Cut(oldWord, true)
		newWords = append(newWords, words...)
	}
	//result := ReplaceAllStringSubmatchFunc(re, s, func(groups []string) string {
	//	//prefix := ""
	//	//if strings.HasSuffix(groups[1], ">") {
	//	//	prefix = " "
	//	//}
	//	//suffix := ""
	//	//if strings.HasPrefix(groups[3], "<") {
	//	//	prefix = " "
	//	//}
	//
	//	return " " + strings.Join(words, " ") + " "
	//})

	//newWords := strings.Split(result, " ")
	return newWords
	//return words
}
