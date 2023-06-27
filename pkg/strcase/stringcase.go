// This package was copied from https://github.com/golang-cz/textcase
package strcase

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Converts input string to "camelCase" (lower camel case) naming convention.
// Removes all whitespace and special characters. Supports Unicode characters.
func ToCamel(input string) string {
	str := markLetterCaseChanges(input)

	var b strings.Builder

	state := idle
	for i := 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		i += size
		state = state.next(r)
		switch state {
		case firstAlphaNum:
			if b.Len() > 0 {
				b.WriteRune(unicode.ToUpper(r))
			} else {
				b.WriteRune(unicode.ToLower(r))
			}
		case alphaNum:
			b.WriteRune(unicode.ToLower(r))
		}
	}

	return b.String()
}

// Converts input string to "PascalCase" (upper camel case) naming convention.
// Removes all whitespace and special characters. Supports Unicode characters.
func ToPascal(input string) string {
	str := markLetterCaseChanges(input)

	var b strings.Builder

	state := idle
	for i := 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		i += size
		state = state.next(r)
		switch state {
		case firstAlphaNum:
			b.WriteRune(unicode.ToUpper(r))
		case alphaNum:
			b.WriteRune(unicode.ToLower(r))
		}
	}

	return b.String()
}

// Converts input string to "snake_case" naming convention.
// Removes all whitespace and special characters. Supports Unicode characters.
func ToSnake(input string) string {
	str := markLetterCaseChanges(input)

	var b bytes.Buffer

	state := idle
	for i := 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		i += size
		state = state.next(r)
		switch state {
		case firstAlphaNum, alphaNum:
			b.WriteRune(unicode.ToLower(r))
		case delimiter:
			b.WriteByte('_')
		}
	}
	if (state == idle || state == delimiter) && b.Len() > 0 {
		b.Truncate(b.Len() - 1)
	}

	return b.String()
}

func IsSnake(input string) bool {
	return input == ToSnake(input)
}

func IsCamel(input string) bool {
	return input == ToCamel(input)
}

func IsPascal(input string) bool {
	return input == ToPascal(input)
}
