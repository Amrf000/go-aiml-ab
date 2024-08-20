package ab

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// Bot struct represents the bot in GoLang
type Bot struct {
	Properties   Properties
	PreProcessor *PreProcessor
	Brain        *Graphmaster
	LearnfGraph  *Graphmaster
	LearnGraph   *Graphmaster
	Name         string
	SetMap       map[string]*AIMLSet
	MapMap       map[string]*AIMLMap
	PronounSet   map[string]bool
	RootPath     string
	BotPath      string
	BotNamePath  string
	AimlifPath   string
	AimlPath     string
	ConfigPath   string
	LogPath      string
	SetsPath     string
	MapsPath     string
}

func NewBot() *Bot {
	return NewBotWithName(DefaultBot)
}

func NewBotWithName(name string) *Bot {
	return NewBotWithPath(name, RootPath)
}

func NewBotWithPath(name, path string) *Bot {
	return NewBotWithAction(name, path, "auto")
}

func GetLastModifiedTime(filePath string) time.Time {
	fileInfo, err := os.Stat(filePath)

	// Checks for the error
	if err != nil {
		log.Fatal(err)
	}
	return fileInfo.ModTime()
}

func NewBotWithAction(name, path, action string) *Bot {
	cnt := 0

	bot := &Bot{
		Name:       name,
		RootPath:   RootPath,
		SetMap:     make(map[string]*AIMLSet),
		MapMap:     make(map[string]*AIMLMap),
		PronounSet: make(map[string]bool),
		Properties: Properties{},
	}

	bot.SetAllPaths(path, name)
	bot.Brain = NewGraphmaster(bot)
	bot.LearnfGraph = NewGraphmasterWithName(bot, "learnf")
	bot.LearnGraph = NewGraphmasterWithName(bot, "learn")
	bot.PreProcessor = NewPreProcessor(bot)
	bot.AddProperties()

	cnt = bot.AddAIMLSets()
	if TraceMode {
		fmt.Printf("Loaded %d set elements.\n", cnt)
	}
	cnt = bot.AddAIMLMaps()
	if TraceMode {
		fmt.Printf("Loaded %d map elements.\n", cnt)
	}
	bot.PronounSet = bot.GetPronouns()

	number := NewAIMLSet(NaturalNumberSetName, bot)
	bot.SetMap[NaturalNumberSetName] = number
	successor := NewAIMLMap(MapSuccessor, bot)
	bot.MapMap[MapSuccessor] = successor
	predecessor := NewAIMLMap(MapPredecessor, bot)
	bot.MapMap[MapPredecessor] = predecessor
	singular := NewAIMLMap(MapSingular, bot)
	bot.MapMap[MapSingular] = singular
	plural := NewAIMLMap(MapPlural, bot)
	bot.MapMap[MapPlural] = plural
	aimlDate := GetLastModifiedTime(bot.AimlPath)
	aimlIFDate := GetLastModifiedTime(bot.AimlifPath)
	if TraceMode {
		fmt.Println("AIML modified " + aimlDate.String() + " AIMLIF modified " + aimlIFDate.String())
	}
	//readUnfinishedIFCategories();
	PannousApiKey = GetPannousAPIKey(bot)
	PannousLogin = GetPannousLogin(bot)

	if action == "aiml2csv" {
		bot.AddCategoriesFromAIML()
	} else if action == "csv2aiml" {
		bot.AddCategoriesFromAIMLIF()
	} else if action == "chat-app" {
		if TraceMode {
			fmt.Println("Loading only AIMLIF files")
		}
		cnt = bot.AddCategoriesFromAIMLIF()
	} else {
		aimlDate := GetLastModifiedTime(bot.AimlPath)
		aimlIFDate := GetLastModifiedTime(bot.AimlifPath)
		if aimlDate.After(aimlIFDate) {
			if TraceMode {
				fmt.Println("AIML modified after AIMLIF")
			}
			cnt = bot.AddCategoriesFromAIML()
			bot.WriteAIMLIFFiles()
		} else {
			bot.AddCategoriesFromAIMLIF()
			if len(bot.Brain.GetCategories()) == 0 {
				fmt.Println("No AIMLIF Files found. Looking for AIML")
				cnt = bot.AddCategoriesFromAIML()
			}
		}
	}

	b := NewCategory(0, "PROGRAM VERSION", "*", "*", ProgramNameVersion, "update.aiml")
	// Input: "*", That: "*", Topic: "*", Template: ProgramNameVersion, Filename: "update.aiml"
	bot.Brain.AddCategory(b)
	bot.Brain.NodeStats()
	bot.LearnfGraph.NodeStats()

	return bot
}

