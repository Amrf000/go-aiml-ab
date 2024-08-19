package ab

import (
	"fmt"
	"time"
)

// Test demonstrates the usage of the interval utility functions.
func Test() {
	date1 := "23:59:59.00"
	date2 := "12:00:00.00"
	format := "15:04:05.00"
	hours := GetHoursBetween(date2, date1, format)
	fmt.Printf("Hours = %d\n", hours)

	date1 = "January 30, 2013"
	date2 = "August 2, 1960"
	format = "January 2, 2006"
	years := GetYearsBetween(date2, date1, format)
	fmt.Printf("Years = %d\n", years)
}

// ParseDate parses a date string with the given format.
func ParseDate(dateStr, format string) (time.Time, error) {
	return time.ParseInLocation(format, dateStr, time.Local)
}

// GetHoursBetween calculates the number of hours between two dates.
func GetHoursBetween(date1, date2, format string) int {
	time1, err1 := ParseDate(date1, format)
	time2, err2 := ParseDate(date2, format)

	if err1 != nil || err2 != nil {
		if err1 != nil {
			fmt.Println("Error parsing date1:", err1)
		}
		if err2 != nil {
			fmt.Println("Error parsing date2:", err2)
		}
		return 0
	}

	duration := time2.Sub(time1)
	return int(duration.Hours())
}

// GetYearsBetween calculates the number of years between two dates.
func GetYearsBetween(date1, date2, format string) int {
	time1, err1 := ParseDate(date1, format)
	time2, err2 := ParseDate(date2, format)

	if err1 != nil || err2 != nil {
		if err1 != nil {
			fmt.Println("Error parsing date1:", err1)
		}
		if err2 != nil {
			fmt.Println("Error parsing date2:", err2)
		}
		return 0
	}

	years := time2.Year() - time1.Year()

	// Adjust for months and days
	if time2.YearDay() < time1.YearDay() {
		years--
	}
	return years
}

// GetMonthsBetween calculates the number of months between two dates.
func GetMonthsBetween(date1, date2, format string) int {
	time1, err1 := ParseDate(date1, format)
	time2, err2 := ParseDate(date2, format)

	if err1 != nil || err2 != nil {
		if err1 != nil {
			fmt.Println("Error parsing date1:", err1)
		}
		if err2 != nil {
			fmt.Println("Error parsing date2:", err2)
		}
		return 0
	}

	months := (time2.Year()-time1.Year())*12 + int(time2.Month()-time1.Month())

	// Adjust for days
	if time2.Day() < time1.Day() {
		months--
	}
	return months
}

// GetDaysBetween calculates the number of days between two dates.
func GetDaysBetween(date1, date2, format string) int {
	time1, err1 := ParseDate(date1, format)
	time2, err2 := ParseDate(date2, format)

	if err1 != nil || err2 != nil {
		if err1 != nil {
			fmt.Println("Error parsing date1:", err1)
		}
		if err2 != nil {
			fmt.Println("Error parsing date2:", err2)
		}
		return 0
	}

	duration := time2.Sub(time1)
	return int(duration.Hours() / 24)
}
