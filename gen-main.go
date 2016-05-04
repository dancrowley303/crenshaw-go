//+build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

const tmpl = `//+build ignore

package main

import (
	pkg "github.com/dcw303/crenshaw-go/%s"
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
`

func main() {
	for _, chapter := range []string{
		"chapter01",
		"chapter02",
		"chapter03",
		"chapter04",
		"chapter05",
		"chapter06",
		"chapter06b",
		"chapter07",
		"chapter07b",
		"chapter09",
		"chapter09b",
		"chapter10",
		"chapter11",
		"chapter12",
		"chapter12b",
		"chapter13",
		"chapter13b",
		"chapter14",
		"chapter15",
		"chapter16",
	} {
		fname := filepath.Join(chapter, "main.go")

		err := ioutil.WriteFile(fname, []byte(fmt.Sprintf(tmpl, chapter)), 0644)

		if err != nil {
			log.Fatal(err)
		}
	}
}
