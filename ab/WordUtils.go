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

func SplitWords(s string) []string {
	re, err := regexp.Compile(`\p{Han}*`)
	if err != nil {
		fmt.Println(err)
		return strings.Split(s, " ")
	}
	if !re.MatchString(s) {
		return strings.Split(s, " ")
	}
	words := jieba.Cut(s, true)
	return words
}
