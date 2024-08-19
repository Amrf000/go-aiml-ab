package ab

import (
	"strconv"
	"strings"
)

var tensNames = [...]string{
	"",
	" ten",
	" twenty",
	" thirty",
	" forty",
	" fifty",
	" sixty",
	" seventy",
	" eighty",
	" ninety",
}

var numNames = [...]string{
	"",
	" one",
	" two",
	" three",
	" four",
	" five",
	" six",
	" seven",
	" eight",
	" nine",
	" ten",
	" eleven",
	" twelve",
	" thirteen",
	" fourteen",
	" fifteen",
	" sixteen",
	" seventeen",
	" eighteen",
	" nineteen",
}

func convertLessThanOneThousand(number int) string {
	var soFar string
	if number%100 < 20 {
		soFar = numNames[number%100]
		number /= 100
	} else {
		soFar = numNames[number%10]
		number /= 10
		soFar = tensNames[number%10] + soFar
		number /= 10
	}
	if number == 0 {
		return soFar
	}
	return numNames[number] + " hundred" + soFar
}

func convert(number int64) string {
	if number == 0 {
		return "zero"
	}
	var result string
	var sNumber = strconv.FormatInt(number, 10)
	mask := "000000000000"
	sNumber = mask[:len(mask)-len(sNumber)] + sNumber

	billions, _ := strconv.Atoi(sNumber[0:3])
	millions, _ := strconv.Atoi(sNumber[3:6])
	hundredThousands, _ := strconv.Atoi(sNumber[6:9])
	thousands, _ := strconv.Atoi(sNumber[9:12])

	if billions != 0 {
		result += convertLessThanOneThousand(billions) + " billion "
	}
	if millions != 0 {
		result += convertLessThanOneThousand(millions) + " million "
	}
	if hundredThousands != 0 {
		if hundredThousands == 1 {
			result += "one thousand "
		} else {
			result += convertLessThanOneThousand(hundredThousands) + " thousand "
		}
	}
	if thousands != 0 {
		result += convertLessThanOneThousand(thousands)
	}
	return strings.TrimSpace(result)
}
