/* Sample test
abc
123
if
<
*
~
END
*/

package kiss

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/dcw303/crenshaw-go/util"
)

// Look is a Lookahead character
var Look rune

// SymType is an enumeration of Symbol Types
type SymType int

const (
	//IfSym If Symbol
	IfSym SymType = iota
	//ElseSym Else Symbol
	ElseSym
	//EndIfSym End If Symbol
	EndIfSym
	//EndSym End Symbol
	EndSym
	//Ident Identity
	Ident
	//Number Number
	Number
	//Operator Operator
	Operator
)

// Token is a Token
var Token SymType

// Value is a String Token of Lookahead
var Value string

// SymTab is a Table of Strings
var SymTab []string

// KWList is a Keyword List
var KWList = []string{"IF", "ELSE", "ENDIF", "END"}

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

// IsOp Recognizes Any Operator
func IsOp(r rune) bool {
	return strings.ContainsRune("+-*/<>:=", r)
}

// GetName Gets an Identifier
func GetName() {
	Value = ""
	if !IsAlpha(Look) {
		Expected("Name")
	}
	for IsAlNum(Look) {
		Value += string(unicode.ToUpper(Look))
		GetChar()
	}
	k := Lookup(KWList, Value, 4)
	if k == -1 {
		Token = Ident
	} else {
		Token = SymType(k)
	}
}

// GetNum Gets a Number
func GetNum() {
	Value = ""
	if !IsDigit(Look) {
		Expected("Integer")
	}
	for IsDigit(Look) {
		Value += string(Look)
		GetChar()
	}
	Token = Number
}

// GetOp Gets an Operator
func GetOp() {
	Value = ""
	if !IsOp(Look) {
		Expected("Operator")
	}
	for IsOp(Look) {
		Value += string(Look)
		GetChar()
	}
	Token = Operator
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

// IsWhite Recognizes White Space
func IsWhite(r rune) bool {
	return r == 0x20 || r == 0x09
}

// SkipComma Skips Over a Comma
func SkipComma() {
	SkipWhite()
	if Look == ',' {
		GetChar()
		SkipWhite()
	}
}

// SkipWhite Skips Over Leading White Space
func SkipWhite() {
	for IsWhite(Look) {
		GetChar()
	}
}

// Scan Is a Lexical Scanner
func Scan() {
	for Look == 0x0D {
		Fin()
	}
	switch {
	case IsAlpha(Look):
		GetName()
	case IsDigit(Look):
		GetNum()
	case IsOp(Look):
		GetOp()
	default:
		Value = string(Look)
		Token = Operator
		GetChar()
	}
	SkipWhite()
}

// Fin Skips a CRLF
func Fin() {
	if Look == 0x0D { //CR
		GetChar()
	}
	if Look == 0x0A { // LF - Does this get hit from stdin?
		GetChar()
	}
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

// Init Initializes
func Init() {
	GetChar()
}

// Go starts the execution of this chapter
func Go() {
	Init()
	for Token != EndSym {
		Scan()
		switch Token {
		case Ident:
			util.Write("Ident ")
		case Number:
			util.Write("Number ")
		case Operator:
			util.Write("Operator ")
		case IfSym, ElseSym, EndIfSym, EndSym:
			util.Write("Keyword ")
		}
		util.WriteLine(Value)
	}
}
