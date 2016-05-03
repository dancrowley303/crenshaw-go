/* Sample test
3+2*3/5
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

// Match Matches a Specific Input Character
func Match(x rune) {
	if Look == x {
		GetChar()
	} else {
		Expected(strconv.QuoteRuneToASCII(x))
	}
}

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

// GetName Gets an Identifier
func GetName() (r rune) {
	if !IsAlpha(Look) {
		Expected("Name")
	}
	r = unicode.ToUpper(Look)
	GetChar()
	return
}

// GetNum Gets a Number
func GetNum() (r rune) {
	if !IsDigit(Look) {
		Expected("Integer")
	}
	r = Look
	GetChar()
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

// Init Initializes
func Init() {
	GetChar()
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
		default:
			Expected("AddOp")
		}
	}
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
		default:
			Expected("MulOp")
		}
	}
}

// Factor Parses and Translates a Math Factor
func Factor() {
	switch {
	case Look == '(':
		Match('(')
		Expression()
		Match(')')
	default:
		EmitLn("MOVE #" + string(GetNum()) + ",D0")
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
	EmitLn("DIVS D1,D0")
}

// Go starts the execution of this chapter
func Go() {
	Init()
	Expression()
}
