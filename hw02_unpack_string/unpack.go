package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var builder strings.Builder
	var resultBuilder strings.Builder

	for _, char := range s {
		switch {
		case unicode.IsLetter(char):
			if builder.Len() > 0 {
				resultBuilder.WriteString(builder.String())
				builder.Reset()
			}
			builder.WriteRune(char)
		case unicode.IsDigit(char):
			if builder.Len() == 0 {
				return "", ErrInvalidString
			}
			count := int(char - '0')
			resultBuilder.WriteString(strings.Repeat(builder.String(), count))
			builder.Reset()
		default:
			return "", ErrInvalidString
		}
	}

	if builder.Len() > 0 {
		resultBuilder.WriteString(builder.String())
	}

	return resultBuilder.String(), nil
}
