/* Sample test
abc123=20+60
*/

package parse

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/dcw303/crenshaw-go/util"
)

// Look is a Lookahead character
var Look rune

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

// Match Matches a Specific INput Character
func Match(x rune) {
	if Look == x {
		GetChar()
	} else {
		Expected(strconv.QuoteRuneToASCII(x))
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
func GetNum() (num string) {
	if !IsDigit(Look) {
		Expected("Integer")
	}
	for IsDigit(Look) {
		num += string(Look)
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

// Ident Parses and Translates an Identifier
func Ident() {
	name := GetName()
	if Look == '(' {
		Match('(')
		Match(')')
		EmitLn("BSR " + name)
	} else {
		EmitLn("MOVE " + name + "(PC),D0")
	}
}

// Factor Parses and Translates a Math Factor
func Factor() {
	switch {
	case Look == '(':
		Match('(')
		Expression()
		Match(')')
	case IsAlpha(Look):
		Ident()
	default:
		EmitLn("MOVE #" + GetNum() + ",D0")
	}
}

// Multiply Recognizes and Translates a Multiply
func Multiply() {
	Match('*')
	Factor()
	EmitLn("MULS (SP)+,D0")
}

// Divide Recognizes and Translates a Divide
func Divide() {
	Match('/')
	Factor()
	EmitLn("MOVE (SP)+,D1")
	EmitLn("EXS.L D0")
	EmitLn("DIVS D1,D0")
}

// Term Parses and Translates a Math Term
func Term() {
	Factor()
	for strings.ContainsRune("*/", Look) {
		EmitLn("MOVE D0,-(SP)")
		switch Look {
		case '*':
			Multiply()
		case '/':
			Divide()
		}
	}
}

// Add Recognizes and Translates an Add
func Add() {
	Match('+')
	Term()
	EmitLn("ADD (SP)+,D0")
}

// Subtract Recognizes and Translates a Subtract
func Subtract() {
	Match('-')
	Term()
	EmitLn("SUB (SP)+,D0")
	EmitLn("NEG D0")
}

// Expression Parses and Translates a Math Expression
func Expression() {
	if IsAddOp(Look) {
		EmitLn("CLR D0")
	} else {
		Term()
	}
	for strings.ContainsRune("+-", Look) {
		EmitLn("MOVE D0,-(SP)")
		switch Look {
		case '+':
			Add()
		case '-':
			Subtract()
		}
	}
}

// Assignment Parses and Translates an Assignment Statement
func Assignment() {
	name := GetName()
	Match('=')
	Expression()
	EmitLn("LEA " + name + "(PC),A0")
	EmitLn("MOVE D0,(A0)")
}

// Init Initializes
func Init() {
	GetChar()
	SkipWhite()
}

// Go starts the execution of this chapter
func Go() {
	Init()
	Assignment()
	if Look != 0x0D {
		Expected("Newline")
	}
}