// setAllPaths sets the paths based on root and bot name
//func (bot *Bot) SetAllPaths(name string) {
//	bot.BotPath = filepath.Join(bot.RootPath, "bots")
//	bot.BotNamePath = filepath.Join(bot.BotPath, name)
//	fmt.Printf("Name = %s Path = %s\n", name, bot.BotNamePath)
//	bot.AimlPath = filepath.Join(bot.BotNamePath, "aiml")
//	bot.AimlifPath = filepath.Join(bot.BotNamePath, "aimlif")
//	bot.ConfigPath = filepath.Join(bot.BotNamePath, "config")
//	bot.LogPath = filepath.Join(bot.BotNamePath, "logs")
//	bot.SetsPath = filepath.Join(bot.BotNamePath, "sets")
//	bot.MapsPath = filepath.Join(bot.BotNamePath, "maps")
//	if TraceMode {
//		fmt.Println(bot.RootPath)
//		fmt.Println(bot.BotPath)
//		fmt.Println(bot.BotNamePath)
//		fmt.Println(bot.AimlPath)
//		fmt.Println(bot.AimlifPath)
//		fmt.Println(bot.ConfigPath)
//		fmt.Println(bot.LogPath)
//		fmt.Println(bot.SetsPath)
//		fmt.Println(bot.MapsPath)
//	}
//}

// getPronouns retrieves pronouns from a file
//func (bot *Bot) GetPronouns() map[string]bool {
//	pronounSet := make(map[string]bool)
//	pronouns := ReadFileContents(filepath.Join(bot.ConfigPath, "pronouns.txt"))
//	splitPronouns := strings.Split(pronouns, "\n")
//	for _, p := range splitPronouns {
//		p = strings.TrimSpace(p)
//		if len(p) > 0 {
//			pronounSet[p] = true
//		}
//	}
//	if TraceMode {
//		fmt.Println("Read pronouns:", pronounSet)
//	}
//	return pronounSet
//}

// addMoreCategories adds more categories based on file and moreCategories
func (bot *Bot) AddMoreCategories(file string, moreCategories []*Category) {
	if strings.Contains(file, DeletedAimlFile) {
		// Handle deleted AIML file logic if needed
	} else if strings.Contains(file, LearnfAimlFile) {
		if TraceMode {
			fmt.Println("Reading Learnf file")
		}
		for _, c := range moreCategories {
			bot.Brain.AddCategory(c)
			bot.LearnfGraph.AddCategory(c)
		}
	} else {
		for _, c := range moreCategories {
			bot.Brain.AddCategory(c)
		}
	}
}

// addCategoriesFromAIML adds categories from AIML files
//func (bot *Bot) AddCategoriesFromAIML() int {
//	timer := NewTimer()
//	timer.Start()
//	cnt := 0
//	files, err := ioutil.ReadDir(bot.AimlPath)
//	if err != nil {
//		fmt.Println("Error reading AIML directory:", err)
//		return cnt
//	}
//	fmt.Println("Loading AIML files from", bot.AimlPath)
//	for _, file := range files {
//		if file.IsDir() {
//			continue
//		}
//		if strings.HasSuffix(strings.ToLower(file.Name()), ".aiml") {
//			fmt.Println(file.Name())
//			moreCategories := AIMLToCategories(filepath.Join(bot.AimlPath, file.Name()))
//			if err != nil {
//				fmt.Println("Problem loading", file.Name())
//				fmt.Println(err)
//				continue
//			}
//			bot.AddMoreCategories(file.Name(), moreCategories)
//			cnt += len(moreCategories)
//		}
//	}
//	fmt.Printf("Loaded %d categories in %.2f sec\n", cnt, timer.ElapsedTimeSecs())
//	return cnt
//}

