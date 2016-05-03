/* Sample test
program
var abc
var def=200
begin
abc=50
read(abc)
if abc = def
def=20+abc
else
def=5+def
endif
write(def)
end.
*/

package tiny

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/dcw303/crenshaw-go/util"
)

// LCount is a Label Counter
var LCount int

// NEntry is the Next Entry in the Symbol Table
var NEntry int

// Look is a Lookahead character
var Look rune

// Token is an Encoded Token
var Token rune

// Value is an Unencoded Token
var Value string

// MaxEntry is the number of Entries allowed in the Symbol Table
const MaxEntry = 100

// ST is the Symbol Table
var ST []string

// SType is the Symbol Type Table
var SType []rune

// Definition of Keywords and Token Types

// NKW is the Number of Keywords
const NKW = 11

// NKW1 is the Number of Keywords + 1 (?)
const NKW1 = 12

// KWList is the Keyword List
var KWList = []string{"IF", "ELSE", "ENDIF", "WHILE", "ENDWHILE", "READ",
	"WRITE", "VAR", "BEGIN", "END", "PROGRAM"}

// KWCode is the Keyword Code
const KWCode string = "xileweRWvbep"

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

// Undefined Reports an Undefined Identifier
func Undefined(n string) {
	Abort("Undefined Identifier " + n)
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

// IsMulOp Recognizes a MulOp
func IsMulOp(r rune) bool {
	return strings.ContainsRune("*/", r)
}

// IsOrOp Recognizes a Boolean OrOp
func IsOrOp(r rune) bool {
	return strings.ContainsRune("|~", r)
}

// IsRelOp Recognizes a RelOp
func IsRelOp(r rune) bool {
	return strings.ContainsRune("=#<>", r)
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

// NewLine Skips Over an End-of-Line
func NewLine() {
	for Look == 0x0D { // CR
		GetChar()
		if Look == 0x0A { // LF
			GetChar()
		}
		SkipWhite()
	}
}

// Match Matches a Specific Input Character
func Match(x rune) {
	if Look == x {
		GetChar()
	} else {
		Expected(strconv.QuoteRuneToASCII(x))
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

// Locate Locates a Symbol in Table
func Locate(n string) int {
	return Lookup(ST, n, MaxEntry)
}

// InTable Looks for Symbol in Table
func InTable(n string) bool {
	// Original code calls Lookup, but I'm using Locate as otherwise it has
	// no use
	return Locate(n) != -1
}

// AddEntry Adds a New Entry to Symbol Table
func AddEntry(n string, t rune) {
	if InTable(n) {
		Abort("Duplicate Identifier " + n)
	}
	if NEntry == MaxEntry {
		Abort("Symbol Table Full")
	}
	NEntry++
	ST[NEntry] = n
	SType[NEntry] = t
}

// GetName Gets an Identifier
func GetName() {
	NewLine()
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
func GetNum() (val int) {
	if !IsDigit(Look) {
		Expected("Integer")
	}
	NewLine()
	val = 0
	for IsDigit(Look) {
		digit, err := strconv.Atoi(string(Look))
		if err != nil {
			panic(err)
		}
		val = 10*val + digit
		GetChar()
	}
	SkipWhite()
	return
}

// Scan Gets an Identifier and Scans it for Keywords
func Scan() {
	GetName()
	Token = rune(KWCode[Lookup(KWList, Value, NKW)+1])
}

// MatchString Matches a Specific Input String
func MatchString(x string) {
	if Value != x {
		Expected(x)
	}
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

// Clear Clears the Primary Register
func Clear() {
	EmitLn("CLR D0")
}

// Negate Negates the Primary Register
func Negate() {
	EmitLn("NEG D0")
}

// NotIt Complements the Primary Register
func NotIt() {
	EmitLn("NOT D0")
}

// LoadConst Loads a Constant Value to Primary Register
func LoadConst(n int) {
	Emit("MOVE #")
	util.WriteLine(strconv.Itoa(n) + ",D0")
}

// LoadVar Loads a Variable to Primary Register
func LoadVar(name string) {
	if !InTable(name) {
		Undefined(name)
	}
	EmitLn("MOVE " + name + "(PC),D0")
}

// Push Pushes Primary onto Stack
func Push() {
	EmitLn("MOVE D0,-(SP)")
}

// PopAdd Adds Top of Stack to Primary
func PopAdd() {
	EmitLn("ADD (SP)+,D0")
}

// PopSub Subtracts Primary from Top of Stack
func PopSub() {
	EmitLn("SUB (SP)+,D0")
	EmitLn("NEG D0")
}

// PopMul Multiplies Top of Stack by Primary
func PopMul() {
	EmitLn("MULS (SP)+,D0")
}

// PopDiv Divides Top of Stack by Primary
func PopDiv() {
	EmitLn("MOVE (SP)+,D7")
	EmitLn("EXT.L D7")
	EmitLn("DIVS D0,D7")
	EmitLn("MOVE D7,D0")
}

// PopAnd ANDs Top of Stack with Primary
func PopAnd() {
	EmitLn("AND (SP)+,D0")
}

// PopOr ORs Top of Stack with Primary
func PopOr() {
	EmitLn("OR (SP)+,D0")
}

// PopXor XORs Top of Stack with Primary
func PopXor() {
	EmitLn("EOR (SP)+,D0")
}

// PopCompare Compares Top of Stack with Primary
func PopCompare() {
	EmitLn("CMP (SP)+,D0")
}

// SetEqual Sets D0 if Compare was =
func SetEqual() {
	EmitLn("SEQ D0")
	EmitLn("EXT D0")
}

// SetNEqual Sets D0
func SetNEqual() {
	EmitLn("SNE D0")
	EmitLn("EXT D0")
}

// SetGreater Sets D0 If Compare was >
func SetGreater() {
	EmitLn("SLT D0")
	EmitLn("EXT D0")
}

// SetLess Sets D0 if Compare was <
func SetLess() {
	EmitLn("SGT D0")
	EmitLn("EXT D0")
}

// SetLessOrEqual Sets D0 if Compare was <= 0
func SetLessOrEqual() {
	EmitLn("SGE D0")
	EmitLn("EXT D0")
}

// SetGreaterOrEqual Sets D0 if Compare was >= 0
func SetGreaterOrEqual() {
	EmitLn("SLE D0")
	EmitLn("EXT D0")
}

// Store Stores Primary to Variable
func Store(name string) {
	if !InTable(name) {
		Undefined(name)
	}
	EmitLn("LEA " + name + "(PC),A0")
	EmitLn("MOVE D0,(A0)")
}

// Branch Branches Unconditional
func Branch(l string) {
	EmitLn("BRA " + l)
}

// BranchFalse Branches false
func BranchFalse(l string) {
	EmitLn("TST D0")
	EmitLn("BEQ " + l)
}

// ReadVar Reads Variable to Primary Register
func ReadVar() {
	EmitLn("BSR READ")
	Store(Value)
}

// WriteVar Writes Variable from Primary Register
func WriteVar() {
	EmitLn("BSR WRITE")
}

// Header Writes Header Info
func Header() {
	util.WriteLine("WARMST\t'EQU $A01E'")
}

// Prolog Writes the Prolog
func Prolog() {
	PostLabel("MAIN")
}

// Epilog Writes the Epilog
func Epilog() {
	util.WriteLine("DC WARMST")
	util.WriteLine("END MAIN")
}

// Factor Parses and Translates a Math Factor
func Factor() {
	switch {
	case Look == '(':
		Match('(')
		BoolExpression()
		Match(')')
	case IsAlpha(Look):
		GetName()
		LoadVar(Value)
	default:
		LoadConst(GetNum())
	}
}

// NegFactor Parses and Translates a Negative Factor
func NegFactor() {
	Match('-')
	if IsDigit(Look) {
		LoadConst(-GetNum())
	} else {
		Factor()
		Negate()
	}
}

// FirstFactor Parses and Translates a Leading Factor
func FirstFactor() {
	switch Look {
	case '+':
		Match('+')
		Factor()
	case '-':
		NegFactor()
	default:
		Factor()
	}
}

// Multiply Recognizes and Translates a Multiply
func Multiply() {
	Match('*')
	Factor()
	PopMul()
}

// Divide Recognizes and Translates a Divide
func Divide() {
	Match('/')
	Factor()
	PopDiv()
}

// Term1 Is Common Code Used by Term and FirstTerm
func Term1() {
	for IsMulOp(Look) {
		Push()
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

// FirstTerm Parses and Translates a Leading Term
func FirstTerm() {
	FirstFactor()
	Term1()
}

// Add Recognizes and Translates an Add
func Add() {
	Match('+')
	Term()
	PopAdd()
}

// Subtract Recognizes and Translates a Subtract
func Subtract() {
	Match('-')
	Term()
	PopSub()
}

// Expression Parses and Translates a Math Expression
func Expression() {
	FirstTerm()
	for IsAddOp(Look) {
		Push()
		switch Look {
		case '+':
			Add()
		case '-':
			Subtract()
		}
	}
}

// Equals Recognizes and Translates a Relational "Equals"
func Equals() {
	Match('=')
	Expression()
	PopCompare()
	SetEqual()
}

// LessOrEqual Recognizes and Translates a Relational "Less Than or Equal"
func LessOrEqual() {
	Match('=')
	Expression()
	PopCompare()
	SetLessOrEqual()
}

// NotEqual Recognizes and Translates a Relational "Not Equals"
func NotEqual() {
	Match('>')
	Expression()
	PopCompare()
	SetNEqual()
}

// Less Recognizes and Translates a Relational "Less Than"
func Less() {
	Match('<')
	switch Look {
	case '=':
		LessOrEqual()
	case '>':
		NotEqual()
	default:
		Expression()
		PopCompare()
		SetLess()
	}
}

// Greater Recognizes and Translates a Relational "Greater Than"
func Greater() {
	Match('>')
	if Look == '=' {
		Match('=')
		Expression()
		PopCompare()
		SetGreaterOrEqual()
	} else {
		Expression()
		PopCompare()
		SetGreater()
	}
}

// Relation Parses and  Translates a Relation
func Relation() {
	Expression()
	if IsRelOp(Look) {
		Push()
		switch Look {
		case '=':
			Equals()
		case '<':
			Less()
		case '>':
			Greater()
		}
	}
}

// NotFactor Parses and Translates a Boolean Factor with Leading NOT
func NotFactor() {
	if Look == '!' {
		Match('!')
		Relation()
		NotIt()
	} else {
		Relation()
	}
}

// BoolTerm Parses and Translates a Boolean Term
func BoolTerm() {
	NotFactor()
	for Look == '&' {
		Push()
		Match('&')
		NotFactor()
		PopAdd()
	}
}

// BoolOr Recognizes and Translates a Boolean OR
func BoolOr() {
	Match('|')
	BoolTerm()
	PopOr()
}

// BoolXor Recognizes and Translates an Exclusive OR
func BoolXor() {
	Match('~')
	BoolTerm()
	PopXor()
}

// BoolExpression Parses and Translates a Boolean Expression
func BoolExpression() {
	BoolTerm()
	for IsOrOp(Look) {
		Push()
		switch Look {
		case '|':
			BoolOr()
		case '~':
			BoolXor()
		}
	}
}

// Assignment Parses and Translates an Assignment Statement
func Assignment() {
	name := Value
	Match('=')
	BoolExpression()
	Store(name)
}

// DoIf Recognizes and Translates an IF Construct
func DoIf() {
	BoolExpression()
	l1 := NewLabel()
	l2 := l1
	BranchFalse(l1)
	Block()
	if Token == 'l' {
		l2 = NewLabel()
		Branch(l2)
		PostLabel(l1)
		Block()
	}
	PostLabel(l2)
	MatchString("ENDIF")
}

// DoWhile Parses and Translates a WHILE Statement
func DoWhile() {
	l1 := NewLabel()
	l2 := NewLabel()
	PostLabel(l1)
	BoolExpression()
	BranchFalse(l2)
	Block()
	MatchString("ENDWHILE")
	Branch(l1)
	PostLabel(l2)
}

// DoRead Processes a Read Statement
func DoRead() {
	Match('(')
	GetName()
	ReadVar()
	for Look == ',' {
		Match(',')
		GetName()
		ReadVar()
	}
	Match(')')
}

// DoWrite Processes a Write Statement
func DoWrite() {
	Match('(')
	Expression()
	WriteVar()
	for Look == ',' {
		Match(',')
		Expression()
		WriteVar()
	}
	Match(')')
}

// Block Parses and Translates a Block of Statements
func Block() {
	Scan()
	for Token != 'e' && Token != 'l' {
		switch Token {
		case 'i':
			DoIf()
		case 'w':
			DoWhile()
		case 'R':
			DoRead()
		case 'W':
			DoWrite()
		default:
			Assignment()
		}
		Scan()
	}
}

// Alloc Allocates Storage for a Variable
func Alloc(n string) {
	if InTable(n) {
		Abort("Duplicate Variable Name " + n)
	}
	AddEntry(n, 'v')
	util.Write(n + ":\tDC ")
	if Look == '=' {
		Match('=')
		if Look == '-' {
			util.Write(string(Look))
			Match('-')
		}
		util.WriteLine(strconv.Itoa(GetNum()))
	} else {
		util.WriteLine("0")
	}
}

// Decl Processes a Data Declaration
func Decl() {
	GetName()
	Alloc(Value)
	for Look == ',' {
		Match(',')
		GetName()
		Alloc(Value)
	}
}

// TopDecls Parses and Translates Global Declarations
func TopDecls() {
	Scan()
	for Token != 'b' {
		switch Token {
		case 'v':
			Decl()
		default:
			Abort("Unrecognized Keyword '" + Value + "'")
		}
		Scan()
	}
}

// Main Parses and Translates a Main PROGRAM
func Main() {
	MatchString("BEGIN")
	Prolog()
	Block()
	MatchString("END")
	Epilog()
}

// Prog Parses and Translates a Program
func Prog() {
	MatchString("PROGRAM")
	Header()
	TopDecls()
	Main()
	Match('.')
}

// Init Initializes
func Init() {
	ST = make([]string, MaxEntry)
	SType = make([]rune, MaxEntry)
	for i := 0; i < MaxEntry; i++ {
		//this is not necessary in Go as empty string is default val for string
		//slice anyway
		ST[i] = ""
		SType[i] = ' '
	}
	GetChar()
	Scan()
}

// Go starts the execution of this chapter
func Go() {
	Init()
	Prog()
	if Look != 0x0D {
		Abort("Unexpected data after '.'")
	}
}
