package codegen

import "github.com/dcw303/crenshaw-go/chapter16/output"

// LoadConstant Loads the Primary Register with a Constant
func LoadConstant(n string) {
	output.EmitLn("MOVE #" + n + ",D0")
}

// LoadVariable Loads a Variable to the Primary Register
func LoadVariable(name string) {
	output.EmitLn("MOVE " + name + "(PC),D0")
}

// Negate Negates Primary
func Negate() {
	output.EmitLn("NEG D0")
}

// Push Pushes Primary to Stack
func Push() {
	output.EmitLn("MOVE D0,-(SP)")
}

// PopAdd Adds TOS to Primary
func PopAdd() {
	output.EmitLn("ADD (SP)+,D0")
}

// PopSub Subtracts TOS from Primary
func PopSub() {
	output.EmitLn("SUB (SP)+,D0")
}

// PopMul Multiples TOS by Primary
func PopMul() {
	output.EmitLn("MULS (SP)+,D0")
}

// PopDiv Divides Primary by TOS
func PopDiv() {
	output.EmitLn("MOVE (SP)+,D7")
	output.EmitLn("EXT.L D7")
	output.EmitLn("DIVS D0,D7")
	output.EmitLn("MOVE D7,D0")
}

// StoreVariable Stores the Primary Register to a Variable
func StoreVariable(name string) {
	output.EmitLn("LEA " + name + "(PC),A0")
	output.EmitLn("MOVE D0,(A0)")
}

// PopOr Ors TOS with Primary
func PopOr() {
	output.EmitLn("OR (SP)+,D0")
}

// PopXor Exclusive-Ors TOS with Primary
func PopXor() {
	output.EmitLn("EOR (SP)+,D0")
}

// PopAnd Ands Primary with TOS
func PopAnd() {
	output.EmitLn("AND (SP)+,D0")
}

// NotIt Bitwise Nots Primary
func NotIt() {
	output.EmitLn("EOR #-1,D0")
}
