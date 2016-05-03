// Second version of the package created for the section of the tutorial
// labelled "Returning a Character"

/* Sample test
abc=123
if
abc=500
else
abc=22/2
endif
jkl=abc
end
*/

package kiss

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/dcw303/crenshaw-go/util"
)

// SymTab is a Table of Strings
var SymTab []string

// Look is a Lookahead character
var Look rune

// Token is a Token
var Token rune

// Value is a String Token of Look
var Value string

// LCount is a Label Counter
var LCount int

// KWList is a Keyword List
var KWList = []string{"IF", "ELSE", "ENDIF", "END"}

// KWCode is a Keyword Code
const KWCode string = "xilee"

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

// IsAlNum Recognizes an Alphanumeric Character
func IsAlNum(r rune) bool {
	return IsAlpha(r) || IsDigit(r)
}

// IsAddOp Recognizes an AddOp
func IsAddOp(r rune) bool {
	return strings.ContainsRune("+-", r)
}

// IsMulOp Recognizes a Mulop
func IsMulOp(r rune) bool {
	return strings.ContainsRune("*/", r)
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
	if Look != x {
		Expected(strconv.QuoteRuneToASCII(x))
	} else {
		GetChar()
		SkipWhite()
	}
}

// Fin Skips a CRLF
func Fin() {
	if Look == 0x0D { //CR
		GetChar()
	}
	if Look == 0x0A { // LF - Does this get hit from stdin?
		GetChar()
	}
	SkipWhite()
}

// Lookup Looks Up Tokens in the Keyword Table
func Lookup(table []string, s string, n int) int {
	found := false
	// the following two lines differ from the original tutorial because
	// pascal has 1-based arrays
	i := n - 1
	for i >= 0 && !found {
		if s == table[i] {
			found = true
		} else {
			i--
		}
	}
	return i
}

// GetName Gets an Identifier
func GetName() {
	for Look == 0x0D {
		Fin()
	}
	if !IsAlpha(Look) {
		Expected("Name")
	}
	Value = ""
	for IsAlNum(Look) {
		Value += string(unicode.ToUpper(Look))
		GetChar()
	}
	SkipWhite()
}

// GetNum Gets a Number
func GetNum() {
	if !IsDigit(Look) {
		Expected("Integer")
	}
	Value = ""
	for IsDigit(Look) {
		Value += string(Look)
		GetChar()
	}
	Token = '#'
	SkipWhite()
}

// Scan Is a Lexical Scanner
func Scan() {
	GetName()
	Token = rune(KWCode[Lookup(KWList, Value, 4)+1])
}

// MatchString Matches a Specific Input String
func MatchString(x string) {
	if Value != x {
		Expected(x)
	}
}

// NewLabel Generates a Unique label
func NewLabel() (out string) {
	out = "L" + strconv.Itoa(LCount)
	LCount++
	return
}

// PostLabel Posts a Label to Outputs
func PostLabel(l string) {
	util.WriteLine(l + ":")
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
	GetName()
	if Look == '(' {
		Match('(')
		Match(')')
		EmitLn("BSR " + Value)
	} else {
		EmitLn("MOVE " + Value + "(PC),D0")
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
		GetNum()
		EmitLn("MOVE #" + Value + ",D0")
	}
}

// SignedFactor Parses and Translates the First Math Factor
func SignedFactor() {
	s := Look == '-'
	if IsAddOp(Look) {
		GetChar()
		SkipWhite()
	}
	Factor()
	if s {
		EmitLn("NEG D0")
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

// Term1 Completes Term Processing (called by Term and FirstTerm)
func Term1() {
	for IsMulOp(Look) {
		EmitLn("MOVE D0,-(SP)")
		switch Look {
		case '*':
			Multiply()
		case '/':
			Divide()
		}
	}
}

// Term Parses and Translates a Math Term
func Term() {
	Factor()
	Term1()
}

// FirstTerm Parses and Translates a Math Term with Possible Leading Sign
func FirstTerm() {
	SignedFactor()
	Term1()
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
	FirstTerm()
	for IsAddOp(Look) {
		EmitLn("MOVE D0,-(SP)")
		switch Look {
		case '+':
			Add()
		case '-':
			Subtract()
		}
	}
}

// Condition Parses and Translates a Boolean Condition
// This version is a dummy
func Condition() {
	EmitLn("<condition>")
}

// DoIf Recognizes and Translates an IF Construct
func DoIf() {
	Condition()
	l1 := NewLabel()
	l2 := l1
	EmitLn("BEQ " + l1)
	Block()
	if Token == 'l' {
		l2 = NewLabel()
		EmitLn("BRA " + l2)
		PostLabel(l1)
		Block()
	}
	PostLabel(l2)
	MatchString("ENDIF")
}

// Assignment Parses and Translates an Assignment Statement
func Assignment() {
	name := Value
	Match('=')
	Expression()
	EmitLn("LEA " + string(name) + "(PC),A0")
	EmitLn("MOVE D0,(A0)")
}

// Block Recognizes and Translates a Statement Block
func Block() {
	Scan()
	for Token != 'e' && Token != 'l' {
		switch Token {
		case 'i':
			DoIf()
		default:
			Assignment()
		}
		Scan()
	}
}

// Init Initializes
func Init() {
	// golang will init it to zero anyway
	LCount = 0
	GetChar()
}

// Go starts the execution of this chapter
func Go() {
	Init()
	Block()
	MatchString("END")
	EmitLn("END")
}
