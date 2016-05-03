//Note: This code covers the C style and semicolon and comment parsing in this
//chapter. The code here matches the following style:
//(1) Semicolons are TERMINATORS, not seperators
//(2) Semicolons are NOT OPTIONAL
//(3) Comments are delimited by /* and */
//(4) /* */ Comments can be nested
//(5) //one sides style comments are also supported

//Sample test (ignore single line comments at start; there is a nested /* */)
//program
//var a,b,c;
//begin
//a=1*(b+4); //single comment
//b=b+3; /* comment
//spanning
//multiple
//lines*/
//c=5;
//end.

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

// TempChar is a Temporary Character
var TempChar = ' '

// ST is the Symbol Table
var ST []string

// SType is the Symbol Type Table
var SType []rune

// Definition of Keywords and Token Types

// NKW is the Number of Keywords
const NKW = 9

// NKW1 is the Number of Keywords + 1 (?)
const NKW1 = 10

// KWList is the Keyword List
var KWList = []string{"IF", "ELSE", "ENDIF", "WHILE", "ENDWHILE", "READ",
	"WRITE", "VAR", "END"}

// KWCode is the Keyword Code
const KWCode string = "xileweRWve"

// GetCharX Reads New Character From Input Stream
func GetCharX() {
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

// CheckIdent Checks to Make Sure the Current Token is an Identifier
func CheckIdent() {
	if Token != 'x' {
		Expected("Identifier")
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
	// SPACE / TAB / CR / LF / comments marker / one sided comment marker
	return r == 0x20 || r == 0x09 || r == 0x0D || r == 0xFF || r == 0xFE
}

// SkipWhite Skips Over Leading White Space
func SkipWhite() {
	for IsWhite(Look) {
		if Look == 0xFF {
			SkipComment()
		} else if Look == 0xFE {
			SkipOneSidedComment()
		} else {
			GetChar()
		}
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

// Locate Locates a Symbol in Table
func Locate(n string) int {
	return Lookup(ST, n, MaxEntry)
}

// InTable Looks for Symbol in Table
func InTable(n string) bool {
	return Locate(n) != -1
}

// CheckTable Checks to See if an Identifier is in the Symbol Table
// Reports an error if it's not.
func CheckTable(n string) {
	if !InTable(n) {
		Undefined(n)
	}
}

// CheckDup Checks the Symbol Table for a Duplicate Identifier
// Reports an error if identifier is already in table.
func CheckDup(n string) {
	if InTable(n) {
		Duplicate(n)
	}
}

// AddEntry Adds a New Entry to Symbol Table
func AddEntry(n string, t rune) {
	CheckDup(n)
	if NEntry == MaxEntry {
		Abort("Symbol Table Full")
	}
	NEntry++
	ST[NEntry] = n
	SType[NEntry] = t
}

// GetName Gets an Identifier
func GetName() {
	SkipWhite()
	if !IsAlpha(Look) {
		Expected("Name")
	}
	Token = 'x'
	Value = ""
	for IsAlNum(Look) {
		Value += string(unicode.ToUpper(Look))
		GetChar()
	}
}

// GetNum Gets a Number
func GetNum() {
	SkipWhite()
	if !IsDigit(Look) {
		Expected("Integer")
	}
	Token = '#'
	Value = ""
	for IsDigit(Look) {
		Value += string(Look)
		GetChar()
	}
}

// GetOp Gets an Operator
func GetOp() {
	SkipWhite()
	Token = Look
	Value = string(Look)
	GetChar()
}

// Next Gets the Next Input Token
func Next() {
	SkipWhite()
	if IsAlpha(Look) {
		GetName()
	} else if IsDigit(Look) {
		GetNum()
	} else {
		GetOp()
	}
}

// Scan Gets an Identifier and Scans it for Keywords
func Scan() {
	if Token == 'x' {
		Token = rune(KWCode[Lookup(KWList, Value, NKW)+1])
	}
}

// MatchString Matches a Specific Input String
func MatchString(x string) {
	if Value != x {
		Expected(x)
	}
	Next()
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
func LoadConst(n string) {
	Emit("MOVE #")
	util.WriteLine(n + ",D0")
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

// ReadIt Reads Variable to Primary Register
func ReadIt(name string) {
	EmitLn("BSR READ")
	Store(name)
}

// WriteIt Writes from Primary Register
func WriteIt() {
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

// Allocate Allocates Storage for a Static Variable
func Allocate(name string, val string) {
	util.WriteLine(name + ":\tDC " + val)
}

// Factor Parses and Translates a Math Factor
func Factor() {
	if Token == '(' {
		Next()
		BoolExpression()
		MatchString(")")
	} else {
		if Token == 'x' {
			LoadVar(Value)
		} else if Token == '#' {
			LoadConst(Value)
		} else {
			Expected("Math Factor")
		}
		Next()
	}
}

// Multiply Recognizes and Translates a Multiply
func Multiply() {
	Next()
	Factor()
	PopMul()
}

// Divide Recognizes and Translates a Divide
func Divide() {
	Next()
	Factor()
	PopDiv()
}

// Term Parses and Translates a Math Term
func Term() {
	Factor()
	for IsMulOp(Token) {
		Push()
		switch Token {
		case '*':
			Multiply()
		case '/':
			Divide()
		}
	}
}

// Add Recognizes and Translates an Add
func Add() {
	Next()
	Term()
	PopAdd()
}

// Subtract Recognizes and Translates a Subtract
func Subtract() {
	Next()
	Term()
	PopSub()
}

// Expression Parses and Translates a Math Expression
func Expression() {
	if IsAddOp(Token) {
		Clear()
	} else {
		Term()
	}
	for IsAddOp(Token) {
		Push()
		switch Token {
		case '+':
			Add()
		case '-':
			Subtract()
		}
	}
}

// CompareExpression Gets Another Expression and Compares
func CompareExpression() {
	Expression()
	PopCompare()
}

// NextExpression Gets the Next Expression and  Compares
func NextExpression() {
	Next()
	CompareExpression()
}

// Equals Recognizes and Translates a Relational "Equals"
func Equals() {
	NextExpression()
	SetEqual()
}

// LessOrEqual Recognizes and Translates a Relational "Less Than or Equal"
func LessOrEqual() {
	NextExpression()
	SetLessOrEqual()
}

// NotEqual Recognizes and Translates a Relational "Not Equals"
func NotEqual() {
	NextExpression()
	SetNEqual()
}

// Less Recognizes and Translates a Relational "Less Than"
func Less() {
	Next()
	switch Token {
	case '=':
		LessOrEqual()
	case '>':
		NotEqual()
	default:
		CompareExpression()
		SetLess()
	}
}

// Greater Recognizes and Translates a Relational "Greater Than"
func Greater() {
	Next()
	if Token == '=' {
		NextExpression()
		SetGreaterOrEqual()
	} else {
		CompareExpression()
		SetGreater()
	}
}

// Relation Parses and  Translates a Relation
func Relation() {
	Expression()
	if IsRelOp(Token) {
		Push()
		switch Token {
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
	if Token == '!' {
		Next()
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
		Next()
		NotFactor()
		PopAdd()
	}
}

// BoolOr Recognizes and Translates a Boolean OR
func BoolOr() {
	Next()
	BoolTerm()
	PopOr()
}

// BoolXor Recognizes and Translates an Exclusive OR
func BoolXor() {
	Next()
	BoolTerm()
	PopXor()
}

// BoolExpression Parses and Translates a Boolean Expression
func BoolExpression() {
	BoolTerm()
	for IsOrOp(Token) {
		Push()
		switch Token {
		case '|':
			BoolOr()
		case '~':
			BoolXor()
		}
	}
}

// Assignment Parses and Translates an Assignment Statement
func Assignment() {
	CheckTable(Value)
	name := Value
	Next()
	MatchString("=")
	BoolExpression()
	Store(name)
}

// DoIf Recognizes and Translates an IF Construct
func DoIf() {
	Next()
	BoolExpression()
	l1 := NewLabel()
	l2 := l1
	BranchFalse(l1)
	Block()
	if Token == 'l' {
		Next()
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
	Next()
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

// ReadVar Reads Variable to Primary Register
func ReadVar() {
	CheckIdent()
	CheckTable(Value)
	ReadIt(Value)
	Next()
}

// DoRead Processes a Read Statement
func DoRead() {
	Next()
	MatchString("(")
	ReadVar()
	for Token == ',' {
		Next()
		ReadVar()
	}
	MatchString(")")
}

// DoWrite Processes a Write Statement
func DoWrite() {
	Next()
	MatchString("(")
	Expression()
	WriteIt()
	for Token == ',' {
		Next()
		Expression()
		WriteIt()
	}
	MatchString(")")
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
		Semi()
		Scan()
	}
}

// Alloc Allocates Storage for a Variable
func Alloc() {
	Next()
	if Token != 'x' {
		Expected("Variable Name")
	}
	CheckDup(Value)
	AddEntry(Value, 'v')
	Allocate(Value, "0")
	Next()
}

// TopDecls Parses and Translates Global Declarations
func TopDecls() {
	Scan()
	for Token == 'v' {
		Alloc()
		for Token == ',' {
			Alloc()
		}
		Semi()
	}
}

// Semi Matches a semicolon
func Semi() {
	MatchString(";")
	/*
		if Token == ';' {
			Next()
		}
	*/
}

// SkipComment Skips a Comment Field
func SkipComment() {
	for Look != '/' {
		for Look != '*' {
			GetCharX()
			// Note: Tutorial suggests that nested C-style comments only need 1 line
			// of code change in SkipComment, but I could only get it to work by
			// testing for both / and * before recursing
			if Look == '/' {
				GetCharX()
				if Look == '*' {
					SkipComment()
				}
			}
		}
		GetCharX()
	}
	GetCharX()
}

// SkipOneSidedComment Skips a One Sided Comment Field
func SkipOneSidedComment() {
	for Look != 0x0D {
		GetCharX()
	}
	GetChar()
}

// GetChar Reads New Character. Intercepts '/*'
func GetChar() {
	if TempChar != ' ' {
		Look = TempChar
		TempChar = ' '
	} else {
		GetCharX()
		if Look == '/' {
			TempChar = util.Read()
			if TempChar == '*' {
				Look = 0xFF
				TempChar = ' '
			} else if TempChar == '/' {
				Look = 0xFE
				TempChar = ' '
			}
		}
	}
}

// Init Initializes
func Init() {
	ST = make([]string, MaxEntry)
	SType = make([]rune, MaxEntry)
	GetChar()
	Next()
}

// Go starts the execution of this chapter
func Go() {
	Init()
	MatchString("PROGRAM")
	Header()
	TopDecls()
	MatchString("BEGIN")
	Prolog()
	Block()
	MatchString("END")
	Epilog()
}
