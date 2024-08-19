package ab

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// fixCSV removes trailing semicolons and handles quoted CSV values
func FixCSV(line string) string {
	for strings.HasSuffix(line, ";") {
		line = line[:len(line)-1]
	}
	if strings.HasPrefix(line, "\"") {
		line = line[1:]
	}
	if strings.HasSuffix(line, "\"") {
		line = line[:len(line)-1]
	}
	line = strings.ReplaceAll(line, "\"\"", "\"")
	return line
}

// tagTrim removes XML tags from xmlExpression
func TagTrim(xmlExpression, tagName string) string {
	stag := "<" + tagName + ">"
	etag := "</" + tagName + ">"
	if len(xmlExpression) >= len(stag+etag) {
		xmlExpression = xmlExpression[len(stag):]
		xmlExpression = xmlExpression[:len(xmlExpression)-len(etag)]
	}
	return xmlExpression
}

// stringSet creates a HashSet of strings from variadic input
func StringSet(strings ...string) map[string]bool {
	set := make(map[string]bool)
	for _, s := range strings {
		set[s] = true
	}
	return set
}

// getFileFromInputStream reads from an input stream and returns contents as a string
func GetFileFromInputStream(in *os.File) string {
	scanner := bufio.NewScanner(in)
	var contents string
	for scanner.Scan() {
		strLine := scanner.Text()
		if !strings.HasPrefix(strLine, TextCommentMark) {
			if len(strLine) == 0 {
				contents += "\n"
			} else {
				contents += strLine + "\n"
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return strings.TrimSpace(contents)
}

// getFile reads contents of a file and returns as a string
func GetFile(filename string) string {
	var contents string
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return contents
	}
	defer file.Close()

	contents = GetFileFromInputStream(file)
	return contents
}

// getCopyrightFromInputStream reads copyright information from an input stream
func GetCopyrightFromInputStream(in *os.File, bot *Bot) string {
	scanner := bufio.NewScanner(in)
	var copyright string
	for scanner.Scan() {
		strLine := scanner.Text()
		if len(strLine) == 0 {
			copyright += "\n"
		} else {
			copyright += "<!-- " + strLine + " -->\n"
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return copyright
}

// getCopyright generates copyright information based on bot and filename
func GetCopyright(bot *Bot, AIMLFilename string) string {
	var copyright string
	year := Utils_Year()
	date := Utils_Date()

	copyright = GetFile(bot.ConfigPath + "/copyright.txt")
	splitCopyright := strings.Split(strings.TrimSpace(copyright), "\n")
	copyright = ""
	for _, line := range splitCopyright {
		copyright += "<!-- " + line + " -->\n"
	}
	replacements := map[string]string{
		"[url]":          bot.Properties.Get("url"),
		"[date]":         date,
		"[YYYY]":         year,
		"[version]":      bot.Properties.Get("version"),
		"[botname]":      strings.ToUpper(bot.Name),
		"[filename]":     AIMLFilename,
		"[botmaster]":    bot.Properties.Get("botmaster"),
		"[organization]": bot.Properties.Get("organization"),
	}
	for key, value := range replacements {
		copyright = strings.ReplaceAll(copyright, key, value)
	}
	copyright += "<!--  -->\n"
	return copyright
}

// getPannousAPIKey retrieves Pannous API key from file or default
func GetPannousAPIKey(bot *Bot) string {
	apiKey := GetFile(bot.ConfigPath + "/pannous-apikey.txt")
	if apiKey == "" {
		apiKey = PannousApiKey
	}
	return apiKey
}

// getPannousLogin retrieves Pannous login information from file or default
func GetPannousLogin(bot *Bot) string {
	login := GetFile(bot.ConfigPath + "/pannous-login.txt")
	if login == "" {
		login = PannousLogin
	}
	return login
}

// isCharCJK checks if a character is in CJK Unicode block
func IsCharCJK(c rune) bool {
	switch {
	case c >= '\u4e00' && c <= '\u9fff': // CJK Unified Ideographs
		return true
	case c >= '\u3400' && c <= '\u4dbf': // CJK Unified Ideographs Extension A
		return true
	case c >= '\uf900' && c <= '\ufaff': // CJK Compatibility Ideographs
		return true
	case c >= '\u2e80' && c <= '\u2eff': // CJK Radicals Supplement
		return true
	case c >= '\u3000' && c <= '\u303f': // CJK Symbols and Punctuation
		return true
	case c >= '\u3200' && c <= '\u32ff': // Enclosed CJK Letters and Months
		return true
	}
	return false
}
