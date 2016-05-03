package codegen

import "github.com/dcw303/crenshaw-go/chapter15/output"

// LoadConstant Loads the Primary Register with a Constant
func LoadConstant(n string) {
	output.EmitLn("MOVE #" + n + ",D0")
}

// LoadVariable Loads a Variable to the Primary Register
func LoadVariable(name string) {
	output.EmitLn("MOVE " + name + "(PC),D0")
}
