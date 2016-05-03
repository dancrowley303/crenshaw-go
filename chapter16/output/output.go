package output

import "github.com/dcw303/crenshaw-go/util"

// Emit Emits an Instruction
func Emit(s string) {
	util.Write("\t" + s)
}

// EmitLn Emits an Instruction, Followed By a Newline
func EmitLn(s string) {
	Emit(s)
	util.WriteBlankLine()
}
