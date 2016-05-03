/* Sample test
ba
wb
lc
B
a=5
b=a*100
c=b/a
..etc..

.
*/

package types

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

// go doesn't have a built in abs function for int64
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

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

// DumpTable Dumps the Sybol Table
func DumpTable() {
	for i := 'A'; i <= 'Z'; i++ {
		if ST[i] != '?' {
			util.WriteLine(string(i) + " " + string(ST[i]))
		}
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
	return
}

// GetNum Gets a Number
func GetNum() (val int64) {
	if !IsDigit(Look) {
		Expected("Integer")
	}
	val = 0
	for IsDigit(Look) {
		digit, err := strconv.Atoi(string(Look))
		if err != nil {
			panic(err)
		}
		val = 10*val + int64(digit)
		GetChar()
	}
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

// TypeOf Reports Type of a Variable
func TypeOf(n rune) rune {
	return ST[n]
}

// InTable Reports if a Variable is in the Table
func InTable(n rune) bool {
	return TypeOf(n) != '?'
}

// CheckDup Checks for a Duplicate Variable Name
func CheckDup(n rune) {
	if InTable(n) {
		Abort("Duplicate Name " + string(n))
	}
}

// AddEntry Adds Entry to Table
func AddEntry(n rune, t rune) {
	CheckDup(n)
	ST[n] = t
}

// Alloc Allocates Storage for a Variable
func Alloc(n, t rune) {
	AddEntry(n, t)
	AllocVar(n, t)
}

// Decl Parses and Translates a Data Declaration
func Decl() {
	typ := GetName()
	Alloc(GetName(), typ)
}

// TopDecls Parses and Translates Global Declarations
func TopDecls() {
	for Look != 'B' {
		switch Look {
		case 'b', 'w', 'l':
			Decl()
		default:
			Abort("Unrecogized Keyword " + string(Look))
		}
		Fin()
	}
}

// AllocVar Generates Code for Allocation of a Variable
func AllocVar(n, t rune) {
	util.WriteLine(string(n) + ":\tDC." + string(t) + " 0")
}

// LoadVar Loads a Variable to Primary Register
func LoadVar(name, typ rune) {
	Move(typ, string(name)+"(PC)", "D0")
}

// Move Generates a Move Instruction
func Move(size rune, source, dest string) {
	EmitLn("MOVE." + string(size) + " " + source + "," + dest)
}

// IsVarType Recognizes a Legal Variable Type
func IsVarType(r rune) bool {
	return strings.ContainsRune("BWL", r)
}

// VarType Gets a Variable Type from the Symbol Table
func VarType(name rune) rune {
	typ := TypeOf(name)
	if !IsVarType(typ) {
		Abort("Identifier " + string(name) + " is not a variable")
	}
	return typ
}

// Load Loads a Variable to the Primary Register
func Load(name rune) rune {
	typ := VarType(name)
	LoadVar(name, typ)
	return typ
}

// StoreVar Stores Primary to Variable
func StoreVar(name, typ rune) {
	EmitLn("LEA " + string(name) + "(PC),A0")
	Move(typ, "DO", "(A0)")
}

// Store Stores a Variable from the Primary Register
func Store(name, t1 rune) {
	t2 := VarType(name)
	//note tutorial does not have revised form of Convert with the third
	//parameter for register. Assuming this should be D0 as it is the Primary
	//Register
	Convert(t1, t2, "D0")
	StoreVar(name, t2)
}

// Term Parses and Translates a Math Term
func Term() (typ rune) {
	typ = Factor()
	for IsMulOp(Look) {
		Push(typ)
		switch Look {
		case '*':
			typ = Multiply(typ)
		case '/':
			typ = Divide(typ)
		}
	}
	return
}

// Expression Parses and Translates an Expression
func Expression() (typ rune) {
	if IsAddOp(Look) {
		typ = UnOp()
	} else {
		typ = Term()
	}
	for IsAddOp(Look) {
		Push(typ)
		switch Look {
		case '+':
			typ = Add(typ)
		case '-':
			typ = Subtract(typ)
		}
	}
	return
}

// UnOp Processes a Term with Leading Unary Operator
func UnOp() rune {
	Clear()
	return 'W'
}

// Clear Clears the Primary Register
// Note this is not defined in Chapter 14, but assuming the same function as
// from Chapter 12
func Clear() {
	EmitLn("CLR D0")
}

// Push Pushes Primary onto Stack
func Push(size rune) {
	Move(size, "D0", "-(SP)")
}

// Add Recognizes and Translates an Add
func Add(t1 rune) rune {
	Match('+')
	return PopAdd(t1, Term())
}

// Subtract Recognizes and Translates a Subtract
func Subtract(t1 rune) rune {
	Match('-')
	return PopSub(t1, Term())
}

// Pop Pops Stack into Secondary Register
func Pop(size rune) {
	Move(size, "(SP)+", "D7")
}

// Convert Convers a Data Item from One Type to Another
func Convert(source, dest rune, reg string) {
	if source != dest {
		if source == 'B' {
			EmitLn("AND.W #$FF," + reg)
		}
		if dest == 'L' {
			EmitLn("EXT.L " + reg)
		}
	}
}

// Promote Promotes the Size of a Register Value
func Promote(t1, t2 rune, reg string) (typ rune) {
	typ = t1
	if t1 != t2 {
		if t1 == 'B' || (t1 == 'W' && t2 == 'L') {
			Convert(t1, t2, reg)
			typ = t2
		}
	}
	return
}

// SameType Forces both Arguments to Same Type
func SameType(t1, t2 rune) rune {
	t1 = Promote(t1, t2, "D7")
	return Promote(t2, t1, "D0")
}

// PopAdd Generates Code to Add Primary to the Stack
func PopAdd(t1, t2 rune) rune {
	Pop(t1)
	t2 = SameType(t1, t2)
	GenAdd(t2)
	return t2
}

// PopSub Generates Code to Subtract Primary from the Stack
func PopSub(t1, t2 rune) rune {
	Pop(t1)
	t2 = SameType(t1, t2)
	GenSub(t2)
	return t2
}

// GenAdd Adds Top of Stack to Primary
func GenAdd(size rune) {
	EmitLn("ADD." + string(size) + " D7,D0")
}

// GenSub Subtracts Primary from Top of Stack
func GenSub(size rune) {
	EmitLn("SUB." + string(size) + " D7,D0")
	EmitLn("NEG." + string(size) + " D0")
}

// Factor Parses and Translates a Factor
func Factor() (typ rune) {
	if Look == '(' {
		Match('(')
		typ = Expression()
		Match(')')
	} else if IsAlpha(Look) {
		typ = Load(GetName())
	} else {
		typ = LoadNum(GetNum())
	}
	return
}

// Multiply Recognizes and Translates a Multiply
func Multiply(t1 rune) rune {
	Match('*')
	return PopMul(t1, Factor())
}

// Divide Recognizes and Translates a Divide
func Divide(t1 rune) rune {
	Match('/')
	return PopDiv(t1, Factor())
}

// PopMul Generates Code to Multiply Primary by Stack
func PopMul(t1, t2 rune) rune {
	Pop(t1)
	t := SameType(t1, t2)
	Convert(t, 'W', "D7")
	Convert(t, 'W', "D0")
	if t == 'L' {
		GenLongMult()
	} else {
		GenMult()
	}
	if t == 'B' {
		return 'W'
	}
	return 'L'
}

// PopDiv Generates Code to Divide Stack by the Primary
func PopDiv(t1, t2 rune) rune {
	Pop(t1)
	Convert(t1, t2, "D7")
	if t1 == 'L' || t2 == 'L' {
		Convert(t2, 'L', "D0")
		GenLongDiv()
		return 'L'
	}
	Convert(t2, 'W', "D0")
	GenDiv()
	return t1
}

// GenDiv Divides Top of Stack by Primary
func GenDiv() {
	EmitLn("DIVS D0,D7")
	Move('W', "D7", "D0")
}

// GenLongDiv Divides Top of Stack by Primary
func GenLongDiv() {
	EmitLn("JSR DIV32")
}

// GenMult Multiplies Top of Stack by Primary (Word)
func GenMult() {
	EmitLn("MULS D7,D0")
}

// GenLongMult Multiplies Top of Stack by Primary (Long)
func GenLongMult() {
	EmitLn("JSR MUL32")
}

// Assignment Parses and Translates an Assignment Statement
func Assignment() {
	name := GetName()
	Match('=')
	Store(name, Expression())
}

// Block Parses and Transaltes a Block of Statements
func Block() {
	for Look != '.' {
		Assignment()
		Fin()
	}
}

// LoadNum Loads a Constant to the Primary Register
func LoadNum(n int64) rune {
	var typ rune
	if abs(n) <= 127 {
		typ = 'B'
	} else if abs(n) <= 32627 {
		typ = 'W'
	} else {
		typ = 'L'
	}
	LoadConst(n, typ)
	return typ
}

// LoadConst Loads a Constant to the Primary Register
func LoadConst(n int64, typ rune) {
	temp := strconv.FormatInt(n, 10)
	Move(typ, "#"+temp, "D0")
}

// Init Initializes
func Init() {
	ST = make(map[rune]rune)
	for i := 'A'; i <= 'Z'; i++ {
		ST[i] = '?'
	}
	GetChar()
	SkipWhite()
}

// Go starts the execution of this chapter
func Go() {
	Init()
	TopDecls()
	Match('B')
	Fin()
	Block()
	DumpTable()
}