// Utility function to get modification date of a file
func GetModificationDate(path string) time.Time {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return time.Time{}
	}
	return fileInfo.ModTime()
}

// Utility function to read file contents
func ReadFileContents(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}
	return string(content)
}

func (bot *Bot) SetAllPaths(root, name string) {
	bot.BotPath = fmt.Sprintf("%s/bots", root)
	bot.BotNamePath = fmt.Sprintf("%s/%s", bot.BotPath, name)
	if TraceMode {
		fmt.Printf("Name = %s Path = %s\n", name, bot.BotNamePath)
	}
	bot.AimlPath = fmt.Sprintf("%s/aiml", bot.BotNamePath)
	bot.AimlifPath = fmt.Sprintf("%s/aimlif", bot.BotNamePath)
	bot.ConfigPath = fmt.Sprintf("%s/config", bot.BotNamePath)
	bot.LogPath = fmt.Sprintf("%s/logs", bot.BotNamePath)
	bot.SetsPath = fmt.Sprintf("%s/sets", bot.BotNamePath)
	bot.MapsPath = fmt.Sprintf("%s/maps", bot.BotNamePath)
	if TraceMode {
		fmt.Println(bot.RootPath)
		fmt.Println(bot.BotPath)
		fmt.Println(bot.BotNamePath)
		fmt.Println(bot.AimlPath)
		fmt.Println(bot.AimlifPath)
		fmt.Println(bot.ConfigPath)
		fmt.Println(bot.LogPath)
		fmt.Println(bot.SetsPath)
		fmt.Println(bot.MapsPath)
	}
}

func (bot *Bot) GetPronouns() map[string]bool {
	pronounSet := make(map[string]bool)
	pronouns := GetFile(fmt.Sprintf("%s/pronouns.txt", bot.ConfigPath))
	splitPronouns := strings.Split(pronouns, "\n")
	for _, p := range splitPronouns {
		p = strings.TrimSpace(p)
		if len(p) > 0 {
			pronounSet[p] = false
		}
	}
	if TraceMode {
		fmt.Println("Read pronouns:", pronounSet)
	}
	return pronounSet
}

func (bot *Bot) AddCategoriesFromAIML() int {
	timer := NewTimer()
	timer.Start()
	cnt := 0
	defer func() {
		if TraceMode {
			fmt.Printf("Loaded %d categories in %.2f sec\n", cnt, timer.ElapsedTimeSecs())
		}
	}()
	files, err := os.ReadDir(bot.AimlPath)
	if err != nil {
		fmt.Println("addCategoriesFromAIML:", err)
		return cnt
	}
	if TraceMode {
		fmt.Println("Loading AIML files from", bot.AimlPath)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(file.Name()), ".AIML") {
			if TraceMode {
				fmt.Println(file.Name())
			}
			//if file.Name() == "client_profile.aiml" {
			//	nn := 0
			//	nn++
			//}
			moreCategories := AIMLToCategories(bot.AimlPath, file.Name())
			if err != nil {
				fmt.Println("Problem loading", file.Name())
				fmt.Println(err)
				continue
			}
			bot.AddMoreCategories(file.Name(), moreCategories)
			cnt += len(moreCategories)
		}
	}
	return cnt
}

