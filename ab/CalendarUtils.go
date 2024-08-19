package ab

import (
	"fmt"
	"time"
)

func FormatTime(formatString string, msSinceEpoch int64) string {
	t := time.Unix(0, msSinceEpoch*int64(time.Millisecond))
	return t.Format(formatString)
}

func TimeZoneOffset() int {
	_, offset := time.Now().Zone()
	return offset / 60
}

func Utils_Year() string {
	return fmt.Sprintf("%d", time.Now().Year())
}

func Utils_Date() string {
	return time.Now().Format("January 02, 2006")
}

func DateCustom(jformat, locale, timezone string) string {
	if jformat == "" {
		jformat = "Mon Jan 02 15:04:05 MST 2006"
	}
	if locale == "" {
		locale = "USA"
	}
	if timezone == "" {
		timezone = time.Now().Location().String()
	}
	t := time.Now()
	dateAsString := t.Format(jformat)
	return dateAsString
}
