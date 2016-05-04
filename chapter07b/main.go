//+build ignore

package main

import (
	pkg "github.com/dcw303/crenshaw-go/chapter07b"
	"github.com/dcw303/crenshaw-go/util"
	"github.com/nsf/termbox-go"
)

func main() {
	defer termbox.Close()
	defer closeLoop()
	pkg.Go()
}

func closeLoop() {
	util.WriteLine("*** Execution Complete - Hit <Enter> to exit ***")
	for r := util.Read(); r != 0x0D; r = util.Read() {
	}

}
