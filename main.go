package main

import (
	"github.com/dcw303/crenshaw-go/chapter16"
	"github.com/dcw303/crenshaw-go/util"
	"github.com/nsf/termbox-go"
)

//1. Import the chapter you want to run

//2. Execute the Go() func on the package for the chapter. This is:

// chapter 01: cradle
// chapter 02/03/06/06b/09/09b: parse
// chapter 04: interpret
// chapter 05: branch
// chapter 07/07b: kiss
// chapter 10/11/12/12b: tiny
// chapter 13/13b: calls
// chapter 14: types
// chapter 15/16: test

func main() {
	defer termbox.Close()
	defer closeLoop()
	test.Go()
}

func closeLoop() {
	util.WriteLine("*** Execution Complete - Hit <Enter> to exit ***")
	for r := util.Read(); r != 0x0D; r = util.Read() {
	}

}
