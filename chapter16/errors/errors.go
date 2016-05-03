package errors

import "github.com/dcw303/crenshaw-go/util"

// Error Writes Error Message and Halts
func Error(s string) {
	util.WriteBlankLine()
	util.WriteLine("Error: " + s + ".")
	panic("Halted")
}

// Expected Writes "<something> Expected"
func Expected(s string) {
	Error(s + " Expected")
}