func (bot *Bot) AddCategoriesFromAIMLIF() int {
	timer := NewTimer()
	timer.Start()
	cnt := 0
	defer func() {
		if TraceMode {
			fmt.Printf("Loaded %d categories in %.2f sec\n", cnt, timer.ElapsedTimeSecs())
		}
	}()
	files, err := os.ReadDir(bot.AimlifPath)
	if err != nil {
		fmt.Println("addCategoriesFromAIMLIF:", err)
		return 0
	}
	if TraceMode {
		fmt.Println("Loading AIML files from", bot.AimlifPath)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(file.Name()), AimlifFileSuffix) {
			if TraceMode {
				fmt.Println(file.Name())
			}
			categories, err := bot.ReadIFCategories(bot.AimlifPath + "/" + file.Name())
			if err != nil {
				fmt.Println("Problem loading", file.Name())
				fmt.Println(err)
				continue
			}
			cnt += len(categories)
			bot.AddMoreCategories(file.Name(), categories)
		}
	}
	return cnt
}

//func (bot *Bot) AddMoreCategories(file string, moreCategories []*Category) {
//	if strings.Contains(file, DeletedAimlFile) {
//		return
//	}
//	for _, c := range moreCategories {
//		bot.Brain.AddCategory(c)
//		if strings.Contains(file, LearnfAimlFile) {
//			bot.LearnfGraph.AddCategory(c)
//		}
//	}
//}

func (bot *Bot) WriteQuit() {
	bot.WriteAIMLIFFiles()
	bot.WriteAIMLFiles()
}

func (bot *Bot) ReadCertainIFCategories(graph *Graphmaster, fileName string) int {
	filePath := fmt.Sprintf("%s/%s%s", bot.AimlifPath, fileName, AimlifFileSuffix)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("No", filePath, "file found")
		return 0
	}
	certainCategories, err := bot.ReadIFCategories(filePath)
	if err != nil {
		fmt.Println("Problem loading", fileName)
		fmt.Println(err)
		return 0
	}
	for _, d := range certainCategories {
		graph.AddCategory(d)
	}
	fmt.Printf("readCertainIFCategories %d categories from %s%s\n", len(certainCategories), fileName, AimlifFileSuffix)
	return len(certainCategories)
}

func (bot *Bot) WriteCertainIFCategories(graph *Graphmaster, file string) {
	if TraceMode {
		fmt.Printf("writeCertainIFCategories %s size=%d\n", file, len(graph.GetCategories()))
	}
	bot.WriteIFCategories(graph.GetCategories(), file+AimlifFileSuffix)
	dir := bot.AimlifPath
	dirInfo, err := os.Stat(dir)
	if err == nil {
		dirInfo.ModTime()
	}
}

func (bot *Bot) WriteLearnfIFCategories() {
	bot.WriteCertainIFCategories(bot.LearnfGraph, LearnfAimlFile)
}

func (bot *Bot) WriteIFCategories(cats []*Category, filename string) {
	filePath := fmt.Sprintf("%s/%s", bot.AimlifPath, filename)
	existsPath, err := os.Stat(bot.AimlifPath)
	if err == nil {
		_ = existsPath
	}
	bw, err := os.Open(filePath)
	if err != nil {
		fmt.Println("writeIFCategories", err)
		return
	}
	for _, category := range cats {
		_, _ = fmt.Fprintln(bw, CategoryToIF(category))
	}
}

