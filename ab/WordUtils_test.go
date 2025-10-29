package ab

import (
	"fmt"
	"testing"
)

func testAliceFunctions(t *testing.T) {
	// 创建分词器实例
	tokenizer := NewTokenizerUtil()

	// 测试中英文混合文本
	text := "自然语言处理(Natural Language Processing)是人工智能的重要方向。Hello世界！"

	// 分词
	tokens := tokenizer.Tokenize(text)
	fmt.Printf("分词结果: %v\n", tokens)

	// 获取分词数量
	count := tokenizer.GetTokenCount(text)
	fmt.Printf("分词数量: %d\n", count)

	// 检查是否包含中文
	hasChinese := tokenizer.ContainsChinese(text)
	fmt.Printf("包含中文: %t\n", hasChinese)

	// 使用简单分词（不依赖外部库）
	simpleTokens := tokenizer.SimpleTokenize(text)
	fmt.Printf("简单分词结果: %v\n", simpleTokens)

	// 测试Alice机器人相关的字符串处理
	testAliceFunctions2(tokenizer)
}

// 模拟Alice机器人中的相关函数
func testAliceFunctions2(tokenizer *TokenizerUtil) {
	fmt.Println("\n=== Alice机器人功能测试 ===")

	// 模拟contains函数中的分词逻辑
	testString := "你好世界 hello world"
	tokens := tokenizer.Tokenize(testString)
	fmt.Printf("contains函数分词: %v, 数量: %d\n", tokens, len(tokens))

	// 模拟sentenceToPath函数
	sentence := "这是一个测试句子 this is a test sentence"
	pathTokens := tokenizer.Tokenize(sentence)
	fmt.Printf("sentenceToPath分词: %v\n", pathTokens)

	// 模拟instantiateSets函数
	pattern := "<SET>测试集合 这是一个模式"
	patternTokens := tokenizer.Tokenize(pattern)
	fmt.Printf("instantiateSets分词: %v\n", patternTokens)
}
