/* Sample test
a<1*3|c=d|e>f|g=5
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

// IsBoolean Recognizes a Boolean Literal
func IsBoolean(r rune) bool {
	return strings.ContainsRune("TF", unicode.ToUpper(r))
}

// GetBoolean Gets a Boolean Literal
func GetBoolean() (out bool) {
	if !IsBoolean(Look) {
		Expected("Boolean Literal")
	}
	out = unicode.ToUpper(Look) == 'T'
	GetChar()
	return
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

// BoolExpression Parses and Translates a Boolean Expression
func BoolExpression() {
	BoolTerm()
	for IsOrOp(Look) {
		EmitLn("MOVE D0,-(SP)")
		switch Look {
		case '|':
			BoolOr()
		case '~':
			BoolXor()
		}
	}
}

// BoolOr Recognizes and Translates a Boolean OR
func BoolOr() {
	Match('|')
	BoolTerm()
	EmitLn("OR (SP)+,D0")
}

// BoolXor Recognizes and Translates an Exclusive OR
func BoolXor() {
	Match('~')
	BoolTerm()
	EmitLn("EOR (SP)+,D0")
}

// IsOrOp Recognizes a Boolean OrOp
func IsOrOp(r rune) bool {
	return strings.ContainsRune("|~", r)
}

// BoolFactor Parses and Translates a Boolean Factor
func BoolFactor() {
	if IsBoolean(Look) {
		if GetBoolean() {
			EmitLn("MOVE #-1,D0")
		} else {
			EmitLn("CLR D0")
		}
	} else {
		Relation()
	}
}

// BoolTerm Parses and Translates a Boolean Term
func BoolTerm() {
	NotFactor()
	for Look == '&' {
		EmitLn("MOVE D0, -(SP)")
		Match('&')
		NotFactor()
		EmitLn("AND (SP)+,D0")
	}
}

// NotFactor Parses and Translates a Boolean Factor with NOT
func NotFactor() {
	if Look == '!' {
		Match('!')
		BoolFactor()
		EmitLn("EOR #-1,Do")
	} else {
		BoolFactor()
	}
}

// IsRelOp Recognizes a RelOp
func IsRelOp(r rune) bool {
	return strings.ContainsRune("=#<>", r)
}

// Equals Recognizes and Translates a Relational "Equals"
func Equals() {
	Match('=')
	Expression()
	EmitLn("CMP (SP)+,D0")
	EmitLn("SEQ D0")
}

// NotEquals Recognizes and Translates a Relational "Not Equals"
func NotEquals() {
	Match('#')
	Expression()
	EmitLn("CMP (SP)+,D0")
	EmitLn("SNE D0")
}

// Less Recognizes and Translates a Relational "Less Than"
func Less() {
	Match('<')
	Expression()
	EmitLn("CMP (SP)+,D0")
	EmitLn("SGE D0")
}

// Greater Recognizes and Translates a Relational "Greater Than"
func Greater() {
	Match('>')
	Expression()
	EmitLn("CMP (SP)+,D0")
	EmitLn("SLE D0")
}

// Relation Parses and  Translates a Relation
func Relation() {
	Expression()
	if IsRelOp(Look) {
		EmitLn("MOVE D0,-(SP)")
		switch Look {
		case '=':
			Equals()
		case '#':
			NotEquals()
		case '<':
			Less()
		case '>':
			Greater()
		}
		EmitLn("TST D0")
	}
}

// Expression Parses and Translates a Math Expression
func Expression() {
	Term()
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

// IsAddOp Recognizes an AddOp
func IsAddOp(r rune) bool {
	return strings.ContainsRune("+-", r)
}

// Term Parses and Translates a Math Term
func Term() {
	SignedFactor()
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
		EmitLn("MOVE #" + string(GetNum()) + ",D0")
	}
}

// SignedFactor Parses and Translates the First Math Factor
func SignedFactor() {
	switch Look {
	case '+':
		GetChar()
	case '-':
		GetChar()
		if IsDigit(Look) {
			EmitLn("MOVE #-" + string(GetNum()) + ",D0")
		} else {
			Factor()
			EmitLn("NEG D0")
		}
	default:
		Factor()
	}
}

// Ident Parses and Translates an Identifier
func Ident() {
	name := GetName()
	if Look == '(' {
		Match('(')
		Match(')')
		EmitLn("BSR " + string(name))
	} else {
		EmitLn("MOVE " + string(name) + "(PC),D0")
	}
}

// Init Initializes
func Init() {
	GetChar()
}

// Go starts the execution of this chapter
func Go() {
	Init()
	BoolExpression()
}
