/* Sample test:
pabe.
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

// Prog Parses and Translates a Program
func Prog() {
	Match('p')
	name := GetName()
	Prolog()
	DoBlock(name)
	Match('.')
	Epilog(name)
}

// Prolog Writes the Prolog
func Prolog() {
	EmitLn("WARMST EQU $A01E")
}

// Epilog Writes the Eiplog
func Epilog(name rune) {
	EmitLn("DC WARMST")
	EmitLn("END " + string(name))
}

// DoBlock Parses and Translates a Pascal Block
func DoBlock(name rune) {
	Declarations()
	PostLabel(name)
	Statements()
}

// Declarations Parses and Translates the Declaration Part
func Declarations() {
	for strings.ContainsRune("lctvpf", Look) {
		switch Look {
		case 'l':
			Labels()
		case 'c':
			Constants()
		case 't':
			Types()
		case 'v':
			Variables()
		case 'p':
			DoProcedure()
		case 'f':
			DoFunction()
		}
	}
}

// Labels Processes Label Statement
func Labels() {
	Match('l')
}

// Constants Processes Const Statement
func Constants() {
	Match('c')
}

// Types Processes Type Statement
func Types() {
	Match('t')
}

// Variables Processes Var Statement
func Variables() {
	Match('v')
}

// DoProcedure Processes Procedure Definition
func DoProcedure() {
	Match('p')
}

// DoFunction Processes Function Definition
func DoFunction() {
	Match('f')
}

// Statements Writes the Statements
func Statements() {
	Match('b')
	for Look != 'e' {
		GetChar()
	}
	Match('e')
}

// PostLabel Posts a Label to Outputs
func PostLabel(l rune) {
	util.WriteLine(string(l) + ":")
}

// Go starts the execution of this chapter
func Go() {
	Init()
	Prog()
}
