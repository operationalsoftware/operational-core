package format

import "strings"

func DecimalWithCommas(input string) string {

	isNegative := false
	if strings.Index(input, "-") == 0 {
		input = input[1:]
		isNegative = true
	}

	// Split the input on the decimal point
	parts := strings.Split(input, ".")

	// Format the integer part with commas
	integerPart := formatIntWithCommas(parts[0])

	// Combine the formatted integer and decimal parts
	if len(parts) == 1 {
		return integerPart
	}

	if isNegative {
		return "-" + integerPart + "." + parts[1]
	} else {
		return integerPart + "." + parts[1]
	}
}

func formatIntWithCommas(s string) string {
	// Add commas to the string
	var formattedStr string
	for i, digit := range reverseString(s) {
		if i > 0 && i%3 == 0 {
			formattedStr = "," + formattedStr
		}
		formattedStr = string(digit) + formattedStr
	}

	return formattedStr
}

func reverseString(s string) string {
	var reversed string
	for _, char := range s {
		reversed = string(char) + reversed
	}
	return reversed
}
