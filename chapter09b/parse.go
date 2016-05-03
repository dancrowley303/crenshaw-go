// Second version of the package created for the C compiler

/* Sample test
ia;ub;lc;cd;sie;xif;ug(){}<ctrl-z>
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

// Class is a Storage Class Specifier
var Class rune

// Sign is Sign Specifier
var Sign rune

// Typ is a Type Specifier
var Typ rune

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

// GetClass Gets a Storage Class
func GetClass() {
	if strings.ContainsRune("axs", Look) {
		Class = Look
		GetChar()
	} else {
		Class = 'a'
	}
}

// GetType Gets a Type Specifier
func GetType() {
	Typ = ' '
	if Look == 'u' {
		Sign = 'u'
		Typ = 'i'
		GetChar()
		return
	}
	Sign = 's'
	if strings.ContainsRune("ilc", Look) {
		Typ = Look
		GetChar()
	}
}

// TopDecl Processes a Top-Level Declaration
func TopDecl() {
	name := GetName()
	if Look == '(' {
		DoFunc(name)
	} else {
		DoData(name)
	}
}

// DoFunc Processes a Function Definition
func DoFunc(n rune) {
	Match('(')
	Match(')')
	Match('{')
	Match('}')
	if Typ == ' ' {
		Typ = 'i'
	}
	util.WriteLine(string(Class) + string(Sign) + string(Typ) + " function " +
		string(n))
}

// DoData Processes a Data Declaration
func DoData(n rune) {
	if Typ == ' ' {
		Expected("Type declaration")
	}
	util.WriteLine(string(Class) + string(Sign) + string(Typ) + " data " +
		string(n))
	for Look == ',' {
		n = GetName()
		util.WriteLine(string(Class) + string(Sign) + string(Typ) + " data " +
			string(n))
	}
	Match(';')
}

// Go starts the execution of this chapter
func Go() {
	Init()
	//0x1A is ctrl-z
	for Look != 0x1A {
		GetClass()
		GetType()
		TopDecl()
	}
}
