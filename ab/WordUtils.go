package ab

import (
	"strings"
	"unicode"

	// 导入中文分词库，这里以gse为例
	"github.com/go-ego/gse"
)

// TokenizerUtil 分词工具结构体
type TokenizerUtil struct {
	segmenter *gse.Segmenter
}

// NewTokenizerUtil 创建新的分词工具实例
func NewTokenizerUtil(dictPath ...string) *TokenizerUtil {
	tokenizer := &TokenizerUtil{
		segmenter: &gse.Segmenter{},
	}

	// 加载词典
	if len(dictPath) > 0 && dictPath[0] != "" {
		tokenizer.segmenter.LoadDict(dictPath[0])
	} else {
		tokenizer.segmenter.LoadDict()
	}

	return tokenizer
}

// Tokenize 对文本进行分词，返回词汇数组
func (t *TokenizerUtil) Tokenize(text string) []string {
	if strings.TrimSpace(text) == "" {
		return []string{}
	}

	// 使用精确模式进行分词
	segments := t.segmenter.Cut(text, true)

	// 过滤空字符串和空格
	var tokens []string
	for _, word := range segments {
		trimmed := strings.TrimSpace(word)
		if trimmed != "" {
			tokens = append(tokens, trimmed)
		}
	}

	return tokens
}

// GetTokenCount 获取文本的分词数量
func (t *TokenizerUtil) GetTokenCount(text string) int {
	return len(t.Tokenize(text))
}

// ContainsChinese 检查文本是否包含中文字符
func (t *TokenizerUtil) ContainsChinese(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

// TokenizeWithPos 分词并返回词汇及其位置信息（可选功能）
func (t *TokenizerUtil) TokenizeWithPos(text string) []struct {
	Word string
	Pos  int
} {
	tokens := t.Tokenize(text)
	var result []struct {
		Word string
		Pos  int
	}

	currentPos := 0
	for _, token := range tokens {
		// 在原文中找到token的位置
		pos := strings.Index(text[currentPos:], token)
		if pos != -1 {
			actualPos := currentPos + pos
			result = append(result, struct {
				Word string
				Pos  int
			}{
				Word: token,
				Pos:  actualPos,
			})
			currentPos = actualPos + len(token)
		}
	}

	return result
}

// SimpleTokenize 简单分词（不依赖外部库，适用于简单场景）
func (t *TokenizerUtil) SimpleTokenize(text string) []string {
	if strings.TrimSpace(text) == "" {
		return []string{}
	}

	var tokens []string
	var currentToken strings.Builder

	for _, r := range text {
		// 如果是空格或标点符号，结束当前token
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			continue
		}

		// 中文字符单独成词
		if unicode.Is(unicode.Han, r) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(r))
		} else {
			// 英文字符积累到当前token
			currentToken.WriteRune(r)
		}
	}

	// 添加最后一个token
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}
