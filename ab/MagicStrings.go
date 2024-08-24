package ab

var (
	ProgramNameVersion       = "Program AB 0.0.6.26 beta -- AI Foundation Reference AIML 2.1 implementation"
	Comment                  = "Added repetition detection."
	AimlifSplitChar          = ","
	DefaultBot               = "alice2"
	DefaultLanguage          = "EN"
	AimlifSplitCharName      = "\\#Comma"
	AimlifFileSuffix         = ".csv"
	AbSampleFile             = "sample.txt"
	TextCommentMark          = ";;"
	PannousApiKey            = "guest"
	PannousLogin             = "test-user"
	SraixFailed              = "SRAIXFAILED"
	RepetitionDetected       = "REPETITIONDETECTED"
	SraixNoHint              = "nohint"
	SraixEventHint           = "event"
	SraixPicHint             = "pic"
	SraixShoppingHint        = "shopping"
	UnknownAimlFile          = "unknownAimlFile.xml"
	DeletedAimlFile          = "deleted.xml"
	LearnfAimlFile           = "learnf.xml"
	NullAimlFile             = "null.xml"
	InappropriateAimlFile    = "inappropriate.xml"
	ProfanityAimlFile        = "profanity.xml"
	InsultAimlFile           = "insults.xml"
	ReductionsUpdateAimlFile = "reductions_update.xml"
	PredicatesAimlFile       = "client_profile.xml"
	UpdateAimlFile           = "update.xml"
	PersonalityAimlFile      = "personality.xml"
	SraixAimlFile            = "sraix.xml"
	OobAimlFile              = "oob.xml"
	UnfinishedAimlFile       = "unfinished.xml"
	InappropriateFilter      = "FILTER INAPPROPRIATE"
	ProfanityFilter          = "FILTER PROFANITY"
	InsultFilter             = "FILTER INSULT"
	DeletedTemplate          = "deleted"
	UnfinishedTemplate       = "unfinished"
	BadJavascript            = "JSFAILED"
	JsEnabled                = "true"
	UnknownHistoryItem       = "unknown"
	DefaultBotResponse       = "I have no answer for that."
	ErrorBotResponse         = "Something is wrong with my brain."
	ScheduleError            = "I'm unable to schedule that event."
	SystemFailed             = "Failed to execute system command."
	DefaultGet               = "unknown"
	DefaultProperty          = "unknown"
	DefaultMap               = "unknown"
	DefaultCustomerId        = "unknown"
	DefaultBotName           = "unknown"
	DefaultThat              = "unknown"
	DefaultTopic             = "unknown"
	DefaultListItem          = "NIL"
	UndefinedTriple          = "NIL"
	UnboundVariable          = "unknown"
	TemplateFailed           = "Template failed."
	TooMuchRecursion         = "Too much recursion in AIML"
	TooMuchLooping           = "Too much looping in AIML"
	BlankTemplate            = "blank template"
	NullInput                = "NORESP"
	NullStar                 = "nullstar"
	SetMemberString          = "ISA"
	RemoteMapKey             = "external"
	RemoteSetKey             = "external"
	NaturalNumberSetName     = "number"
	MapSuccessor             = "successor"
	MapPredecessor           = "predecessor"
	MapSingular              = "singular"
	MapPlural                = "plural"
	RootPath                 = "c:/ab"
)

func SetRootPath(newRootPath string) {
	RootPath = newRootPath
}

func SetRootPathFromSystem() {
	// Simulate System.getProperty("user.dir")
	// Not directly translatable in Go, use current directory
	SetRootPath(".")
}
