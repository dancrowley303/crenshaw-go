package parser

import (
	"github.com/dcw303/crenshaw-go/chapter16/codegen"
	"github.com/dcw303/crenshaw-go/chapter16/errors"
	"github.com/dcw303/crenshaw-go/chapter16/input"
	"github.com/dcw303/crenshaw-go/chapter16/scanner"
)

// Factor Parses and Translates a Factor
func Factor() {
	if input.Look == '(' {
		scanner.Match('(')
		Expression()
		scanner.Match(')')
	} else if scanner.IsDigit(input.Look) {
		codegen.LoadConstant(scanner.GetNumber())
	} else if scanner.IsAlpha(input.Look) {
		codegen.LoadVariable(scanner.GetName())
	} else {
		errors.Error("Unrecognized character " + string(input.Look))
	}
}

// SignedTerm Parses and Translates a Term with Optonal Leading SignedTerm
func SignedTerm() {
	sign := input.Look
	if scanner.IsAddOp(input.Look) {
		input.GetChar()
	}
	Term()
	if sign == '-' {
		codegen.Negate()
	}
}

// Expression Parses and Translates an Expression
func Expression() {
	SignedTerm()
	for scanner.IsAddOp(input.Look) {
		switch input.Look {
		case '+':
			Add()
		case '-':
			Subtract()
		case '|':
			Or()
		case '~':
			Xor()
		}
	}
}

// Add Parses and Translates an Addition Operator
func Add() {
	scanner.Match('+')
	codegen.Push()
	Term()
	codegen.PopAdd()
}

// Subtract Parses and Translates a Subtraction Operation
func Subtract() {
	scanner.Match('-')
	codegen.Push()
	Term()
	codegen.PopSub()
}

// Term Parses and Translates a Term
func Term() {
	NotFactor()
	for scanner.IsMulOp(input.Look) {
		switch input.Look {
		case '*':
			Multiply()
		case '/':
			Divide()
		case '&':
			And()
		}
	}
}

// Multiply Parses and Translates a Multiplication Operation
// Note this function is not documented in the tutorial but assumed to work
// the same as add/subtract
func Multiply() {
	scanner.Match('*')
	codegen.Push()
	NotFactor()
	codegen.PopMul()
}

// Divide Parses and Translates a Division Operation
// Note this function is not documented in the tutorial but assumed to work
// the same as add/subtract
func Divide() {
	scanner.Match('/')
	codegen.Push()
	NotFactor()
	codegen.PopDiv()
}

// Assignment Parses and Translates an Assignment Statement
func Assignment() {
	name := scanner.GetName()
	scanner.Match('=')
	Expression()
	codegen.StoreVariable(name)
}

// Or Parses and Translates a Boolean Or Operation
func Or() {
	scanner.Match('|')
	codegen.Push()
	Term()
	codegen.PopOr()
}

// Xor Pars and Translates a Boolean Xor Operation
func Xor() {
	scanner.Match('~')
	codegen.Push()
	Term()
	codegen.PopXor()
}

// And Parses and Translates a Boolean And Operation
func And() {
	scanner.Match('&')
	codegen.Push()
	NotFactor()
	codegen.PopAnd()
}

// NotFactor Parses and Translates a Factor with Optional "NOT"
func NotFactor() {
	if input.Look == '!' {
		scanner.Match('!')
		Factor()
		codegen.NotIt()
	} else {
		Factor()
	}
}
