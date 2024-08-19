package ab

import (
	"regexp"
	"strings"
	"sync"
)

type Inflector struct {
	plurals      []Rule
	singulars    []Rule
	uncountables map[string]struct{}
}

type Rule struct {
	expression  *regexp.Regexp
	replacement string
}

var instance *Inflector
var once sync.Once

func GetInstance() *Inflector {
	once.Do(func() {
		instance = &Inflector{
			uncountables: make(map[string]struct{}),
		}
		instance.Initialize()
	})
	return instance
}

func (r *Rule) Apply(input string) string {
	if r == nil {
		return ""
	}
	if r.expression == nil {
		return input
	}
	return r.expression.ReplaceAllString(input, r.replacement)
}

func (inf *Inflector) Pluralize(word interface{}) string {
	wordStr := strings.TrimSpace(word.(string))
	if len(wordStr) == 0 || inf.IsUncountable(wordStr) {
		return wordStr
	}
	for _, rule := range inf.plurals {
		result := rule.Apply(wordStr)
		if result != wordStr {
			return result
		}
	}
	return wordStr
}

func (inf *Inflector) Singularize(word interface{}) string {
	wordStr := strings.TrimSpace(word.(string))
	if len(wordStr) == 0 || inf.IsUncountable(wordStr) {
		return wordStr
	}
	for _, rule := range inf.singulars {
		result := rule.Apply(wordStr)
		if result != wordStr {
			return result
		}
	}
	return wordStr
}

func (inf *Inflector) IsUncountable(word string) bool {
	_, exists := inf.uncountables[strings.ToLower(strings.TrimSpace(word))]
	return exists
}

func (inf *Inflector) AddPluralize(rule, replacement string) {
	inf.plurals = append([]Rule{{expression: regexp.MustCompile("(?i)" + rule), replacement: replacement}}, inf.plurals...)
}

func (inf *Inflector) AddSingularize(rule, replacement string) {
	inf.singulars = append([]Rule{{expression: regexp.MustCompile("(?i)" + rule), replacement: replacement}}, inf.singulars...)
}

func (inf *Inflector) AddIrregular(singular, plural string) {
	singularRemainder := ""
	if len(singular) > 1 {
		singularRemainder = singular[1:]
	}
	pluralRemainder := ""
	if len(plural) > 1 {
		pluralRemainder = plural[1:]
	}
	inf.AddPluralize("("+string(singular[0])+")"+singularRemainder+"$", "$1"+pluralRemainder)
	inf.AddSingularize("("+string(plural[0])+")"+pluralRemainder+"$", "$1"+singularRemainder)
}

func (inf *Inflector) AddUncountable(words ...string) {
	for _, word := range words {
		inf.uncountables[strings.ToLower(strings.TrimSpace(word))] = struct{}{}
	}
}

func (inf *Inflector) Initialize() {
	inf.AddPluralize("$", "s")
	inf.AddPluralize("s$", "s")
	inf.AddPluralize("(ax|test)is$", "$1es")
	inf.AddPluralize("(octop|vir)us$", "$1i")
	inf.AddPluralize("(octop|vir)i$", "$1i")
	inf.AddPluralize("(alias|status)$", "$1es")
	inf.AddPluralize("(bu)s$", "$1ses")
	inf.AddPluralize("(buffal|tomat)o$", "$1oes")
	inf.AddPluralize("([ti])um$", "$1a")
	inf.AddPluralize("([ti])a$", "$1a")
	inf.AddPluralize("sis$", "ses")
	inf.AddPluralize("(?:([^f])fe|([lr])f)$", "$1$2ves")
	inf.AddPluralize("(hive)$", "$1s")
	inf.AddPluralize("([^aeiouy]|qu)y$", "$1ies")
	inf.AddPluralize("(x|ch|ss|sh)$", "$1es")
	inf.AddPluralize("(matr|vert|ind)ix|ex$", "$1ices")
	inf.AddPluralize("([m|l])ouse$", "$1ice")
	inf.AddPluralize("([m|l])ice$", "$1ice")
	inf.AddPluralize("^(ox)$", "$1en")
	inf.AddPluralize("(quiz)$", "$1zes")
	inf.AddPluralize("(people|men|children|sexes|moves|stadiums)$", "$1")
	inf.AddPluralize("(oxen|octopi|viri|aliases|quizzes)$", "$1")

	inf.AddSingularize("s$", "")
	inf.AddSingularize("(s|si|u)s$", "$1s")
	inf.AddSingularize("(n)ews$", "$1ews")
	inf.AddSingularize("([ti])a$", "$1um")
	inf.AddSingularize("((a)naly|(b)a|(d)iagno|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$", "$1$2sis")
	inf.AddSingularize("(^analy)ses$", "$1sis")
	inf.AddSingularize("(^analy)sis$", "$1sis")
	inf.AddSingularize("([^f])ves$", "$1fe")
	inf.AddSingularize("(hive)s$", "$1")
	inf.AddSingularize("(tive)s$", "$1")
	inf.AddSingularize("([lr])ves$", "$1f")
	inf.AddSingularize("([^aeiouy]|qu)ies$", "$1y")
	inf.AddSingularize("(s)eries$", "$1eries")
	inf.AddSingularize("(m)ovies$", "$1ovie")
	inf.AddSingularize("(x|ch|ss|sh)es$", "$1")
	inf.AddSingularize("([m|l])ice$", "$1ouse")
	inf.AddSingularize("(bus)es$", "$1")
	inf.AddSingularize("(o)es$", "$1")
	inf.AddSingularize("(shoe)s$", "$1")
	inf.AddSingularize("(cris|ax|test)is$", "$1is")
	inf.AddSingularize("(cris|ax|test)es$", "$1is")
	inf.AddSingularize("(octop|vir)i$", "$1us")
	inf.AddSingularize("(octop|vir)us$", "$1us")
	inf.AddSingularize("(alias|status)es$", "$1")
	inf.AddSingularize("(alias|status)$", "$1")
	inf.AddSingularize("^(ox)en", "$1")
	inf.AddSingularize("(vert|ind)ices$", "$1ex")
	inf.AddSingularize("(matr)ices$", "$1ix")
	inf.AddSingularize("(quiz)zes$", "$1")

	inf.AddIrregular("person", "people")
	inf.AddIrregular("man", "men")
	inf.AddIrregular("child", "children")
	inf.AddIrregular("sex", "sexes")
	inf.AddIrregular("move", "moves")
	inf.AddIrregular("stadium", "stadiums")
	inf.AddUncountable("equipment", "information", "rice", "money", "species", "series", "fish", "sheep")
}
