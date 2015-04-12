package bencoding

// Converts a CamelCase string to a lower case string with spaces
import "unicode"

func ToLowerCaseWithSpaces(s string) (lower string) {
	for i, r := range s {
		if unicode.IsLower(r) {
			lower += string(r)
		} else if unicode.IsUpper(r) {
			if i != 0 {
				lower += " "
			}
			lower += string(unicode.ToLower(r))
		} else if string(r) == "_" {
			lower += "-"
		}
	}
	return lower
}

func ToLowerCaseWithUnderscores(s string) (lower string) {
	for i, r := range s {
		if unicode.IsLower(r) {
			lower += string(r)
		} else if unicode.IsUpper(r) {
			if i != 0 {
				lower += "_"
			}
			lower += string(unicode.ToLower(r))
		} else if string(r) == "_" {
			lower += "-"
		}
	}
	return lower
}

func ToCamelCase(s string) (camel string) {
	previousRuneWasSpace := true
	for _, r := range s {
		if previousRuneWasSpace {
			camel += string(unicode.ToUpper(r))
			previousRuneWasSpace = false
		} else if unicode.IsSpace(r) {
			previousRuneWasSpace = true
		} else if string(r) == "-" {
			camel += "_"
			previousRuneWasSpace = true
		} else {
			camel += string(r)
		}
	}
	return camel
}
