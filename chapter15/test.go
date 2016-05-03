package test

//import (
//"github.com/dcw303/crenshaw-go/chapter15/input"
//"github.com/dcw303/crenshaw-go/chapter15/output"
//"github.com/dcw303/crenshaw-go/chapter15/errors"
//"github.com/dcw303/crenshaw-go/chapter15/scanner1"
//)

// Go is equivalent to the program Test / program Main entry point defined in
// the tutorial
func Go() {

	// This Go function departs from previous chapters. As there is a test for
	// each of the turbo pascal units being written, these are coded and commented
	// out. Uncomment a section, and then read the comments to run each test.

	// Early on, there is a test of the input/ouput packages that runs the
	// following. You'll need to add imports for input/output packages.
	// Also, there is some discussion about difficulties pretty-printing the
	// labels with the Turbo Pascal library WinCRT; this is not a problem in Go.

	/*
		util.Write("MAIN:\t")
		output.EmitLn("Hello, world!")
		util.WriteLine(string(input.Look))
	*/

	// Then, there is a test program for the Errors unit. You'll need to add
	//just the errors import for that

	/*
		errors.Expected("Integer")
	*/

	// Next up is a test for the single character scanner unit, scanner1.
	// Import package as expected, and also fmt to do the Printlns

	/*
		util.Write(string(scanner1.GetName()))
		scanner1.Match('=')
		util.Write(string(scanner1.GetNumber()))
		scanner1.Match('+')
		util.WriteLine(string(scanner1.GetName()))
	*/

	// This test is for the multi character scanner unit.
	// The same package includes as for scanner1 are required

	/*
		util.Write(scanner.GetName())
		scanner.Match('=')
		util.Write(scanner.GetNumber())
		scanner.Match('+')
		util.WriteLine(scanner.GetName())
	*/

	// The penultimate and final test both test parsing a factor. The initial
	// test can only test a constant, whereas the final can also test a variable.
	// As the factor func is modified for the final test, I've skipped the
	// constant-only version.

	/*
		parser.Factor()
	*/
}
