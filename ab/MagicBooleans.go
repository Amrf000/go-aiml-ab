package ab

import "fmt"

var (
	TraceMode               = false
	EnableExternalSets      = true
	EnableExternalMaps      = true
	JpTokenize              = false
	FixExcelCsv             = true
	EnableNetworkConnection = true
	CacheSraix              = false
	QaTestMode              = false
	MakeVerbsSetsMapsFlag   = false
)

func trace(traceString string) {
	if TraceMode {
		fmt.Println(traceString)
	}
}
