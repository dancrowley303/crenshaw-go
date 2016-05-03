package parser

import (
	"github.com/dcw303/crenshaw-go/chapter15/codegen"
	"github.com/dcw303/crenshaw-go/chapter15/errors"
	"github.com/dcw303/crenshaw-go/chapter15/input"
	"github.com/dcw303/crenshaw-go/chapter15/scanner"
)

// Factor Parses and Translates a Factor
func Factor() {
	if scanner.IsDigit(input.Look) {
		codegen.LoadConstant(scanner.GetNumber())
	} else if scanner.IsAlpha(input.Look) {
		codegen.LoadVariable(scanner.GetName())
	} else {
		errors.Error("Unrecognized character " + string(input.Look))
	}
}
