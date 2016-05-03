package scanner

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/dcw303/crenshaw-go/chapter16/errors"
	"github.com/dcw303/crenshaw-go/chapter16/input"
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
	return strings.ContainsRune("+-|~", r)
}

// IsMulOp Recognizes a MulOp
func IsMulOp(r rune) bool {
	return strings.ContainsRune("*/&", r)
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
func GetName() string {
	n := ""
	if !IsAlpha(input.Look) {
		errors.Expected("Name")
	}
	for IsAlNum(input.Look) {
		//note that the ToUpper is not documented in the tutorial, but added in
		//to conform with the case-insensitive variable theme of this compiler
		n += string(unicode.ToUpper(input.Look))
		input.GetChar()
	}
	return n
}

// GetNumber Gets a Number
func GetNumber() string {
	n := ""
	if !IsDigit(input.Look) {
		errors.Expected("Integer")
	}
	for IsDigit(input.Look) {
		n += string(input.Look)
		input.GetChar()
	}
	return n
}

// GetNumberAsLongInt Gets a Number (integer version)
// Note that this function is not actively called in the tutorial, and is
// merely documentation for a possible implementation as the tutorial suggests.
func GetNumberAsLongInt() (n int64) {
	n = 0
	if !IsDigit(input.Look) {
		errors.Expected("Integer")
	}
	for IsDigit(input.Look) {
		digit, err := strconv.Atoi(string(input.Look))
		if err != nil {
			panic(err)
		}
		n = 10*n + int64(digit)
		input.GetChar()
	}
	return
}
