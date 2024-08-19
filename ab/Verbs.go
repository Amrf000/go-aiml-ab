package ab

import (
	"fmt"
	"strings"
)

var (
	Es        = []string{"sh", "ch", "th", "ss", "x"}
	Ies       = []string{"ly", "ry", "ny", "fy", "dy", "py"}
	Ring      = []string{"be", "me", "re", "se", "ve", "de", "le", "ce", "ze", "ke", "te", "ge", "ne", "pe", "ue"}
	Bing      = []string{"ab", "at", "op", "el", "in", "ur", "op", "er", "un", "in", "it", "et", "ut", "im", "id", "ol", "ig"}
	NotBing   = []string{"der", "eat", "ber", "ain", "sit", "ait", "uit", "eet", "ter", "lop", "ver", "wer", "aim", "oid", "eel", "out", "oin", "fer", "vel", "mit"}
	Irregular = make(map[string][]string)
	Be2was    = make(map[string]string)
	Be2been   = make(map[string]string)
	Be2is     = make(map[string]string)
	Be2being  = make(map[string]string)
	AllVerbs  []string
)

func EndsWith(verb string, endings []string) string {
	for _, x := range endings {
		if strings.HasSuffix(verb, x) {
			return x
		}
	}
	return ""
}

func Is(verb string) string {
	var ending string
	if val, ok := Irregular[verb]; ok {
		return val[2]
	}
	if strings.HasSuffix(verb, "go") {
		return verb + "Es"
	}
	if ending = EndsWith(verb, Es); ending != "" {
		return verb + "Es"
	}
	if ending = EndsWith(verb, Ies); ending != "" {
		return verb[:len(verb)-1] + "Ies"
	}
	return verb + "s"
}

func was(verb string) string {
	var ending string
	verb = strings.TrimSpace(verb)
	switch verb {
	case "admit":
		return "admitted"
	case "commit":
		return "committed"
	case "die":
		return "died"
	case "agree":
		return "agreed"
	}
	if strings.HasSuffix(verb, "efer") {
		return verb + "red"
	}
	if val, ok := Irregular[verb]; ok {
		return val[1]
	}
	if ending = EndsWith(verb, Ies); ending != "" {
		return verb[:len(verb)-1] + "ied"
	}
	if ending = EndsWith(verb, Ring); ending != "" {
		return verb + "d"
	}
	if ending = EndsWith(verb, Bing); ending != "" {
		if EndsWith(verb, NotBing) == "" {
			return verb + ending[1:2] + "ed"
		}
	}
	return verb + "ed"
}

func being(verb string) string {
	var ending string
	if val, ok := Irregular[verb]; ok {
		return val[3]
	}
	switch verb {
	case "admit":
		return "admitting"
	case "commit":
		return "committing"
	case "quit":
		return "quitting"
	case "die":
		return "dying"
	case "lie":
		return "lying"
	}
	if strings.HasSuffix(verb, "efer") {
		return verb + "Ring"
	}
	if ending = EndsWith(verb, Ring); ending != "" {
		return verb[:len(verb)-1] + "ing"
	}
	if ending = EndsWith(verb, Bing); ending != "" {
		if EndsWith(verb, NotBing) == "" {
			return verb + ending[1:2] + "ing"
		}
	}
	return verb + "ing"
}

func been(verb string) string {
	if val, ok := Irregular[verb]; ok {
		return val[2]
	}
	return was(verb)
}

func getIrregulars() {
	// Replace with actual file reading logic
	// irrFile := Utilities.getFile("c:/ab/data/irrverbs.txt")
	// triples := strings.Split(irrFile, "\n")
}

func MakeVerbSetsMaps(bot *Bot) {
	getIrregulars()
	verbFile := GetFile("c:/ab/data/verb300.txt")
	verbs := strings.Split(verbFile, "\n")

	allVerbs := make([]string, 0, len(verbs))
	for _, verb := range verbs {
		if verb != "" {
			allVerbs = append(allVerbs, verb)
		}
	}

	be := NewAIMLSet("be", bot)
	isSet := NewAIMLSet("is", bot)
	wasSet := NewAIMLSet("was", bot)
	beenSet := NewAIMLSet("been", bot)
	beingSet := NewAIMLSet("being", bot)
	is2be := NewAIMLMap("is2be", bot)
	be2is := NewAIMLMap("be2is", bot)
	was2be := NewAIMLMap("was2be", bot)
	be2was := NewAIMLMap("be2was", bot)
	been2be := NewAIMLMap("been2be", bot)
	be2been := NewAIMLMap("be2been", bot)
	be2being := NewAIMLMap("be2being", bot)
	being2be := NewAIMLMap("being2be", bot)

	for _, verb := range allVerbs {
		beForm := verb
		isForm := Is(verb)
		wasForm := was(verb)
		beenForm := been(verb)
		beingForm := being(verb)
		fmt.Println(verb + "," + isForm + "," + wasForm + "," + beingForm + "," + beenForm)
		be.Add(beForm)
		isSet.Add(isForm)
		wasSet.Add(wasForm)
		beenSet.Add(beenForm)
		beingSet.Add(beingForm)
		be2is.Put(beForm, isForm)
		is2be.Put(isForm, beForm)
		be2was.Put(beForm, wasForm)
		was2be.Put(wasForm, beForm)
		be2been.Put(beForm, beenForm)
		been2be.Put(beenForm, beForm)
		be2being.Put(beForm, beingForm)
		being2be.Put(beingForm, beForm)
	}

	be.WriteAIMLSet()
	isSet.WriteAIMLSet()
	wasSet.WriteAIMLSet()
	beenSet.WriteAIMLSet()
	beingSet.WriteAIMLSet()
	be2is.WriteAIMLMap()
	is2be.WriteAIMLMap()
	be2was.WriteAIMLMap()
	was2be.WriteAIMLMap()
	be2been.WriteAIMLMap()
	been2be.WriteAIMLMap()
	be2being.WriteAIMLMap()
	being2be.WriteAIMLMap()
}