func (bot *Bot) WriteAIMLIFFiles() {
	if TraceMode {
		fmt.Println("writeAIMLIFFiles")
	}
	fileMap := make(map[string]*os.File)
	b := NewCategory(0, "BRAIN BUILD", "*", "*", time.Now().String(), "update.aiml")
	// {Input: "*", That: "*", Topic: "*", Template: time.Now().String(), Filename: "update.aiml"}
	bot.Brain.AddCategory(b)
	brainCategories := bot.Brain.GetCategories()
	sort.Slice(brainCategories, func(i, j int) bool {
		return brainCategories[i].CategoryNumber < brainCategories[j].CategoryNumber
	})
	for _, c := range brainCategories {
		fileName := c.Filename
		var bw *os.File
		if _, existsPath := fileMap[fileName]; existsPath {
			bw = fileMap[fileName]
		} else {
			bw, err := os.Create(bot.AimlifPath + "/" + fileName + AimlifFileSuffix)
			if err != nil {
				fmt.Println("writeAIMLIFFiles", err)
				continue
			}
			fileMap[fileName] = bw
			_, _ = fmt.Fprintln(bw, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<aiml>")
			copyright := GetCopyright(bot, fileName)
			fmt.Fprintln(bw, copyright)
		}
		_, _ = fmt.Fprintln(bw, CategoryToIF(c))
	}
	for _, bw := range fileMap {
		_ = bw
	}
	dir, err := os.Stat(bot.AimlifPath)
	if err == nil {
		_ = dir.ModTime()
	}
}

func (bot *Bot) WriteAIMLFiles() {
	if TraceMode {
		fmt.Println("writeAIMLFiles")
	}
	fileMap := make(map[string]*os.File)
	b := NewCategory(0, "BRAIN BUILD", "*", "*", time.Now().String(), "update.aiml")
	// {Input: "*", That: "*", Topic: "*", Template: time.Now().String(), Filename: "update.aiml"}
	bot.Brain.AddCategory(b)
	brainCategories := bot.Brain.GetCategories()
	sort.Slice(brainCategories, func(i, j int) bool {
		return brainCategories[i].CategoryNumber < brainCategories[j].CategoryNumber
	})
	for _, c := range brainCategories {
		if c.Filename != NullAimlFile {
			var bw *os.File
			if _, existsPath := fileMap[c.Filename]; existsPath {
				bw = fileMap[c.Filename]
			} else {
				copyright := GetCopyright(bot, c.Filename)
				var err error
				bw, err = os.Create(bot.AimlPath + "/" + c.Filename)
				if err != nil {
					fmt.Println("writeAIMLFiles", err)
					continue
				}
				fileMap[c.Filename] = bw
				_, _ = fmt.Fprintln(bw, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<aiml>")
				fmt.Fprintln(bw, copyright)
			}
			_, _ = fmt.Fprintf(bw, CategoryToAIML(c)+"\n")
		}
	}
	for _, bw := range fileMap {
		_, _ = fmt.Fprintf(bw, "</aiml>\n")
	}
	currentTime := time.Now().Local()
	err := os.Chtimes(bot.AimlPath, currentTime, currentTime)
	if err != nil {
		fmt.Println(err)
	}
}

func (bot *Bot) AddProperties() {
	bot.Properties.GetProperties(bot.ConfigPath + "/properties.txt")
	//if err != nil {
	//	fmt.Println("addProperties:", err)
	//}
}

func (bot *Bot) ReadIFCategories(filename string) ([]*Category, error) {
	var categories []*Category
	fstream, err := os.Open(filename)
	if err != nil {
		fmt.Println("readIFCategories: Error opening file", err)
		return categories, err
	}
	defer fstream.Close()

	br := bufio.NewScanner(fstream)
	for br.Scan() {
		strLine := br.Text()
		c := IFToCategory(strLine)
		//if err != nil {
		//	fmt.Println("Invalid AIMLIF in", filename, "line", strLine)
		//	continue
		//}
		categories = append(categories, c)
	}
	if err := br.Err(); err != nil {
		fmt.Println("readIFCategories: Error reading file", err)
	}
	return categories, nil
}

func (bot *Bot) AddAIMLSets() int {
	timer := NewTimer()
	timer.Start()
	defer func() {
		if TraceMode {
			fmt.Printf("Loaded %d AIML Set files in %.2f sec\n", len(bot.SetMap), timer.ElapsedTimeSecs())
		}
	}()
	files, err := ioutil.ReadDir(bot.SetsPath)
	if err != nil {
		fmt.Println("addAIMLSets:", err)
		return 0
	}
	if TraceMode {
		fmt.Println("Loading AIML Set files from", bot.SetsPath)
	}
	cnt := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(file.Name()), ".TXT") {
			if TraceMode {
				fmt.Println(file.Name())
			}
			setName := strings.TrimSuffix(file.Name(), ".txt")
			if TraceMode {
				fmt.Println("Read AIML Set", setName)
			}
			aimlSet := NewAIMLSet(setName, bot)
			cnt += aimlSet.ReadAIMLSet(bot)
			bot.SetMap[setName] = aimlSet
		}
	}
	return cnt
}

