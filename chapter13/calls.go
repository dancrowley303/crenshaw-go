//this is the version of the package with call-by-value semantics

/* Sample test
va
vb
vc
pd(e,f)
vh
vi
vj
b
h=e
i=f
j=a
e
Px
b
d(b,c)
e.
*/

package calls

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/dcw303/crenshaw-go/util"
)

// Look is a Lookahead character
var Look rune

// ST is a Symbol Table
var ST map[rune]rune

// Params is a Table of Function Parameters
var Params map[rune]int

// NumParams is the Number of Parameters
var NumParams int

// Base is Used to Compute Stack Offsets
var Base int

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

// Duplicate Reports a Duplicate Identifier
func Duplicate(n string) {
	Abort("Duplicate Identifier " + n)
}

// TypeOf Gets Type of Symbol
func TypeOf(n rune) rune {
	if IsParam(n) {
		return 'f'
	}
	return ST[n]
}

// InTable Looks for Symbol in Table
func InTable(n rune) bool {
	return ST[n] != ' '
}

// AddEntry Adds a New Entry to Symbol Table
func AddEntry(name rune, t rune) {
	if InTable(name) {
		Duplicate(string(name))
	}
	ST[name] = t
}

// CheckVar Checks an Entry to Make Sure It's a Variable
func CheckVar(name rune) {
	if !InTable(name) {
		Undefined(string(name))
	}
	if TypeOf(name) != 'v' {
		Abort(string(name) + " is not a variable")
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
	// SPACE / TAB
	return r == 0x20 || r == 0x09
}

// SkipWhite Skips Over Leading White Space
func SkipWhite() {
	for IsWhite(Look) {
		GetChar()
	}
}

// Fin Skips Over an End-of-Line
func Fin() {
	if Look == 0x0D {
		GetChar()
		if Look == 0x0A {
			GetChar()
		}
	}
}

// Match Matches a Specific Input Character
func Match(x rune) {
	if Look == x {
		GetChar()
	} else {
		Expected(strconv.QuoteRuneToASCII(x))
	}
}

// GetName Gets an Identifier
func GetName() (r rune) {
	if !IsAlpha(Look) {
		Expected("Name")
	}
	r = unicode.ToUpper(Look)
	GetChar()
	SkipWhite()
	return
}

// GetNum Gets a Number
func GetNum() (r rune) {
	if !IsDigit(Look) {
		Expected("Integer")
	}
	r = Look
	GetChar()
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

// PostLabel Posts a Label to Outputs
func PostLabel(l string) {
	util.WriteLine(l + ":")
}

// LoadVar Loads a Variable to Primary Register
func LoadVar(name rune) {
	CheckVar(name)
	EmitLn("MOVE " + string(name) + "(PC),D0")
}

// StoreVar Stores the Primary Register
func StoreVar(name rune) {
	CheckVar(name)
	EmitLn("LEA " + string(name) + "(PC),A0")
	EmitLn("MOVE D0,(A0)")
}

// Expression Parses and Translates an Expression
// Vestigal Version
func Expression() {
	name := GetName()
	if IsParam(name) {
		LoadParam(ParamNumber(name))
	} else {
		LoadVar(name)
	}
}

// Assignment Parses and Translates an Assignment Statement
func Assignment(name rune) {
	Match('=')
	Expression()
	if IsParam(name) {
		StoreParam(ParamNumber(name))
	} else {
		StoreVar(name)
	}
}

// DoBlock Parses and Translates a Block of Statements
func DoBlock() {
	for Look != 'e' {
		AssignOrProc()
		Fin()
	}
}

// BeginBlock Parses and Translates a Begin-Block
func BeginBlock() {
	Match('b')
	Fin()
	DoBlock()
	Match('e')
	Fin()
}

// Alloc Allocates Storage for a Variable
func Alloc(n rune) {
	if InTable(n) {
		Duplicate(string(n))
	}
	ST[n] = 'v'
	util.WriteLine(string(n) + ":\tDC 0")
}

// Decl Parses and Translates a Data Declaration
func Decl() {
	Match('v')
	Alloc(GetName())
}

// TopDecls Parses and Translates Global Declarations
func TopDecls() {
	for Look != '.' {
		switch Look {
		case 'v':
			Decl()
		case 'p':
			DoProc()
		case 'P':
			DoMain()
		default:
			Abort("Unrecognized Keyword " + string(Look))
		}
		Fin()
	}
}

// Return Emits an RTS Instruction
func Return() {
	EmitLn("RTS")
}

// DoProc Parses and Translates a Procedure Declaration
func DoProc() {
	Match('p')
	n := GetName()
	if InTable(n) {
		Duplicate(string(n))
	}
	ST[n] = 'p'
	FormalList()
	k := LocDecls()
	ProcProlog(n, k)
	BeginBlock()
	ProcEpilog()
	ClearParams()
}

// DoMain Parses and Translates a Main Program
func DoMain() {
	Match('P')
	n := GetName()
	Fin()
	if InTable(n) {
		Duplicate(string(n))
	}
	Prolog()
	BeginBlock()
}

// AssignOrProc Decides if a Statement is an Assignment or Procedure call
func AssignOrProc() {
	name := GetName()
	switch TypeOf(name) {
	case ' ':
		Undefined(string(name))
	case 'v', 'f':
		Assignment(name)
	case 'p':
		CallProc(name)
	default:
		Abort("Identifier " + string(name) + " Cannot Be Used Here")
	}
}

// CallProc Processes a Procedure Call
func CallProc(name rune) {
	n := ParamList()
	Call(name)
	CleanStack(n)
}

// Call Generates code to Emit BSR instruction
func Call(name rune) {
	EmitLn("BSR " + string(name))
}

// FormalList Processes the Formal Parameter List of a Procedure
func FormalList() {
	Match('(')
	if Look != ')' {
		FormalParam()
		for Look == ',' {
			Match(',')
			FormalParam()
		}
	}
	Match(')')
	Fin()
	Base = NumParams
	NumParams += 4
}

// FormalParam Processes a Formal Parameter
func FormalParam() {
	AddParam(GetName())
}

// Param Processes an Actual Parameter
func Param() {
	Expression()
	Push()
}

// ParamList Processes the Parameter List for a Procedure Call
func ParamList() int {
	n := 0
	Match('(')
	if Look != ')' {
		Param()
		n++
		for Look == ',' {
			Match(',')
			Param()
			n++
		}
	}
	Match(')')
	return 2 * n
}

// ClearParams Initializes Parameter Table to Null
func ClearParams() {
	for i := 'A'; i <= 'Z'; i++ {
		Params[i] = 0
	}
	NumParams = 0
}

// ParamNumber Finds the Parameter Number
func ParamNumber(n rune) int {
	return Params[n]
}

// IsParam Sees if an Identifer is a Parameter
func IsParam(n rune) bool {
	return Params[n] != 0
}

// AddParam Adds a New Parameter to Table
func AddParam(name rune) {
	if IsParam(name) {
		Duplicate(string(name))
	}
	NumParams++
	Params[name] = NumParams
}

// LoadParam Loads a Parameter to the Primary Register
func LoadParam(n int) {
	offset := 8 + 2*(Base-n)
	Emit("MOVE ")
	util.WriteLine(strconv.Itoa(offset) + "(A6),D0")
}

// StoreParam Stores a Parameter from the Primary Register
func StoreParam(n int) {
	offset := 8 + 2*(Base-n)
	Emit("MOVE D0,")
	util.WriteLine(strconv.Itoa(offset) + "(A6)")
}

// Push Pushes the Primary Register to the Stack
func Push() {
	EmitLn("MOVE D0,-(SP)")
}

// CleanStack Adjusts the Stack Pointer Upwards by N bytes
func CleanStack(n int) {
	if n > 0 {
		Emit("ADD #")
		util.WriteLine(strconv.Itoa(n) + ",SP")
	}
}

// ProcProlog Writes the Prolog for a Procedure
func ProcProlog(n rune, k int) {
	PostLabel(string(n))
	Emit("LINK A6,#")
	util.WriteLine(strconv.Itoa(-2 * k))
}

// ProcEpilog Writes the Epilog for a Procedure
func ProcEpilog() {
	EmitLn("UNLK A6")
	EmitLn("RTS")
}

// Prolog Writes the Prolog
func Prolog() {
	PostLabel("MAIN")
}

// Epilog Writes the Epilog
func Epilog() {
	EmitLn("RTS")
}

// LocDecl Parses and Translates a Local Data Declaration
func LocDecl() {
	Match('v')
	AddParam(GetName())
	Fin()
}

// LocDecls Parses and Translates Local Declarations
func LocDecls() int {
	n := 0
	for Look == 'v' {
		LocDecl()
		n++
	}
	return n
}

// Init Initializes
func Init() {
	GetChar()
	SkipWhite()
	ST = make(map[rune]rune)
	Params = make(map[rune]int)
	for i := 'A'; i <= 'Z'; i++ {
		ST[i] = ' '
	}
	ClearParams()
}

// Go starts the execution of this chapter
func Go() {
	Init()
	TopDecls()
	Epilog()
}
