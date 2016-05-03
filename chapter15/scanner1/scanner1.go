package scanner1

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/dcw303/crenshaw-go/chapter15/errors"
	"github.com/dcw303/crenshaw-go/chapter15/input"
)

// IsAlpha Recognizes an Alpha Character
func IsAlpha(r rune) bool {
	return unicode.IsLetter(r)
}

// IsDigit Recognizes a Decimal Digit
func IsDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// IsAlNum Recognizes an Alphanumeric Character
func IsAlNum(r rune) bool {
	return IsAlpha(r) || IsDigit(r)
}

// IsAddOp Recognizes an AddOp
func IsAddOp(r rune) bool {
	return strings.ContainsRune("+-", r)
}

// IsMulOp Recognizes a MulOp
func IsMulOp(r rune) bool {
	return strings.ContainsRune("*/", r)
}

// Match Matches a Specific Input Character
func Match(x rune) {
	if input.Look == x {
		input.GetChar()
	} else {
		errors.Expected(strconv.QuoteRuneToASCII(x))
	}
}

// GetName Gets an Identifier
func GetName() (r rune) {
	if !IsAlpha(input.Look) {
		errors.Expected("Name")
	}
	r = unicode.ToUpper(input.Look)
	input.GetChar()
	return
}

// GetNumber Gets a Number
func GetNumber() (r rune) {
	if !IsDigit(input.Look) {
		errors.Expected("Integer")
	}
	r = input.Look
	input.GetChar()
	return
}
