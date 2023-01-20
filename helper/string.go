package helper

import "unicode"

func CapitalizeFirstLetter(strToCapitalize string) string {
	strRune := []rune(strToCapitalize)
	strRune[0] = unicode.ToUpper(strRune[0])
	capitalizedStr := string(strRune)

	return capitalizedStr
}