func (bot *Bot) AddAIMLMaps() int {
	timer := NewTimer()
	timer.Start()
	defer func() {
		if TraceMode {
			fmt.Printf("Loaded %d AIML Map files in %.2f sec\n", len(bot.MapMap), timer.ElapsedTimeSecs())
		}
	}()
	files, err := os.ReadDir(bot.MapsPath)
	if err != nil {
		fmt.Println("addAIMLMaps:", err)
		return 0
	}
	if TraceMode {
		fmt.Println("Loading AIML Map files from", bot.MapsPath)
	}
	cnt := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(file.Name()), ".TXT") {
			if TraceMode {
				fmt.Println(file.Name())
			}
			mapName := strings.TrimSuffix(file.Name(), ".txt")
			if TraceMode {
				fmt.Println("Read AIML Map", mapName)
			}
			aimlMap := NewAIMLMap(mapName, bot)
			cnt += aimlMap.ReadAIMLMap(bot)
			bot.MapMap[mapName] = aimlMap
		}
	}
	return cnt
}

func (bot *Bot) DeleteLearnfCategories() {
	learnfCategories := bot.LearnfGraph.GetCategories()
	for _, c := range learnfCategories {
		n := bot.Brain.FindNode(c)
		if n != nil {
			fmt.Println("Found node", n, "for", c.InputThatTopic())
			n.Category = nil
		}
	}
	bot.LearnfGraph = NewGraphmaster(bot)
}

func (bot *Bot) DeleteLearnCategories() {
	learnCategories := bot.LearnGraph.GetCategories()
	for _, c := range learnCategories {
		n := bot.Brain.FindNode(c)
		if n != nil {
			fmt.Println("Found node", n, "for", c.InputThatTopic())
			n.Category = nil
		}
	}
	bot.LearnGraph = NewGraphmaster(bot)
}

func (bot *Bot) ShadowChecker() {
	bot.ShadowCheckerWithNode(bot.Brain.Root)
}

func (bot *Bot) ShadowCheckerWithNode(node *Nodemapper) {
	if IsLeaf(node) {
		input := node.Category.GetPattern()
		input = bot.Brain.ReplaceBotProperties(input)
		input = strings.NewReplacer("*", "XXX", "_", "XXX", "^", "", "#", "").Replace(input)
		that := strings.NewReplacer("*", "XXX", "_", "XXX", "^", "", "#", "").Replace(node.Category.That)
		topic := strings.NewReplacer("*", "XXX", "_", "XXX", "^", "", "#", "").Replace(node.Category.Topic)
		input = bot.InstantiateSets(input)
		fmt.Println("shadowChecker: input=", input)
		match := bot.Brain.MatchRaw(input, that, topic)
		if match != node {
			fmt.Println(InputThatTopic(input, that, topic))
			fmt.Println("MATCHED:     ", match.Category.InputThatTopic())
			fmt.Println("SHOULD MATCH:", node.Category.InputThatTopic())
		}
	} else {
		for _, key := range KeySet(node) {
			bot.ShadowCheckerWithNode(Get(node, key))
		}
	}
}

func (bot *Bot) InstantiateSets(pattern string) string {
	splitPattern := strings.Split(pattern, " ")
	for i, x := range splitPattern {
		if strings.HasPrefix(x, "<SET>") {
			setName := TrimTag(x, "SET")
			if _, exists := bot.SetMap[setName]; exists {
				x = "FOUNDITEM"
			} else {
				x = "NOTFOUND"
			}
			splitPattern[i] = x
		}
	}
	return strings.Join(splitPattern, " ")
}
