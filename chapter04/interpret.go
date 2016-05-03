/* Sample test
?cat=1
?dog=2
!cat
!dog
pets=cat+dog
!pets
*/

package interpret

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/dcw303/crenshaw-go/util"
)

// Look is a Lookahead character
var Look rune

// Table is used to store variables
var Table map[string]int

// GetChar Reads New Character From Input Stream
func GetChar() {
	Look = util.Read()
}

// Error Reports an Error
func Error(s string) {
	util.WriteBlankLine()
	util.WriteLine("Error: " + s)
}

// Abort Reports Error and Halts
func Abort(s string) {
	Error(s)
	panic("Aborted")
}

// Expected Reports What Was Expected
func Expected(s string) {
	Abort(s + " Expected")
}

// IsAlpha Recognizes an Alpha Character
func IsAlpha(r rune) bool {
	return unicode.IsLetter(r)
}

// IsDigit Recognizes a Decimal Digit
func IsDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// IsAlNum Recognizes an Alphanumeric
func IsAlNum(r rune) bool {
	return IsAlpha(r) || IsDigit(r)
}

// IsAddOp Recognizes an AddOp
func IsAddOp(r rune) bool {
	return strings.ContainsRune("+-", r)
}

// IsWhite Recognizes White Space
func IsWhite(r rune) bool {
	return r == 0x20 || r == 0x09
}

// SkipWhite Skips Over Leading White Space
func SkipWhite() {
	for IsWhite(Look) {
		GetChar()
	}
}

// Match Matches a Specific Input Character
func Match(x rune) {
	if Look == x {
		GetChar()
	} else {
		Expected(strconv.QuoteRuneToASCII(x))
	}
}

// NewLine Recognizes and Skips Over a Newline
func NewLine() {
	if Look == 0x0D {
		GetChar()
	}
}

// GetName Gets an Identifier
func GetName() (token string) {
	if !IsAlpha(Look) {
		Expected("Name")
	}
	for IsAlNum(Look) {
		token += string(unicode.ToUpper(Look))
		GetChar()
	}
	SkipWhite()
	return
}

// GetNum Gets a Number
func GetNum() (value int) {
	if !IsDigit(Look) {
		Expected("Integer")
	}
	for IsDigit(Look) {
		digit, err := strconv.Atoi(string(Look))
		if err != nil {
			panic(err)
		}
		value = 10*value + digit
		GetChar()
	}
	SkipWhite()
	return
}

// Emit Ouputs a String with Tab
func Emit(s string) {
	util.Write("\t " + s)
}

// EmitLn Ouputs a String with Tab and CRLF
func EmitLn(s string) {
	Emit(s)
	util.WriteBlankLine()
}

// Factor Parses and Translates a Math Factor
func Factor() (value int) {
	switch {
	case Look == '(':
		Match('(')
		value = Expression()
		Match(')')
	case IsAlpha(Look):
		value = Table[GetName()]
	default:
		value = GetNum()
	}
	return
}

// Term Parses and Translates a Math Term
func Term() (value int) {
	value = Factor()
	for strings.ContainsRune("*/", Look) {
		switch Look {
		case '*':
			Match('*')
			value = value * Factor()
		case '/':
			Match('/')
			value = value / Factor()
		}
	}
	return
}

// Expression Parses and Translates a Math Expression
func Expression() (value int) {
	if IsAddOp(Look) {
		value = 0
	} else {
		value = Term()
	}
	for IsAddOp(Look) {
		switch Look {
		case '+':
			Match('+')
			value = value + Term()
		case '-':
			Match('-')
			value = value - Term()
		}
	}
	return
}

// Assignment Parses and Translates an Assignment Statement
func Assignment() {
	name := GetName()
	Match('=')
	Table[name] = Expression()
}

// InitTable Initializes the Table of variables
func InitTable() {
	Table = make(map[string]int)
}

// Init Initializes
func Init() {
	InitTable()
	GetChar()
	SkipWhite()
}

// Input Inputs Routine
func Input() {
	Match('?')
	index := GetName()
	Match('=')
	number := GetNum()
	Table[index] = number
}

// Output Outputs Routine
func Output() {
	Match('!')
	util.WriteLine(strconv.Itoa(Table[GetName()]))
}

// Go starts the execution of this chapter
func Go() {
	Init()
	for Look != '.' {
		switch Look {
		case '?':
			Input()
		case '!':
			Output()
		default:
			Assignment()
		}
		NewLine()
	}
}
