/* Sample tests
if else: aiceje
nested if else: aicijekeme
while: awjeke
loop: paibejeke
repeat until: rauke
for: afi=bece
do: dajke
*/

package branch

import (
	"strconv"
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

// Condition Parses and Translates a Boolean Condition
// This version is a dummy
func Condition() {
	EmitLn("<condition>")
}

// Expression Parses and Translates an Expression
// This version is a dummy
func Expression() {
	EmitLn("<expr>")
}

// Block Recognizes and Translates a Statement Block
func Block(l string) {
	for Look != 'e' && Look != 'l' && Look != 'u' {
		switch Look {
		case 'i':
			DoIf(l)
		case 'w':
			DoWhile()
		case 'p':
			DoLoop()
		case 'r':
			DoRepeat()
		case 'f':
			DoFor()
		case 'd':
			DoDo()
		case 'b':
			DoBreak(l)
		default:
			Other()
		}
	}
}

// DoIf Recognizes and Translates an IF Construct
func DoIf(l string) {
	Match('i')
	Condition()
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
	Condition()
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
	Condition()
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

// DoProgram Parses and Translates a Program
func DoProgram() {
	Block("")
	if Look != 'e' {
		Expected("End")
	}
	EmitLn("END")
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
