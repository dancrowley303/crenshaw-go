// Second version of the package created for the section of the tutorial
// labelled "Merging With Control Constructs"

/* Sample test
a=1
j=2
k=3
ia+j=k
z=9
e
e
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

// LCount is a Label Counter
var LCount int

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

// NewLabel Generates a Unique Label
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

// Block Recognizes and Translates a Statement Block
func Block(l string) {
	for Look != 'e' && Look != 'l' && Look != 'u' {
		Fin()
		switch Look {
		case 'i':
			DoIf(l)
		case 'w':
			DoWhile()
		case 'p':
			DoLoop()
		case 'r':
			DoRepeat()
		case 'd':
			DoDo()
		case 'b':
			DoBreak(l)
		case 0X0D:
			//Do Nothing - This is not in the tutorial but is needed to stop 3 or
			//more CRs feeding a CR into Other()
		default:
			Assignment()
		}
		Fin()
	}
}

// DoIf Recognizes and Translates an IF Construct
func DoIf(l string) {
	Match('i')
	BoolExpression()
	l1 := NewLabel()
	l2 := l1
	EmitLn("BEQ " + l1)
	Block(l)
	if Look == 'l' {
		Match('l')
		l2 = NewLabel()
		EmitLn("BRA " + l2)
		PostLabel(l1)
		Block(l)
	}
	Match('e')
	PostLabel(l2)
}

// DoWhile Parses and Translates a WHILE Statement
func DoWhile() {
	Match('w')
	l1 := NewLabel()
	l2 := NewLabel()
	PostLabel(l1)
	BoolExpression()
	EmitLn("BEQ " + l2)
	Block(l2)
	Match('e')
	EmitLn("BRA " + l1)
	PostLabel(l2)
}

// DoLoop Parses and Translates a LOOP Statement
func DoLoop() {
	Match('p')
	l1 := NewLabel()
	l2 := NewLabel()
	PostLabel(l1)
	Block(l2)
	Match('e')
	EmitLn("BRA " + l1)
	PostLabel(l2)
}

// DoRepeat Parses and Translates a REPEAT Statement
func DoRepeat() {
	Match('r')
	l1 := NewLabel()
	l2 := NewLabel()
	PostLabel(l1)
	Block(l2)
	Match('u')
	BoolExpression()
	EmitLn("BEQ " + l1)
	PostLabel(l2)
}

// DoFor Parses and Translates a FOR Statement
func DoFor() {
	Match('f')
	l1 := NewLabel()
	l2 := NewLabel()
	name := GetName()
	Match('=')
	Expression()
	EmitLn("SUBQ #1,D0")
	EmitLn("LEA " + string(name) + "(PC),A0")
	EmitLn("MOVE D0,(A0)")
	Expression()
	EmitLn("MOVE D0,-(SP)")
	PostLabel(l1)
	EmitLn("LEA " + string(name) + "(PC),A0")
	EmitLn("MOVE (A0),D0")
	EmitLn("ADDQ #1,D0")
	EmitLn("MOVE D0,(A0)")
	EmitLn("CMP (SP),D0")
	EmitLn("BGT " + l2)
	Block(l2)
	Match('e')
	EmitLn("BRA " + l1)
	PostLabel(l2)
	EmitLn("ADDQ #2,SP")
}

// DoDo Parses and Translates a DO Statement
func DoDo() {
	Match('d')
	l1 := NewLabel()
	l2 := NewLabel()
	Expression()
	EmitLn("SUBQ #1,D0")
	PostLabel(l1)
	EmitLn("MOVE D0,-(SP)")
	Block(l2)
	EmitLn("MOVE (SP)+,D0")
	EmitLn("DBRA D0," + l1)
	EmitLn("SUBQ #2,SP")
	PostLabel(l2)
	EmitLn("ADDQ #2,SP")
}

// DoBreak Recognizes and Translates a BREAK
func DoBreak(l string) {
	Match('b')
	if l != "" {
		EmitLn("BRA " + l)
	} else {
		Abort("No loop to break from")
	}
}

// Other Recognizes and Translates an "Other"
func Other() {
	EmitLn(string(GetName()))
}

// DoProgram Recognizes and Translates a Program
func DoProgram() {
	Block("")
	if Look != 'e' {
		Expected("END")
	}
	EmitLn("END")
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

// Fin Skips a CRLF
func Fin() {
	if Look == 0x0D { //CR
		GetChar()
	}
	if Look == 0x0A { //LF - Note this never gets hit from stdin on Windows
		GetChar()
	}
}

// Assignment Parses and Translates an Assignment Statement
func Assignment() {
	name := GetName()
	Match('=')
	BoolExpression()
	EmitLn("LEA " + string(name) + "(PC),A0")
	EmitLn("MOVE D0,(A0)")
}

// Init Initializes
func Init() {
	LCount = 1
	GetChar()
}

// Go starts the execution of this chapter
func Go() {
	Init()
	DoProgram()
}
